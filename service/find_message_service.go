package service

import (
	"context"
	"fmt"
	"giligili/conf"
	"giligili/model/ws"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

// 插入聊天记录
func InsertMsg(database string, id string, content string, read uint, expire int64) (err error) {
	collection := conf.MongoDBClient.Database(database).Collection(id)

	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}

	// 如果不知道该使用什么 context，可以通过 context.TODO() 产生 context
	_, err = collection.InsertOne(context.TODO(), comment)
	return
}

// findInCollection 从集合中查找结果
func findInCollection(collection *mongo.Collection, pageSize int) ([]ws.Trainer, error) {
	var results []ws.Trainer

	// 定义查找操作的选项
	findOptions := options.Find().SetSort(bson.D{{"startTime", -1}}).SetLimit(int64(pageSize))
	// findOptions := options.Find().SetSort(bson.M{"startTime": -1}).SetLimit(int64(pageSize))

	// 使用指定的集合和选项执行查找操作
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// 将结果从游标解码到 'results' 切片中
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// FindMany 找多个结果
func FindMany(database string, sendId string, id string, time int64, pageSize int) (results []ws.Result, err error) {
	// 根据提供的数据库和集合名称获取集合
	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)

	// 从 sendIdCollection 中查找结果
	resultsYou, err := findInCollection(sendIdCollection, pageSize)
	if err != nil {
		return nil, err
	}

	// 从 idCollection 中查找结果
	resultsMe, err := findInCollection(idCollection, pageSize)
	if err != nil {
		return nil, err
	}

	// 合并并排序来自两个集合的结果
	results, _ = appendAndSort(resultsMe, resultsYou)
	return
}

func FirsFindtMsg(database string, sendId string, id string) (results []ws.Result, err error) {
	// 首次查询(把对方发来的所有未读都取出来)
	var resultsMe []ws.Trainer
	var resultsYou []ws.Trainer

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendId)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)

	filter := bson.M{"read": bson.M{
		"&all": []uint{0},
	}}

	sendIdCursor, err := sendIdCollection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{
		"startTime", 1}}), options.Find().SetLimit(1))

	if err != nil {
		log.Println("sendIdCursor err", err)
	}

	if sendIdCursor == nil {
		return
	}

	var unRead []ws.Trainer
	err = sendIdCursor.All(context.TODO(), &unRead)

	if err != nil {
		log.Println("sendIdCursor err", err)
	}

	if len(unRead) > 0 {
		timeFilter := bson.M{
			"startTime": bson.M{
				"$gte": unRead[0].StartTime,
			},
		}

		sendIdTimeCursor, err := sendIdCollection.Find(context.TODO(), timeFilter)
		if err != nil {
			log.Println("sendIdTimeCursor err", err)
		}

		idTimeCursor, err := idCollection.Find(context.TODO(), timeFilter)
		if err != nil {
			log.Println("idTimeCursor err", err)
		}

		err = sendIdTimeCursor.All(context.TODO(), &resultsYou)
		err = idTimeCursor.All(context.TODO(), &resultsMe)
		results, err = appendAndSort(resultsMe, resultsYou)
	} else {
		results, err = FindMany(database, sendId, id, 9999999999, 10)
	}

	overTimeFilter := bson.D{
		{"$and", bson.A{
			bson.D{{"endTime", bson.M{"&lt": time.Now().Unix()}}},
			bson.D{{"read", bson.M{"$eq": 1}}},
		}},
	}

	_, _ = sendIdCollection.DeleteMany(context.TODO(), overTimeFilter)
	_, _ = idCollection.DeleteMany(context.TODO(), overTimeFilter)
	// 将所有的维度设置为已读
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{"read": 1},
	})
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"&set": bson.M{"ebdTime": time.Now().Unix() + int64(time.Hour*24*30*3)},
	})

	return
}

// 负责将 ws.Trainer 类型转换为 ws.Result 类型，并接受一个额外的参数 from 来指定消息是来自 "me" 还是 "you"。这样我们就可以在 appendAndSort 函数中复用 trainerToResult 函数，减少了代码重复，并且使得 appendAndSort 函数更加简洁。
func trainerToResult(trainer ws.Trainer, from string) ws.Result {
	sendsort := SendSortMsg{
		Content:  trainer.Content,
		Read:     trainer.Read,
		CreateAt: trainer.StartTime,
	}

	return ws.Result{
		StartTime: trainer.StartTime,
		Msg:       fmt.Sprintf("%v", sendsort),
		From:      from,
	}
}

func appendAndSort(resultsMe, resultsYou []ws.Trainer) (results []ws.Result, err error) {
	for _, resultMe := range resultsMe {
		results = append(results, trainerToResult(resultMe, "me"))
	}

	for _, resultYou := range resultsYou {
		results = append(results, trainerToResult(resultYou, "you"))
	}

	// 最后进行排序
	sort.Slice(results, func(i, j int) bool { return results[i].StartTime < results[j].StartTime })
	return results, nil
}
