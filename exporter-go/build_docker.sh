#!/bin/bash

IMAGE=${1:-registry.dso.mil/platform-one/big-bang/apps/core/cluster-auditor/opa-exporter:v0.0.5}


docker build --tag="${IMAGE}" .
docker push "${IMAGE}"
