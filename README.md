本项目已经预先创建了一系列文件夹划分出下列模块:

1. api文件夹就是MVC框架的controller，负责协调各部件完成任务
2. model文件夹负责存储数据库模型和数据库操作相关的代码
3. service负责处理比较复杂的业务，把业务代码模型化可以有效提高业务代码的质量（比如用户注册，充值，下单等）
4. serializer储存通用的json模型，把model得到的数据库模型转换成api需要的json对象
5. cache负责redis缓存相关的代码
6. auth权限控制文件夹
7. util一些通用的小工具
8. conf放一些静态存放的配置文件，其中locales内放置翻译相关的配置文件

## MySQL、Redis、Mongo
```shell
docker pull mysql
docker pull redis
docker pull mongo
```

## Godotenv

项目在启动的时候依赖以下环境变量，但是在也可以在项目根目录创建`.env`文件设置环境变量便于使用(建议开发环境使用)

```shell
MYSQL_DSN="db_user@/db_name?charset=utf8&parseTime=True&loc=Local"
REDIS_ADDR="127.0.0.1:6379"
REDIS_PW=""
REDIS_DB=""
SESSION_SECRET="setOnProducation"
GIN_MODE="debug"
LOG_LEVEL="debug"

TENCENT_OSS_APP_ID=""
TENCENT_OSS_REGION=""
TENCENT_OSS_BUCKET=""
TENCENT_OSS_SECRET_ID=""
TENCENT_OSS_SECRET_KEY=""


MONGO_DB_DSN="mongodb://mongoadmin:mongopassword@localhost:27017"
MONGO_DB_NAME="giligili"
```

```shell
# 下载所需镜像
docker pull mysql
docker pull redis
docker pull mongo

# 新建文件夹，位置随便
# 主要是将上面的镜像数据挂载到文件夹下，保证停止/删除后还能找到（持久化）
mkdir -p ~/docker/data/mysql
mkdir -p ~/docker/data/redis
mkdir -p ~/docker/data/mongo
```

```shell
# 启动 mysql
docker run -d --name mysql \
  -v ~/data/mysql:/var/lib/mysql \
  -e MYSQL_ROOT_PASSWORD=mysql \
  -e MYSQL_DATABASE=mysql \
  -e MYSQL_USER=mysql \
  -e MYSQL_PASSWORD=mysql \
  -p 3306:3306
  mysql

# 创建数据库
docker exec -it mysql mysql -u root -pmysql -e "CREATE DATABASE giligili;"
```

```shell
# 启动 redis
docker run --name redis -p 6379:6379 -v ~/data/redis:/data -d redis
```

```shell
# 启动 mongo
docker run -d --name mongodb \
  -v ~/data/mongodb:/data/db \
  -e MONGO_INITDB_ROOT_USERNAME=mongoadmin \
  -e MONGO_INITDB_ROOT_PASSWORD=mongopassword \
  -p 27017:27017 \
  mongo

```

## Go Mod

本项目使用[Go Mod](https://github.com/golang/go/wiki/Modules)管理依赖。

```shell
go run main.go // 自动安装
```

## 运行

```shell
go run main.go
```

项目运行后启动在3000端口（可以修改，参考gin文档）