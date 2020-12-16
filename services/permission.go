package services

import (
	"github.com/apulis/AIArtsBackend/configs"
)

type PermissionListResult struct {
	PermissionList []string `json:"permissionList"`
}

func CanSubmitPrivilegedJob(token string, bypassCode string) (bool, error) {
	hasPermission, err := HasPermission(token, "SUBMIT_PRIVILEGE_JOB")
	if err != nil {
		return false, err
	}

	if !hasPermission {
		return false, nil
	}

	setting, err := GetPrivilegedSetting()
	if err != nil {
		return false, err
	}

	if !setting.IsEnable || setting.BypassCode == "" {
		return false, nil
	}

	canSubmit := (bypassCode == setting.BypassCode)

	return canSubmit, nil
}

func HasPermission(token string, permission string) (bool, error) {
	permissionList, err := getCurrentPermissionList(token)
	if err != nil {
		return false, err
	}

	return contains(permissionList, permission), nil
}

func getCurrentPermissionList(token string) ([]string, error) {
	var resp PermissionListResult
	requestUrl := configs.Config.Auth.Url + "/custom-user-dashboard-backend/auth/currentUser"

	headers := map[string]string{
		"Authorization": token,
		"Content-Type":  "application/json",
	}

	if err := DoRequest(requestUrl, "GET", headers, nil, &resp); err != nil {
		return nil, err
	}

	return resp.PermissionList, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
