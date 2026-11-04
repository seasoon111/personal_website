package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"personal-website/backend/internal/config"
	"personal-website/backend/internal/handler"
	"personal-website/backend/internal/middleware"
	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
	"personal-website/backend/internal/storage"
)

// New 组装全部路由。
func New(cfg *config.Config, db *gorm.DB, store storage.Storage) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.CORSOrigins))
	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "接口不存在")
	})

	// 本地图片静态访问：/uploads/2026/11/xxx.jpg
	r.Static("/uploads", cfg.Upload.Dir)

	articleSvc := service.NewArticleService(db)
	taxonomySvc := service.NewTaxonomyService(db)
	projectSvc := service.NewProjectService(db)
	momentSvc := service.NewMomentService(db)
	aboutSvc := service.NewAboutService(db)
	messageSvc := service.NewMessageService(db)
	imageSvc := service.NewImageService(db, store, cfg.Upload.MaxMB)
	dashboardSvc := service.NewDashboardService(db)

	authH := handler.NewAuthHandler(db, cfg)
	articleH := handler.NewArticleHandler(articleSvc)
	taxonomyH := handler.NewTaxonomyHandler(taxonomySvc)
	projectH := handler.NewProjectHandler(projectSvc)
	momentH := handler.NewMomentHandler(momentSvc)
	siteH := handler.NewSiteHandler(cfg, aboutSvc, messageSvc)
	imageH := handler.NewImageHandler(imageSvc)
	messageH := handler.NewMessageHandler(messageSvc)
	dashboardH := handler.NewDashboardHandler(dashboardSvc)

	api := r.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	// —— 公开接口 ——
	api.GET("/site", siteH.Site)
	api.GET("/about", siteH.About)
	api.GET("/articles", articleH.PublicList)
	api.GET("/articles/:slug", articleH.PublicDetail)
	api.GET("/archives", articleH.Archives)
	api.GET("/categories", taxonomyH.Categories)
	api.GET("/tags", taxonomyH.Tags)
	api.GET("/projects", projectH.PublicList)
	api.GET("/projects/:slug", projectH.PublicDetail)
	api.GET("/moments", momentH.PublicList)
	api.POST("/contact", middleware.RateLimit(3, time.Hour), siteH.Contact)

	// —— 登录 ——
	admin := api.Group("/admin")
	admin.POST("/login", middleware.LoginLimiter(5), authH.Login)

	// —— 需要登录的后台接口 ——
	authed := admin.Group("")
	authed.Use(middleware.Auth(cfg.JWT.Secret))
	{
		authed.GET("/me", authH.Me)
		authed.PUT("/password", authH.ChangePassword)
		authed.GET("/dashboard", dashboardH.Stats)

		authed.GET("/articles", articleH.AdminList)
		authed.POST("/articles", articleH.AdminCreate)
		authed.GET("/articles/:id", articleH.AdminGet)
		authed.PUT("/articles/:id", articleH.AdminUpdate)
		authed.DELETE("/articles/:id", articleH.AdminDelete)
		authed.PATCH("/articles/:id/status", articleH.AdminSetStatus)
		authed.PATCH("/articles/:id/pin", articleH.AdminSetPinned)

		authed.POST("/categories", taxonomyH.CreateCategory)
		authed.PUT("/categories/:id", taxonomyH.UpdateCategory)
		authed.DELETE("/categories/:id", taxonomyH.DeleteCategory)

		authed.POST("/tags", taxonomyH.CreateTag)
		authed.PUT("/tags/:id", taxonomyH.UpdateTag)
		authed.DELETE("/tags/:id", taxonomyH.DeleteTag)

		authed.GET("/projects", projectH.AdminList)
		authed.POST("/projects", projectH.AdminCreate)
		authed.GET("/projects/:id", projectH.AdminGet)
		authed.PUT("/projects/:id", projectH.AdminUpdate)
		authed.DELETE("/projects/:id", projectH.AdminDelete)

		authed.GET("/moments", momentH.AdminList)
		authed.POST("/moments", momentH.AdminCreate)
		authed.GET("/moments/:id", momentH.AdminGet)
		authed.PUT("/moments/:id", momentH.AdminUpdate)
		authed.DELETE("/moments/:id", momentH.AdminDelete)

		authed.GET("/about", siteH.About)
		authed.PUT("/about", siteH.UpdateAbout)

		authed.GET("/images", imageH.List)
		authed.POST("/images", imageH.Upload)
		authed.DELETE("/images/:id", imageH.Delete)

		authed.GET("/messages", messageH.List)
		authed.PATCH("/messages/:id/read", messageH.SetRead)
		authed.DELETE("/messages/:id", messageH.Delete)
	}

	return r
}
