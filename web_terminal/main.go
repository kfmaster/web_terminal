package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/infraboard/mcube/v2/http/response"
	"gitlab.com/go-course-project/public-project/web_terminal/k8s"
	"gitlab.com/go-course-project/public-project/web_terminal/terminal"
)

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 60 * time.Second,
		ReadBufferSize:   8192,
		WriteBufferSize:  8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//go:embed ui
var ui embed.FS

func main() {
	http.HandleFunc("/ws/pod/terminal/log", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			response.Failed(w, err)
			return
		}

		// Construct Terminal object based on Websocket
		term := terminal.NewWebSocketTerminal(ws)

		// Get k8sClient
		kubeConf := k8s.MustReadContentFile("k8s/kube_config.yml")
		k8sClient, err := k8s.NewClient(kubeConf)
		if err != nil {
			term.Failed(err)
			return
		}

		// After establish Websocket, read user request from Terminal object
		// Read specific pod log based on requests sent from socket
		// TO-DO, authentication information can be provided
		req := k8s.NewWatchConainterLogRequest()
		if err = term.ReadReq(req); err != nil {
			term.Failed(err)
			return
		}

		// Get pod log
		podReader, err := k8sClient.WatchConainterLog(r.Context(), req)
		if err != nil {
			term.Failed(err)
			return
		}

		//Copy read stream to term, data writen by term will be send to user
		_, err = io.Copy(term, podReader)
		if err != nil {
			term.Failed(err)
			return
		}
	})

	http.HandleFunc("/ws/pod/terminal/login", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			response.Failed(w, err)
			return
		}

		// Construct Terminal object based on Websocket
		term := terminal.NewWebSocketTerminal(ws)

		// Get k8sClient
		kubeConf := k8s.MustReadContentFile("k8s/kube_config.yml")
		k8sClient, err := k8s.NewClient(kubeConf)
		if err != nil {
			term.Failed(err)
			return
		}

		req := k8s.NewLoginContainerRequest(term)
		if err = term.ReadReq(req); err != nil {
			term.Failed(err)
			return
		}

		// Login to pod
		err = k8sClient.LoginContainer(r.Context(), req)
		if err != nil {
			term.Failed(err)
			return
		}
	})

	web, _ := fs.Sub(ui, "ui")
	http.Handle("/", http.FileServer(http.FS(web)))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
