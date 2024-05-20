package k8s

import (
	"context"
	"io"
	"os"

	"github.com/infraboard/mcube/v2/tools/pretty"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

var (
	shellCmd = []string{
		"sh",
		"-c",
		`TERM=xterm-256color; export TERM; [ -x /bin/bash ] && ([ -x /usr/bin/script ] && /usr/bin/script -q -c "/bin/bash" /dev/null || exec /bin/bash) || exec /bin/sh`,
	}
)

func NewLoginContainerRequest(ce ContainerTerminal) *LoginContainerRequest {
	return &LoginContainerRequest{
		Command:  shellCmd,
		Executor: ce,
	}
}

type LoginContainerRequest struct {
	Namespace     string            `json:"namespace" validate:"required"`
	PodName       string            `json:"pod_name" validate:"required"`
	ContainerName string            `json:"container_name"`
	Command       []string          `json:"command"`
	Executor      ContainerTerminal `json:"-"`
}

func (req *LoginContainerRequest) String() string {
	return pretty.ToJSON(req)
}

type ContainerTerminal interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// Login to pod, simulate docker exec to get command executed in container
func (c *Client) LoginContainer(ctx context.Context, req *LoginContainerRequest) error {
	// Construct container login request
	restReq := c.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(req.PodName).
		Namespace(req.Namespace).
		SubResource("exec")

	// Construct login container parameters
	restReq.VersionedParams(&v1.PodExecOptions{
		Container: req.ContainerName,
		Command:   req.Command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	// Login to container
	executor, err := remotecommand.NewSPDYExecutor(c.restconf, "POST", restReq.URL())
	if err != nil {
		return err
	}

	return executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             req.Executor,
		Stdout:            req.Executor,
		Stderr:            req.Executor,
		Tty:               true,
		TerminalSizeQueue: req.Executor,
	})
}

type MockContainerTerminal struct {
	In io.Reader
}

func (t *MockContainerTerminal) Read(p []byte) (n int, err error) {
	return t.In.Read(p)
}

func (t *MockContainerTerminal) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (t *MockContainerTerminal) Next() *remotecommand.TerminalSize {
	return &remotecommand.TerminalSize{
		Width:  100,
		Height: 100,
	}
}
