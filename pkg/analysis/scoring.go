package analysis

import (
	"analysis-engine/pkg/api/score"
)

type AnalysisInterface interface {
	RunNodeScoringPlugins(analysisScore *score.AnalysisScore, metricCache *MetricCache)
	RunGPUScoringPlugins(analysisScore *score.AnalysisScore, metricCache *MetricCache)
}

type AnalysisFramework struct {
	NodeScoring []ScoringPlugin
	GPUScoring  []ScoringPlugin
}

type ScoringPlugin interface {
	Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache)
}

func GetAnalysisFramework() AnalysisInterface {
	return &AnalysisFramework{
		NodeScoring: []ScoringPlugin{
			NodeCPUCore{},
			NodeMemory{},
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

func (af AnalysisFramework) RunNodeScoringPlugins(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	for _, afn := range af.NodeScoring {
		afn.Scoring(analysisScore, metricCache)
	}
}

func (af AnalysisFramework) RunGPUScoringPlugins(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	for _, afg := range af.GPUScoring {
		afg.Scoring(analysisScore, metricCache)
	}
}

type NodeCPUCore struct{}
type NodeMemory struct{}
type NodeStorage struct{}

func (s NodeCPUCore) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		nodeScore := (1 - (float32(multiMetric.NodeMetric.MemoryUsage) / float32(multiMetric.NodeMetric.MemoryTotal))) * 100
		analysisScore.Scores[nodeName].NodeScore += nodeScore
	}
}
func (s NodeMemory) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		nodeScore := (1 - (float32(multiMetric.NodeMetric.MilliCpuUsage) / float32(multiMetric.NodeMetric.MilliCpuTotal))) * 100
		analysisScore.Scores[nodeName].NodeScore += nodeScore
	}
}
func (s NodeStorage) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		nodeScore := (1 - (float32(multiMetric.NodeMetric.StorageUsage) / float32(multiMetric.NodeMetric.StorageTotal))) * 100
		analysisScore.Scores[nodeName].NodeScore += nodeScore
	}
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

func (s GPUFlops) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
func (s GPUPodCount) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
func (s GPUUtilization) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
func (s GPUMemory) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
func (s GPUTemperature) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
func (s GPUPower) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	var total_gpu_power = float32(0.0)
	for _, multiMetric := range metricCache.MultiMetrics {
		for _, gpu_metric := range multiMetric.GpuMetrics {
			total_gpu_power += float32(gpu_metric.PowerUsed)
		}
	}

	for nodeName, multiMetric := range metricCache.MultiMetrics {
		for gpuName, gpu_metric := range multiMetric.GpuMetrics {
			analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += total_gpu_power / float32(gpu_metric.PowerUsed)
		}
	}
}
func (s GPUBandwidth) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
func (s GPUDirectStorage) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
func (s GPUProcessTypeBalance) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
