package routers

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			logger.Info("User not found in session, check token")
			fmt.Println("<<<<<<<<<<<<<<<<<<<<")
			fmt.Println(c.Request.Header.Get("token"))
			fmt.Println(c.Request.Header)
			fmt.Println("<<<<<<<<<<<<<<<<<<<<")
		}
	}
}
