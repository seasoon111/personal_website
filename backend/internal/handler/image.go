package handler

import (
	"io"

	"github.com/gin-gonic/gin"

	"personal-website/backend/internal/dto"
	"personal-website/backend/internal/pkg/response"
	"personal-website/backend/internal/service"
)

// ImageHandler 图片库接口。
type ImageHandler struct{ svc *service.ImageService }

func NewImageHandler(svc *service.ImageService) *ImageHandler {
	return &ImageHandler{svc: svc}
}

// Upload POST /api/v1/admin/images  （multipart，字段名 files，可多张）
func (h *ImageHandler) Upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "请使用 multipart/form-data 上传，字段名为 files")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		files = form.File["file"]
	}
	if len(files) == 0 {
		response.BadRequest(c, "没有收到文件")
		return
	}
	if len(files) > 9 {
		response.BadRequest(c, "一次最多上传 9 张图片")
		return
	}

	uploaded := make([]interface{}, 0, len(files))
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			response.BadRequest(c, "读取上传文件失败")
			return
		}
		data, err := io.ReadAll(io.LimitReader(f, 32*1024*1024))
		_ = f.Close()
		if err != nil {
			response.BadRequest(c, "读取上传文件失败")
			return
		}

		img, err := h.svc.Upload(fh.Filename, data)
		if err != nil {
			handleErr(c, err)
			return
		}
		uploaded = append(uploaded, img)
	}

	response.OK(c, uploaded)
}

// List GET /api/v1/admin/images
func (h *ImageHandler) List(c *gin.Context) {
	var q dto.Pagination
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, "分页参数不正确")
		return
	}
	res, err := h.svc.List(q)
	if err != nil {
		handleErr(c, err)
		return
	}
	response.OK(c, res)
}

// Delete DELETE /api/v1/admin/images/:id
func (h *ImageHandler) Delete(c *gin.Context) {
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
