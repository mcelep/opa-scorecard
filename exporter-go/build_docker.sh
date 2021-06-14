#!/bin/bash

IMAGE=${1:-mcelep/opa_scorecard_exporter:v0.0.3}


docker build --tag="${IMAGE}" .
docker push "${IMAGE}"
