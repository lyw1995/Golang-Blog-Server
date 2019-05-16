#!/bin/bash

#当前时间
current_time=` date +"%Y%m%d%H%M%S"`

# 删除上一次的 可执行文件
rm blog_server-linux-amd64-* >& /dev/null

# 编译linux可执行文件, 进行制作镜像
export GOOS=linux
go build  -ldflags="-s -w" -o blog_server-linux-amd64-${current_time} ../cmd/main.go
export GOOS=darwin

# Dockerfile生成
echo "FROM iron/base" > Dockerfile
echo "MAINTAINER Track \"24245@163.com\"" >> Dockerfile
echo "LABEL build_time=\"$current_time\"" >> Dockerfile
echo "RUN mkdir -p /home/blogserver/api_docs && mkdir -p /var/log/blogserver && mkdir -p /home/blogserver/config && mkdir -p /home/blogserver/images" >> Dockerfile
echo "COPY ./docs/swaggerui/ /home/blogserver/api_docs/" >> Dockerfile
echo "COPY ./config/config.toml /home/blogserver/config/" >> Dockerfile
echo "ENV APP_CONFIG_PATH /home/blogserver/config/config.toml" >> Dockerfile
echo "EXPOSE 8888" >> Dockerfile
echo "ADD  ./build/blog_server-linux-amd64-$current_time /" >> Dockerfile
echo "ENTRYPOINT [\"./blog_server-linux-amd64-$current_time\"]" >> Dockerfile

EXPOSE 8888

# build image , 之前的用脚本删除之类的..
docker image rm -f yinguiw/blog_server

docker build -t yinguiw/blog_server:1.0 -f ../build/Dockerfile ..

# clean
rm blog_server-linux-amd64-* >& /dev/null