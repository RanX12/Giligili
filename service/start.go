package service

import (
	"encoding/json"
	"fmt"
	"giligili/conf"
	"giligili/pkg/e"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var mutex sync.Mutex

func (manager *ClientManager) Start() {
	for {
		log.Println("<------监听管道通信----->")
		select {
		case conn := <-Manager.Register: // 建立连接
			log.Printf("建立新连接: %v", conn.ID)
			Manager.Clients[conn.ID] = conn
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已连接至服务器",
			}

			msg, err := json.Marshal(replyMsg)
			if err != nil {
				log.Println("Manager.Register json.Marshal err: ", err)
			}

			err = conn.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Manager.Register conn.Socket.WriteMessage err: ", err)
			}
		case conn := <-Manager.Unregister:
			log.Printf("连接失败:%v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "连接断开",
				}

				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		// 广播信息
		case broadcast := <-Manager.Broadcast:
			message := broadcast.Message
			sendId := broadcast.Client.SendID
			flag := false // 默认对方不在线

			log.Println("Manager.Clients : ", sendId)

			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}

				flag = true

				select {
				case conn.Send <- message:
				default:
					log.Println("close 了 conn")
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}

			id := broadcast.Client.ID
			if flag {
				log.Println("对方在线应答")
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答",
				}

				msg, err := json.Marshal(replyMsg)
				if err != nil {
					log.Printf("flag true Error marshalling reply message: %v\n", err)
				} else {
					err = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						log.Printf("flag true Error sending reply message: %v\n", err)
					}

				}
				err = InsertMsg(conf.MongoDBName, id, string(message), 1, int64(time.Hour*24*30*3))
				if err != nil {
					fmt.Println("flag true InsertOneMsg Err", err)
				}
			} else {
				log.Println("对方不在线")
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方不在线",
				}
				msg, err := json.Marshal(replyMsg)
				if err != nil {
					log.Printf("Error marshalling reply message: %v\n", err)
				} else {
					if broadcast == nil || broadcast.Client == nil || broadcast.Client.Socket == nil {
						log.Println("WebSocket connection or client is nil")
						return // 或者进行一些错误处理
					}

					// 现在我们确定 broadcast.Client.Socket 不是 nil，可以安全地调用方法
					err = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						log.Printf("Error sending reply message: %v\n", err)
						// 处理错误，可能是关闭连接或者重试发送
					}
				}

				err = InsertMsg(conf.MongoDBName, id, string(message), 0, int64(time.Hour*24*30*3))
				if err != nil {
					fmt.Println("flag false InsertOneMsg Err", err)
				}
			}
		}
	}
}
