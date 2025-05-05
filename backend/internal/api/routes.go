package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/tastygo/internal/auth"
	"github.com/yourusername/tastygo/internal/models"
)

func SetupRoutes(router *gin.Engine) {
    // Public routes
    router.POST("/api/auth/login", auth.HandleLogin)
    
    // Protected routes
    authRoutes := router.Group("/api")
    authRoutes.Use(auth.AuthMiddleware())
    {
        authRoutes.POST("/auth/logout", auth.HandleLogout)
        authRoutes.GET("/profile", auth.HandleGetProfile)
        
        // Admin routes
        adminRoutes := authRoutes.Group("/admin")
        adminRoutes.Use(auth.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
        {
            // Routes accessible by both admin and superadmin
            adminRoutes.GET("/dashboard", func(c *gin.Context) {
                c.JSON(200, gin.H{"message": "Admin dashboard"})
            })
        }
        
        // SuperAdmin routes
        superAdminRoutes := authRoutes.Group("/admin")
        superAdminRoutes.Use(auth.RoleMiddleware(models.RoleSuperAdmin))
        {
            // Routes accessible only by superadmin
            superAdminRoutes.POST("/users", auth.HandleCreateAdmin)
            superAdminRoutes.GET("/users/admins", auth.HandleListAdmins)
            superAdminRoutes.POST("/users/reset-password", auth.HandleResetPassword)
            superAdminRoutes.POST("/users/update-status", auth.HandleUpdateUserStatus)
            superAdminRoutes.POST("/users/unlock-account", auth.HandleUnlockAccount) // Thêm route mới
            superAdminRoutes.GET("/logs", auth.HandleGetActivityLogs)
        }
    }
}
