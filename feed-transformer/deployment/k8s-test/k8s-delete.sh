#!/bin/bash

export RELATIVE_PATH=${RELATIVE_PATH:-}

# Delete Feed Transformer Service
kubectl delete -f "${RELATIVE_PATH}feed-transformer-service.yaml"

# Delete Feed Transformer Deployment
kubectl delete -f "${RELATIVE_PATH}feed-transformer-deployment.yaml"

# Delete RabbitMQ Service
kubectl delete -f "${RELATIVE_PATH}rabbitmq-service.yaml"

# Delete RabbitMQ Deployment
kubectl delete -f "${RELATIVE_PATH}rabbitmq-deployment.yaml"

# Delete RabbitMQ ConfigMap
kubectl delete -f "${RELATIVE_PATH}rabbitmq-configmap.yaml"

# Delete MongoDB Service
kubectl delete -f "${RELATIVE_PATH}mongodb-service.yaml"

# Delete MongoDB Deployment
kubectl delete -f "${RELATIVE_PATH}mongodb-deployment.yaml"

# Delete MongoDB ConfigMap
kubectl delete -f "${RELATIVE_PATH}mongodb-configmap.yaml"

# Delete Redis Service
kubectl delete -f "${RELATIVE_PATH}redis-service.yaml"

# Delete Redis Deployment
kubectl delete -f "${RELATIVE_PATH}redis-deployment.yaml"

# Delete Feed Transformer Namespace
kubectl delete -f "${RELATIVE_PATH}feed-transformer-namespace.yaml"
