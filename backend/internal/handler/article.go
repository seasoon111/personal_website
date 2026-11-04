package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// ArticleHandler 文章的前台与后台接口。
type ArticleHandler struct{ svc *service.ArticleService }

func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// PublicList GET /api/v1/articles
func (h *ArticleHandler) PublicList(c *gin.Context) {
	var q dto.ArticleQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "查询参数不正确")
		return
	}
	res, err := h.svc.PublicList(q)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, res)
}

// PublicDetail GET /api/v1/articles/:slug
func (h *ArticleHandler) PublicDetail(c *gin.Context) {
	res, err := h.svc.PublicDetail(c.Param("slug"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, res)
}

// Archives GET /api/v1/archives
func (h *ArticleHandler) Archives(c *gin.Context) {
	res, err := h.svc.Archives()
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, res)
}

// AdminList GET /api/v1/admin/articles
func (h *ArticleHandler) AdminList(c *gin.Context) {
	var q dto.AdminArticleQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "查询参数不正确")
		return
	}
	res, err := h.svc.AdminList(q)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, res)
}

// AdminGet GET /api/v1/admin/articles/:id
func (h *ArticleHandler) AdminGet(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	a, err := h.svc.GetByID(id)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, a)
}

// AdminCreate POST /api/v1/admin/articles
func (h *ArticleHandler) AdminCreate(c *gin.Context) {
	var in dto.ArticleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "标题不能为空，且不能超过 200 字")
		return
	}
	a, err := h.svc.Create(in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, a)
}

// AdminUpdate PUT /api/v1/admin/articles/:id
func (h *ArticleHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in dto.ArticleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "标题不能为空，且不能超过 200 字")
		return
	}
	a, err := h.svc.Update(id, in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, a)
}

// AdminDelete DELETE /api/v1/admin/articles/:id
func (h *ArticleHandler) AdminDelete(c *gin.Context) {
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

type statusReq struct {
	Status string `json:"status" binding:"required,oneof=draft published"`
}

// AdminSetStatus PATCH /api/v1/admin/articles/:id/status
func (h *ArticleHandler) AdminSetStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req statusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "状态只能是 draft 或 published")
		return
	}
	if err := h.svc.SetStatus(id, req.Status); err != nil {
		handleErr(c, err)
		return
	}
	response.OKMsg(c, "状态已更新")
}

type pinReq struct {
	IsPinned bool `json:"is_pinned"`
}

// AdminSetPinned PATCH /api/v1/admin/articles/:id/pin
func (h *ArticleHandler) AdminSetPinned(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req pinReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数不正确")
		return
	}
	if err := h.svc.SetPinned(id, req.IsPinned); err != nil {
		handleErr(c, err)
		return
	}
	response.OKMsg(c, "已更新")
}

// parseID 解析路径参数 :id，失败时已写好响应。
func parseID(c *gin.Context) (uint, bool) {
	v, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || v == 0 {
		response.BadRequest(c, "ID 不合法")
		return 0, false
	}
	return uint(v), true
}
