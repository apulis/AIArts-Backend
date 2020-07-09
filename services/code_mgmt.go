package services

import (
	"time"
	"math/rand"
	"github.com/apulis/AIArtsBackend/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetCodeset(page, size int) ([] *models.CodesetItem, int, int, error) {

	rand.Seed(time.Now().Unix())
	item := &models.CodesetItem{
		Name: RandStringRunes(16),
		Status: "started",
		Engine: "tf_1.15",
		CodePath: "/home/bifeng.peng/",
		CreateTime: "2020-07-01T08:46:42.075952+00:00",
		Desc: "test test test",
	}

	codes := make([] *models.CodesetItem, 0)
	codes = append(codes, item)

	return codes, 1, 1, nil
}

func CreateCodeset(name, description string, framework models.AIFrameworkItem) (string, error) {
	return RandStringRunes(16), nil
}

func DeleteCodeset(id string) error {
	return nil
}
