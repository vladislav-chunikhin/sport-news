#!/bin/bash

kubectl --namespace=feed-fetcher get pods
kubectl get service -n feed-fetcher
kubectl --namespace=feed-fetcher get nodes -o wide

# curl -i http://192.168.65.3:31001