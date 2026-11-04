package service

import (
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/model"
)

// AboutService 管理「关于我」单例内容。
type AboutService struct{ db *gorm.DB }

func NewAboutService(db *gorm.DB) *AboutService { return &AboutService{db: db} }

func (s *AboutService) Get() (*model.AboutProfile, error) {
	var p model.AboutProfile
	err := s.db.First(&p, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p = model.AboutProfile{ID: 1, Skills: datatypes.JSON("[]"), Timeline: datatypes.JSON("[]"), Socials: datatypes.JSON("[]")}
		if err := s.db.Create(&p).Error; err != nil {
			return nil, err
		}
		return &p, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *AboutService) Update(in dto.AboutInput) (*model.AboutProfile, error) {
	p, err := s.Get()
	if err != nil {
		return nil, err
	}
	p.AvatarURL = in.AvatarURL
	p.IntroMD = in.IntroMD
	p.Skills = marshalJSON(in.Skills)
	p.Timeline = marshalJSON(in.Timeline)
	p.Socials = marshalJSON(in.Socials)
	if err := s.db.Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func marshalJSON(v interface{}) datatypes.JSON {
	if v == nil {
		return datatypes.JSON("[]")
	}
	b, err := json.Marshal(v)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}
