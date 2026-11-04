package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 项目状态。
const (
	ProjectStatusDraft     = "draft"
	ProjectStatusPublished = "published"
)

// Project 作品集项目。
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Slug        string         `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	Summary     string         `gorm:"size:500" json:"summary"`
	ContentMD   string         `gorm:"type:longtext" json:"content_md"`
	ContentHTML string         `gorm:"type:longtext" json:"content_html,omitempty"`
	CoverURL    string         `gorm:"size:500" json:"cover_url"`
	TechStack   datatypes.JSON `gorm:"type:json" json:"tech_stack"`
	RepoURL     string         `gorm:"size:500" json:"repo_url"`
	DemoURL     string         `gorm:"size:500" json:"demo_url"`
	Sort        int            `gorm:"default:0" json:"sort"`
	Status      string         `gorm:"size:16;not null;default:draft;index" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Project) TableName() string { return "projects" }
