#!/bin/bash

export RELATIVE_PATH=${RELATIVE_PATH:-}

# Delete Feed Fetcher Service
kubectl delete -f "${RELATIVE_PATH}feed-service-service.yaml"

# Delete Feed Fetcher Deployment
kubectl delete -f "${RELATIVE_PATH}feed-service-deployment.yaml"

# Delete RabbitMQ Service
kubectl delete -f "${RELATIVE_PATH}rabbitmq-service.yaml"

# Delete RabbitMQ Deployment
kubectl delete -f "${RELATIVE_PATH}rabbitmq-deployment.yaml"

# Delete RabbitMQ ConfigMap
kubectl delete -f "${RELATIVE_PATH}rabbitmq-configmap.yaml"

# Delete Feed Fetcher Namespace
kubectl delete -f "${RELATIVE_PATH}feed-service-namespace.yaml"
