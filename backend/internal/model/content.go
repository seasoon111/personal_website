package model

import (
	"time"

	"gorm.io/gorm"
)

// Category 文章分类。
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Slug      string    `gorm:"size:64;uniqueIndex;not null" json:"slug"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Category) TableName() string { return "categories" }

// Tag 文章标签。
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Slug      string    `gorm:"size:64;uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Tag) TableName() string { return "tags" }

// 文章状态。
const (
	ArticleStatusDraft     = "draft"
	ArticleStatusPublished = "published"
)

// Article 文章。content_md 是 Markdown 原文，content_html 是渲染结果（列表页不返回）。
type Article struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Slug        string         `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	Summary     string         `gorm:"size:500" json:"summary"`
	ContentMD   string         `gorm:"type:longtext" json:"content_md"`
	ContentHTML string         `gorm:"type:longtext" json:"content_html,omitempty"`
	CoverURL    string         `gorm:"size:500" json:"cover_url"`
	CategoryID  *uint          `gorm:"index" json:"category_id"`
	Status      string         `gorm:"size:16;not null;default:draft;index" json:"status"`
	IsPinned    bool           `gorm:"default:false" json:"is_pinned"`
	PublishedAt *time.Time     `gorm:"index" json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag     `gorm:"many2many:article_tags;" json:"tags,omitempty"`
}

func (Article) TableName() string { return "articles" }
