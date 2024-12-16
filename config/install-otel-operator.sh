#!/bin/bash

# Define the URLs for Cert-Manager and OpenTelemetry Operator
CERT_MANAGER_URL="https://github.com/cert-manager/cert-manager/releases/download/v1.8.2/cert-manager.yaml"
OTEL_OPERATOR_URL="https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml"

# Apply Cert-Manager manifest
echo "Deploying Cert-Manager..."
kubectl apply -f $CERT_MANAGER_URL
if [ $? -ne 0 ]; then
  echo "Failed to deploy Cert-Manager. Please check the output for errors."
  exit 1
fi

# Wait for Cert-Manager components to be ready
echo "Waiting for Cert-Manager to become ready..."
kubectl wait --namespace cert-manager \
  --for=condition=Ready pods \
  --all \
  --timeout=120s

# Apply OpenTelemetry Operator manifest
echo "Deploying OpenTelemetry Operator..."
kubectl apply -f $OTEL_OPERATOR_URL
if [ $? -ne 0 ]; then
  echo "Failed to deploy OpenTelemetry Operator. Please check the output for errors."
  exit 1
fi

# Wait for OpenTelemetry Operator components to be ready
echo "Waiting for OpenTelemetry Operator to become ready..."
kubectl wait --namespace opentelemetry-operator-system \
  --for=condition=Ready pods \
  --all \
  --timeout=120s

echo "Cert-Manager and OpenTelemetry Operator have been successfully deployed!"

# https://github.com/isItObservable/Otel-Collector-Observek8s/blob/master/README.md
