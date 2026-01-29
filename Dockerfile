# --- 阶段 1: 后端构建 ---
FROM golang:1.25-alpine AS builder
# 设置 Go 代理（中国镜像）
ENV GOPROXY=https://goproxy.cn,direct
# 安装 CGO 所需构建依赖
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download
# 复制其余后端代码
COPY . .
# 构建开启 CGO 的后端二进制文件
RUN CGO_ENABLED=1 GOOS=linux go build -o server cmd/server/main.go

# --- 阶段 2: 生产环境运行 ---
FROM alpine:latest
# 安装运行时依赖：ffmpeg 用于视频处理，exiftool 用于图像元数据获取
RUN apk add --no-cache \
    ffmpeg \
    exiftool \
    ca-certificates \
    tzdata \
    libc6-compat

WORKDIR /app

# 复制构建好的后端二进制文件
COPY --from=builder /app/server .
# 复制配置文件
COPY config ./config

# 暴露应用程序端口
EXPOSE 53123

# 设置默认环境为生产环境
ENTRYPOINT ["./server", "-env", "prod"]
