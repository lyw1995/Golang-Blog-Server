# 构建docker镜像
.PHONY: build
build:
	cd build && bash build.sh

# 启动容器
.PHONY: up
up:
	cd build && docker-compose up -d

# 停止并删除容器
.PHONY: down
down:
	cd build && docker-compose down

# 重新生成api文档
.PHONY: api
api:
	cd api && swagger generate spec -o  ../docs/swaggerui/swagger.json --scan-models && swagger serve -F=swagger ../docs/swaggerui/swagger.json