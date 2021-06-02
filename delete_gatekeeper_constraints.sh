#!/bin/bash

set -eux

constraints=$(find gatekeeper-library/library/general -name constraint.yaml)

for c in $constraints
do
  ytt -f $c -f config.yaml | kubectl delete -f -
done

kustomize build gatekeeper-library/library/general | kubectl delete -n styra-system -f -

