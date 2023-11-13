package api

import (
	"analysis-engine/pkg/api/k8s"
	"analysis-engine/pkg/api/metric"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type ClientManager struct {
	KubeClient *kubernetes.Clientset
}

func NewClientManager() *ClientManager {
	var err error
	result := &ClientManager{}

	result.KubeClient, err = k8s.NewClient()
	if err != nil {
		klog.Errorln(err)
	}

	return result
}

func GetMultiMetric(ip string) (*metric.MultiMetric, error) {
	host := ip + ":9323"
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	grpcClient := metric.NewMetricCollectorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	res, err := grpcClient.GetMultiMetric(ctx, &metric.Request{})
	if err != nil {
		cancel()
		return nil, err
	}

	cancel()

	fmt.Println(res)

	return res, nil
}
