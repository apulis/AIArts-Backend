package router

import (
	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code int
	Msg  string
	Data gin.H
}

func NewRouter() *gin.Engine {
	r := gin.Default()

	AddGroupCode(r)
	AddGroupModel(r)
	AddGroupTraining(r)
	AddGroupDataset(r)
	AddGroupAnnotation(r)
	AddGroupInference(r)

	return r
}
