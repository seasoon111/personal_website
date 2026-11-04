package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 是全站统一的响应结构：code=0 表示成功。
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK 返回 200 与数据。
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "ok", Data: data})
}

// OKMsg 返回 200 与自定义提示（无数据）。
func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: msg})
}

// Fail 返回指定状态码与业务错误码。
func Fail(c *gin.Context, status, code int, msg string) {
	c.JSON(status, Body{Code: code, Message: msg})
}

// BadRequest 参数错误。
func BadRequest(c *gin.Context, msg string) {
	Fail(c, http.StatusBadRequest, 40000, msg)
}

// Unauthorized 未登录或 token 失效。
func Unauthorized(c *gin.Context, msg string) {
	Fail(c, http.StatusUnauthorized, 40100, msg)
}

// Forbidden 无权限。
func Forbidden(c *gin.Context, msg string) {
	Fail(c, http.StatusForbidden, 40300, msg)
}

// NotFound 资源不存在。
func NotFound(c *gin.Context, msg string) {
	Fail(c, http.StatusNotFound, 40400, msg)
}

// Conflict 资源冲突（如 slug 重复、图片被引用）。
func Conflict(c *gin.Context, msg string) {
	Fail(c, http.StatusConflict, 40900, msg)
}

// Internal 服务器内部错误，细节只进日志，不回传给客户端。
func Internal(c *gin.Context, msg string) {
	if msg == "" {
		msg = "服务器开小差了，请稍后再试"
	}
	Fail(c, http.StatusInternalServerError, 50000, msg)
}
