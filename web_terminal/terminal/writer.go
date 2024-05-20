package terminal

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// 4K
	DefaultWriteBuf = 4 * 1024
)

func (i *WebSocketTerminal) ReadReq(req any) error {
	mt, data, err := i.ws.ReadMessage()
	if err != nil {
		return err
	}
	if mt != websocket.TextMessage {
		return fmt.Errorf("req must be TextMessage, but now not, is %d", mt)
	}
	if !json.Valid(data) {
		return fmt.Errorf("req must be json data, but %s", string(data))
	}

	return json.Unmarshal(data, req)
}

func (i *WebSocketTerminal) WriteTo(r io.Reader) (err error) {
	_, err = io.CopyBuffer(i, r, i.writeBuf)
	if err != nil {
		return err
	}
	defer i.ResetWriteBuf()

	_, err = i.Write(i.writeBuf)
	return
}

func (i *WebSocketTerminal) Write(p []byte) (n int, err error) {
	err = i.ws.WriteMessage(websocket.BinaryMessage, p)
	n = len(p)
	return
}

func (i *WebSocketTerminal) WriteTextln(format string, a ...any) {
	i.WriteTextf(format, a...)
	i.WriteText("\r\n")
}

func (i *WebSocketTerminal) WriteText(msg string) {
	err := i.ws.WriteMessage(websocket.BinaryMessage, []byte(msg))
	if err != nil {
		fmt.Printf("write message error, %s\n", err)
	}
}

func (i *WebSocketTerminal) WriteTextf(format string, a ...any) {
	i.WriteText(fmt.Sprintf(format, a...))
}

func (i *WebSocketTerminal) Failed(err error) {
	i.close(websocket.CloseGoingAway, err.Error())
}

func (i *WebSocketTerminal) Success(msg string) {
	i.close(websocket.CloseNormalClosure, msg)
}

func (i *WebSocketTerminal) ResetWriteBuf() {
	i.writeBuf = make([]byte, DefaultWriteBuf)
}

func (i *WebSocketTerminal) close(code int, msg string) {
	fmt.Printf("close code: %d, msg: %s\n", code, msg)
	err := i.ws.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(code, msg),
		time.Now().Add(i.timeout),
	)
	if err != nil {
		fmt.Printf("close error, %s\n", err)
		i.WriteText("\n" + msg)
	}
}
