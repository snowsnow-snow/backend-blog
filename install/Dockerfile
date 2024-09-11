FROM --platform=linux/amd64 alpine:latest

MAINTAINER snowsnowsnow

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 日志
RUN mkdir -p /var/log/go/backend-blog/logs
# 上传文件
RUN mkdir -p /home/snowsnowsnow/file/public
# 数据库文件
RUN mkdir -p /var/db/backend-blog
# 主程序
RUN mkdir -p /usr/bin/backend-blog

# 安装 libstdc++6
RUN apk add --no-cache libstdc++
# 安装 perl
RUN apk add perl

WORKDIR /usr/bin/backend-blog/

COPY backend_blog ./
COPY config-prod.yaml ./

EXPOSE 800

ENTRYPOINT ["./backend_blog","-env=prod"]

