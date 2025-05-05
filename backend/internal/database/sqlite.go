package database

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"os"

	"github.com/yourusername/tastygo/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error
	
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return err
	}
	
	// Migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.UserProfile{}, &models.Session{}, &models.ActivityLog{})
	if err != nil {
		return err
	}
	
	// Check if superadmin exists, if not create one
	var count int64
	DB.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count)
	
	if count == 0 {
		// Tạo mật khẩu ngẫu nhiên nếu không phải môi trường development
		defaultPassword := "admin123"
		isProduction := os.Getenv("GIN_MODE") == "release"
		
		if isProduction {
			// Tạo mật khẩu ngẫu nhiên cho môi trường production
			randomBytes := make([]byte, 12)
			if _, err := rand.Read(randomBytes); err == nil {
				defaultPassword = base64.URLEncoding.EncodeToString(randomBytes)
				log.Printf("Generated random password for superadmin: %s", defaultPassword)
				log.Println("IMPORTANT: Save this password securely. It will not be shown again.")
			}
		}
		
		superAdmin := models.User{
			Email:    "superadmin@tastygo.com",
			Username: "superadmin",
			Role:     models.RoleSuperAdmin,
			Active:   true,
			Profile: models.UserProfile{
				FullName: "Super Admin",
			},
		}
		
		err = superAdmin.SetPassword(defaultPassword)
		if err != nil {
			return err
		}
		
		result := DB.Create(&superAdmin)
		if result.Error != nil {
			return result.Error
		}
		
		log.Println("Created default superadmin account")
	}
	
	return nil
}
