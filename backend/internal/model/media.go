package model

import "time"

// Image 图片库中的一张图。path 是存储层内部路径，url 是对外访问地址。
type Image struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Filename     string    `gorm:"size:255;not null" json:"filename"`
	OriginalName string    `gorm:"size:255" json:"original_name"`
	Path         string    `gorm:"size:500;not null" json:"path"`
	URL          string    `gorm:"size:500;not null" json:"url"`
	Mime         string    `gorm:"size:100" json:"mime"`
	Size         int64     `json:"size"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Image) TableName() string { return "images" }
