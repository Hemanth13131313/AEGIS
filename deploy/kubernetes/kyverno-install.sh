#!/usr/bin/env bash
# Installs Kyverno and applies AEGIS policies
set -euo pipefail
helm repo add kyverno https://kyverno.github.io/kyverno/
helm repo update
helm install kyverno kyverno/kyverno -n kyverno --create-namespace \
  --set admissionController.replicas=3 \
  --set backgroundController.replicas=2
echo "Waiting for Kyverno to be ready..."
kubectl wait --for=condition=available --timeout=120s deployment/kyverno-admission-controller -n kyverno
kubectl apply -f deploy/kubernetes/kyverno-policies.yaml
echo "Kyverno installed and AEGIS policies applied."
