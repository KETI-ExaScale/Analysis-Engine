package analysis

import (
	"analysis-engine/pkg/api"
	"analysis-engine/pkg/api/score"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

type Engine struct {
	score.UnimplementedMetricGRPCServer
	Client   *api.ClientManager
	Analysis AnalysisInterface
	MasterIP string
	// NodeScore     map[string]float32
	// DeploymentMap map[string]bool
}

func InitEngine() *Engine {
	client := api.NewClientManager()
	ip := "10.0.5.20"

	return &Engine{
		Client:   client,
		Analysis: GetAnalysisFramework(),
		MasterIP: ip,
		// NodeScore:     make(map[string]float32),
		// DeploymentMap: make(map[string]bool),
	}
}

func (e *Engine) Work() {
	e.StartGRPCServer()
}

func (e *Engine) RunMetricAnalyzer() (map[string]*Score, error) {
	multi_metric, err := api.GetMultiMetric(e.MasterIP)
	if err != nil {
		return nil, fmt.Errorf("metric collector gRPC error: %s", err)
	}

	scores := make(map[string]*Score)
	for node_name, node_metric := range multi_metric.NodeMetrics {
		score := &Score{
			NodeScore: 0.0,
			GPUScores: make(map[string]float32),
		}
		for gpu_name, _ := range node_metric.GpuMetrics {
			score.GPUScores[gpu_name] = 0.0
		}
		scores[node_name] = score
	}

	e.Analysis.RunNodeScoringPlugins(scores, multi_metric)
	e.Analysis.RunGPUScoringPlugins(scores, multi_metric)

	return scores, nil
}

func (e *Engine) GetScore(context.Context, *score.Request) (*score.Response, error) {
	scores, err := e.RunMetricAnalyzer()
	if err != nil {
		return nil, err
	}

	response := &score.Response{
		Message: make(map[string]*score.Score),
	}

	for node_name, node_score := range scores {
		ns := &score.Score{
			Nodescore: node_score.NodeScore,
			Gpuscore:  make(map[string]float32),
		}
		for gpu_name, gpu_score := range node_score.GPUScores {
			ns.Gpuscore[gpu_name] = gpu_score
		}
		response.Message[node_name] = ns
	}

	sc1 := &score.Score{
		Nodescore: 50,
		Gpuscore:  make(map[string]float32),
	}
	response.Message["c1-master"] = sc1

	sc2 := &score.Score{
		Nodescore: 52,
		Gpuscore:  make(map[string]float32),
	}
	response.Message["cpu-node1"] = sc2

	sc3 := &score.Score{
		Nodescore: 60,
		Gpuscore:  make(map[string]float32),
	}
	sc3.Gpuscore["GPU-de67e9b5-fd88-d5b3-133f-8c839735ab87"] = 37.0
	response.Message["gpu-node1"] = sc3

	sc4 := &score.Score{
		Nodescore: 30,
		Gpuscore:  make(map[string]float32),
	}
	sc4.Gpuscore["GPU-476b21ae-61f7-8fe1-fc92-a8e15c6ee0bc"] = 40.0
	response.Message["gpu-node2"] = sc4

	fmt.Println(response)
	fmt.Println("# Analysis score transmission completed")

	return response, nil
}

func (e *Engine) StartGRPCServer() {
	lis, err := net.Listen("tcp", ":9322")
	if err != nil {
		klog.Fatalf("failed to listen: %v", err)
	}
	scoreServer := grpc.NewServer()
	score.RegisterMetricGRPCServer(scoreServer, e)
	fmt.Println("-----:: Analysis Engine Server Running... ::-----")
	if err := scoreServer.Serve(lis); err != nil {
		klog.Fatalf("failed to serve: %v", err)
	}
}
