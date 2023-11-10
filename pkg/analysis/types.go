package analysis

import (
	"fmt"
	"os"
	"sync"
)

const (
	LEVEL1 = "LEVEL1"
	LEVEL2 = "LEVEL2"
	LEVEL3 = "LEVEL3"
)

var ANALYSIS_ENGINE_DEBUGG_LEVEL = os.Getenv("DEBUGG_LEVEL")

type MetricCollector struct {
	MultiMetric map[string]*MultiMetric
}

type MultiMetric struct {
	mutex      sync.Mutex
	nodeName   string
	gpuCount   int32
	nvlinkList []*NVLink
	gpuUUID    []string
	nodeMetric *NodeMetric
	gpuMetric  map[string]*GPUMetric
}

type NodeMetric struct {
	milliCPUTotal int64
	milliCPUFree  int64
	memoryTotal   int64
	memoryFree    int64
	storageTotal  int64
	storageFree   int64
	networkRX     int64
	networkTX     int64
}

type GPUMetric struct {
	index            int32
	gpuName          string
	architecture     int32
	maxClock         int64
	cudacore         int64
	bandwidth        float32
	flops            int64
	maxOperativeTemp int64
	slowdownTemp     int64
	shutdownTemp     int64
	memoryTotal      int64
	memoryUsed       int64
	powerUsed        int64
	pciRX            int64
	pciTX            int64
	memoryGauge      int64
	memoryCounter    int64
	temperature      int64
	utilization      int64
	fanSpeed         int64
	podCount         int64
	runningPods      map[string]PodMetric
}

type PodMetric struct {
	milliCPUUsed int64
	memoryUsed   int64
	storageUsed  int64
	networkRX    int64
	networkTX    int64
}

type NVLink struct {
	GPU1UUID  string
	GPU2UUID  string
	LaneCount int32
}

type NVLinkStatus struct {
	UUID          string
	BusID         string
	Lanes         map[string]int
	P2PUUID       []string
	P2PDeviceType []int //0 GPU, 1 IBMNPU, 2 SWITCH, 255 = UNKNOWN
	P2PBusID      []string
}

func NewMetricCollector() *MetricCollector {

}

func NewMultiMetric() *MultiMetric {
	return &MultiMetric{
		nodeName:   "",
		gpuCount:   0,
		nvlinkList: make([]*NVLink, 0),
		gpuUUID:    make([]string, 0),
		nodeMetric: NewNodeMetric(),
		gpuMetric:  make(map[string]*GPUMetric),
	}
}

func NewNodeMetric() *NodeMetric {
	return &NodeMetric{
		milliCPUTotal: 0,
		milliCPUFree:  0,
		memoryTotal:   0,
		memoryFree:    0,
		storageTotal:  0,
		storageFree:   0,
		networkRX:     0,
		networkTX:     0,
	}
}

func NewGPUMetric(i int32) *GPUMetric {
	return &GPUMetric{
		index:            i,
		gpuName:          "",
		architecture:     0,
		maxClock:         0,
		cudacore:         0,
		bandwidth:        0,
		flops:            0,
		maxOperativeTemp: 0,
		slowdownTemp:     0,
		shutdownTemp:     0,
		memoryTotal:      0,
		memoryUsed:       0,
		powerUsed:        0,
		pciRX:            0,
		pciTX:            0,
		memoryGauge:      0,
		memoryCounter:    0,
		temperature:      0,
		utilization:      0,
		fanSpeed:         0,
		podCount:         0,
		runningPods:      make(map[string]PodMetric),
	}
}

func NewPodMetric() *PodMetric {
	return &PodMetric{
		milliCPUUsed: 0,
		memoryUsed:   0,
		storageUsed:  0,
		networkRX:    0,
		networkTX:    0,
	}
}

func NewNVLink(s1 string, s2 string, l int32) *NVLink {
	return &NVLink{
		GPU1UUID:  s1,
		GPU2UUID:  s2,
		LaneCount: l,
	}
}

