package models

type PrivilegedSettings struct {
	Id         string `json:"id"` `gorm:"primaryKey"`
	IsEnable   bool   `json:"isEnable" binding:"required"`
	BypassCode string `json:"bypassCode" binding:"required"`
}
