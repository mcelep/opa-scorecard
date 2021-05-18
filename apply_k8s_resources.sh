#!/bin/bash

set -eux

kustomize build . | kubectl apply -n styra-system -f -

