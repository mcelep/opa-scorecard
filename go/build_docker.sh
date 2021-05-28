#!/bin/bash

IMAGE=mcelep/opa_scorecard_exporter:v0.0.1
docker build --tag="${IMAGE}" .
docker push "${IMAGE}"