package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/yourusername/tastygo/config"
	"github.com/yourusername/tastygo/internal/api"
	"github.com/yourusername/tastygo/internal/auth"
	"github.com/yourusername/tastygo/internal/database"
	"github.com/yourusername/tastygo/internal/logging"
)

func main() {
	// Tải cấu hình ứng dụng
	appConfig := config.LoadAppConfig()
	dbConfig := config.LoadDBConfig()
	
	// Kiểm tra và đảm bảo JWT secret đủ mạnh
	if len(appConfig.JWTSecret) < 32 {
		logging.Warn("JWT_SECRET should be at least 32 characters long for security", nil)
		if appConfig.JWTSecret == "change_this_in_production" || appConfig.JWTSecret == "" {
			logging.Error("Using default JWT_SECRET in production is not secure!", nil)
		}
	}

	// Khởi tạo JWT secret
	auth.InitJWTSecret()

	// Khởi tạo database
	err := database.InitDB(dbConfig.Path)
	if err != nil {
		logging.Fatal("Failed to initialize database", map[string]interface{}{
			"error": err.Error(),
			"path":  dbConfig.Path,
		})
	}

	// Khởi tạo server
	server := api.NewServer()

	// Xử lý graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Khởi động server
		port := appConfig.Port
		portStr := os.Getenv("PORT")
		if portStr != "" {
			port = appConfig.Port
		}

		logging.Info("Server starting", map[string]interface{}{
			"port": port,
			"mode": appConfig.GinMode,
		})
		
		if err := server.Run(":" + strconv.Itoa(port)); err != nil {
			logging.Fatal("Failed to start server", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	<-quit
	logging.Info("Shutting down server...", nil)

	// Thực hiện các tác vụ cleanup nếu cần
	// ...

	logging.Info("Server exited properly", nil)
}
