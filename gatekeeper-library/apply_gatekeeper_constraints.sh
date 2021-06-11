#!/bin/bash

set -eux

kustomize build library/library/general | kubectl apply -n styra-system -f -

# wait a couple of seconds for all the constraint kinds to become ready
sleep 4

constraints=$(find library/library/general -name constraint.yaml)

for c in $constraints
do
  ytt -f $c -f config.yaml | kubectl apply -f - 
done

