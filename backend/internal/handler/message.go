package handler

import (
	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// MessageHandler 后台留言管理。
type MessageHandler struct{ svc *service.MessageService }

func NewMessageHandler(svc *service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

// List GET /api/v1/admin/messages?unread=1
func (h *MessageHandler) List(c *gin.Context) {
	var q dto.Pagination
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确")
		return
	}
	unreadOnly := c.Query("unread") == "1"
	res, err := h.svc.List(q, unreadOnly)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, res)
}

type readReq struct {
	IsRead bool `json:"is_read"`
}

// SetRead PATCH /api/v1/admin/messages/:id/read
func (h *MessageHandler) SetRead(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req readReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数不正确")
		return
	}
	if err := h.svc.SetRead(id, req.IsRead); err != nil {
		handleErr(c, err)
		return
	}
	response.OKMsg(c, "已更新")
}

// Delete DELETE /api/v1/admin/messages/:id
func (h *MessageHandler) Delete(c *gin.Context) {
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
