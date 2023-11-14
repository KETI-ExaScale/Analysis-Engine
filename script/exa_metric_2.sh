#!/bin/bash

echo "create pod exa_metric_nbody_1.yaml"
kubectl create -f deployment/exa_metric_nbody_1.yaml

echo "create pod exa_metric_nbody_2.yaml"
kubectl create -f deployment/exa_metric_nbody_2.yaml

# echo "create pod exa_metric_nbody_3.yaml"
# kubectl create -f deployment/exa_metric_nbody_3.yaml

# echo "create pod exa_metric_nbody_4.yaml"
# kubectl create -f deployment/exa_metric_nbody_4.yaml

# echo "create pod exa_metric_nbody_5.yaml"
# kubectl create -f deployment/exa_metric_nbody_5.yaml

# echo "create pod exa_metric_nbody_6.yaml"
# kubectl create -f deployment/exa_metric_nbody_6.yaml

# echo "create pod exa_metric_nbody_7.yaml"
# kubectl create -f deployment/exa_metric_nbody_7.yaml

# echo "create pod exa_metric_nbody_8.yaml"
# kubectl create -f deployment/exa_metric_nbody_8.yaml

# echo "create pod exa_metric_nbody_9.yaml"
# kubectl create -f deployment/exa_metric_nbody_9.yaml

# echo "create pod exa_metric_nbody_10.yaml"
# kubectl create -f deployment/exa_metric_nbody_10.yaml

sleep 10

echo "check metric"
curl http://10.0.5.20:32555/analysis/metric