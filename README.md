### K8s WebTerminal based on https://github.com/infraboard/mpaas
### All credits to the original author: yumaojun03

### This is a shrink down version based on above repo
### It mainly implement a web_terminal on local browser to login to a Kubernetes Pod container and show its logs

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

## To run this project

```
Get your kubernetes config file (usually in ~/.kub/config) and copy to web_terminal/terminal subfolder as kube_config.yaml (which is ignored in .gitignore)

cd web_terminal
go run main.go
```

## Then go to local brower: http://127.0.0.1:8080
## You should see 2 tabs: "ShowContainerLog" and "LoginContainer"

## The "ShowContainerlog" will ask you for the namespace, pod and container information then show container log

## The "LoginContainer" will ask you for the namespace, pod and container information then you can login to the container and do interactive session on the container


