#!/bin/bash

kubectl krew install kubectl-replay
cd ./examples

make default
kubectl apply -f examples/new-configmap.yaml

# Example 1. Replay audit log to create a ConfigMap in a non-existent namespace
kubectl-replay -f missing-configmap.yaml --audit-log-file=/var/log/kubernetes/audit.log --replay-speed=10x &

# Run 10 times...
for i in {1..10}
do
   kubectl get configmap new-configmap --namespace=non-existent-namespace || echo "ConfigMap not found"
   sleep 1
done

# Example 2. Replay event logs from event log of kubernetes cluster
kubectl-replay --from-event-log --replay-speed=10x &

# Run 10 times...
for i in {1..10}
do
   kubectl get configmap new-configmap --namespace=default || echo "ConfigMap not found"
   sleep 1
done

# Deploy nginx
kubectl run nginx --image=nginx

# Cleanup
make cleanup
