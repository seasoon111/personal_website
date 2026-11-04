package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/model"
)

// MomentService 负责个人动态（文字 + 最多 9 图）。
type MomentService struct{ db *gorm.DB }

func NewMomentService(db *gorm.DB) *MomentService { return &MomentService{db: db} }

// PublicList 前台动态时间线。
func (s *MomentService) PublicList(q dto.Pagination) (*dto.PageResult, error) {
	q.Normalize()
	tx := s.db.Model(&model.Moment{}).Where("is_public = ?", true)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []model.Moment
	if err := tx.Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC, id ASC")
	}).Preload("Images.Image").
		Order("created_at DESC, id DESC").
		Offset(q.Offset()).Limit(q.PageSize).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// AdminList 后台动态列表（含不公开的）。
func (s *MomentService) AdminList(q dto.Pagination) (*dto.PageResult, error) {
	q.Normalize()
	var total int64
	if err := s.db.Model(&model.Moment{}).Count(&total).Error; err != nil {
		return nil, err
	}
	var items []model.Moment
	if err := s.db.Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC, id ASC")
	}).Preload("Images.Image").
		Order("created_at DESC, id DESC").
		Offset(q.Offset()).Limit(q.PageSize).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *MomentService) GetByID(id uint) (*model.Moment, error) {
	var m model.Moment
	err := s.db.Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC, id ASC")
	}).Preload("Images.Image").First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *MomentService) Create(in dto.MomentInput) (*model.Moment, error) {
	if err := validateMoment(in); err != nil {
		return nil, err
	}
	isPublic := true
	if in.IsPublic != nil {
		isPublic = *in.IsPublic
	}
	m := model.Moment{Content: strings.TrimSpace(in.Content), IsPublic: isPublic}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&m).Error; err != nil {
			return err
		}
		return s.replaceImages(tx, m.ID, in.ImageIDs)
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(m.ID)
}

func (s *MomentService) Update(id uint, in dto.MomentInput) (*model.Moment, error) {
	if err := validateMoment(in); err != nil {
		return nil, err
	}
	m, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	m.Content = strings.TrimSpace(in.Content)
	if in.IsPublic != nil {
		m.IsPublic = *in.IsPublic
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		return s.replaceImages(tx, m.ID, in.ImageIDs)
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s *MomentService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&model.Moment{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrNotFound
		}
		return tx.Where("moment_id = ?", id).Delete(&model.MomentImage{}).Error
	})
}

func (s *MomentService) replaceImages(tx *gorm.DB, momentID uint, imageIDs []uint) error {
	if err := tx.Where("moment_id = ?", momentID).Delete(&model.MomentImage{}).Error; err != nil {
		return err
	}
	if len(imageIDs) == 0 {
		return nil
	}
	rows := make([]model.MomentImage, 0, len(imageIDs))
	for i, id := range imageIDs {
		rows = append(rows, model.MomentImage{MomentID: momentID, ImageID: id, Sort: i})
	}
	return tx.Create(&rows).Error
}

func validateMoment(in dto.MomentInput) error {
	if strings.TrimSpace(in.Content) == "" && len(in.ImageIDs) == 0 {
		return ErrBadInput
	}
	if len(in.ImageIDs) > 9 {
		return ErrBadInput
	}
	return nil
}
