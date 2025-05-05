package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
    RoleSuperAdmin Role = "superadmin"
    RoleAdmin     Role = "admin"
    RoleCustomer  Role = "customer"
)

type User struct {
    ID                uint           `gorm:"primarykey" json:"id"`
    Email             string         `gorm:"uniqueIndex;not null" json:"email"`
    Username          string         `gorm:"uniqueIndex;not null" json:"username"`
    PasswordHash      string         `gorm:"not null" json:"-"`
    Role              Role           `gorm:"not null" json:"role"`
    Active            bool           `gorm:"default:true" json:"active"`
    CreatedAt         time.Time      `json:"created_at"`
    UpdatedAt         time.Time      `json:"updated_at"`
    LastLogin         *time.Time     `json:"last_login"`
    FailedLoginCount  int            `gorm:"default:0" json:"-"`
    LastFailedLogin   *time.Time     `json:"-"`
    LockedUntil       *time.Time     `json:"-"`
    DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
    Profile           UserProfile    `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

type UserProfile struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
    FullName  string    `json:"full_name"`
    Phone     string    `json:"phone"`
    Address   string    `json:"address"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) SetPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.PasswordHash = string(hashedPassword)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
    return err == nil
}
