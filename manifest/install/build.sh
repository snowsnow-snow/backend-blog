#!/bin/bash

#read -p "请输入要打包的镜像名称: " -r NAME
read -p "请输入要打包的镜像版本号: " -r VERSION

# 定义变量
#IMAGE_NAME="$NAME"
IMAGE_NAME="backend-blog"
TAG="$VERSION"

# 指定编译版本为 linux/amd64
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build -o backend_blog ../../.
# 将配置文件拷贝到同级目录
cp ../../config/config-prod.yaml .

echo "项目打包完成"

echo "镜像 $IMAGE_NAME:$TAG 开始构建"


# 判断当前镜像是否存在，如果存在就删除，构建新的镜像
if docker image inspect "$IMAGE_NAME:$TAG" &> /dev/null; then
  docker image rm "$IMAGE_NAME:$TAG"
fi

# 构建镜像
docker build -t "$IMAGE_NAME:$TAG" .


echo "镜像 $IMAGE_NAME:$TAG 构建完成"

echo "开始打包镜像..."

# 打包镜像
docker save -o "$IMAGE_NAME.tar" "$IMAGE_NAME:$TAG"

echo "镜像打包完成，镜像文件：$IMAGE_NAME.tar"

rm config-prod.yaml


cp "$IMAGE_NAME.tar" /Users/snowsnowsnow/DockerTar