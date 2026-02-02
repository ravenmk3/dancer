# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates 用于 HTTPS
RUN apk --no-cache add ca-certificates

# 从 builder 复制二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/config.toml .

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./server"]
