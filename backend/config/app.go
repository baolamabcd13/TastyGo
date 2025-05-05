package config

import (
    "os"
    "strconv"
)

// AppConfig chứa cấu hình ứng dụng
type AppConfig struct {
    Port      int
    JWTSecret string
    GinMode   string
    LogLevel  string
}

// LoadAppConfig tải cấu hình từ biến môi trường
func LoadAppConfig() AppConfig {
    port, _ := strconv.Atoi(getEnvOrDefault("PORT", "8080"))
    
    return AppConfig{
        Port:      port,
        JWTSecret: getEnvOrDefault("JWT_SECRET", ""),
        GinMode:   getEnvOrDefault("GIN_MODE", "debug"),
        LogLevel:  getEnvOrDefault("LOG_LEVEL", "INFO"),
    }
}

// getEnvOrDefault lấy giá trị từ biến môi trường hoặc giá trị mặc định
func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
