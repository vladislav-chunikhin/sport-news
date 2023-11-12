#!/bin/bash

# Declare environment variables
export MONGO_USERNAME=${MONGO_USERNAME:-vlad}
export MONGO_PASSWORD=${MONGO_PASSWORD:-sport}
export TEST_IMAGE=${TEST_IMAGE:-sport-news/feed-service:local}
export RELATIVE_PATH=${RELATIVE_PATH:-}

# Create Namespace
kubectl apply -f "${RELATIVE_PATH}feed-app-namespace.yaml"

# Prepare ConfigMap YAML with variable substitution
envsubst < "${RELATIVE_PATH}mongodb-configmap.yaml" | kubectl apply -f -

# Deploy MongoDB
kubectl apply -f "${RELATIVE_PATH}mongodb-deployment.yaml"
kubectl apply -f "${RELATIVE_PATH}mongodb-service.yaml"

# Wait for MongoDB to be ready
echo "Waiting for MongoDB to become ready..."
if kubectl wait --namespace feed-service --for=condition=ready pod -l app=mongodb --timeout=180s; then
    echo "MongoDB is ready."
    # Deploy Feed Service
    envsubst < "${RELATIVE_PATH}feed-app-deployment.yaml" | kubectl apply -f -
    kubectl apply -f "${RELATIVE_PATH}feed-app-service.yaml"
else
    echo "MongoDB is not ready within the timeout period. Exiting..."
    ./k8s-delete.sh
    exit 1
fi
