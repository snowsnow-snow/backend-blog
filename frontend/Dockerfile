# --- 阶段 1: 构建 ---
FROM node:20-alpine AS builder
WORKDIR /app
# 先复制 package 文件以利用缓存
COPY package*.json ./
RUN npm install
# 复制前端源码
COPY . .
# 构建项目
RUN npm run build

# --- 阶段 2: 运行时 (Nginx) ---
FROM nginx:alpine
# 复制自定义 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf
# 从构建阶段复制构建好的静态资源
COPY --from=builder /app/dist /usr/share/nginx/html
# 暴露 80 端口
EXPOSE 80
ENTRYPOINT ["nginx", "-g", "daemon off;"]
