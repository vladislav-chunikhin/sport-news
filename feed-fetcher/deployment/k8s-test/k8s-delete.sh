#!/bin/bash

# Delete Feed Fetcher Service
kubectl delete -f feed-fetcher-service.yaml

# Delete Feed Fetcher Deployment
kubectl delete -f feed-fetcher-deployment.yaml

# Delete RabbitMQ Service
kubectl delete -f rabbitmq-service.yaml

# Delete RabbitMQ Deployment
kubectl delete -f rabbitmq-deployment.yaml

# Delete ConfigMap
kubectl delete -n feed-fetcher configmap rabbitmq-config

# Delete the 'feed-fetcher' namespace
kubectl delete namespace feed-fetcher
