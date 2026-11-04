package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Storage 抽象文件存储。首版只有本地实现，后续接阿里云 OSS 时新增一个实现即可。
type Storage interface {
	// Save 保存文件并返回内部相对路径与对外 URL。
	Save(r io.Reader, ext string) (relPath string, url string, err error)
	// Delete 删除内部相对路径对应的文件；文件不存在视为成功。
	Delete(relPath string) error
	// URL 把内部相对路径转换为对外 URL。
	URL(relPath string) string
}

// Local 把文件存在服务器磁盘上，通过静态路由对外提供。
type Local struct {
	root      string
	urlPrefix string
}

// NewLocal 创建本地存储，root 是磁盘目录（如 ./uploads），urlPrefix 是对外前缀（如 /uploads）。
func NewLocal(root, urlPrefix string) (*Local, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}
	return &Local{root: root, urlPrefix: strings.TrimRight(urlPrefix, "/")}, nil
}

func (l *Local) Save(r io.Reader, ext string) (string, string, error) {
	now := time.Now()
	dir := filepath.Join(now.Format("2006"), now.Format("01"))
	rel := filepath.ToSlash(filepath.Join(dir, randomName()+ext))
	full := filepath.Join(l.root, filepath.FromSlash(rel))

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %w", err)
	}
	f, err := os.Create(full)
	if err != nil {
		return "", "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", "", fmt.Errorf("写入文件失败: %w", err)
	}
	return rel, l.URL(rel), nil
}

func (l *Local) Delete(relPath string) error {
	full := filepath.Join(l.root, filepath.FromSlash(relPath))
	cleanRoot := filepath.Clean(l.root) + string(os.PathSeparator)
	if !strings.HasPrefix(filepath.Clean(full), cleanRoot) {
		return fmt.Errorf("非法文件路径")
	}
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

func (l *Local) URL(relPath string) string {
	return l.urlPrefix + "/" + strings.TrimLeft(filepath.ToSlash(relPath), "/")
}

func randomName() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), hex.EncodeToString(b))
}
