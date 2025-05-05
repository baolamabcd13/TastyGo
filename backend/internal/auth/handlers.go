package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/tastygo/internal/cache"
	"github.com/yourusername/tastygo/internal/database"
	"github.com/yourusername/tastygo/internal/models"
	"github.com/yourusername/tastygo/internal/pagination"
)

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type UserResponse struct {
    ID       uint        `json:"id"`
    Email    string      `json:"email"`
    Username string      `json:"username"`
    Role     models.Role `json:"role"`
}

type ResetPasswordRequest struct {
    UserID   uint   `json:"user_id" binding:"required"`
    Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserStatusRequest struct {
    UserID uint `json:"user_id" binding:"required"`
    Active bool `json:"active"`
}

type ActivityLogResponse struct {
    ID           uint               `json:"id"`
    UserID       uint               `json:"user_id"`
    Username     string             `json:"username"`
    ActivityType models.ActivityType `json:"activity_type"`
    Description  string             `json:"description"`
    IPAddress    string             `json:"ip_address"`
    UserAgent    string             `json:"user_agent"`
    CreatedAt    time.Time          `json:"created_at"`
}

type UnlockAccountRequest struct {
    UserID uint `json:"user_id" binding:"required"`
}

func HandleLogin(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    token, err := Login(req.Email, req.Password, c.ClientIP(), c.GetHeader("User-Agent"))
    if err != nil {
        status := http.StatusUnauthorized
        if err == ErrUserNotFound {
            status = http.StatusNotFound
        }
        c.JSON(status, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"token": token})
}

func HandleLogout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    parts := strings.Split(authHeader, " ")
    if len(parts) == 2 {
        tokenString := parts[1]
        Logout(tokenString)
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func HandleGetProfile(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    // Tạo cache key
    cacheKey := "profile_" + strconv.FormatUint(uint64(userID.(uint)), 10)
    
    // Kiểm tra cache
    if cachedProfile, found := cache.Get(cacheKey); found {
        c.JSON(http.StatusOK, cachedProfile)
        return
    }
    
    var user models.User
    result := database.DB.Preload("Profile").First(&user, userID)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    
    response := gin.H{
        "id":       user.ID,
        "email":    user.Email,
        "username": user.Username,
        "role":     user.Role,
        "profile": gin.H{
            "full_name": user.Profile.FullName,
        },
    }
    
    // Lưu vào cache trong 5 phút
    cache.Set(cacheKey, response, 5*time.Minute)
    
    c.JSON(http.StatusOK, response)
}

func HandleCreateAdmin(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Only superadmin can create admin accounts
    user.Role = models.RoleAdmin
    
    if err := user.SetPassword(c.PostForm("password")); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    result := database.DB.Create(&user)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    
    // Ghi log tạo admin mới
    creatorID, _ := c.Get("user_id")
    LogActivity(creatorID.(uint), models.ActivityCreateUser, 
        fmt.Sprintf("Created admin user: %s (ID: %d)", user.Username, user.ID),
        c.ClientIP(), c.GetHeader("User-Agent"))
    
    c.JSON(http.StatusCreated, UserResponse{
        ID:       user.ID,
        Email:    user.Email,
        Username: user.Username,
        Role:     user.Role,
    })
}

func HandleResetPassword(c *gin.Context) {
    // Chỉ SuperAdmin mới có quyền reset password
    role, _ := c.Get("role")
    if role != models.RoleSuperAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can reset passwords"})
        return
    }
    
    var req ResetPasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    var user models.User
    result := database.DB.First(&user, req.UserID)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    
    // Không cho phép reset password của SuperAdmin khác
    if user.Role == models.RoleSuperAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "cannot reset password of another superadmin"})
        return
    }
    
    if err := user.SetPassword(req.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    database.DB.Save(&user)
    
    // Ghi log reset password
    adminID, _ := c.Get("user_id")
    LogActivity(adminID.(uint), models.ActivityResetPassword, 
        fmt.Sprintf("Reset password for user ID: %d", req.UserID),
        c.ClientIP(), c.GetHeader("User-Agent"))
    
    c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func HandleUpdateUserStatus(c *gin.Context) {
    // Chỉ SuperAdmin mới có quyền kích hoạt/vô hiệu hóa
    role, _ := c.Get("role")
    if role != models.RoleSuperAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can update user status"})
        return
    }
    
    var req UpdateUserStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    var user models.User
    result := database.DB.First(&user, req.UserID)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    
    // Không cho phép vô hiệu hóa SuperAdmin
    if user.Role == models.RoleSuperAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "cannot update status of a superadmin"})
        return
    }
    
    user.Active = req.Active
    database.DB.Save(&user)
    
    status := "activated"
    if !req.Active {
        status = "deactivated"
    }
    
    // Ghi log cập nhật trạng thái
    adminID, _ := c.Get("user_id")
    LogActivity(adminID.(uint), models.ActivityUpdateStatus, 
        fmt.Sprintf("User ID %d %s", req.UserID, status),
        c.ClientIP(), c.GetHeader("User-Agent"))
    
    c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user %s successfully", status)})
}

