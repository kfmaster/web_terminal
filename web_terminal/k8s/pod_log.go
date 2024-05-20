package k8s

import (
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
)

func NewWatchConainterLogRequest() *WatchConainterLogRequest {
	return &WatchConainterLogRequest{
		PodLogOptions: &v1.PodLogOptions{
			Follow:                       true,
			Previous:                     false,
			InsecureSkipTLSVerifyBackend: true,
		},
	}
}

type WatchConainterLogRequest struct {
	Namespace string `json:"namespace" validate:"required"`
	PodName   string `json:"pod_name" validate:"required"`
	*v1.PodLogOptions
}

// 查看容器日志
func (c *Client) WatchConainterLog(
	ctx context.Context,
	req *WatchConainterLogRequest) (
	io.ReadCloser, error) {
	restReq := c.client.CoreV1().
		Pods(req.Namespace).
		GetLogs(req.PodName, req.PodLogOptions)
	return restReq.Stream(ctx)
}
