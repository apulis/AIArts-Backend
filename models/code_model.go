package models

type CodesetItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"   binding:"required"`
	Status      string `json:"status" binding:"required"`
	Engine      string `json:"engine"  binding:"required"`
	CodePath 	string `json:"codePath"`
	CodeUrl 	string `json:"codeUrl"`
	CreateTime  int64 `json:"createTime"`
	Desc		string `json:"desc"`
}

type DeviceItem struct {
	DeviceType	string 	`json:"deviceType"`
	Avail		int 	`json:"avail"`
}

