package model

import (
	"time"

	"gorm.io/datatypes"
)

// AboutProfile 是单例表（id 恒为 1），存「关于我」页面的可编辑内容。
// Skills / Timeline / Socials 用 JSON 存数组，后台以结构化表单编辑。
type AboutProfile struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	AvatarURL string         `gorm:"size:500" json:"avatar_url"`
	IntroMD   string         `gorm:"type:longtext" json:"intro_md"`
	Skills    datatypes.JSON `gorm:"type:json" json:"skills"`
	Timeline  datatypes.JSON `gorm:"type:json" json:"timeline"`
	Socials   datatypes.JSON `gorm:"type:json" json:"socials"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (AboutProfile) TableName() string { return "about_profile" }
