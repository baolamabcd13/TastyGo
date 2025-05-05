# TastyGo Admin Dashboard

Đây là dự án dashboard quản trị cho ứng dụng TastyGo, được xây dựng bằng Next.js, TypeScript, Tailwind CSS v3 và DaisyUI.

## Công nghệ sử dụng

- Next.js
- TypeScript
- Tailwind CSS v3
- DaisyUI
- React Hook Form
- Axios
- React Icons
- Zustand

## Cài đặt và chạy

### Yêu cầu

- Node.js (phiên bản 18 trở lên)
- Backend TastyGo đang chạy (mặc định ở cổng 8080)

### Các bước cài đặt

1. Clone repository
2. Cài đặt dependencies:

```bash
npm install
```

3. Chạy ứng dụng ở môi trường phát triển:

```bash
npm run dev
```

4. Mở [http://localhost:3000](http://localhost:3000) để xem kết quả.

### Tài khoản mặc định

- SuperAdmin:
  - Email: superadmin@tastygo.com
  - Password: admin123

## Cấu trúc dự án

```
dashboard/
├── public/              # Static files
│   └── tastygo.png      # Logo
├── src/
│   ├── app/             # App router pages
│   │   ├── login/       # Login page
│   │   └── dashboard/   # Dashboard pages
│   └── components/      # Reusable components
│       └── layout/      # Layout components
└── ...
```

## Tính năng

- Đăng nhập/đăng xuất
- Xem thông tin cá nhân
- Quản lý tài khoản Admin
- Xem nhật ký hoạt động
- Responsive design (desktop và mobile)

## Lưu ý

- Đảm bảo backend TastyGo đang chạy ở cổng 8080 trước khi sử dụng dashboard
- Thay đổi file `tastygo.png` trong thư mục `public` để sử dụng logo của riêng bạn
