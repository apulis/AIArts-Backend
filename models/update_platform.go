package models

type VersionInfoSet struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`

	Description string `gorm:"type:text" json:"description"`
	Version     string `gorm:"not null" json:"version"`
}

func UploadVersionInfoSet(versionInfo VersionInfoSet) error {
	return db.Create(&versionInfo).Error
}

func GetCurrentVersion() (VersionInfoSet, error) {
	var currentVersion VersionInfoSet
	res := db.Order("created_at desc").First(&currentVersion)
	if res.Error != nil {
		return currentVersion, res.Error
	}
	return currentVersion, nil
}

func GetVersionLogs() ([]VersionInfoSet, error) {
	var versionLog []VersionInfoSet
	res := db.Limit(10).Find(&versionLog)
	if res.Error != nil {
		return versionLog, res.Error
	}
	return versionLog, nil
}
