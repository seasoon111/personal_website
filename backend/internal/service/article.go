package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/model"
	"personal-website/backend/internal/pkg/markdown"
	"personal-website/backend/internal/pkg/slug"
)

const articleFields = "articles.id, articles.title, articles.slug, articles.summary, articles.cover_url, articles.category_id, articles.status, articles.is_pinned, articles.published_at, articles.created_at, articles.updated_at"

// ArticleService 负责文章的读写。
type ArticleService struct{ db *gorm.DB }

func NewArticleService(db *gorm.DB) *ArticleService { return &ArticleService{db: db} }

// ArticleDetail 详情页返回当前文章与相邻文章。
type ArticleDetail struct {
	Article model.Article  `json:"article"`
	Prev    *model.Article `json:"prev"`
	Next    *model.Article `json:"next"`
}

// ArchiveGroup 归档页按年份分组。
type ArchiveGroup struct {
	Year     int             `json:"year"`
	Articles []model.Article `json:"articles"`
}

// PublicList 前台文章列表，只返回已发布文章。
func (s *ArticleService) PublicList(q dto.ArticleQuery) (*dto.PageResult, error) {
	q.Normalize()

	tx := s.db.Model(&model.Article{}).Where("articles.status = ?", model.ArticleStatusPublished)
	if q.Category != "" {
		tx = tx.Joins("JOIN categories ON categories.id = articles.category_id").
			Where("categories.slug = ?", q.Category)
	}
	if q.Tag != "" {
		tx = tx.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.slug = ?", q.Tag)
	}
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		tx = tx.Where("articles.title LIKE ? OR articles.summary LIKE ?", like, like)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []model.Article
	if err := tx.Select(articleFields).
		Preload("Category").Preload("Tags").
		Order("articles.is_pinned DESC, articles.published_at DESC, articles.id DESC").
		Offset(q.Offset()).Limit(q.PageSize).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// PublicDetail 按 slug 取已发布文章，并附带上一篇（更旧）/下一篇（更新）。
func (s *ArticleService) PublicDetail(slugStr string) (*ArticleDetail, error) {
	var a model.Article
	err := s.db.Preload("Category").Preload("Tags").
		Where("slug = ? AND status = ?", slugStr, model.ArticleStatusPublished).
		First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	detail := &ArticleDetail{Article: a}
	if a.PublishedAt != nil {
		var prev, next model.Article
		_ = s.db.Select(articleFields).Where("status = ? AND published_at < ?", model.ArticleStatusPublished, a.PublishedAt).
			Order("published_at DESC").First(&prev).Error
		_ = s.db.Select(articleFields).Where("status = ? AND published_at > ?", model.ArticleStatusPublished, a.PublishedAt).
			Order("published_at ASC").First(&next).Error
		if prev.ID != 0 {
			detail.Prev = &prev
		}
		if next.ID != 0 {
			detail.Next = &next
		}
	}
	return detail, nil
}

// AdminList 后台文章列表，包含草稿。
func (s *ArticleService) AdminList(q dto.AdminArticleQuery) (*dto.PageResult, error) {
	q.Normalize()

	tx := s.db.Model(&model.Article{})
	if q.Status != "" {
		tx = tx.Where("articles.status = ?", q.Status)
	}
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		tx = tx.Where("articles.title LIKE ? OR articles.summary LIKE ?", like, like)
	}
	if q.Category != "" {
		tx = tx.Joins("JOIN categories ON categories.id = articles.category_id").
			Where("categories.slug = ?", q.Category)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []model.Article
	if err := tx.Select(articleFields).
		Preload("Category").Preload("Tags").
		Order("articles.is_pinned DESC, articles.updated_at DESC").
		Offset(q.Offset()).Limit(q.PageSize).
		Find(&items).Error; err != nil {
		return nil, err
	}

	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

// GetByID 后台按 ID 取文章（含正文）。
func (s *ArticleService) GetByID(id uint) (*model.Article, error) {
	var a model.Article
	err := s.db.Preload("Category").Preload("Tags").First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Create 新建文章。
func (s *ArticleService) Create(in dto.ArticleInput) (*model.Article, error) {
	a := model.Article{
		Title:      in.Title,
		Slug:       in.Slug,
		Summary:    in.Summary,
		ContentMD:  in.ContentMD,
		CoverURL:   in.CoverURL,
		CategoryID: in.CategoryID,
		Status:     in.Status,
		IsPinned:   in.IsPinned,
	}
	if a.Status == "" {
		a.Status = model.ArticleStatusDraft
	}
	if a.Slug == "" {
		a.Slug = slug.Make(a.Title)
	}
	a.Slug = s.uniqueSlug(a.Slug, 0)
	a.ContentHTML = markdown.Render(a.ContentMD)
	if a.Status == model.ArticleStatusPublished {
		now := time.Now()
		a.PublishedAt = &now
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&a).Error; err != nil {
			return err
		}
		return s.replaceTags(tx, &a, in.TagIDs)
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(a.ID)
}

// Update 更新文章；首次从草稿变为发布时补发布时间。
func (s *ArticleService) Update(id uint, in dto.ArticleInput) (*model.Article, error) {
	a, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if in.Slug == "" {
		in.Slug = slug.Make(in.Title)
	}
	a.Slug = s.uniqueSlug(in.Slug, id)
	a.Title = in.Title
	a.Summary = in.Summary
	a.ContentMD = in.ContentMD
	a.ContentHTML = markdown.Render(in.ContentMD)
	a.CoverURL = in.CoverURL
	a.CategoryID = in.CategoryID
	a.IsPinned = in.IsPinned

	if in.Status == "" {
		in.Status = model.ArticleStatusDraft
	}
	if in.Status == model.ArticleStatusPublished && a.PublishedAt == nil {
		now := time.Now()
		a.PublishedAt = &now
	}
	a.Status = in.Status

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(a).Error; err != nil {
			return err
		}
		return s.replaceTags(tx, a, in.TagIDs)
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

// Delete 软删除文章。
func (s *ArticleService) Delete(id uint) error {
	res := s.db.Delete(&model.Article{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// SetStatus 发布 / 下架。
func (s *ArticleService) SetStatus(id uint, status string) error {
	if status != model.ArticleStatusDraft && status != model.ArticleStatusPublished {
		return ErrBadInput
	}
	a, err := s.GetByID(id)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{"status": status}
	if status == model.ArticleStatusPublished && a.PublishedAt == nil {
		updates["published_at"] = time.Now()
	}
	return s.db.Model(&model.Article{}).Where("id = ?", id).Updates(updates).Error
}

// SetPinned 置顶 / 取消置顶。
func (s *ArticleService) SetPinned(id uint, pinned bool) error {
	res := s.db.Model(&model.Article{}).Where("id = ?", id).Update("is_pinned", pinned)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Archives 按年份倒序归档。
func (s *ArticleService) Archives() ([]ArchiveGroup, error) {
	var items []model.Article
	if err := s.db.Select(articleFields).
		Where("status = ?", model.ArticleStatusPublished).
		Order("published_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	groups := make([]ArchiveGroup, 0, 8)
	for _, a := range items {
		year := a.CreatedAt.Year()
		if a.PublishedAt != nil {
			year = a.PublishedAt.Year()
		}
		if len(groups) == 0 || groups[len(groups)-1].Year != year {
			groups = append(groups, ArchiveGroup{Year: year})
		}
		groups[len(groups)-1].Articles = append(groups[len(groups)-1].Articles, a)
	}
	return groups, nil
}

func (s *ArticleService) replaceTags(tx *gorm.DB, a *model.Article, ids []uint) error {
	tags := make([]model.Tag, 0, len(ids))
	if len(ids) > 0 {
		if err := tx.Find(&tags, ids).Error; err != nil {
			return err
		}
	}
	return tx.Model(a).Association("Tags").Replace(tags)
}

// uniqueSlug 保证 slug 唯一；软删除的记录也算占用，避免数据库唯一索引冲突。
func (s *ArticleService) uniqueSlug(base string, excludeID uint) string {
	candidate := base
	for i := 2; i < 200; i++ {
		var count int64
		q := s.db.Unscoped().Model(&model.Article{}).Where("slug = ?", candidate)
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
