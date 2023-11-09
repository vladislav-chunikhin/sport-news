#!/bin/bash

# Create Namespace
kubectl create namespace feed-fetcher

# Deploy RabbitMQ
kubectl apply -f rabbitmq-deployment.yaml
kubectl apply -f rabbitmq-service.yaml

# Wait for RabbitMQ to be ready
echo "Waiting for RabbitMQ to become ready..."
while true; do
    # Check if RabbitMQ pods are ready
    if kubectl wait --namespace feed-fetcher --for=condition=ready pod -l app=rabbitmq --timeout=300s; then
        echo "RabbitMQ is ready."
        break
    else
        echo "RabbitMQ is not ready yet. Retrying in 10 seconds..."
        sleep 10
    fi
done

# Deploy Feed Fetcher
kubectl apply -f feed-fetcher-deployment.yaml
kubectl apply -f feed-fetcher-service.yaml
