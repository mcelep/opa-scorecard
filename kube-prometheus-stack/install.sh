#!/bin/bash

kubectl create namespace prometheus

kubectl apply -f cm-custom-dashboard.yaml -n prometheus

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
helm install prometheus prometheus-community/kube-prometheus-stack --atomic --namespace prometheus -f values.yaml