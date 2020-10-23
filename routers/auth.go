package routers

import (
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/crewjam/saml/samlsp"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type Claim struct {
	jwt.StandardClaims
	Uid      int    `json:"uid"`
	UserName string `json:"userName"`
}

var JwtSecret = configs.Config.Auth.Key

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
	return func(c *gin.Context) {
		r := c.Request

		auth := r.Header.Get("Authorization")
		samlSession, err := samlValidator.Session.GetSession(r)

		if len(auth) == 0 && (samlSession == nil || err == samlsp.ErrNoSession) {
			c.Abort()
			c.JSON(http.StatusUnauthorized, UnAuthorizedError("Cannot authorize"))
			c.Next()

			return
		}

		// 1. judge authentication token
		if len(auth) > 0 {
			auth = strings.Fields(auth)[1]

			// Check token
			claim, err := parseToken(auth)
			if err != nil {
				c.Abort()
				c.JSON(http.StatusUnauthorized, UnAuthorizedError(err.Error()))
			} else {
				if time.Now().Unix() > claim.ExpiresAt {
					c.Abort()
					c.JSON(http.StatusUnauthorized, UnAuthorizedError("Token expired"))
				}
				c.Set("uid", claim.Uid)
				c.Set("userName", claim.UserName)
				c.Set("userId", claim.Uid)
			}
		} else {
			r = r.WithContext(samlsp.ContextWithSession(r.Context(), samlSession))
		}

		c.Next()
	}
}
