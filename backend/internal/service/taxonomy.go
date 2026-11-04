package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/model"
	"personal-website/backend/internal/pkg/slug"
)

// TaxonomyService 管理分类与标签。
type TaxonomyService struct{ db *gorm.DB }

func NewTaxonomyService(db *gorm.DB) *TaxonomyService { return &TaxonomyService{db: db} }

// CategoryItem 分类 + 已发布文章数。
type CategoryItem struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Sort         int    `json:"sort"`
	ArticleCount int64  `json:"article_count"`
}

// TagItem 标签 + 已发布文章数。
type TagItem struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ArticleCount int64  `json:"article_count"`
}

func (s *TaxonomyService) Categories() ([]CategoryItem, error) {
	var items []CategoryItem
	err := s.db.Model(&model.Category{}).
		Select("categories.id, categories.name, categories.slug, categories.sort, COUNT(articles.id) AS article_count").
		Joins("LEFT JOIN articles ON articles.category_id = categories.id AND articles.deleted_at IS NULL AND articles.status = ?", model.ArticleStatusPublished).
		Group("categories.id, categories.name, categories.slug, categories.sort").
		Order("categories.sort ASC, categories.id ASC").
		Scan(&items).Error
	return items, err
}

func (s *TaxonomyService) Tags() ([]TagItem, error) {
	var items []TagItem
	err := s.db.Model(&model.Tag{}).
		Select("tags.id, tags.name, tags.slug, COUNT(articles.id) AS article_count").
		Joins("LEFT JOIN article_tags ON article_tags.tag_id = tags.id").
		Joins("LEFT JOIN articles ON articles.id = article_tags.article_id AND articles.deleted_at IS NULL AND articles.status = ?", model.ArticleStatusPublished).
		Group("tags.id, tags.name, tags.slug").
		Order("article_count DESC, tags.id ASC").
		Scan(&items).Error
	return items, err
}

func (s *TaxonomyService) CreateCategory(in dto.CategoryInput) (*model.Category, error) {
	c := model.Category{Name: in.Name, Slug: in.Slug, Sort: in.Sort}
	if c.Slug == "" {
		c.Slug = slug.Make(c.Name)
	}
	c.Slug = s.uniqueCategorySlug(c.Slug, 0)
	if err := s.db.Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *TaxonomyService) UpdateCategory(id uint, in dto.CategoryInput) (*model.Category, error) {
	var c model.Category
	if err := s.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if in.Slug == "" {
		in.Slug = slug.Make(in.Name)
	}
	c.Name = in.Name
	c.Sort = in.Sort
	c.Slug = s.uniqueCategorySlug(in.Slug, id)
	if err := s.db.Save(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// DeleteCategory 删除分类，并把原分类下的文章置为未分类。
func (s *TaxonomyService) DeleteCategory(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&model.Category{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrNotFound
		}
		return tx.Model(&model.Article{}).Where("category_id = ?", id).Update("category_id", nil).Error
	})
}

func (s *TaxonomyService) CreateTag(in dto.TagInput) (*model.Tag, error) {
	t := model.Tag{Name: in.Name, Slug: in.Slug}
	if t.Slug == "" {
		t.Slug = slug.Make(t.Name)
	}
	t.Slug = s.uniqueTagSlug(t.Slug, 0)
	if err := s.db.Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TaxonomyService) UpdateTag(id uint, in dto.TagInput) (*model.Tag, error) {
	var t model.Tag
	if err := s.db.First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if in.Slug == "" {
		in.Slug = slug.Make(in.Name)
	}
	t.Name = in.Name
	t.Slug = s.uniqueTagSlug(in.Slug, id)
	if err := s.db.Save(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// DeleteTag 删除标签并清理文章关联。
func (s *TaxonomyService) DeleteTag(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Delete(&model.Tag{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrNotFound
		}
		return tx.Exec("DELETE FROM article_tags WHERE tag_id = ?", id).Error
	})
}

func (s *TaxonomyService) uniqueCategorySlug(base string, excludeID uint) string {
	candidate := base
	for i := 2; i < 200; i++ {
		var count int64
		q := s.db.Unscoped().Model(&model.Category{}).Where("slug = ?", candidate)
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

func (s *TaxonomyService) uniqueTagSlug(base string, excludeID uint) string {
	candidate := base
	for i := 2; i < 200; i++ {
		var count int64
		q := s.db.Unscoped().Model(&model.Tag{}).Where("slug = ?", candidate)
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
