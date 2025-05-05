FROM golang:1.23 AS builder

WORKDIR /app

# Sao chép go.mod và go.sum
COPY go.mod go.sum ./
RUN go mod download

# Sao chép mã nguồn
COPY . .

# Build ứng dụng
RUN CGO_ENABLED=1 GOOS=linux go build -a -o tastygo ./cmd/server

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Sao chép binary từ builder
COPY --from=builder /app/tastygo .

# Tạo thư mục cho database
RUN mkdir -p /data

# Thiết lập biến môi trường
ENV DB_PATH=/data/tastygo.db
ENV JWT_SECRET=change_this_in_production
ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./tastygo"]
