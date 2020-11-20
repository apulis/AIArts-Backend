package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/apulis/AIArtsBackend/utils"
	"io/ioutil"
	"strings"
	"time"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"golang.org/x/crypto/ssh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	RecordInterval = 86400 * 24 * time.Second
)

var (
	saveElapsed = int64(0)
	saveSize    = uint64(0)
	// only statistic the latest records(e.g. 24 hour)
	lastSavePoint = time.Now()
)

func ListSavedImages(page, count int, orderBy, order, imageName, username string) ([]models.SavedImage, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListSavedImages(offset, limit, orderBy, order, imageName, username)
}

func getContainerSize(hostIp, containerId string) (uint64, error) {
	// docker ps -f "id={containerId}" --format '{{ .Size }}'
	cmd := fmt.Sprintf("docker ps -f \"id=%s\" --format '{{ .Size }}'", containerId)

	// output format: 5.55MB (virtual 916MB)
	output, err := runSShCmd(hostIp, cmd)
	if err != nil {
		logger.Errorf("run \"%s\" occurs error: %s", cmd, err.Error())
		return 0, err
	}

	parts := strings.Split(output, "(virtual")
	cSize, err := utils.ToBytes(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, err
	}
	iSize, err := utils.ToBytes(strings.TrimSpace(strings.TrimRight(parts[1], ")")))
	if err != nil {
		return 0, err
	}

	return cSize + iSize, nil
}

func CreateSavedImage(name, version, description, jobId, username string, isPrivate bool) (int64, error) {
	hostIp, containerId, err := getPodStatus(jobId)
	if err != nil {
		return 0, err
	}
	logger.Info("Job " + jobId + " is running at: " + hostIp + ". Container id is: " + containerId)

	// 执行docker commit
	name = trimImageName(name)
	version = trimImageName(version)
	imageName := trimImageName(username) + "_" + name + ":" + version
	fullName := configs.Config.PrivateRegistry + imageName
	commitCmd := "docker commit " + containerId + " " + fullName
	pushCmd := "docker push " + fullName

	start := time.Now()
	// reset the statistic of save image action
	if start.Sub(lastSavePoint) > RecordInterval {
		saveElapsed = 0
		saveSize = 0
		lastSavePoint = start
	}
	size, sizeError := getContainerSize(hostIp, containerId)

	// save image by goroutine
	go func() {
		logger.Info("Running docker-commit command: ", commitCmd)
		res, err := runSShCmd(hostIp, commitCmd)
		if err != nil {
			logger.Errorf(fmt.Sprintf("Run %s error: %s", commitCmd, err.Error()))
			return
		}
		logger.Info("Run ssh command docker-commit result: ", res)

		// 执行docker push
		logger.Info("Running docker-push command: ", pushCmd)
		res, err = runSShCmd(hostIp, pushCmd)
		if err != nil {
			logger.Errorf(fmt.Sprintf("Run %s error: %s", pushCmd, err.Error()))
			return
		}
		logger.Info("Run ssh command docker-push result: ", res)

		strList := strings.Split(res, "sha256:")
		imageId := strings.TrimPrefix(strList[1], "sha256:")

		end := time.Now()
		if sizeError == nil {
			saveElapsed += end.Sub(start).Milliseconds()
			saveSize += size
		}

		savedImage := models.SavedImage{
			Name:        name,
			Version:     version,
			Description: description,
			Creator:     username,
			FullName:    imageName,
			IsPrivate:   isPrivate,
			ContaindrId: containerId[0:10],
			ImageId:     imageId[0:10],
		}

		// 删除重名镜像
		duplicatedImages, err := models.ListSavedImageByName(imageName)
		if err != nil {
			logger.Errorf("list images %s occurs error: %s", imageName, err.Error())
			return
		}
		for _, img := range duplicatedImages {
			if err := models.DeleteSavedImage(&img); err != nil {
				logger.Errorf("delete image %s occurs error: %s", img.ImageId, err.Error())
			}
		}

		if err := models.CreateSavedImage(savedImage); err != nil {
			logger.Errorf("create image %s occurs error: %s", savedImage.ImageId, err.Error())
		}
	}()

	if saveSize == 0 {
		return 0, nil
	}

	return saveElapsed * int64(size) / int64(saveSize * 1000), nil
}

func UpdateSavedImage(id int, description string) error {
	savedImage, err := models.GetSavedImage(id)
	if err != nil {
		return err
	}
	savedImage.Description = description
	return models.UpdateSavedImage(&savedImage)
}

func GetSavedImage(id int) (models.SavedImage, error) {
	return models.GetSavedImage(id)
}

func DeleteSavedImage(id int) error {
	savedImage, err := models.GetSavedImage(id)
	if err != nil {
		return err
	}
	return models.DeleteSavedImage(&savedImage)
}

func trimImageName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, ":", "_")
	return name
}

func getPodStatus(podId string) (string, string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", configs.Config.Imagesave.K8sconf)
	if err != nil {
		return "", "", err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", "", err
	}
	pod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), podId, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}

	hostIp := pod.Status.HostIP
	containers := pod.Status.ContainerStatuses
	containerId := ""
	for _, container := range containers {
		if container.Name == podId {
			containerId = strings.TrimPrefix(container.ContainerID, "docker://")
		}
	}
	if containerId == "" {
		return "", "", errors.New("container not found in pod")
	}
	return hostIp, containerId, nil
}

func runSShCmd(address, command string) (string, error) {
	user := configs.Config.Imagesave.Sshuser
	port := configs.Config.Imagesave.Sshport

	key, err := ioutil.ReadFile(configs.Config.Imagesave.Sshkey)
	if err != nil {
		return "", err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", address+":"+port, config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	sshSession, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer sshSession.Close()

	output, err := sshSession.CombinedOutput(command)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
