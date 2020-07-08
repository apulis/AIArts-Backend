package models

type CodesetItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"   binding:"required"`
	Status      string `json:"status" binding:"required"`
	Image       string `json:"image"  binding:"required"`
	Creator     string `json:"creator" binding:"required"`
	CodePath 	string `json:"code_path" binding:"required"`
	CreateTime  string `json:"string" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type AIFrameworkItem struct {
	Name 		string `json:"name"`
	Image 		string `json:"image"`
}

type DeviceItem struct {

}

