#!/bin/bash

# Delete Feed Fetcher Service
kubectl delete -f feed-fetcher-service.yaml

# Delete Feed Fetcher Deployment
kubectl delete -f feed-fetcher-deployment.yaml

# Delete RabbitMQ Service
kubectl delete -f rabbitmq-service.yaml

# Delete RabbitMQ Deployment
kubectl delete -f rabbitmq-deployment.yaml

# Delete RabbitMQ ConfigMap
kubectl delete -f rabbitmq-configmap.yaml

# Delete Feed Fetcher Namespace
kubectl delete -f feed-fetcher-namespace.yaml
