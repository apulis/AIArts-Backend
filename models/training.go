package models


type Training struct {
	Id          string `json:"id"`
	Name        string `json:"name"   binding:"required"`
	Status      string `json:"status" binding:"required"`
	Engine      string `json:"engine"  binding:"required"`
	CodePath 	string `json:"codePath"`
	CreateTime  int64 `json:"createTime"`
	Desc		string `json:"desc"`
}

