package handler

import (
	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// TaxonomyHandler 分类与标签。
type TaxonomyHandler struct{ svc *service.TaxonomyService }

func NewTaxonomyHandler(svc *service.TaxonomyService) *TaxonomyHandler {
	return &TaxonomyHandler{svc: svc}
}

// Categories GET /api/v1/categories
func (h *TaxonomyHandler) Categories(c *gin.Context) {
	items, err := h.svc.Categories()
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, items)
}

// Tags GET /api/v1/tags
func (h *TaxonomyHandler) Tags(c *gin.Context) {
	items, err := h.svc.Tags()
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, items)
}

// CreateCategory POST /api/v1/admin/categories
func (h *TaxonomyHandler) CreateCategory(c *gin.Context) {
	var in dto.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "分类名称不能为空，且不能超过 64 字")
		return
	}
	item, err := h.svc.CreateCategory(in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// UpdateCategory PUT /api/v1/admin/categories/:id
func (h *TaxonomyHandler) UpdateCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in dto.CategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "分类名称不能为空，且不能超过 64 字")
		return
	}
	item, err := h.svc.UpdateCategory(id, in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// DeleteCategory DELETE /api/v1/admin/categories/:id
func (h *TaxonomyHandler) DeleteCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteCategory(id); err != nil {
		handleErr(c, err)
		return
	}
	response.OKMsg(c, "已删除，原分类下的文章已变为未分类")
}

// CreateTag POST /api/v1/admin/tags
func (h *TaxonomyHandler) CreateTag(c *gin.Context) {
	var in dto.TagInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "标签名称不能为空，且不能超过 64 字")
		return
	}
	item, err := h.svc.CreateTag(in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// UpdateTag PUT /api/v1/admin/tags/:id
func (h *TaxonomyHandler) UpdateTag(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in dto.TagInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "标签名称不能为空，且不能超过 64 字")
		return
	}
	item, err := h.svc.UpdateTag(id, in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// DeleteTag DELETE /api/v1/admin/tags/:id
func (h *TaxonomyHandler) DeleteTag(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteTag(id); err != nil {
		handleErr(c, err)
		return
	}
	response.OKMsg(c, "已删除")
}
