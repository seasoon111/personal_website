package service

import (
	"gorm.io/gorm"

	"personal-website/backend/internal/model"
)

// DashboardService 汇总后台首页统计。
type DashboardService struct{ db *gorm.DB }

func NewDashboardService(db *gorm.DB) *DashboardService { return &DashboardService{db: db} }

// Stats 后台仪表盘数据。
type Stats struct {
	Articles       int64           `json:"articles"`
	Drafts         int64           `json:"drafts"`
	Projects       int64           `json:"projects"`
	Moments        int64           `json:"moments"`
	Images         int64           `json:"images"`
	UnreadMessages int64           `json:"unread_messages"`
	RecentDrafts   []model.Article `json:"recent_drafts"`
}

func (s *DashboardService) Stats() (*Stats, error) {
	st := &Stats{}

	if err := s.db.Model(&model.Article{}).Where("status = ?", model.ArticleStatusPublished).Count(&st.Articles).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Article{}).Where("status = ?", model.ArticleStatusDraft).Count(&st.Drafts).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Project{}).Count(&st.Projects).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Moment{}).Count(&st.Moments).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Image{}).Count(&st.Images).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.ContactMessage{}).Where("is_read = ?", false).Count(&st.UnreadMessages).Error; err != nil {
		return nil, err
	}
	if err := s.db.Select(articleFields).
		Where("status = ?", model.ArticleStatusDraft).
		Order("updated_at DESC").Limit(5).
		Find(&st.RecentDrafts).Error; err != nil {
		return nil, err
	}
	return st, nil
}
