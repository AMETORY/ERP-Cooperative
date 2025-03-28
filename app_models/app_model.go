package app_models

type AppModel struct {
	ID      int64  `gorm:"primaryKey;" json:"id"`
	AppName string `gorm:"type:varchar(255)" json:"app_name"`
	AppKey  string `gorm:"type:varchar(255)" json:"app_key"`
	Version string `gorm:"type:varchar(255)" json:"version"`
	Build   string `gorm:"type:varchar(255)" json:"build"`
}

func (AppModel) TableName() string {
	return "apps"
}
