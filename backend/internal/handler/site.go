package handler

import (
	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/config"
	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// SiteHandler 站点信息、关于我、联系表单。
type SiteHandler struct {
	cfg     *config.Config
	about   *service.AboutService
	message *service.MessageService
}

func NewSiteHandler(cfg *config.Config, about *service.AboutService, message *service.MessageService) *SiteHandler {
	return &SiteHandler{cfg: cfg, about: about, message: message}
}

// Site GET /api/v1/site
func (h *SiteHandler) Site(c *gin.Context) {
	profile, err := h.about.Get()
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, gin.H{
		"name":        h.cfg.SiteName,
		"description": h.cfg.SiteDescription,
		"icp":         h.cfg.SiteICP,
		"avatar_url":  profile.AvatarURL,
		"socials":     profile.Socials,
	})
}

// About GET /api/v1/about
func (h *SiteHandler) About(c *gin.Context) {
	profile, err := h.about.Get()
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, profile)
}

// UpdateAbout PUT /api/v1/admin/about
func (h *SiteHandler) UpdateAbout(c *gin.Context) {
	var in dto.AboutInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "关于我内容格式不正确")
		return
	}
	profile, err := h.about.Update(in)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, profile)
}

// Contact POST /api/v1/contact
func (h *SiteHandler) Contact(c *gin.Context) {
	var in dto.ContactInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "请填写称呼、正确的邮箱和至少 5 个字的留言")
		return
	}
	if _, err := h.message.Create(in); err != nil {
		handleErr(c, err)
		return
	}
	response.OKMsg(c, "留言已收到，谢谢你的来信")
}
