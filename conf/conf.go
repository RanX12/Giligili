package conf

import (
	"context"
	"giligili/cache"
	"giligili/model"
	"giligili/util"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoDBClient *mongo.Client

	MongoDBName string
)

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	MongoDBName = os.Getenv("MONGO_DB_NAME")

	// 读取翻译文件
	if err := LoadLocales("conf/locales/zh-cn.yaml"); err != nil {
		util.Log().Panic("翻译文件加载失败", err)
	}

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
	MongoDB()
}

func MongoDB() {
	// 设置mongoDB客户端连接信息
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_DB_DSN"))
	var err error

	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		util.Log().Panic("MongoDBClient加载失败", err)
	}
	err = MongoDBClient.Ping(context.TODO(), nil)
	if err != nil {
		util.Log().Panic("MongoDBClient加载失败", err)
	}
	util.Log().Info("MongoDB Connect")
}
