#!/bin/bash

# Declare environment variables
export RABBITMQ_USERNAME=${RABBITMQ_USERNAME:-vlad}
export RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD:-sport}

# Create Namespace
kubectl create namespace feed-fetcher

# Create or update the ConfigMap with the environment variables
kubectl create configmap rabbitmq-config --namespace=feed-fetcher \
  --from-literal=RABBITMQ_DEFAULT_USER="$RABBITMQ_USERNAME" \
  --from-literal=RABBITMQ_DEFAULT_PASS="$RABBITMQ_PASSWORD" \
  --dry-run=client -o yaml | kubectl apply -f -

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