func HandleListAdmins(c *gin.Context) {
    // Chỉ SuperAdmin mới có quyền xem danh sách Admin
    role, _ := c.Get("role")
    if role != models.RoleSuperAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can view admin list"})
        return
    }
    
    // Lấy tham số phân trang
    params := pagination.Extract(c)
    
    var admins []models.User
    var total int64
    
    // Đếm tổng số admin
    query := database.DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin)
    countResult := query.Count(&total)
    if countResult.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": countResult.Error.Error()})
        return
    }
    
    // Lấy danh sách admin có phân trang
    result := database.DB.Preload("Profile").
        Where("role = ?", models.RoleAdmin).
        Offset((params.Page - 1) * params.PageSize).
        Limit(params.PageSize).
        Find(&admins)
    
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    
    // Cập nhật tổng số bản ghi
    params.Total = total
    
    // Chuyển đổi sang response format
    var adminResponses []UserResponse
    for _, admin := range admins {
        adminResponses = append(adminResponses, UserResponse{
            ID:       admin.ID,
            Email:    admin.Email,
            Username: admin.Username,
            Role:     admin.Role,
        })
    }
    
    c.JSON(http.StatusOK, pagination.NewResponse(adminResponses, params))
}

func HandleGetActivityLogs(c *gin.Context) {
    // Chỉ SuperAdmin mới có quyền xem logs
    role, _ := c.Get("role")
    if role != models.RoleSuperAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can view activity logs"})
        return
    }
    
    // Lấy tham số phân trang
    params := pagination.Extract(c)
    
    var logs []models.ActivityLog
    var total int64
    
    // Hỗ trợ lọc theo user_id
    userID := c.Query("user_id")
    query := database.DB.Model(&models.ActivityLog{}).Order("created_at DESC")
    
    if userID != "" {
        query = query.Where("user_id = ?", userID)
    }
    
    // Đếm tổng số logs
    countResult := query.Count(&total)
    if countResult.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": countResult.Error.Error()})
        return
    }
    
    // Lấy logs có phân trang
    result := query.Offset((params.Page - 1) * params.PageSize).
        Limit(params.PageSize).
        Find(&logs)
    
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
    
    // Cập nhật tổng số bản ghi
    params.Total = total
    
    c.JSON(http.StatusOK, pagination.NewResponse(logs, params))
}

func HandleUnlockAccount(c *gin.Context) {
    // Chỉ SuperAdmin mới có quyền mở khóa tài khoản
    role, _ := c.Get("role")
    if role != models.RoleSuperAdmin {
        c.JSON(http.StatusForbidden, gin.H{"error": "only superadmin can unlock accounts"})
        return
    }
    
    var req UnlockAccountRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    var user models.User
    result := database.DB.First(&user, req.UserID)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    
    // Kiểm tra xem tài khoản có bị khóa không
    if user.LockedUntil == nil || user.LockedUntil.Before(time.Now()) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "account is not locked"})
        return
    }
    
    // Mở khóa tài khoản
    user.LockedUntil = nil
    user.FailedLoginCount = 0
    database.DB.Save(&user)
    
    // Ghi log mở khóa tài khoản
    adminID, _ := c.Get("user_id")
    LogActivity(adminID.(uint), models.ActivityUnlockAccount, 
        fmt.Sprintf("Unlocked account for user ID: %d", req.UserID),
        c.ClientIP(), c.GetHeader("User-Agent"))
    
    c.JSON(http.StatusOK, gin.H{"message": "account unlocked successfully"})
}