func (m *MultiMetric) DumpMultiMetric() {
	KETI_LOG_L1("\n---:: Dump Metric Collector ::---")

	KETI_LOG_L1("1. [Multi Metric]")
	KETI_LOG_L1(fmt.Sprintf("1-1. node name : %s", m.nodeName))
	KETI_LOG_L1(fmt.Sprintf("1-2. gpu count : %d", m.gpuCount))
	KETI_LOG_L1("1-3. nvlink list ")
	for _, nvlink := range m.nvlinkList {
		KETI_LOG_L1(fmt.Sprintf("%s:%s lane:%d", nvlink.GPU1UUID, nvlink.GPU2UUID, nvlink.LaneCount))
	}
	KETI_LOG_L1("1-4. gpu uuid ")
	for _, uuid := range m.gpuUUID {
		KETI_LOG_L1(fmt.Sprintf("- %s", uuid))
	}

	KETI_LOG_L1("2. [Node Metric]")
	KETI_LOG_L1(fmt.Sprintf("2-1. milli cpu (free/total) : %d/%d", m.nodeMetric.milliCPUFree, m.nodeMetric.milliCPUTotal))
	KETI_LOG_L1(fmt.Sprintf("2-2. memory (free/total) : %d/%d", m.nodeMetric.memoryFree, m.nodeMetric.memoryTotal))
	KETI_LOG_L1(fmt.Sprintf("2-3. storage (free/total) : %d/%d", m.nodeMetric.storageFree, m.nodeMetric.storageTotal))
	KETI_LOG_L1(fmt.Sprintf("2-4. network (rx/tx) : %d/%d", m.nodeMetric.networkRX, m.nodeMetric.networkTX))

	KETI_LOG_L1("3. [GPU Metric]")
	for gpuName, gpuMetric := range m.gpuMetric {
		KETI_LOG_L1(fmt.Sprintf("3-0 GPU UUID : %s", gpuName))
		KETI_LOG_L1(fmt.Sprintf("3-1. index : %d", gpuMetric.index))
		KETI_LOG_L1(fmt.Sprintf("3-2. gpuName : %s", gpuMetric.gpuName))
		KETI_LOG_L1(fmt.Sprintf("3-3. flops : %d", gpuMetric.flops))
		KETI_LOG_L1(fmt.Sprintf("3-4. architecture : %d", gpuMetric.architecture))
		KETI_LOG_L1(fmt.Sprintf("3-5. max operative temperature : %d", gpuMetric.maxOperativeTemp))
		KETI_LOG_L1(fmt.Sprintf("3-6. slow down temperature : %d", gpuMetric.slowdownTemp))
		KETI_LOG_L1(fmt.Sprintf("3-7. shut dowm temperature : %d", gpuMetric.shutdownTemp))
		KETI_LOG_L1(fmt.Sprintf("3-8. memory (used/total) : %d/%d", gpuMetric.memoryUsed, gpuMetric.memoryTotal))
		KETI_LOG_L1(fmt.Sprintf("3-9. power (used) : %d", gpuMetric.powerUsed))
		KETI_LOG_L1(fmt.Sprintf("3-10. network (rx/tx) :  %d/%d", gpuMetric.pciRX, gpuMetric.pciTX))
		KETI_LOG_L1(fmt.Sprintf("3-11. memory gauge : %d", gpuMetric.memoryGauge))
		KETI_LOG_L1(fmt.Sprintf("3-12. memory counter : %d", gpuMetric.memoryCounter))
		KETI_LOG_L1(fmt.Sprintf("3-13. temperature : %d", gpuMetric.temperature))
		KETI_LOG_L1(fmt.Sprintf("3-14. utilization : %d", gpuMetric.utilization))
		KETI_LOG_L1(fmt.Sprintf("3-15. bandwidth : %f", gpuMetric.bandwidth))
		KETI_LOG_L1(fmt.Sprintf("3-16. fan speed : %d", gpuMetric.fanSpeed))
		KETI_LOG_L1(fmt.Sprintf("3-17. pod count : %d", gpuMetric.podCount))

		KETI_LOG_L1("4. [Pod Metric]")
		for podName, podMetric := range gpuMetric.runningPods {
			KETI_LOG_L1(fmt.Sprintf("# Pod Name : %s", podName))
			KETI_LOG_L1(fmt.Sprintf("4-1. milli cpu (used) : %d", podMetric.milliCPUUsed))
			KETI_LOG_L1(fmt.Sprintf("4-2. memory (used) : %d", podMetric.memoryUsed))
			KETI_LOG_L1(fmt.Sprintf("4-3. storage (used) : %d", podMetric.storageUsed))
			KETI_LOG_L1(fmt.Sprintf("4-4. network (rx/tx) :  %d/%d", podMetric.networkRX, podMetric.networkTX))
		}
	}
}

