package service

import (
	"encoding/json"
	"fmt"
	"giligili/conf"
	"giligili/pkg/e"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

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

			msg, _ := json.Marshal(replyMsg)

			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
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

			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}

				select {
				case conn.Send <- message:
					flag = true
				default:
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
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
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
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err = InsertMsg(conf.MongoDBName, id, string(message), 0, int64(time.Hour*24*30*3))
				if err != nil {
					fmt.Println("flag false InsertOneMsg Err", err)
				}
			}
		}
	}
}
