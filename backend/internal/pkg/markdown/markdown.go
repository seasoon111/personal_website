package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// 渲染器开启 GFM（表格、删除线、任务列表、自动链接），并允许 Markdown 中内嵌 HTML。
var engine = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

// Render 把 Markdown 渲染为 HTML；失败时返回空串而不是半成品。
func Render(source string) string {
	if source == "" {
		return ""
	}
	var buf bytes.Buffer
	if err := engine.Convert([]byte(source), &buf); err != nil {
		return ""
	}
	return buf.String()
}
