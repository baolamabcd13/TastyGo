package main

import (
    "log"
    "os"

    "github.com/yourusername/tastygo/internal/database"
    "github.com/yourusername/tastygo/internal/models"
)

func main() {
    // Khởi tạo database
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "tastygo.db"
    }

    err := database.InitDB(dbPath)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    // Tạo SuperAdmin nếu chưa có
    var count int64
    database.DB.Model(&models.User{}).Where("role = ?", models.RoleSuperAdmin).Count(&count)
    
    if count == 0 {
        superAdmin := models.User{
            Email:    "superadmin@tastygo.com",
            Username: "superadmin",
            Role:     models.RoleSuperAdmin,
            Active:   true,
            Profile: models.UserProfile{
                FullName: "Super Admin",
            },
        }
        
        err = superAdmin.SetPassword("admin123")
        if err != nil {
            log.Fatalf("Failed to set password: %v", err)
        }
        
        result := database.DB.Create(&superAdmin)
        if result.Error != nil {
            log.Fatalf("Failed to create superadmin: %v", result.Error)
        }
        
        log.Println("Created default superadmin account")
    }

    // Tạo Admin test nếu chưa có
    database.DB.Model(&models.User{}).Where("email = ?", "testadmin@tastygo.com").Count(&count)
    
    if count == 0 {
        admin := models.User{
            Email:    "testadmin@tastygo.com",
            Username: "testadmin",
            Role:     models.RoleAdmin,
            Active:   true,
            Profile: models.UserProfile{
                FullName: "Test Admin",
                Phone:    "1234567890",
            },
        }
        
        err = admin.SetPassword("admin123")
        if err != nil {
            log.Fatalf("Failed to set password: %v", err)
        }
        
        result := database.DB.Create(&admin)
        if result.Error != nil {
            log.Fatalf("Failed to create test admin: %v", result.Error)
        }
        
        log.Println("Created test admin account")
    }

    log.Println("Seed completed successfully")
}
