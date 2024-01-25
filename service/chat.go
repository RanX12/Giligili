package service

import (
	"context"
	"encoding/json"
	"fmt"
	"giligili/cache"
	"giligili/conf"
	"giligili/pkg/e"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 发送消息的类型
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// 回复的消息
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// 用户类
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

// Message 信息转JSON (包括：发送者、接收者、内容)
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func createId(uid, toUid string) string {
	return uid + "->" + toUid
}

func WsHandler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin 解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成 ws 协议

	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 创建一个用户实例
	client := &Client{
		ID:     createId(uid, toUid),
		SendID: createId(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}

	// 用户注册到用户管理上
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		// c.Socket.PongHandler() 设置一个具体的处理函数，而不是直接调用
		c.Socket.SetPongHandler(func(data string) error {
			log.Println("SetPongHandler data: ", data)
			return nil
		})
		var sendMsg SendMsg

		// 读取消息
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// 如果期望的消息类型是文本消息，则尝试解析 JSON
		if messageType == websocket.TextMessage {
			// 解析 JSON 格式的消息到 SendMsg 结构体
			err := json.Unmarshal(message, &sendMsg)
			if err != nil {
				log.Println("Error unmarshalling message:", err)
				continue // 解析失败时，继续下一次循环
			}

			log.Printf("Received message: %+v\n", sendMsg)
		} else {
			log.Printf("Received binary message, messageType: %d\n", messageType)
			// 如果消息是二进制的，您可以在这里处理它
		}

		// err = c.Socket.ReadJSON(&sendMsg) // 读取json格式，如果不是json格式，会报错
		if err != nil {
			log.Println("数据格式不正确", err)

			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		if sendMsg.Type == 1 {
			ctx := context.Background()
			r1, _ := cache.RedisClient.Get(ctx, c.ID).Result()
			r2, _ := cache.RedisClient.Get(ctx, c.SendID).Result()

			log.Println("是否被限制： ", r1 >= "3" && r2 == "")

			if r1 >= "3" && r2 == "" { // 限制单聊
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "超过 3 条对方未回复，请等对方回复",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				_, _ = cache.RedisClient.Expire(ctx, c.ID, time.Hour*24*30).Result() // 防止重复骚扰，未建立连接刷新过期时间一个月
				continue
			} else {
				cache.RedisClient.Incr(ctx, c.ID)
				_, _ = cache.RedisClient.Expire(ctx, c.ID, time.Hour*24*90).Result() // 防止过快“分手”，建立连接三个月过期
			}

			log.Println(c.ID, "发送消息", sendMsg.Content)
			// 进行消息广播
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 { // 拉取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) // 传送来时间
			if err != nil {
				timeT = 999999999
			}

			results, _ := FindMany(conf.MongoDBName, c.SendID, c.ID, int64(timeT), 10)
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}

				msg, _ := json.Marshal(replyMsg)
				c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}

			// 循环发送历史消息
			// TODO 整合历史记录
			for _, result := range results {
				replyMsg := &ReplyMsg{
					From:    result.From,
					Code:    e.WebsocketHistroy,
					Content: fmt.Sprintf("%s", result.Msg),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		} else if sendMsg.Type == 3 {
			results, err := FirsFindtMsg(conf.MongoDBName, c.SendID, c.ID)
			if err != nil {
				log.Println(err)
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: fmt.Sprintf("%s", result.Msg),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println(c.ID, "接受消息:", string(message))
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
