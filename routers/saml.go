package routers

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/crewjam/saml/samlsp"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

var (
	openSaml      = false
	samlValidator *samlsp.Middleware
)

func initSamlValidator() error {
	authConfig    := configs.Config.Auth
	if len(authConfig.SamlIdpMetadataURL) == 0 {
		return nil
	}

	keyPair, err := tls.LoadX509KeyPair(authConfig.SamlCertificate, authConfig.SamlPrivateKey)
	if err != nil {
		return err
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return err
	}

	idpMetadataURL, err := url.Parse(authConfig.SamlIdpMetadataURL)
	if err != nil {
		return err
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
	if err != nil {
		return err
	}

	rootURL, err := url.Parse(fmt.Sprintf("%s:%d", configs.Config.RootUrl, configs.Config.Port))
	if err != nil {
		return err
	}

	if samlValidator, err = samlsp.New(samlsp.Options{
		URL:               *rootURL,
		AllowIDPInitiated: true,
		Key:               keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:       keyPair.Leaf,
		IDPMetadata:       idpMetadata,
	}); err != nil {
		return err
	}
	openSaml = true

	logger.Infof("Open saml login: %s\n", authConfig.SamlIdpMetadataURL)

	return nil
}

func loginSuccess(w http.ResponseWriter, r *http.Request) {
	authConfig    := configs.Config.Auth
	logger.Info("user login by saml way")

	s := samlsp.SessionFromContext(r.Context())
	if s == nil {
		fmt.Println("session is nil")
	}
	sa, _ := s.(samlsp.SessionWithAttributes)

	data := services.ExtractSamlAttrs(sa.GetAttributes())

	claim := Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Uid:      30000,
		UserName: data["userName"].(string),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, _ := token.SignedString([]byte(authConfig.Key))
	services.CreateSamlUser(tokenStr, map[string]interface{}{
		"openId":   data["uid"],
		"userName": data["userName"],
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func AddSamlInterface(r *gin.Engine) {
	g := r.Group("/ai_arts/")

	g.Any("/saml/*action", gin.WrapH(samlValidator))

	app := samlValidator.RequireAccount(http.HandlerFunc(loginSuccess))
	g.GET("/saml_login", gin.WrapH(app))
}
