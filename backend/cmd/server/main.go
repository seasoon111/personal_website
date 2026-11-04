package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"personal-website/backend/internal/config"
	"personal-website/backend/internal/database"
	"personal-website/backend/internal/router"
	"personal-website/backend/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("加载配置失败", "err", err)
		os.Exit(1)
	}

	db, err := database.Open(cfg.DB)
	if err != nil {
		slog.Error("数据库初始化失败", "err", err)
		os.Exit(1)
	}
	if err := database.Migrate(db); err != nil {
		slog.Error("数据库迁移失败", "err", err)
		os.Exit(1)
	}
	if err := database.Seed(db, cfg.Admin.Username, cfg.Admin.Password); err != nil {
		slog.Error("种子数据写入失败", "err", err)
		os.Exit(1)
	}

	store, err := storage.NewLocal(cfg.Upload.Dir, "/uploads")
	if err != nil {
		slog.Error("初始化存储失败", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router.New(cfg, db, store),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("服务已启动", "addr", "http://127.0.0.1:"+cfg.ServerPort, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP 服务异常退出", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("正在关闭服务…")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("关闭服务失败", "err", err)
	}
	slog.Info("服务已退出")
}
