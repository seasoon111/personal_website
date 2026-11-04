package handler

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// handleErr 把业务错误映射成统一的 HTTP 响应。
func handleErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "内容不存在或已被删除")
	case errors.Is(err, service.ErrInUse):
		response.Conflict(c, "该资源正在被使用，无法删除")
	case errors.Is(err, service.ErrBadInput):
		response.BadRequest(c, err.Error())
	default:
		slog.Error("请求处理失败",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"err", err,
		)
		response.Internal(c, "")
	}
}
