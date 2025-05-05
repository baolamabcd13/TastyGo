package pagination

import (
    "strconv"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// Params chứa các tham số phân trang
type Params struct {
    Page     int   `json:"page"`
    PageSize int   `json:"page_size"`
    Total    int64 `json:"total"`
}

// Response là cấu trúc phản hồi có phân trang
type Response struct {
    Data       interface{} `json:"data"`
    Pagination Params      `json:"pagination"`
}

// Extract trích xuất tham số phân trang từ request
func Extract(c *gin.Context) Params {
    // Giá trị mặc định
    page := 1
    pageSize := 10

    // Lấy tham số từ query string
    pageStr := c.Query("page")
    pageSizeStr := c.Query("page_size")

    // Chuyển đổi sang số
    if pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }

    if pageSizeStr != "" {
        if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
            pageSize = ps
        }
    }

    return Params{
        Page:     page,
        PageSize: pageSize,
    }
}

// Apply áp dụng phân trang cho query GORM
func Apply(db *gorm.DB, params Params) *gorm.DB {
    offset := (params.Page - 1) * params.PageSize
    return db.Offset(offset).Limit(params.PageSize)
}

// NewResponse tạo phản hồi có phân trang
func NewResponse(data interface{}, params Params) Response {
    return Response{
        Data:       data,
        Pagination: params,
    }
}