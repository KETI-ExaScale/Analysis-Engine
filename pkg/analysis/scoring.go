package analysis

import "analysis-engine/pkg/api/metric"

type AnalysisInterface interface {
	RunNodeScoringPlugins(scores map[string]*Score, multi_metric *metric.MultiMetric)
	RunGPUScoringPlugins(scores map[string]*Score, multi_metric *metric.MultiMetric)
}

type AnalysisFramework struct {
	NodeScoring []ScoringPlugin
	GPUScoring  []ScoringPlugin
}

type ScoringPlugin interface {
	Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric)
}

type Score struct {
	NodeScore float32
	GPUScores map[string]float32
}

func GetAnalysisFramework() AnalysisInterface {
	return &AnalysisFramework{
		NodeScoring: []ScoringPlugin{
			// NodeCPUCore{},
			// NodeMemory{},
			// NodeStorage{},
		},
		GPUScoring: []ScoringPlugin{
			// GPUFlops{},
			// GPUPodCount{},
			// GPUUtilization{},
			// GPUTemperature{},
			GPUPower{},
			// GPUBandwidth{},
			// GPUDirectStorage{},
			// GPUProcessTypeBalance{},
		},
	}
}

func (af AnalysisFramework) RunNodeScoringPlugins(scores map[string]*Score, multi_metric *metric.MultiMetric) {
	for _, afn := range af.NodeScoring {
		afn.Scoring(scores, multi_metric)
	}
}

func (af AnalysisFramework) RunGPUScoringPlugins(scores map[string]*Score, multi_metric *metric.MultiMetric) {
	for _, afg := range af.GPUScoring {
		afg.Scoring(scores, multi_metric)
	}
}

type NodeCPUCore struct{}
type NodeMemory struct{}
type NodeStorage struct{}

func (s NodeCPUCore) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s NodeMemory) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s NodeStorage) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}

type GPUFlops struct{}
type GPUPodCount struct{}
type GPUUtilization struct{}
type GPUMemory struct{}
type GPUTemperature struct{}
type GPUPower struct{}
type GPUBandwidth struct{}
type GPUDirectStorage struct{}
type GPUProcessTypeBalance struct{}

func (s GPUFlops) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s GPUPodCount) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s GPUUtilization) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s GPUMemory) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s GPUTemperature) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s GPUPower) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {
	var total_gpu_power = float32(0.0)
	for _, node_metric := range multi_metric.NodeMetrics {
		for _, gpu_metric := range node_metric.GpuMetrics {
			total_gpu_power += gpu_metric.GpuPower
		}
	}

	for node_name, node_metric := range multi_metric.NodeMetrics {
		for gpu_name, gpu_metric := range node_metric.GpuMetrics {
			scores[node_name].GPUScores[gpu_name] = total_gpu_power / gpu_metric.GpuPower
		}
	}
}
func (s GPUBandwidth) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s GPUDirectStorage) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
func (s GPUProcessTypeBalance) Scoring(scores map[string]*Score, multi_metric *metric.MultiMetric) {

}
