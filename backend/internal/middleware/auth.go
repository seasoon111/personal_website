package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/pkg/jwtx"
	"personal-website/backend/internal/pkg/response"
)

const (
	CtxUserID   = "user_id"
	CtxUsername = "username"
)

// Auth 校验 Bearer token，并把用户信息写进 gin.Context。
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		claims, err := jwtx.Parse(secret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Unauthorized(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Next()
	}
}

// CurrentUserID 从上下文取当前用户 ID。
func CurrentUserID(c *gin.Context) uint {
	if v, ok := c.Get(CtxUserID); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}
