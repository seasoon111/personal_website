package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config 是应用的全部运行配置，来自环境变量与 .env 文件。
type Config struct {
	AppEnv          string
	ServerPort      string
	SiteName        string
	SiteDescription string
	SiteICP         string
	DB              DBConfig
	JWT         JWTConfig
	Upload      UploadConfig
	Admin       AdminConfig
	CORSOrigins []string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// DSN 返回 GORM 使用的 MySQL 连接串。
func (d DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// TTL 返回 token 有效期。
func (j JWTConfig) TTL() time.Duration {
	return time.Duration(j.ExpireHours) * time.Hour
}

type UploadConfig struct {
	Dir   string
	MaxMB int64
}

type AdminConfig struct {
	Username string
	Password string
}

// Load 读取 .env（存在时）并组装配置，同时做最基本的合法性校验。
func Load() (*Config, error) {
	// 从 backend/ 目录运行；也兼容从仓库根目录运行的情况。
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	cfg := &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		SiteName:        getEnv("SITE_NAME", "个人网站"),
		SiteDescription: getEnv("SITE_DESCRIPTION", "记录文章、项目和日常的地方"),
		SiteICP:         getEnv("SITE_ICP", ""),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "personal_website"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", ""),
			ExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 168),
		},
		Upload: UploadConfig{
			Dir:   getEnv("UPLOAD_DIR", "./uploads"),
			MaxMB: int64(getEnvInt("UPLOAD_MAX_MB", 10)),
		},
		Admin: AdminConfig{
			Username: getEnv("ADMIN_USERNAME", "admin"),
			Password: getEnv("ADMIN_PASSWORD", "admin123"),
		},
		CORSOrigins: splitCSV(getEnv("CORS_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173")),
	}

	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD 未配置，请检查 backend/.env")
	}
	if len(cfg.JWT.Secret) < 16 {
		return nil, fmt.Errorf("JWT_SECRET 至少需要 16 位，请检查 backend/.env")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
