package routers

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claim struct {
	jwt.StandardClaims
	Uid      int    `json:"uid"`
	UserName string `json:"userName"`
}

var JwtSecret string = "Sign key for JWT"

func parseToken(token string) (*Claim, error) {

	jwtToken, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(JwtSecret), nil
	})

	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*Claim); ok && jwtToken.Valid {
			return claim, nil
		}
	}

	return nil, err
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		result := APIException{
			StatusCode: http.StatusUnauthorized,
			Msg:        "无法认证，重新登录",
		}

		auth := context.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			context.Abort()
			context.JSON(http.StatusUnauthorized, gin.H{
				"result": result,
			})
		} else {
			auth = strings.Fields(auth)[1]

			// 校验token
			claim, err := parseToken(auth)
			if err != nil {
				context.Abort()
				result.Msg = "token 过期" + err.Error()
				context.JSON(http.StatusUnauthorized, gin.H{
					"result": result,
				})
			} else {
				println("token 正确: ", claim.UserName)
				context.Set("userName", claim.UserName)
			}
		}

		context.Next()
	}
}
