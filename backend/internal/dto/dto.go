package dto

// Pagination 通用分页参数。
type Pagination struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// Normalize 把非法分页参数收敛到合法范围。
func (p *Pagination) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Offset 计算 SQL OFFSET。
func (p Pagination) Offset() int { return (p.Page - 1) * p.PageSize }

// PageResult 统一分页返回。
type PageResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// ArticleQuery 前台文章筛选。
type ArticleQuery struct {
	Pagination
	Category string `form:"category"`
	Tag      string `form:"tag"`
	Keyword  string `form:"keyword"`
}

// AdminArticleQuery 后台文章筛选。
type AdminArticleQuery struct {
	Pagination
	Status   string `form:"status"`
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
}

// ArticleInput 新建 / 更新文章。
type ArticleInput struct {
	Title      string `json:"title" binding:"required,max=200"`
	Slug       string `json:"slug" binding:"max=200"`
	Summary    string `json:"summary" binding:"max=500"`
	ContentMD  string `json:"content_md"`
	CoverURL   string `json:"cover_url" binding:"max=500"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
	Status     string `json:"status" binding:"omitempty,oneof=draft published"`
	IsPinned   bool   `json:"is_pinned"`
}

// ProjectInput 新建 / 更新项目。
type ProjectInput struct {
	Title     string   `json:"title" binding:"required,max=200"`
	Slug      string   `json:"slug" binding:"max=200"`
	Summary   string   `json:"summary" binding:"max=500"`
	ContentMD string   `json:"content_md"`
	CoverURL  string   `json:"cover_url" binding:"max=500"`
	TechStack []string `json:"tech_stack"`
	RepoURL   string   `json:"repo_url" binding:"max=500"`
	DemoURL   string   `json:"demo_url" binding:"max=500"`
	Sort      int      `json:"sort"`
	Status    string   `json:"status" binding:"omitempty,oneof=draft published"`
}

// MomentInput 新建 / 更新动态，ImageIDs 最多 9 个。
type MomentInput struct {
	Content  string `json:"content" binding:"max=2000"`
	ImageIDs []uint `json:"image_ids" binding:"max=9"`
	IsPublic *bool  `json:"is_public"`
}

// CategoryInput 分类。
type CategoryInput struct {
	Name string `json:"name" binding:"required,max=64"`
	Slug string `json:"slug" binding:"max=64"`
	Sort int    `json:"sort"`
}

// TagInput 标签。
type TagInput struct {
	Name string `json:"name" binding:"required,max=64"`
	Slug string `json:"slug" binding:"max=64"`
}

// TimelineItem 关于我页的经历条目。
type TimelineItem struct {
	Period      string `json:"period"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// SocialItem 关于我页的社交链接。
type SocialItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}

// AboutInput 关于我页面内容。
type AboutInput struct {
	AvatarURL string         `json:"avatar_url" binding:"max=500"`
	IntroMD   string         `json:"intro_md"`
	Skills    []string       `json:"skills"`
	Timeline  []TimelineItem `json:"timeline"`
	Socials   []SocialItem   `json:"socials"`
}

// ContactInput 前台联系表单。
type ContactInput struct {
	Name    string `json:"name" binding:"required,min=1,max=64"`
	Email   string `json:"email" binding:"required,email,max=128"`
	Subject string `json:"subject" binding:"max=200"`
	Content string `json:"content" binding:"required,min=5,max=5000"`
}
