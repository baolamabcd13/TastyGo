package models

import (
    "time"
)

type Session struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    uint      `gorm:"index;not null" json:"user_id"`
    Token     string    `gorm:"uniqueIndex;not null" json:"token"`
    ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
}