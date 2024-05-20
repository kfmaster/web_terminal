package terminal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func NewWebSocketTerminal(conn *websocket.Conn) *WebSocketTerminal {
	return &WebSocketTerminal{
		ws:              conn,
		TerminalResizer: NewTerminalSize(),
		timeout:         3 * time.Second,
		writeBuf:        make([]byte, DefaultWriteBuf),
	}
}

type WebSocketTerminal struct {
	ws *websocket.Conn
	*TerminalResizer

	// Write需要属性
	timeout  time.Duration
	writeBuf []byte
}

func (t *WebSocketTerminal) Close() error {
	return t.ws.Close()
}

// 通过 websocket读取用户数据, 可能是指令也有可能是用户输入的数据
// 然后 Terminal 实现Reader接口, 通过Read 方法为 k8s client excutor 提供输入 输入到 容器
// ws.Read --> k8s Pod
func (t *WebSocketTerminal) Read(p []byte) (n int, err error) {
	mt, m, err := t.ws.ReadMessage()
	if err != nil {
		return 0, err
	}

	// 注意文本消息和关闭消息专门被设计为了指令通道
	switch mt {
	case websocket.TextMessage:
		t.HandleCmd(m)
	case websocket.CloseMessage:
		fmt.Printf("receive client close: %s\n", m)
	default:
		n = copy(p, m)
	}

	return n, nil
}

func (t *WebSocketTerminal) HandleCmd(m []byte) {
	resp := NewResponse()
	defer t.Response(resp)

	req, err := ParseRequest(m)
	if err != nil {
		resp.Message = err.Error()
		return
	}
	resp.Request = req

	// 单独处理Resize请求
	switch req.Command {
	case "resize":
		payload := NewTerminalSzie()
		err := json.Unmarshal(req.Params, payload)
		if err != nil {
			resp.Message = err.Error()
			return
		}
		t.SetSize(*payload)
		fmt.Printf("resize add to queue success: %s\n", req)
		return
	}

	// 处理自定义指令
	fn := GetCmdHandleFunc(req.Command)
	if fn == nil {
		resp.Message = "command not found"
		return
	}

	fn(req, resp)
}

// 命令的返回
func (i *WebSocketTerminal) Response(resp *Response) {
	if resp.Message != "" {
		fmt.Printf("response error, %s\n", resp.Message)
	}

	if err := i.ws.WriteJSON(resp); err != nil {
		fmt.Printf("write message error, %s\n", err)
	}
}
