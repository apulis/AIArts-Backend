package routers

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
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
	authConfig    = configs.Config.Auth
)

func initSamlValidator() error {
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
	if err := services.CreateSamlUser(tokenStr, map[string]interface{}{
		"openId":   data["uid"],
		"userName": data["userName"],
	}); err != nil {
		resp := map[string]interface{}{"code": -1, "msg": "save new saml user error"}
		bResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(bResp)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func getSamlUser(w http.ResponseWriter, r *http.Request)  {
	s := samlsp.SessionFromContext(r.Context())
	if s == nil {
		fmt.Println("session is nil")
	}
	sa, _ := s.(samlsp.SessionWithAttributes)

	data := services.ExtractSamlAttrs(sa.GetAttributes())
	resp := map[string]interface{}{
		"code":           0,
		"userName":       data["userName"],
		"permissionList": []string{"AI_ARTS_ALL"},
	}
	bResp, _ := json.Marshal(resp)
	logger.Info("saml user:", string(bResp))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bResp)
}

func AddSamlInterface(r *gin.Engine) {
	g := r.Group("/ai_arts/")

	g.Any("/saml/*action", gin.WrapH(samlValidator))
	g.GET("/saml_login", gin.WrapH(samlValidator.RequireAccount(http.HandlerFunc(loginSuccess))))
	g.GET("/saml_role", gin.WrapH(samlValidator.RequireAccount(http.HandlerFunc(getSamlUser))))
}
