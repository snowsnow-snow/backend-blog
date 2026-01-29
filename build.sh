#!/bin/bash

# 配置变量
BACKEND_IMAGE="my-blog-backend"
FRONTEND_IMAGE="my-blog-frontend"
TAG="latest"

echo "=== 博客系统多容器构建脚本 ==="

# 提示：对于 1GB 内存的服务器，强烈建议在本地机器上运行此脚本，
# 然后将镜像推送到仓库或使用 'docker save' 进行传输。
# 直接在服务器上构建可能会导致 OOM（内存溢出）错误。

echo "正在构建后端镜像: $BACKEND_IMAGE:$TAG"
docker build -t $BACKEND_IMAGE:$TAG -f Dockerfile .
if [ $? -ne 0 ]; then
    echo "后端构建失败！"
    exit 1
fi

echo "正在构建前端镜像: $FRONTEND_IMAGE:$TAG"
docker build -t $FRONTEND_IMAGE:$TAG -f frontend/Dockerfile frontend/
if [ $? -ne 0 ]; then
    echo "前端构建失败！"
    exit 1
fi

echo "=== 构建成功 ==="
echo "要启动系统，请运行："
echo "  docker-compose up -d"
