
## Design

We use two key features from OPA:
- https://open-policy-agent.github.io/gatekeeper/website/docs/violations/#dry-run-enforcement-action
- https://open-policy-agent.github.io/gatekeeper/website/docs/audit#audit-using-kinds-specified-in-the-constraints-only



## Action

### Required tools

- Kubectl
- Ytt
- Kustomize

### Git submodule
Run ```git submodule update --init --recursive``` to get Gatekeeper dependency

### Install gatekeeper example constrints

Use script:```./apply_gatekeeper_constraints``` to create example constraints from gatekeeper-library.
Bear in mind that this script uses ytt to inject ```  enforcementAction: dryrun``` in order to not enforce any actions.



### Install exporter 

```kubectl -n opa-exporter apply -f expoter-k8s-resources```

### Install kube-stack

```cd kube-stack && ./install.sh```

This helm chart comes with some configuration to set up a Grafana dashboard

### Log on to Grafana


!TODO add screenshot