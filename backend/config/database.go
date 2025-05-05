package config

// DBConfig chứa cấu hình database
type DBConfig struct {
    Path string
}

// LoadDBConfig tải cấu hình database từ biến môi trường
func LoadDBConfig() DBConfig {
    return DBConfig{
        Path: getEnvOrDefault("DB_PATH", "tastygo.db"),
    }
}
