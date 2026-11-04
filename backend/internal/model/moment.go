package model

import (
	"time"

	"gorm.io/gorm"
)

// Moment 个人动态（朋友圈式）：一段文字 + 最多 9 张图。
type Moment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"type:text" json:"content"`
	IsPublic  bool           `gorm:"default:true" json:"is_public"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Images []MomentImage `gorm:"foreignKey:MomentID" json:"images"`
}

func (Moment) TableName() string { return "moments" }

// MomentImage 动态与图片的关联，sort 决定九宫格顺序。
type MomentImage struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	MomentID uint    `gorm:"index;not null" json:"moment_id"`
	ImageID  uint    `gorm:"index;not null" json:"image_id"`
	Sort     int     `gorm:"default:0" json:"sort"`
	Image    *Image  `gorm:"foreignKey:ImageID" json:"image,omitempty"`
}

func (MomentImage) TableName() string { return "moment_images" }
