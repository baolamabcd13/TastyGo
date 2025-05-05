package utils

import (
    "crypto/rand"
    "encoding/base64"
)

// GenerateRandomString tạo chuỗi ngẫu nhiên với độ dài xác định
func GenerateRandomString(length int) (string, error) {
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
