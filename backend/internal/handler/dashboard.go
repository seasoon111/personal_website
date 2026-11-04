package handler

import (
	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// DashboardHandler 后台仪表盘。
type DashboardHandler struct{ svc *service.DashboardService }

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// Stats GET /api/v1/admin/dashboard
func (h *DashboardHandler) Stats(c *gin.Context) {
	st, err := h.svc.Stats()
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, st)
}
