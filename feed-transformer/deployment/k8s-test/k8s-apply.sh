#!/bin/bash

# Declare environment variables
export RABBITMQ_USERNAME=${RABBITMQ_USERNAME:-vlad}
export RABBITMQ_PASSWORD=${RABBITMQ_PASSWORD:-sport}
export MONGO_USERNAME=${MONGO_USERNAME:-vlad}
export MONGO_PASSWORD=${MONGO_PASSWORD:-sport}
export TEST_IMAGE=${TEST_IMAGE:-sport-news/feed-transformer:local}
export RELATIVE_PATH=${RELATIVE_PATH:-}

# Create Namespace
kubectl apply -f "${RELATIVE_PATH}feed-transformer-namespace.yaml"

# Prepare ConfigMap YAML with variable substitution
envsubst < "${RELATIVE_PATH}rabbitmq-configmap.yaml" | kubectl apply -f -
envsubst < "${RELATIVE_PATH}mongodb-configmap.yaml" | kubectl apply -f -

# Deploy RabbitMQ
kubectl apply -f "${RELATIVE_PATH}rabbitmq-deployment.yaml"
kubectl apply -f "${RELATIVE_PATH}rabbitmq-service.yaml"

# Deploy MongoDB
kubectl apply -f "${RELATIVE_PATH}mongodb-deployment.yaml"
kubectl apply -f "${RELATIVE_PATH}mongodb-service.yaml"

# Deploy Redis
kubectl apply -f "${RELATIVE_PATH}redis-deployment.yaml"
kubectl apply -f "${RELATIVE_PATH}redis-service.yaml"

# Wait for RabbitMQ to be ready
echo "Waiting for RabbitMQ to become ready..."
if kubectl wait --namespace feed-transformer --for=condition=ready pod -l app=rabbitmq --timeout=180s; then
    echo "RabbitMQ is ready."

    # Wait for MongoDB to be ready
    echo "Waiting for MongoDB to become ready..."
    if kubectl wait --namespace feed-transformer --for=condition=ready pod -l app=mongodb --timeout=180s; then
        echo "MongoDB is ready."

        # Wait for Redis to be ready
        echo "Waiting for Redis to become ready..."
        if kubectl wait --namespace feed-transformer --for=condition=ready pod -l app=redis --timeout=180s; then

            # Deploy Feed transformer
            envsubst < "${RELATIVE_PATH}feed-transformer-deployment.yaml" | kubectl apply -f -
            kubectl apply -f "${RELATIVE_PATH}feed-transformer-service.yaml"

        else
            echo "Redis is not ready within the timeout period. Exiting..."
            ./k8s-delete.sh
            exit 1
        fi
    else
        echo "MongoDB is not ready within the timeout period. Exiting..."
        ./k8s-delete.sh
        exit 1
    fi

else
    echo "RabbitMQ is not ready within the timeout period. Exiting..."
    ./k8s-delete.sh
    exit 1
fi
