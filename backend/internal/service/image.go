package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	_ "golang.org/x/image/webp"
	"gorm.io/gorm"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/model"
	"personal-website/backend/internal/storage"
)

// 允许上传的图片类型与对应扩展名。
var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// ImageService 负责图片库：上传、列表、删除。
type ImageService struct {
	db       *gorm.DB
	store    storage.Storage
	maxBytes int64
}

func NewImageService(db *gorm.DB, store storage.Storage, maxMB int64) *ImageService {
	if maxMB <= 0 {
		maxMB = 10
	}
	return &ImageService{db: db, store: store, maxBytes: maxMB * 1024 * 1024}
}

// Upload 校验并保存一张图片，返回入库后的记录。
func (s *ImageService) Upload(originalName string, data []byte) (*model.Image, error) {
	if int64(len(data)) > s.maxBytes {
		return nil, fmt.Errorf("%w: 图片不能超过 %d MB", ErrBadInput, s.maxBytes/1024/1024)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: 文件为空", ErrBadInput)
	}

	mime := http.DetectContentType(data)
	ext, ok := allowedImageTypes[mime]
	if !ok {
		return nil, fmt.Errorf("%w: 只支持 jpg / png / webp / gif", ErrBadInput)
	}

	width, height := 0, 0
	if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		width, height = cfg.Width, cfg.Height
	}

	relPath, url, err := s.store.Save(bytes.NewReader(data), ext)
	if err != nil {
		return nil, err
	}

	img := model.Image{
		Filename:     relPath,
		OriginalName: originalName,
		Path:         relPath,
		URL:          url,
		Mime:         mime,
		Size:         int64(len(data)),
		Width:        width,
		Height:       height,
	}
	if err := s.db.Create(&img).Error; err != nil {
		_ = s.store.Delete(relPath)
		return nil, err
	}
	return &img, nil
}

func (s *ImageService) List(q dto.Pagination) (*dto.PageResult, error) {
	q.Normalize()
	var total int64
	if err := s.db.Model(&model.Image{}).Count(&total).Error; err != nil {
		return nil, err
	}
	var items []model.Image
	if err := s.db.Order("id DESC").Offset(q.Offset()).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, err
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *ImageService) GetByID(id uint) (*model.Image, error) {
	var img model.Image
	if err := s.db.First(&img, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &img, nil
}

// Delete 删除图片；被动态或文章封面引用时拒绝，避免前台出现裂图。
func (s *ImageService) Delete(id uint) error {
	img, err := s.GetByID(id)
	if err != nil {
		return err
	}

	var momentRefs int64
	if err := s.db.Model(&model.MomentImage{}).Where("image_id = ?", id).Count(&momentRefs).Error; err != nil {
		return err
	}
	var articleRefs int64
	if err := s.db.Model(&model.Article{}).Where("cover_url = ?", img.URL).Count(&articleRefs).Error; err != nil {
		return err
	}
	var projectRefs int64
	if err := s.db.Model(&model.Project{}).Where("cover_url = ?", img.URL).Count(&projectRefs).Error; err != nil {
		return err
	}
	if momentRefs > 0 || articleRefs > 0 || projectRefs > 0 {
		return ErrInUse
	}

	if err := s.db.Delete(&model.Image{}, id).Error; err != nil {
		return err
	}
	// 磁盘删除失败不回滚数据库，只记录（由调用方打日志）。
	_ = s.store.Delete(img.Path)
	return nil
}
