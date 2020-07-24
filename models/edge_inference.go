package models

type FDInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Url      string `json:"url"`
}

type ConvertionTypes struct {
	ConvertionTypes []string `json:"conversionTypes"`
}
