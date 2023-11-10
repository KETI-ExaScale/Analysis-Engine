package main

import "analysis-engine/pkg/analysis"

func main() {
	w := analysis.InitEngine()
	w.Work()
}
