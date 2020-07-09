package services

import (
	"github.com/apulis/AIArtsBackend/models"
	"math/rand"
	"time"
)

func GetAllTraining(page, size int) ([] *models.Training, int, int, error) {

	rand.Seed(time.Now().Unix())
	item := &models.Training{
		Name: RandStringRunes(16),
		Status: "started",
		Engine: "tf_1.15",
		CodePath: "/home/bifeng.peng/",
		CreateTime: time.Now().Unix(),
		Desc: "test test test",
	}

	codes := make([] *models.Training, 0)
	codes = append(codes, item)

	return codes, 1, 1, nil
}

func CreateTraining(name, description string, framework models.AIFrameworkItem) (string, error) {
	return RandStringRunes(16), nil
}

func DeleteTraining(id string) error {
	return nil
}
