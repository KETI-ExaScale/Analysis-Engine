package analysis

import (
	"analysis-engine/pkg/api/score"
	"fmt"
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
			NodeStorage{},
		},
		GPUScoring: []ScoringPlugin{
			GPUFlops{},
			GPUPodCount{},
			GPUUtilization{},
			GPUTemperature{},
			GPUPower{},
			GPUBandwidth{},
			// GPUDirectStorage{},
			// GPUProcessTypeBalance{},
		},
	}
}

func (af AnalysisFramework) RunNodeScoringPlugins(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] run node scoring plugins")
	for _, afn := range af.NodeScoring {
		afn.Scoring(analysisScore, metricCache)
	}
}

func (af AnalysisFramework) RunGPUScoringPlugins(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] run gpu scoring plugins")
	for _, afg := range af.GPUScoring {
		afg.Scoring(analysisScore, metricCache)
	}
}

type NodeCPUCore struct{}
type NodeMemory struct{}
type NodeStorage struct{}

func (s NodeCPUCore) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#1-1. NodeGPUCore")
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		nodeScore := (1 - (float32(multiMetric.NodeMetric.MemoryUsage) / float32(multiMetric.NodeMetric.MemoryTotal))) * 100
		analysisScore.Scores[nodeName].NodeScore += nodeScore
		KETI_LOG_L2(fmt.Sprintf("[debugg] node {%s} score: %f", nodeName, analysisScore.Scores[nodeName].NodeScore))
	}
}
func (s NodeMemory) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#1-2. NodeMemory")
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		nodeScore := (1 - (float32(multiMetric.NodeMetric.MilliCpuUsage) / float32(multiMetric.NodeMetric.MilliCpuTotal))) * 100
		analysisScore.Scores[nodeName].NodeScore += nodeScore
		KETI_LOG_L2(fmt.Sprintf("[debugg] node {%s} score: %f", nodeName, analysisScore.Scores[nodeName].NodeScore))
	}
}
func (s NodeStorage) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#1-3. NodeStorage")
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		nodeScore := (1 - (float32(multiMetric.NodeMetric.StorageUsage) / float32(multiMetric.NodeMetric.StorageTotal))) * 100
		analysisScore.Scores[nodeName].NodeScore += nodeScore
		KETI_LOG_L2(fmt.Sprintf("[debugg] node {%s} score: %f", nodeName, analysisScore.Scores[nodeName].NodeScore))
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
	KETI_LOG_L2("[stage] S#2-1. GPUFlops")
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		for gpuName, gpu_metric := range multiMetric.GpuMetrics {
			analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += float32(gpu_metric.Flops) / 1000
			KETI_LOG_L2(fmt.Sprintf("[debugg] gpu {%s} score: %f", gpuName, analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore))
		}
	}
}
func (s GPUPodCount) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#2-2. GPUPodCount")
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		for gpuName, gpu_metric := range multiMetric.GpuMetrics {
			analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += float32(gpu_metric.GetPodCount())
			KETI_LOG_L2(fmt.Sprintf("[debugg] gpu {%s} score: %f", gpuName, analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore))
		}
	}
}
func (s GPUUtilization) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#2-3. GPUUtilization")

	for nodeName, multiMetric := range metricCache.MultiMetrics {
		for gpuName, gpu_metric := range multiMetric.GpuMetrics {
			if gpu_metric.GetPodCount() == 0 {
				analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += 11
			} else {
				var score = float32(gpu_metric.Utilization) / 100 * 10
				analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += score
			}
			KETI_LOG_L2(fmt.Sprintf("[debugg] gpu {%s} score: %f", gpuName, analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore))
		}
	}
}
func (s GPUMemory) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#2-4. GPUMemory")
	var total_gpu_memory = float32(0.0)
	for _, multiMetric := range metricCache.MultiMetrics {
		for _, gpu_metric := range multiMetric.GpuMetrics {
			total_gpu_memory += float32(gpu_metric.MemoryUsed)
		}
	}

	for nodeName, multiMetric := range metricCache.MultiMetrics {
		for gpuName, gpu_metric := range multiMetric.GpuMetrics {
			analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += total_gpu_memory / float32(gpu_metric.MemoryUsed)
			KETI_LOG_L2(fmt.Sprintf("[debugg] gpu {%s} score: %f", gpuName, analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore))
		}
	}
}
func (s GPUTemperature) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#2-5. GPUTemperature")
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		for gpuName, gpu_metric := range multiMetric.GpuMetrics {
			analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += float32(gpu_metric.Temperature)
			KETI_LOG_L2(fmt.Sprintf("[debugg] gpu {%s} score: %f", gpuName, analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore))
		}
	}
}
func (s GPUPower) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#2-6. GPUPower")
	var total_gpu_power = float32(0.0)
	for _, multiMetric := range metricCache.MultiMetrics {
		for _, gpu_metric := range multiMetric.GpuMetrics {
			total_gpu_power += float32(gpu_metric.PowerUsed)
		}
	}

	for nodeName, multiMetric := range metricCache.MultiMetrics {
		for gpuName, gpu_metric := range multiMetric.GpuMetrics {
			analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += total_gpu_power / float32(gpu_metric.PowerUsed)
			KETI_LOG_L2(fmt.Sprintf("[debugg] gpu {%s} score: %f", gpuName, analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore))
		}
	}
}
func (s GPUBandwidth) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {
	KETI_LOG_L2("[stage] S#2-7. GPUBandwidth")
	for nodeName, multiMetric := range metricCache.MultiMetrics {
		for gpuName, gpu_metric := range multiMetric.GpuMetrics {
			analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore += float32(gpu_metric.Bandwidth) * 100
			KETI_LOG_L2(fmt.Sprintf("[debugg] gpu {%s} score: %f", gpuName, analysisScore.Scores[nodeName].GpuScores[gpuName].GpuScore))
		}
	}
}
func (s GPUDirectStorage) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
func (s GPUProcessTypeBalance) Scoring(analysisScore *score.AnalysisScore, metricCache *MetricCache) {

}
