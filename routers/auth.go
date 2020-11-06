package routers

import (
	"net/http"
	"strings"
	"time"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/crewjam/saml/samlsp"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

		tokenEmpty, samlEmpty := true, true
		if len(auth) > 0 {
			tokenEmpty = false
		}
		if openSaml {
			if s, _ := samlValidator.Session.GetSession(r); s != nil {
				samlEmpty = false
			}
		}

		if tokenEmpty && samlEmpty {
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
				// TODO expiration has been detected in parseToken actually, why do it again
				if time.Now().Unix() > claim.ExpiresAt {
					c.Abort()
					c.JSON(http.StatusUnauthorized, UnAuthorizedError("Token expired"))
				}
				c.Set("uid", claim.Uid)
				c.Set("userName", claim.UserName)
				c.Set("userId", claim.Uid)
			}
		} else if openSaml {
			samlSession, _ := samlValidator.Session.GetSession(r)
			r = r.WithContext(samlsp.ContextWithSession(r.Context(), samlSession))
			sa, _ := samlSession.(samlsp.SessionWithAttributes)

			claim := services.ExtractSamlAttrs(sa.GetAttributes())
			c.Set("uid", claim["uid"])
			c.Set("userName", claim["userName"])
			c.Set("userId", claim["userId"])
		}

		c.Next()
	}
}
