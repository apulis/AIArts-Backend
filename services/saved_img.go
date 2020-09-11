package services

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"golang.org/x/crypto/ssh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func ListSavedImages(page, count int, orderBy, order, name, username string) ([]models.SavedImage, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListSavedImages(offset, limit, orderBy, order, name, username)
}

func CreateSavedImage(name, version, description, jobId, username string, isPrivate bool) error {
	hostIp, containerId, err := getPodStatus(jobId)
	if err != nil {
		return err
	}
	logger.Info("Job " + jobId + " is running at: " + hostIp + ". Container id is: " + containerId)

	// 执行docker commit
	name = trimImageName(name)
	version = trimImageName(version)
	imageName := configs.Config.PrivateRegistry + trimImageName(username) + name + ":" + version
	cmd := "sudo docker commit " + containerId + " " + imageName
	logger.Info("Running docker-commit command: ", cmd)
	res, err := runSShCmd(hostIp, cmd)
	logger.Info("Run ssh command docker-commit result: ", res)
	if err != nil {
		return err
	}

	// 执行docker push
	cmd = "sudo docker push " + imageName
	logger.Info("Running docker-push command: ", cmd)
	res, err = runSShCmd(hostIp, cmd)
	logger.Info("Run ssh command docker-push result: ", res)
	if err != nil {
		return err
	}
	strList := strings.Split(res, "sha256:")
	imageId := strings.TrimPrefix(strList[1], "sha256:")

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

	return models.CreateSavedImage(savedImage)
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
