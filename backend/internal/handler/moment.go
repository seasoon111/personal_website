package handler

import (
	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// MomentHandler 个人动态接口。
type MomentHandler struct{ svc *service.MomentService }

func NewMomentHandler(svc *service.MomentService) *MomentHandler {
	return &MomentHandler{svc: svc}
}

// PublicList GET /api/v1/moments
func (h *MomentHandler) PublicList(c *gin.Context) {
	var q dto.Pagination
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确")
		return
	}
	res, err := h.svc.PublicList(q)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, res)
}

// AdminList GET /api/v1/admin/moments
func (h *MomentHandler) AdminList(c *gin.Context) {
	var q dto.Pagination
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确")
		return
	}
	res, err := h.svc.AdminList(q)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, res)
}

// AdminGet GET /api/v1/admin/moments/:id
func (h *MomentHandler) AdminGet(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.svc.GetByID(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// AdminCreate POST /api/v1/admin/moments
func (h *MomentHandler) AdminCreate(c *gin.Context) {
	var in dto.MomentInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "动态内容过长，或图片超过 9 张")
		return
	}
	item, err := h.svc.Create(in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// AdminUpdate PUT /api/v1/admin/moments/:id
func (h *MomentHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in dto.MomentInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "动态内容过长，或图片超过 9 张")
		return
	}
	item, err := h.svc.Update(id, in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// AdminDelete DELETE /api/v1/admin/moments/:id
func (h *MomentHandler) AdminDelete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		handleErr(c, err)
		return
	}
	response.OKMsg(c, "已删除")
}
