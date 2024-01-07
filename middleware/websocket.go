package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 定义 WebSocket 服务器
type WebSocketServer struct {
	clients   map[*websocket.Conn]bool
	broadcast chan Message
	upgrader  websocket.Upgrader
}

// 定义消息结构
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// NewWebSocketServer 创建一个新的 WebSocketServer 实例。
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		broadcast: make(chan Message),
		clients:   make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许跨域
			},
		},
	}
}

// HandleConnections 处理新的 WebSocket 连接。
func (server *WebSocketServer) HandleConnections(c *gin.Context) {
	// 升级初始 GET 请求到一个 WebSocket 连接
	ws, err := server.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error during connection upgradation:", err)
		return
	}
	defer ws.Close()

	// 注册新的客户端
	server.clients[ws] = true

	for {
		var msg Message
		// 读取新的消息
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error during message reading: %v", err)
			delete(server.clients, ws)
			break
		}
		// 发送消息到广播频道
		server.broadcast <- msg
	}
}

// HandleMessages 处理从 WebSocket 接收到的消息。
// 广播消息到所有客户端
func (server *WebSocketServer) handleMessages() {
	for {
		// 从广播频道中获取消息
		msg := <-server.broadcast
		// 发送消息到所有的客户端
		for client := range server.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("Error during message writing: %v", err)
				client.Close()
				delete(server.clients, client)
			}
		}
	}
}

// WebSocketMiddleware 创建 WebSocket 中间件。
func WebSocketMiddleware(wsServer *WebSocketServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.IsWebsocket() {
			wsServer.HandleConnections(c)
		}
	}
}
