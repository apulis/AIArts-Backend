package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error

type APIException struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Data       gin.H  `json:"data"`
}

func wrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			var apiException *APIException
			if h, ok := err.(*APIException); ok {
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

func (e *APIException) Error() string {
	return e.Msg
}

func newAPIException(statusCode, code int, msg string) *APIException {
	return &APIException{
		StatusCode: statusCode,
		Code:       code,
		Msg:        msg,
	}
}

func ServerError() *APIException {
	return newAPIException(http.StatusInternalServerError, SERVER_ERROR, http.StatusText(http.StatusInternalServerError))
}

func NotFound() *APIException {
	return newAPIException(http.StatusNotFound, NOT_FOUND, http.StatusText(http.StatusNotFound))
}

func UnknownError(msg string) *APIException {
	return newAPIException(http.StatusForbidden, UNKNOWN_ERROR, msg)
}

func ParameterError(msg string) *APIException {
	return newAPIException(http.StatusBadRequest, PARAMETER_ERROR, msg)
}

func HandleNotFound(c *gin.Context) {
	handleErr := NotFound()
	c.JSON(handleErr.StatusCode, handleErr)
	return
}
