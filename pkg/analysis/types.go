package analysis

import (
	"analysis-engine/pkg/api/metric"
	"analysis-engine/pkg/api/score"
	"fmt"
	"os"
)

const (
	LEVEL1 = "LEVEL1"
	LEVEL2 = "LEVEL2"
	LEVEL3 = "LEVEL3"
)

var ANALYSIS_ENGINE_DEBUGG_LEVEL = os.Getenv("DEBUGG_LEVEL")

type MetricCache struct {
	MultiMetrics map[string]*metric.MultiMetric
}

func NewMetricCache() *MetricCache {
	return &MetricCache{
		MultiMetrics: make(map[string]*metric.MultiMetric),
	}
}

func (m *MetricCache) DumpMetricCache() {
	KETI_LOG_L1("\n---:: Dump Multi Metric All ::---")

	for nodeName, multiMetric := range m.MultiMetrics {
		KETI_LOG_L1(fmt.Sprintf("## Node[%s] Metric", nodeName))

		KETI_LOG_L1("1. [Multi Metric]")
		KETI_LOG_L1(fmt.Sprintf("1-1. node name : %s", multiMetric.NodeName))
		KETI_LOG_L1(fmt.Sprintf("1-2. gpu count : %d", multiMetric.GpuCount))
		KETI_LOG_L1("1-3. nvlink list ")
		for _, nvlink := range multiMetric.NvlinkInfo {
			KETI_LOG_L1(fmt.Sprintf("%s:%s lane:%d", nvlink.Gpu1Uuid, nvlink.Gpu2Uuid, nvlink.Lanecount))
		}

		KETI_LOG_L1("2. [Node Metric]")
		KETI_LOG_L1(fmt.Sprintf("2-1. milli cpu (free/total) : %d/%d", multiMetric.NodeMetric.MilliCpuFree, multiMetric.NodeMetric.MilliCpuTotal))
		KETI_LOG_L1(fmt.Sprintf("2-2. memory (free/total) : %d/%d", multiMetric.NodeMetric.MemoryFree, multiMetric.NodeMetric.MemoryTotal))
		KETI_LOG_L1(fmt.Sprintf("2-3. storage (free/total) : %d/%d", multiMetric.NodeMetric.StorageFree, multiMetric.NodeMetric.StprageTotal))
		KETI_LOG_L1(fmt.Sprintf("2-4. network (rx/tx) : %d/%d", multiMetric.NodeMetric.NetworkRx, multiMetric.NodeMetric.NetworkTx))

		KETI_LOG_L1("3. [GPU Metric]")
		for gpuName, gpuMetric := range multiMetric.GpuMetrics {
			KETI_LOG_L1(fmt.Sprintf("3-0 GPU UUID : %s", gpuName))
			KETI_LOG_L1(fmt.Sprintf("3-1. index : %d", gpuMetric.Index))
			KETI_LOG_L1(fmt.Sprintf("3-2. gpu name : %s", gpuMetric.GpuName))
			KETI_LOG_L1(fmt.Sprintf("3-3. architecture : %s", gpuMetric.Architecture))
			KETI_LOG_L1(fmt.Sprintf("3-4. max clock : %d", gpuMetric.MaxClock))
			KETI_LOG_L1(fmt.Sprintf("3-5. cudacore : %d", gpuMetric.Cudacore))
			KETI_LOG_L1(fmt.Sprintf("3-6. bandwidth : %f", gpuMetric.Bandwidth))
			KETI_LOG_L1(fmt.Sprintf("3-7. flops : %d", gpuMetric.Flops))
			KETI_LOG_L1(fmt.Sprintf("3-8. max operative temperature : %d", gpuMetric.MaxOperativeTemp))
			KETI_LOG_L1(fmt.Sprintf("3-9. slow down temperature : %d", gpuMetric.SlowdownTemp))
			KETI_LOG_L1(fmt.Sprintf("3-10. shut dowm temperature : %d", gpuMetric.ShutdownTemp))
			KETI_LOG_L1(fmt.Sprintf("3-11. memory (used/total) : %d/%d", gpuMetric.MemoryUsed, gpuMetric.MemoryTotal))
			KETI_LOG_L1(fmt.Sprintf("3-12. power (used) : %d", gpuMetric.PowerUsed))
			KETI_LOG_L1(fmt.Sprintf("3-13. pci (rx/tx) :  %d/%d", gpuMetric.PciRx, gpuMetric.PciTx))
			KETI_LOG_L1(fmt.Sprintf("3-14. temperature : %d", gpuMetric.Temperature))
			KETI_LOG_L1(fmt.Sprintf("3-15. utilization : %d", gpuMetric.Utilization))
			KETI_LOG_L1(fmt.Sprintf("3-16. fan speed : %d", gpuMetric.FanSpeed))
			KETI_LOG_L1(fmt.Sprintf("3-17. pod count : %d", gpuMetric.PodCount))

			KETI_LOG_L1("4. [GPU Pod Metric]")
			for podName, podMetric := range gpuMetric.PodMetrics {
				KETI_LOG_L1(fmt.Sprintf("# Pod Name : %s", podName))
				KETI_LOG_L1(fmt.Sprintf("4-1. node milli cpu (used) : %d", podMetric.NodeMilliCpuUsed))
				KETI_LOG_L1(fmt.Sprintf("4-2. node memory (used) : %d", podMetric.NodeMemoryUsed))
				KETI_LOG_L1(fmt.Sprintf("4-3. node storage (used) : %d", podMetric.NodeStorageUsed))
				KETI_LOG_L1(fmt.Sprintf("4-4. node network (rx/tx) :  %d/%d", podMetric.NodeNetworkRx, podMetric.NodeNetworkTx))
				KETI_LOG_L1(fmt.Sprintf("4-5. gpu memory :  %d", podMetric.GpuMemoryUsed))
			}
		}
	}
	KETI_LOG_L1("-----------------------------------\n")
}

func DumpScore(score *score.AnalysisScore) {
	KETI_LOG_L1("\n---:: Dump Analysis Score All ::---")

	for nodeName, nodeScore := range score.Scores {
		KETI_LOG_L1(fmt.Sprintf("1. Node [%s] Score: %f", nodeName, nodeScore.NodeScore))

		for gpuName, gpuScore := range nodeScore.GpuScores {
			KETI_LOG_L1(fmt.Sprintf("2. GPU [%s] Score: %f", gpuName, gpuScore.GpuScore))
		}
	}
	KETI_LOG_L1("-----------------------------------\n")
}

func KETI_LOG_L1(log string) { //자세한 출력, DumpClusterInfo DumpNodeInfo
	if ANALYSIS_ENGINE_DEBUGG_LEVEL == LEVEL1 {
		fmt.Println(log)
	}
}

func KETI_LOG_L2(log string) { // 기본출력
	if ANALYSIS_ENGINE_DEBUGG_LEVEL == LEVEL1 || ANALYSIS_ENGINE_DEBUGG_LEVEL == LEVEL2 {
		fmt.Println(log)
	}
}

func KETI_LOG_L3(log string) { //필수출력, 정량용, 에러
	if ANALYSIS_ENGINE_DEBUGG_LEVEL == LEVEL1 || ANALYSIS_ENGINE_DEBUGG_LEVEL == LEVEL2 || ANALYSIS_ENGINE_DEBUGG_LEVEL == LEVEL3 {
		fmt.Println(log)
	}
}
