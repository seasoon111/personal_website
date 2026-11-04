package handler

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"personal-website/backend/internal/config"
	"personal-website/backend/internal/middleware"
	"personal-website/backend/internal/model"
	"personal-website/backend/internal/pkg/jwtx"
	"personal-website/backend/internal/pkg/password"
	"personal-website/backend/internal/pkg/response"
)

// AuthHandler 处理登录、当前用户、修改密码。
type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type loginReq struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type loginResp struct {
	Token     string     `json:"token"`
	ExpiresIn int64      `json:"expires_in"`
	User      model.User `json:"user"`
}

// Login POST /api/v1/admin/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "用户名或密码格式不正确")
		return
	}

	var user model.User
	err := h.db.Where("username = ?", req.Username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}
	if err != nil {
		slog.Error("查询用户失败", "err", err)
		response.Internal(c, "")
		return
	}

	if !password.Verify(user.PasswordHash, req.Password) {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, err := jwtx.Generate(h.cfg.JWT.Secret, user.ID, user.Username, h.cfg.JWT.TTL())
	if err != nil {
		slog.Error("签发 token 失败", "err", err)
		response.Internal(c, "")
		return
	}

	response.OK(c, loginResp{
		Token:     token,
		ExpiresIn: int64(h.cfg.JWT.TTL().Seconds()),
		User:      user,
	})
}

// Me GET /api/v1/admin/me
func (h *AuthHandler) Me(c *gin.Context) {
	var user model.User
	if err := h.db.First(&user, middleware.CurrentUserID(c)).Error; err != nil {
		response.Unauthorized(c, "用户不存在，请重新登录")
		return
	}
	response.OK(c, user)
}

type changePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

// ChangePassword PUT /api/v1/admin/password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req changePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "新密码至少 8 位")
		return
	}
	if req.OldPassword == req.NewPassword {
		response.BadRequest(c, "新密码不能与旧密码相同")
		return
	}

	var user model.User
	if err := h.db.First(&user, middleware.CurrentUserID(c)).Error; err != nil {
		response.Unauthorized(c, "用户不存在，请重新登录")
		return
	}
	if !password.Verify(user.PasswordHash, req.OldPassword) {
		response.BadRequest(c, "旧密码不正确")
		return
	}

	hash, err := password.Hash(req.NewPassword)
	if err != nil {
		slog.Error("生成密码哈希失败", "err", err)
		response.Internal(c, "")
		return
	}
	if err := h.db.Model(&user).Update("password_hash", hash).Error; err != nil {
		slog.Error("更新密码失败", "err", err)
		response.Internal(c, "")
		return
	}

	response.OKMsg(c, "密码已更新")
}
