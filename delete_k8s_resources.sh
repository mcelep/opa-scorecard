#!/bin/bash

set -eux

kustomize build . | kubectl delete -n styra-system -f -

