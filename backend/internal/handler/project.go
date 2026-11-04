package handler

import (
	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// ProjectHandler 作品集接口。
type ProjectHandler struct{ svc *service.ProjectService }

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

// PublicList GET /api/v1/projects
func (h *ProjectHandler) PublicList(c *gin.Context) {
	items, err := h.svc.PublicList()
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, items)
}

// PublicDetail GET /api/v1/projects/:slug
func (h *ProjectHandler) PublicDetail(c *gin.Context) {
	item, err := h.svc.PublicDetail(c.Param("slug"))
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// AdminList GET /api/v1/admin/projects
func (h *ProjectHandler) AdminList(c *gin.Context) {
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

// AdminGet GET /api/v1/admin/projects/:id
func (h *ProjectHandler) AdminGet(c *gin.Context) {
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

// AdminCreate POST /api/v1/admin/projects
func (h *ProjectHandler) AdminCreate(c *gin.Context) {
	var in dto.ProjectInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "项目名称不能为空，且不能超过 200 字")
		return
	}
	item, err := h.svc.Create(in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// AdminUpdate PUT /api/v1/admin/projects/:id
func (h *ProjectHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in dto.ProjectInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "项目名称不能为空，且不能超过 200 字")
		return
	}
	item, err := h.svc.Update(id, in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, item)
}

// AdminDelete DELETE /api/v1/admin/projects/:id
func (h *ProjectHandler) AdminDelete(c *gin.Context) {
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
