package models

import (
	"time"
)

type ActivityType string

const (
    ActivityLogin          ActivityType = "login"
    ActivityLogout         ActivityType = "logout"
    ActivityCreateUser     ActivityType = "create_user"
    ActivityResetPassword  ActivityType = "reset_password"
    ActivityUpdateStatus   ActivityType = "update_status"
    ActivityUnlockAccount  ActivityType = "unlock_account"
)

type ActivityLog struct {
    ID          uint         `gorm:"primarykey" json:"id"`
    UserID      uint         `gorm:"index;not null" json:"user_id"`
    ActivityType ActivityType `gorm:"not null" json:"activity_type"`
    Description string       `json:"description"`
    IPAddress   string       `json:"ip_address"`
    UserAgent   string       `json:"user_agent"`
    CreatedAt   time.Time    `json:"created_at"`
}