func (m *MultiMetric) DumpMultiMetricForTest() {
	KETI_LOG_L3("\n---:: Dump Metric Collector ::---")

	KETI_LOG_L3("1. [Multi Metric]")
	KETI_LOG_L3(fmt.Sprintf("1-1. node name : %s", m.nodeName))
	KETI_LOG_L3(fmt.Sprintf("1-2. gpu count : %d", m.gpuCount))
	KETI_LOG_L3("1-3. nvlink list ")
	for _, nvlink := range m.nvlinkList {
		KETI_LOG_L3(fmt.Sprintf("%s:%s lane:%d", nvlink.GPU1UUID, nvlink.GPU2UUID, nvlink.LaneCount))
	}
	KETI_LOG_L3("1-4. gpu uuid ")
	for _, uuid := range m.gpuUUID {
		KETI_LOG_L3(fmt.Sprintf("- %s", uuid))
	}

	KETI_LOG_L3("2. [Node Metric]")
	KETI_LOG_L3(fmt.Sprintf("2-1. milli cpu (free/total) : %d/%d", m.nodeMetric.milliCPUFree, m.nodeMetric.milliCPUTotal))
	KETI_LOG_L3(fmt.Sprintf("2-2. memory (free/total) : %d/%d", m.nodeMetric.memoryFree, m.nodeMetric.memoryTotal))
	KETI_LOG_L3(fmt.Sprintf("2-3. storage (free/total) : %d/%d", m.nodeMetric.storageFree, m.nodeMetric.storageTotal))
	KETI_LOG_L3(fmt.Sprintf("2-4. network (rx/tx) : %d/%d", m.nodeMetric.networkRX, m.nodeMetric.networkTX))

	KETI_LOG_L3("3. [GPU Metric]")
	for gpuName, gpuMetric := range m.gpuMetric {
		KETI_LOG_L3(fmt.Sprintf("# GPU UUID : %s", gpuName))
		KETI_LOG_L3(fmt.Sprintf("3-1. index : %d", gpuMetric.index))
		KETI_LOG_L3(fmt.Sprintf("3-2. gpuName : %s", gpuMetric.gpuName))
		KETI_LOG_L3(fmt.Sprintf("3-3. flops : %d", gpuMetric.flops))
		KETI_LOG_L3(fmt.Sprintf("3-4. architecture : %d", gpuMetric.architecture))
		KETI_LOG_L3(fmt.Sprintf("3-5. max operative temperature : %d", gpuMetric.maxOperativeTemp))
		KETI_LOG_L3(fmt.Sprintf("3-6. slow down temperature : %d", gpuMetric.slowdownTemp))
		KETI_LOG_L3(fmt.Sprintf("3-7. shut dowm temperature : %d", gpuMetric.shutdownTemp))
		KETI_LOG_L3(fmt.Sprintf("3-8. memory (used/total) : %d/%d", gpuMetric.memoryUsed, gpuMetric.memoryTotal))
		KETI_LOG_L3(fmt.Sprintf("3-9. power (used) : %d", gpuMetric.powerUsed))
		KETI_LOG_L3(fmt.Sprintf("3-10. network (rx/tx) :  %d/%d", gpuMetric.pciRX, gpuMetric.pciTX))
		KETI_LOG_L3(fmt.Sprintf("3-11. memory gauge : %d", gpuMetric.memoryGauge))
		KETI_LOG_L3(fmt.Sprintf("3-12. memory counter : %d", gpuMetric.memoryCounter))
		KETI_LOG_L3(fmt.Sprintf("3-13. temperature : %d", gpuMetric.temperature))
		KETI_LOG_L3(fmt.Sprintf("3-14. utilization : %d", gpuMetric.utilization))
		KETI_LOG_L3(fmt.Sprintf("3-15. bandwidth : %f", gpuMetric.bandwidth))
		KETI_LOG_L3(fmt.Sprintf("3-16. fan speed : %d", gpuMetric.fanSpeed))
		KETI_LOG_L3(fmt.Sprintf("3-17. pod count : %d", gpuMetric.podCount))

		KETI_LOG_L3("4. [Pod Metric]")
		for podName, podMetric := range gpuMetric.runningPods {
			KETI_LOG_L3(fmt.Sprintf("# Pod Name : %s", podName))
			KETI_LOG_L3(fmt.Sprintf("4-1. milli cpu (used) : %d", podMetric.milliCPUUsed))
			KETI_LOG_L3(fmt.Sprintf("4-2. memory (used) : %d", podMetric.memoryUsed))
			KETI_LOG_L3(fmt.Sprintf("4-3. storage (used) : %d", podMetric.storageUsed))
			KETI_LOG_L3(fmt.Sprintf("4-4. network (rx/tx) :  %d/%d", podMetric.networkRX, podMetric.networkTX))
		}
	}
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
