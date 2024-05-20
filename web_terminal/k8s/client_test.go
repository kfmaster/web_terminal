package k8s_test

import (
	"context"
	"io"
	"os"
	"testing"

	"gitlab.com/go-course-project/public-project/web_terminal/k8s"
)

var (
	client *k8s.Client
)

func init() {
	kubeConf := k8s.MustReadContentFile("kube_config.yml")
	c, err := k8s.NewClient(kubeConf)
	if err != nil {
		panic(err)
	}
	client = c
}

func TestServerVersion(t *testing.T) {
	v, err := client.ServerVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestWatchConainterLog(t *testing.T) {
	req := k8s.NewWatchConainterLogRequest()
	req.Namespace = "default"
	req.PodName = "cicd-test-565794f58d-qv52s"
	stream, err := client.WatchConainterLog(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	_, err = io.Copy(os.Stdout, stream)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoginContainer(t *testing.T) {
	reader, writer := io.Pipe()

	term := &k8s.MockContainerTerminal{
		In: reader,
	}

	// Simulate user input
	go func() {
		writer.Write([]byte("ls -al / \n"))
	}()

	req := k8s.NewLoginContainerRequest(term)
	req.Namespace = "default"
	req.PodName = "cicd-test-565794f58d-qv52s"
	err := client.LoginContainer(context.Background(), req)
	if err != nil {
		panic(err)
	}
}
