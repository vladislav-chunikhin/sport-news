#!/bin/bash

export RELATIVE_PATH=${RELATIVE_PATH:-}

# Delete Feed App Service
kubectl delete -f "${RELATIVE_PATH}feed-app-service.yaml"

# Delete Feed App Deployment
kubectl delete -f "${RELATIVE_PATH}feed-app-deployment.yaml"

# Delete MongoDB Service
kubectl delete -f "${RELATIVE_PATH}mongodb-service.yaml"

# Delete MongoDB Deployment
kubectl delete -f "${RELATIVE_PATH}mongodb-deployment.yaml"

# Delete MongoDB ConfigMap
kubectl delete -f "${RELATIVE_PATH}mongodb-configmap.yaml"

# Delete Feed App Namespace
kubectl delete -f "${RELATIVE_PATH}feed-app-namespace.yaml"
