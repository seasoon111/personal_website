package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/model"
	"personal-website/backend/internal/pkg/markdown"
	"personal-website/backend/internal/pkg/slug"
)

const projectFields = "projects.id, projects.title, projects.slug, projects.summary, projects.cover_url, projects.tech_stack, projects.repo_url, projects.demo_url, projects.sort, projects.status, projects.created_at, projects.updated_at"

// ProjectService 负责作品集。
type ProjectService struct{ db *gorm.DB }

func NewProjectService(db *gorm.DB) *ProjectService { return &ProjectService{db: db} }

func (s *ProjectService) PublicList() ([]model.Project, error) {
	var items []model.Project
	err := s.db.Select(projectFields).
		Where("status = ?", model.ProjectStatusPublished).
		Order("sort DESC, id DESC").
		Find(&items).Error
	return items, err
}

func (s *ProjectService) PublicDetail(slugStr string) (*model.Project, error) {
	var p model.Project
	err := s.db.Where("slug = ? AND status = ?", slugStr, model.ProjectStatusPublished).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ProjectService) AdminList(q dto.Pagination) (*dto.PageResult, error) {
	q.Normalize()
	var total int64
	if err := s.db.Model(&model.Project{}).Count(&total).Error; err != nil {
		return nil, err
	}
	var items []model.Project
	if err := s.db.Select(projectFields).Order("sort DESC, id DESC").
		Offset(q.Offset()).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, err
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *ProjectService) GetByID(id uint) (*model.Project, error) {
	var p model.Project
	if err := s.db.First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (s *ProjectService) Create(in dto.ProjectInput) (*model.Project, error) {
	p := model.Project{
		Title:     in.Title,
		Slug:      in.Slug,
		Summary:   in.Summary,
		ContentMD: in.ContentMD,
		CoverURL:  in.CoverURL,
		RepoURL:   in.RepoURL,
		DemoURL:   in.DemoURL,
		Sort:      in.Sort,
		Status:    in.Status,
	}
	if p.Status == "" {
		p.Status = model.ProjectStatusDraft
	}
	if p.Slug == "" {
		p.Slug = slug.Make(p.Title)
	}
	p.Slug = s.uniqueSlug(p.Slug, 0)
	p.ContentHTML = markdown.Render(p.ContentMD)
	p.TechStack = techStackJSON(in.TechStack)
	if err := s.db.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *ProjectService) Update(id uint, in dto.ProjectInput) (*model.Project, error) {
	p, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if in.Slug == "" {
		in.Slug = slug.Make(in.Title)
	}
	if in.Status == "" {
		in.Status = model.ProjectStatusDraft
	}

	p.Title = in.Title
	p.Slug = s.uniqueSlug(in.Slug, id)
	p.Summary = in.Summary
	p.ContentMD = in.ContentMD
	p.ContentHTML = markdown.Render(in.ContentMD)
	p.CoverURL = in.CoverURL
	p.RepoURL = in.RepoURL
	p.DemoURL = in.DemoURL
	p.Sort = in.Sort
	p.Status = in.Status
	p.TechStack = techStackJSON(in.TechStack)

	if err := s.db.Save(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) Delete(id uint) error {
	res := s.db.Delete(&model.Project{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func techStackJSON(items []string) datatypes.JSON {
	if items == nil {
		items = []string{}
	}
	b, err := json.Marshal(items)
	if err != nil {
		return datatypes.JSON("[]")
	}
	return datatypes.JSON(b)
}

func (s *ProjectService) uniqueSlug(base string, excludeID uint) string {
	candidate := base
	for i := 2; i < 200; i++ {
		var count int64
		q := s.db.Unscoped().Model(&model.Project{}).Where("slug = ?", candidate)
		if excludeID > 0 {
			q = q.Where("id <> ?", excludeID)
		}
		if err := q.Count(&count).Error; err != nil {
			return candidate
		}
		if count == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return fmt.Sprintf("%s-%d", base, time.Now().Unix())
}
