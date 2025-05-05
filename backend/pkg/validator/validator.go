package validator

import (
    "regexp"
    "unicode"
)

// IsValidEmail kiểm tra email có hợp lệ không
func IsValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    match, _ := regexp.MatchString(pattern, email)
    return match
}

// IsStrongPassword kiểm tra mật khẩu có đủ mạnh không
func IsStrongPassword(password string) bool {
    if len(password) < 8 {
        return false
    }
    
    var hasUpper, hasLower, hasNumber, hasSpecial bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    return hasUpper && hasLower && hasNumber && hasSpecial
}
