package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/yourusername/tastygo/internal/api"
)

func TestLoginEndpoint(t *testing.T) {
    // Khởi tạo server
    router := api.NewServer()
    
    // Tạo request đăng nhập
    loginData := map[string]string{
        "email":    "superadmin@tastygo.com",
        "password": "admin123",
    }
    jsonData, _ := json.Marshal(loginData)
    
    req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    // Thực hiện request
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Kiểm tra kết quả
    if w.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    if _, exists := response["token"]; !exists {
        t.Error("Expected response to contain token")
    }
}
