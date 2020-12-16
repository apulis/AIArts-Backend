package services

import (
	"github.com/apulis/AIArtsBackend/configs"
)

type PermissionListResult struct {
	PermissionList []string `json:"permissionList"`
}

func CanSubmitPrivilegedJob(token string, bypassCode string) (int, error) {
	setting, err := GetPrivilegedSetting()
	if err != nil {
		return configs.OPERATION_FORBIDDEN, err
	}

	if !setting.IsEnable || setting.BypassCode == "" {
		return configs.PRIVILEGE_JOB_NOT_ENABLE, nil
	}

	canSubmit, err := CanSubmitPrivilegeJob(token)
	if err != nil {
		return configs.OPERATION_FORBIDDEN, err
	}

	if !canSubmit {
		return configs.OPERATION_FORBIDDEN, nil
	}

	if bypassCode != setting.BypassCode {
		return configs.PRIVILEGE_JOB_CODE_INVALID, nil
	}

	return configs.SUCCESS_CODE, nil
}

func CanSubmitPrivilegeJob(token string) (bool, error) {
	return hasPermission(token, "SUBMIT_PRIVILEGE_JOB")
}

func CanManagePrivilegeJob(token string) (bool, error) {
	return hasPermission(token, "MANAGE_PRIVILEGE_JOB")
}

func hasPermission(token string, permission string) (bool, error) {
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
