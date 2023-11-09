#!/bin/bash

# Declare environment variables
export RABBITMQ_USERNAME=${RABBITMQ_USERNAME:-vlad}
export RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD:-sport}

# Create Namespace
kubectl apply -f feed-fetcher-namespace.yaml

# Prepare ConfigMap YAML with variable substitution
envsubst < rabbitmq-configmap.yaml | kubectl apply -f -

# Deploy RabbitMQ
kubectl apply -f rabbitmq-deployment.yaml
kubectl apply -f rabbitmq-service.yaml

# Wait for RabbitMQ to be ready
echo "Waiting for RabbitMQ to become ready..."
if kubectl wait --namespace feed-fetcher --for=condition=ready pod -l app=rabbitmq --timeout=180s; then
    echo "RabbitMQ is ready."
    # Deploy Feed Fetcher
    kubectl apply -f feed-fetcher-deployment.yaml
    kubectl apply -f feed-fetcher-service.yaml
else
    echo "RabbitMQ is not ready within the timeout period. Exiting..."
    ./k8s-delete.sh
    exit 1
fi
