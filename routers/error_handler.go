package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/apulis/AIArtsBackend/configs"
)

type HandlerFunc func(c *gin.Context) error



func wrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			logger.Error(err.Error())
			var apiException *configs.APIException
			if h, ok := err.(*configs.APIException); ok {
				apiException = h
			} else if e, ok := err.(error); ok {
				if gin.Mode() == "debug" {
					apiException = UnknownError(e.Error())
				} else {
					apiException = UnknownError(e.Error())
				}
			} else {
				apiException = ServerError()
			}
			c.JSON(apiException.StatusCode, apiException)
			return
		}
	}
}



func UnAuthorizedError(msg string) *configs.APIException {
	return configs.NewAPIException(http.StatusUnauthorized, configs.AUTH_ERROR_CODE, msg)
}

func ServerError() *configs.APIException {
	return configs.NewAPIException(http.StatusInternalServerError, configs.SERVER_ERROR_CODE, http.StatusText(http.StatusInternalServerError))
}

func NotFound() *configs.APIException {
	return configs.NewAPIException(http.StatusNotFound, configs.NOT_FOUND_ERROR_CODE, http.StatusText(http.StatusNotFound))
}

func UnknownError(msg string) *configs.APIException {
	return configs.NewAPIException(http.StatusForbidden, configs.UNKNOWN_ERROR_CODE, msg)
}
func RemoteServerError(msg string)*configs.APIException{
	return ServeError(configs.REMOTE_SERVE_ERROR_CODE, msg)
}
func ParameterError(msg string) *configs.APIException {
	return configs.NewAPIException(http.StatusBadRequest, configs.PARAMETER_ERROR_CODE, msg)
}

func AppError(errorCode int, msg string) *configs.APIException {
	return configs.NewAPIException(http.StatusBadRequest, errorCode, msg)
}

func ServeError(errorCode int, msg string) *configs.APIException {
	return configs.NewAPIException(http.StatusInternalServerError, errorCode, msg)
}



func HandleNotFound(c *gin.Context) {
	handleErr := NotFound()
	c.JSON(handleErr.StatusCode, handleErr)
	return
}
