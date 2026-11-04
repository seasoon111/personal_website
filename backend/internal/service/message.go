package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/model"
)

// MessageService 负责联系表单留言。
type MessageService struct{ db *gorm.DB }

func NewMessageService(db *gorm.DB) *MessageService { return &MessageService{db: db} }

func (s *MessageService) Create(in dto.ContactInput) (*model.ContactMessage, error) {
	m := model.ContactMessage{
		Name:    strings.TrimSpace(in.Name),
		Email:   strings.TrimSpace(in.Email),
		Subject: strings.TrimSpace(in.Subject),
		Content: strings.TrimSpace(in.Content),
	}
	if err := s.db.Create(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *MessageService) List(q dto.Pagination, unreadOnly bool) (*dto.PageResult, error) {
	q.Normalize()
	tx := s.db.Model(&model.ContactMessage{})
	if unreadOnly {
		tx = tx.Where("is_read = ?", false)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}
	var items []model.ContactMessage
	if err := tx.Order("created_at DESC, id DESC").Offset(q.Offset()).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, err
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *MessageService) SetRead(id uint, read bool) error {
	res := s.db.Model(&model.ContactMessage{}).Where("id = ?", id).Update("is_read", read)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *MessageService) Delete(id uint) error {
	res := s.db.Delete(&model.ContactMessage{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *MessageService) GetByID(id uint) (*model.ContactMessage, error) {
	var m model.ContactMessage
	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}
