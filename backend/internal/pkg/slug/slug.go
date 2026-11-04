package slug

import (
	"strings"
	"unicode"
)

// Make 生成 URL slug：保留字母、数字、汉字，其余字符折叠为单个连字符。
// 例："Hello World! 你好" -> "hello-world-你好"
func Make(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || unicode.IsSpace(r):
			if !lastDash && b.Len() > 0 {
				b.WriteRune('-')
				lastDash = true
			}
		default:
			// 标点等直接丢弃
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "post"
	}
	if len([]rune(out)) > 80 {
		out = string([]rune(out)[:80])
	}
	return strings.Trim(out, "-")
}
