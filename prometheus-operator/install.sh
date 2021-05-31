#!/bin/bash

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
helm install prometheus prometheus-community/prometheus --atomic --namespace prometheus -f values.yaml