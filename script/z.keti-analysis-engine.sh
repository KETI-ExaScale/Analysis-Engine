#!/usr/bin/env bash
dir=$( pwd )

#$1 create/c or delete/d

if [ "$1" == "delete" ] || [ "$1" == "d" ]; then   
    echo kubectl delete -f $dir/../deployments/keti-analysis-engine.yaml
    kubectl delete -f $dir/../deployments/keti-analysis-engine.yaml
else
    echo kubectl create -f $dir/../deployments/keti-analysis-engine.yaml
    kubectl create -f $dir/../deployments/keti-analysis-engine.yaml
fi