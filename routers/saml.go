package routers

import (
	"github.com/crewjam/saml/samlsp"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

var (
	samlValidator, _ = samlsp.New(samlsp.Options{
		EntityID:          "",
		URL:               url.URL{},
		Key:               nil,
		Certificate:       nil,
		Intermediates:     nil,
		AllowIDPInitiated: false,
		IDPMetadata:       nil,
		SignRequest:       false,
		ForceAuthn:        false,
		CookieSameSite:    -1,
	})
)

func loginSuccess(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusOK)
}

func AddSamlInterface(r *gin.Engine) {
	g := r.Group("/ai_arts/")

	g.GET("/saml_login", func(c *gin.Context) {
		samlValidator.RequireAccount(http.HandlerFunc(loginSuccess))
	})
	g.GET("/saml", func(c *gin.Context) {
		samlValidator.ServeHTTP(c.Writer, c.Request)
	})
}
