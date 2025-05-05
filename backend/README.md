# TastyGo Backend API

## Cài đặt và chạy

### Phát triển

1. Cài đặt Go 1.21 hoặc cao hơn
2. Clone repository
3. Cài đặt dependencies:
   ```
   go mod download
   ```
4. Chạy server:
   ```
   go run cmd/server/main.go
   ```

### Triển khai với Docker

1. Build và chạy với Docker Compose:

   ```
   docker-compose up -d
   ```

   Lệnh này sẽ khởi động container API:

   - **tastygo-api**: API backend

2. Hoặc build và chạy thủ công:
   ```
   docker build -t tastygo-api .
   docker run -p 8080:8080 -e JWT_SECRET=your_secure_secret -v /path/to/data:/data tastygo-api
   ```

## Biến môi trường

- `PORT`: Cổng server (mặc định: 8080)
- `DB_PATH`: Đường dẫn đến file SQLite (mặc định: tastygo.db)
- `JWT_SECRET`: Secret key cho JWT (bắt buộc trong môi trường production)
- `GIN_MODE`: Chế độ Gin framework (development/release)

## Tài khoản mặc định

- SuperAdmin:
  - Email: superadmin@tastygo.com
  - Password: admin123 (nên đổi trong môi trường production)

## API Endpoints

### Authentication

- `POST /api/auth/login`: Đăng nhập
- `POST /api/auth/logout`: Đăng xuất

### User Management

- `GET /api/profile`: Xem thông tin cá nhân
- `POST /api/admin/users`: Tạo tài khoản Admin (SuperAdmin only)
- `GET /api/admin/users/admins`: Xem danh sách Admin (SuperAdmin only)
- `POST /api/admin/users/reset-password`: Đặt lại mật khẩu (SuperAdmin only)
- `POST /api/admin/users/update-status`: Kích hoạt/vô hiệu hóa tài khoản (SuperAdmin only)
- `POST /api/admin/users/unlock-account`: Mở khóa tài khoản bị khóa (SuperAdmin only)
- `GET /api/admin/logs`: Xem lịch sử hoạt động (SuperAdmin only)

## Postman Collection

Dự án bao gồm file Postman Collection để dễ dàng test API:

- File: `tastygo_api.postman_collection.json`
- Cách sử dụng:
  1. Import vào Postman
  2. Tạo environment với biến `base_url` (ví dụ: http://localhost:8080)
  3. Chạy request "Login SuperAdmin" để lấy token
  4. Các request khác sẽ tự động sử dụng token này

## Phát triển

### Cấu trúc thư mục

```
backend/
├── cmd/                # Entry points
│   └── server/         # API server
├── internal/           # Private application code
│   ├── api/            # API handlers và routes
│   ├── auth/           # Authentication và authorization
│   ├── database/       # Database setup và migrations
│   ├── models/         # Data models
│   └── pagination/     # Pagination utilities
├── Dockerfile          # Docker build file
├── docker-compose.yml  # Docker Compose configuration
└── go.mod              # Go modules
```

### Tính năng bảo mật

- JWT authentication
- Role-based access control
- Password hashing với bcrypt
- Rate limiting để ngăn chặn brute force
- Activity logging cho audit trail
