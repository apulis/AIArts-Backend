package routers

const (
	SUCCESS_CODE = 0

	NOT_FOUND_ERROR_CODE = 10001
	UNKNOWN_ERROR_CODE   = 10002
	SERVER_ERROR_CODE    = 10003

	// Request error codes
	PARAMETER_ERROR_CODE = 20001
	AUTH_ERROR_CODE      = 20002

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

	//上传大文件template目录满
	UPLOAD_TEMPDIR_FULL_COD = 30009
	// user auth
	NO_USRNAME = 30101

	// training
	INVALID_CODE_PATH = 30201

	// upgrade platform
	UPGRADE_FILE_PATH_DO_NOT_EXIST = 30501
	FILE_READ_ERROR                = 30502

	REMOTE_SERVE_ERROR_CODE = 40000
)
