package models

type PrivilegedSetting struct {
	Id         string `json:"id" gorm:"primaryKey"`
	IsEnable   bool   `json:"isEnable"`
	BypassCode string `json:"bypassCode" binding:"required"`
}
