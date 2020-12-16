package configs

const (
	SUCCESS_CODE = 0

	NOT_FOUND_ERROR_CODE = 10001
	UNKNOWN_ERROR_CODE   = 10002
	SERVER_ERROR_CODE    = 10003

	// Request error codes
	PARAMETER_ERROR_CODE    = 20001
	AUTH_ERROR_CODE         = 20002
	NAME_ALREADY_EXIST_CODE = 20003
	NOT_IMPLEMENT_CODE      = 20004

	// APP error codes
	APP_ERROR_CODE              = 30000
	FILETYPE_NOT_SUPPORTED_CODE = 30001
	SAVE_FILE_ERROR_CODE        = 30002
	EXTRACT_FILE_ERROR_CODE     = 30003
	REMOVE_FILE_ERROR_CODE      = 30004
	FILEPATH_NOT_EXISTS_CODE    = 30005
	FILE_OVERSIZE_CODE          = 30006
	FILEPATH_NOT_VALID_CODE     = 30007
	COMPRESS_PATH_ERROR_CODE    = 30008

	// dataset
	//上传大文件template目录满
	UPLOAD_TEMPDIR_FULL_CODE = 30009
	//无法删除正在使用的数据集
	DATASET_IS_STILL_USE_CODE = 30010
	//已经存在同名的数据集
	DATASET_IS_EXISTED = 30012

	// 数据集文件不存在
	DATASE_NOT_FOUND = 30013
	DATASE_MOVE_FAIL = 30014

	// modelset
	//创建评估训练任务出错
	CREATE_EVALUATION_FAILED_CODE = 30701

	// user auth
	NO_USRNAME = 30101

	// training
	INVALID_CODE_PATH       = 30201
	TEMPLATE_INVALID_PARAMS = 30202
	INVALID_TRAINING_TYPE   = 30203
	FAILED_START_TRAINING   = 30204
	DOCKER_IMAGE_NOT_FOUNT  = 30205

	// code
	COMPLETE_UPLOAD_ERR = 30301

	// upgrade platform
	UPGRADE_FILE_PATH_DO_NOT_EXIST = 30501
	FILE_READ_ERROR                = 30502
	ALREADY_UPGRADING_CODE         = 30503

	// edge inference
	FDINFO_NOT_SET     = 30601
	FDINFO_SET_ERROR   = 30602
	FD_PUSH_ERROR_CODE = 30603

	// vc
	VC_ERROR = 30701

	// permission
	OPERATION_FORBIDDEN = 30801

	REMOTE_SERVE_ERROR_CODE = 40000
)

type APIException struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

func (e *APIException) Error() string {
	return e.Msg
}

func NewAPIException(statusCode, code int, msg string) *APIException {
	return &APIException{
		StatusCode: statusCode,
		Code:       code,
		Msg:        msg,
	}
}
