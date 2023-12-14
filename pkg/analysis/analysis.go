package analysis

import (
	"analysis-engine/pkg/api"
	"analysis-engine/pkg/api/score"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

type Engine struct {
	score.UnimplementedMetricGRPCServer
	Client                   *api.ClientManager
	Analysis                 AnalysisInterface
	GPUMetricCollectorIPList []string
}

func InitEngine() *Engine {
	client := api.NewClientManager()

	gpuMetricCollectorIPList := make([]string, 0)
	pods, _ := client.KubeClient.CoreV1().Pods("gpu").List(context.TODO(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		if strings.HasPrefix(pod.Name, "keti-gpu-metric-collector") && pod.Status.Phase == "Running" {
			internalIP := pod.Status.PodIP
			gpuMetricCollectorIPList = append(gpuMetricCollectorIPList, internalIP)
		}
	}

	return &Engine{
		Client:                   client,
		Analysis:                 GetAnalysisFramework(),
		GPUMetricCollectorIPList: gpuMetricCollectorIPList,
	}
}

func (e *Engine) Work(ctx context.Context, wg *sync.WaitGroup) {
	go e.StartGRPCServer(ctx, wg)

	if ANALYSIS_ENGINE_DEBUGG_LEVEL == LEVEL3 {
		go e.StartHTTPServer(ctx, wg)
	}
}

func (e *Engine) RunMetricAnalyzer() (*score.AnalysisScore, error) {
	KETI_LOG_L2("[analysis] run metric analyzer")

	metricCache, err := e.GetMetricCache()
	if err != nil {
		// return nil, fmt.Errorf("get multi metric error: %s", err)
	}

	KETI_LOG_L2("[analysis] get multi metric sucess")

	metricCache.DumpMetricCache() //테스트용

	analysisScores := &score.AnalysisScore{
		Scores: make(map[string]*score.NodeScore),
	}

	for nodeName, multiMetric := range metricCache.MultiMetrics {
		nodeScore := &score.NodeScore{
			NodeScore: 0.0,
			GpuScores: make(map[string]*score.GPUScore),
		}
		for gpuUUID := range multiMetric.GpuMetrics {
			nodeScore.GpuScores[gpuUUID] = &score.GPUScore{
				GpuScore: 0.0,
			}
		}
		analysisScores.Scores[nodeName] = nodeScore
	}

	e.Analysis.RunNodeScoringPlugins(analysisScores, metricCache)
	e.Analysis.RunGPUScoringPlugins(analysisScores, metricCache)

	KETI_LOG_L2("[analysis] finish scoring")

	DumpScore(analysisScores) //DEBUGG LEVEL = 1 일때 출력

	return analysisScores, nil
}

func (e *Engine) GetMetricCache() (*MetricCache, error) {
	metricCache := NewMetricCache()

	for _, ip := range e.GPUMetricCollectorIPList {
		multiMetric_, err := api.GetMultiMetric(ip)
		if err != nil {
			return nil, fmt.Errorf("metric collector gRPC error: %s", err)
		}

		nodeName := multiMetric_.NodeName
		metricCache.MultiMetrics[nodeName] = multiMetric_
	}

	return metricCache, nil
}

func (e *Engine) GetScore(context.Context, *score.Request) (*score.AnalysisScore, error) {
	KETI_LOG_L3("[gRPC] called get score")

	analysisScores, err := e.RunMetricAnalyzer()
	if err != nil {
		return nil, err
	}

	return analysisScores, nil
}

func (e *Engine) StartGRPCServer(ctx context.Context, wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", ":9322")
	if err != nil {
		klog.Fatalf("[error] start grpc server error: %v", err)
	}
	scoreServer := grpc.NewServer()
	score.RegisterMetricGRPCServer(scoreServer, e)
	KETI_LOG_L3("[gRPC] analysis engine server running...")

	go func() {
		if err := scoreServer.Serve(lis); err != nil {
			klog.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ctx.Done()
	wg.Done()
}

func (e *Engine) StartHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		KETI_LOG_L3("[http] analysis engine server running...")
		if err := http.ListenAndServe(":9595", e); err != nil {
			klog.Fatalf("[error] start http server error: %v", err)
		}
	}()

	<-ctx.Done()
	wg.Done()
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/analysis/metric":
		metricCache, err := e.GetMetricCache()
		if err != nil {
			klog.Fatal("[error] get metric cache error from servehttp")
		}

		metricCache.DumpMultiMetricForTest()

		w.WriteHeader(http.StatusOK)
	default:
		http.NotFound(w, r)
	}
}
