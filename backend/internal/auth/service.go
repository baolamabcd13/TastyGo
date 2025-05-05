package auth

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/tastygo/internal/database"
	"github.com/yourusername/tastygo/internal/logging"
	"github.com/yourusername/tastygo/internal/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	jwtSecret             = []byte(os.Getenv("JWT_SECRET")) // Lấy từ biến môi trường
)

// Thêm hàm khởi tạo JWT secret
func InitJWTSecret() {
	secretEnv := os.Getenv("JWT_SECRET")
	if secretEnv == "" {
		// Tạo secret ngẫu nhiên nếu không có trong biến môi trường
		randomBytes := make([]byte, 32)
		if _, err := rand.Read(randomBytes); err != nil {
			log.Fatal("Failed to generate random JWT secret")
		}
		jwtSecret = randomBytes
		log.Println("WARNING: Using randomly generated JWT secret. Set JWT_SECRET environment variable for production.")
	} else {
		jwtSecret = []byte(secretEnv)
	}
}

type TokenClaims struct {
	UserID uint        `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

func LogActivity(userID uint, activityType models.ActivityType, description string, ipAddress, userAgent string) {
	log := models.ActivityLog{
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	}
	
	database.DB.Create(&log)
}

func Login(email, password string, ipAddress, userAgent string) (string, error) {
	var user models.User
	
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		logging.Warn("Login attempt failed: user not found", map[string]interface{}{
			"email": email,
			"ip":    ipAddress,
		})
		return "", ErrUserNotFound
	}
	
	// Kiểm tra tài khoản có bị khóa tạm thời không
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		remainingTime := time.Until(*user.LockedUntil).Minutes()
		logging.Warn("Login attempt on locked account", map[string]interface{}{
			"user_id":     user.ID,
			"email":       user.Email,
			"ip":          ipAddress,
			"locked_until": user.LockedUntil,
		})
		return "", fmt.Errorf("account is temporarily locked. Try again in %.0f minutes", remainingTime)
	}
	
	// Kiểm tra mật khẩu
	if !user.CheckPassword(password) {
		// Tăng số lần đăng nhập sai
		now := time.Now()
		user.LastFailedLogin = &now
		user.FailedLoginCount++
		
		// Nếu sai 5 lần liên tiếp, khóa tài khoản 30 phút
		if user.FailedLoginCount >= 5 {
			lockTime := time.Now().Add(30 * time.Minute)
			user.LockedUntil = &lockTime
			user.FailedLoginCount = 0
			logging.Warn("Account locked due to multiple failed login attempts", map[string]interface{}{
				"user_id":     user.ID,
				"email":       user.Email,
				"ip":          ipAddress,
				"locked_until": lockTime,
			})
		}
		
		database.DB.Save(&user)
		logging.Warn("Login failed: invalid password", map[string]interface{}{
			"user_id":     user.ID,
			"email":       user.Email,
			"ip":          ipAddress,
			"failed_count": user.FailedLoginCount,
		})
		return "", ErrInvalidCredentials
	}
	
	// Kiểm tra tài khoản có active không
	if !user.Active {
		return "", errors.New("account is disabled")
	}
	
	// Reset số lần đăng nhập sai
	user.FailedLoginCount = 0
	user.LockedUntil = nil
	
	// Update last login
	now := time.Now()
	user.LastLogin = &now
	database.DB.Save(&user)
	
	// Generate token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &TokenClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	
	// Save session
	session := models.Session{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: expirationTime,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	
	database.DB.Create(&session)
	
	// Ghi log đăng nhập thành công
	LogActivity(user.ID, models.ActivityLogin, "Successful login", ipAddress, userAgent)
	
	logging.Info("User logged in successfully", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"ip":      ipAddress,
	})
	
	return tokenString, nil
}

func ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		// Check if token is in sessions table
		var session models.Session
		result := database.DB.Where("token = ? AND expires_at > ?", tokenString, time.Now()).First(&session)
		if result.Error != nil {
			return nil, errors.New("invalid or expired session")
		}
		
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

func Logout(tokenString string) error {
	// Lấy thông tin user từ token
	claims, _ := ValidateToken(tokenString)
	if claims != nil {
		// Ghi log đăng xuất
		LogActivity(claims.UserID, models.ActivityLogout, "User logged out", "", "")
	}
	
	return database.DB.Where("token = ?", tokenString).Delete(&models.Session{}).Error
}
