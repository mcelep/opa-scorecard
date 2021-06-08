#!/bin/bash

set -eux

constraints=$(find library/library/general -name constraint.yaml)

for c in $constraints
do
  ytt -f $c -f config.yaml | kubectl delete -f -
done

kustomize build library/library/general | kubectl delete -n styra-system -f -

