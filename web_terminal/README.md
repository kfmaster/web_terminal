# K8s WebTerminal based on https://github.com/infraboard/mpaas
# All credits to the original author: yumaojun03

## Initialize Project

```sh
# Initialize project
$ mdkir web_terminal
$ cd web_terminal
$ go mod init gitlab.com/go-course-project/public-project/web_terminal
# Open project using vscode
$ code .
```

## Folder structures

+ main.go： entry point program, WebSocket Server is implemented here
+ ui: web frontend 
+ terminal: WebSocket Terminal implementation
+ k8s: k8s related functions

## k8s package

Package k8s related functions and test

## Data vs Command

+ BinaryMessage: data

+ TextMessage: command


## UI

```sh
cd ui
npm install @xterm/xterm
code .
```

## Package

