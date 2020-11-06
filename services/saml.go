package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"strings"
)

var (
	NameIDSuffix = "nameidentifier"
	NameSuffix   = "claims/name"
)

func ExtractSamlAttrs(attrs map[string][]string) map[string]interface{} {
	data := make(map[string]interface{})

	for key, value := range attrs {
		if len(value) == 0 {
			continue
		}

		if strings.HasSuffix(key, NameIDSuffix) {
			data["userId"] = value[0]
			data["uid"] = value[0]
		}
		if strings.HasSuffix(key, NameSuffix) {
			data["userName"] = value[0]
		}
	}

	logger.Infof("extract saml attributes: %s", data)

	return data
}

func CreateSamlUser(token string, data map[string]interface{}) error {
	resp := struct {
		Success bool `json:"success"`
	}{}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type": "application/json",
	}
	requestUrl := configs.Config.Auth.Url + "/custom-user-dashboard-backend/open/user"
	if err := DoRequest(requestUrl, "POST", headers, data, &resp); err != nil {
		logger.WithError(err).Errorf("%s create saml user error", requestUrl)
		return err
	}
	return nil
}
