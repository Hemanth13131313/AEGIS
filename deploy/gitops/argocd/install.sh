#!/usr/bin/env bash
# Bootstraps ArgoCD and registers the AEGIS applications
set -euo pipefail
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
echo "Waiting for ArgoCD..."
kubectl wait --for=condition=available --timeout=180s deployment/argocd-server -n argocd
kubectl apply -f deploy/gitops/argocd/project.yaml
kubectl apply -f deploy/gitops/argocd/apps/
echo "ArgoCD bootstrapped. Access via: kubectl port-forward svc/argocd-server -n argocd 8080:443"
