package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"personal-website/backend/internal/config"
	"personal-website/backend/internal/model"
	"personal-website/backend/internal/pkg/password"
)

// Open 建立 MySQL 连接池并确认可连通。
func Open(cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time { return time.Now().Local() },
	})
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接池失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("MySQL Ping 失败: %w", err)
	}
	return db, nil
}

// Migrate 建表 / 补字段。首版用 AutoMigrate，后续结构变更可换成迁移文件。
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.Project{},
		&model.Image{},
		&model.Moment{},
		&model.MomentImage{},
		&model.AboutProfile{},
		&model.ContactMessage{},
	); err != nil {
		return fmt.Errorf("自动建表失败: %w", err)
	}
	return nil
}

// Seed 在首次启动时创建管理员账号与「关于我」单例记录；已有数据时不做任何事。
func Seed(db *gorm.DB, adminUsername, adminPassword string) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("统计用户失败: %w", err)
	}
	if count == 0 {
		hash, err := password.Hash(adminPassword)
		if err != nil {
			return fmt.Errorf("生成管理员密码哈希失败: %w", err)
		}
		user := model.User{Username: adminUsername, PasswordHash: hash}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("创建管理员失败: %w", err)
		}
		slog.Info("已创建初始管理员，请登录后立即修改密码", "username", adminUsername)
	}

	profile := model.AboutProfile{
		ID:       1,
		AvatarURL: "",
		IntroMD:  "你好，我是这个网站的主人。这里记录我的文章、项目和日常。",
		Skills:   datatypes.JSON(`[]`),
		Timeline: datatypes.JSON(`[]`),
		Socials:  datatypes.JSON(`[]`),
	}
	if err := db.FirstOrCreate(&profile, 1).Error; err != nil {
		return fmt.Errorf("初始化关于我失败: %w", err)
	}
	return nil
}
