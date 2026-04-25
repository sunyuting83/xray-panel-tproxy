package websocket

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	config "xpanel/Config"

	utils "xpanel/utils"

	"github.com/google/uuid"
	"github.com/lxzan/gws"
)

func TestTCPing(ID string) {
	current_path, _ := utils.GetCurrentPath()
	jsonFile := filepath.Join(current_path, "data", "data.json")
	data, _ := os.ReadFile(jsonFile)
	list := utils.ListToJsons(data)
	if len(*list) > 0 {
		for _, item := range *list {
			port := strconv.Itoa(item.Port)
			index := strconv.Itoa(item.Index)
			elapsedTime, err := TCPing(item.Address, port)
			speed := "0"
			if err == nil {
				speed = Float64ToStringWithPrecision(elapsedTime.Seconds()*1000, 2)
			}
			speeData := &config.Message{
				Type: "tcping",
				UUID: ID,
				Data: strings.Join([]string{index, speed}, "||||"),
			}
			sedData, _ := json.Marshal(speeData)
			Manager.Broadcast(sedData)
		}
	}
}

type WebSocketManager struct {
	// sessions: ID -> *gws.Conn (用于根据 ID 发消息)
	sessions sync.Map
	// infoMap: *gws.Conn -> *config.GopherInfo (用于在接收消息时识别是谁发的)
	infoMap sync.Map
}

// Manager 全局单例
var Manager = &WebSocketManager{}

// OnOpen 实现握手后的逻辑
func (m *WebSocketManager) OnOpen(c *gws.Conn) {
	id := uuid.NewString()
	info := &config.GopherInfo{ID: id}

	// 将连接与 ID 双向绑定在自己的 Map 中
	m.sessions.Store(id, c)
	m.infoMap.Store(c, info)

	// 1. 同步逻辑：将已在线的其他用户信息发送给当前新连接
	m.sessions.Range(func(key, value any) bool {
		conn := value.(*gws.Conn)
		// 通过连接反查信息
		if val, ok := m.infoMap.Load(conn); ok {
			otherInfo := val.(*config.GopherInfo)
			m.sendToConn(c, "message", otherInfo.ID, "first connect")
		}
		return true
	})

	// 2. 发送新连接自己的 ID 消息
	m.sendToConn(c, "message", id, "first connect")
}

// OnMessage 处理接收到的消息
func (m *WebSocketManager) OnMessage(c *gws.Conn, msg *gws.Message) {
	// 1. 必须操作：延迟关闭消息，将内存还给 gws 的复用池
	defer msg.Close()
	msgBytes := msg.Bytes()

	// 2. 反查身份：利用双向 Map 找到当前连接是谁
	val, ok := m.infoMap.Load(c)
	if !ok {
		return
	}
	info := val.(*config.GopherInfo)

	// 3. 解析协议：解析前端发来的 JSON
	var message config.Message
	if err := json.Unmarshal(msgBytes, &message); err != nil {
		return
	}

	// 4. 业务分发逻辑

	// 处理特殊的指令（即使 UUID 不匹配也可能需要处理的逻辑，如公共指令）
	if message.Type == "download" {
		// 传入当前 Manager 实例和用户 ID
		MakeWsData(info.ID, message.Data)
	}

	// 校验 UUID 匹配，确保安全性（防止伪造他人 ID 发指令）
	if message.UUID == info.ID {
		switch message.Type {
		case "message":
			// 广播原始字节流，性能最高
			m.Broadcast(msgBytes)

		case "testspeed":
			// 假设 datafactory 已适配全局 Manager
			SyncCheckData(message.UUID)

		case "tcping":
			// 假设 utils 已适配全局 Manager
			TestTCPing(message.UUID)

		case "active":
			// 心跳包，通常直接 break 即可。gws 内部会自动处理协议级的 Ping/Pong
			break

		default:
			// 默认行为：广播未知类型的消息
			m.Broadcast(msgBytes)
		}
	}
}

// SendMessageToWs 定向发送消息 (供 utils/download.go 等外部调用)
func (m *WebSocketManager) SendMessageToWs(id string, msgType string, data string) {
	if value, ok := m.sessions.Load(id); ok {
		conn := value.(*gws.Conn)
		m.sendToConn(conn, msgType, id, data)
	}
}

// Broadcast 广播给所有人
func (m *WebSocketManager) Broadcast(data []byte) {
	m.sessions.Range(func(key, value any) bool {
		conn := value.(*gws.Conn)
		_ = conn.WriteMessage(gws.OpcodeText, data)
		return true
	})
}

// sendToConn 内部辅助发送函数
func (m *WebSocketManager) sendToConn(c *gws.Conn, t, id, d string) {
	resp := config.Message{
		Type: t,
		UUID: id,
		Data: d,
	}
	payload, _ := json.Marshal(resp)
	_ = c.WriteMessage(gws.OpcodeText, payload)
}

// OnClose 连接断开时的清理工作
func (m *WebSocketManager) OnClose(c *gws.Conn, err error) {
	// 找到该连接对应的 Info，清理两个 Map
	if val, ok := m.infoMap.Load(c); ok {
		info := val.(*config.GopherInfo)
		m.sessions.Delete(info.ID)
		m.infoMap.Delete(c)
	}
}

// 必须实现的其他接口方法
func (m *WebSocketManager) OnPing(c *gws.Conn, p []byte) { _ = c.WritePong(p) }
func (m *WebSocketManager) OnPong(c *gws.Conn, p []byte) {}
