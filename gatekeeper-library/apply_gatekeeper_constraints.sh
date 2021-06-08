#!/bin/bash

set -eux

kustomize build library/library/general | kubectl apply -n styra-system -f -


constraints=$(find library/library/general -name constraint.yaml)

for c in $constraints
do
  ytt -f $c -f config.yaml | kubectl apply -f - 
done

