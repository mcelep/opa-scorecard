---
title: Expose Open Policy Agent/Gatekeeper Violations for Kubernetes Applications with Prometheus and Grafana
tags: ['kubernetes','open policy agent', 'opa','violations','scorecard']
status: draft
---

*TL;DR: In this blog post, we talk about a solution which gives platform users a succint view about which Gatekeeper constraints are violated by using Prometheus & Grafana. *


# Expose Open Policy Agent/Gatekeeper Violations for Kubernetes Applications with Prometheus and Grafana

Application teams that just start to use Kubernetes might find it a bit difficult to get into as Kubernetes is a quiet complex & large ecosystem([CNCF ecosystem ](https://landscape.cncf.io/)). Moreover, although Kubernetes is starting to mature, it's still being developed very actively and it keeps getting new featuers at a facer pace than many other enterprise software out there. On top of that, Kubernetes deployments due to the integration requirements into the rest of a company's ecosystem (Authenticaton, Authorization, Security, Network,storage) are tailored specifically for each company. So even for a seasoned Kubernetes expert there are usually many things to consider to deploy an application in a way that it fulfills security, resiliency, performance requirements. How can you assure that applications that run on Kubernetes keep fulfilling those requirements?

## Enter OPA/Gatekeeper

[Open Policy Agent](https://www.openpolicyagent.org/) and its Kubernetes(K8S) targeting component [Gatekeeper](https://github.com/open-policy-agent/gatekeeper) gives you means to enforce policies on Kubernetes clusters. What we mean by policies here, is a formal definition of rules & best practices & behavior that you want to see in your company's Kubernetes clusters. When using OPA, you use a Domain Specific Language called [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) and you write down policies in *Rego* language. By doing this, you leave no room for misinterpretations that would maybe if you tried to explain a policy in free text on your company's internal wiki.

Moreover, when using Gatekeeper, different policies can have different enforcement actions. There might be certain policies that are treated as **MUST** whereas there might be policies which are **NICE-TO-HAVE**. A **MUST** policy will stop a Kubernetes resource being admissioned onto a cluster and a **NICE-TO-HAVE** policy will only cause warning messages which should be noted by platform users.


In this blog post, we talk about about how you can:
- Apply example Gatekeeper constraints to K8S clusters
- Expose prometheus metrics from Gatekeeper constraint violations
- Create a Grafana dashboard to display key information about violations

If you want to read more about enforcing policies in Kubernetes, check out [this article](https://itnext.io/enforcing-policies-in-kubernetes-c0f6192bd5ca).

## Design

!TODO Diagram 

### OPA Constraints

Every company has its own set of requirements for applications running on Kubernetes. You might have heard about the *production readiness checklist* concept(in a nutshell, you want to create a checklist of items for your platform users to use before they deploy an an application in production). You want to have your own *production readiness checklist* based on Rego and these links below might give you a good starting point for creating your own list:
- [Application Readiness Checklist on Tanzu Developer Center](https://tanzu.vmware.com/developer/guides/kubernetes/app-enhancements-checklist/)
- [Production best practices on learnk8s.io](https://learnk8s.io/production-best-practices) 
- [Kubernetes in Production: Readiness Checklist and Best Practices on replex.io](https://www.replex.io/blog/kubernetes-in-production-readiness-checklist-and-best-practices) 

Bear in mind that, you will need to create an OPA-ready of your production readiness checklist and you might not be able to cover all of the concerns you might have in your checklist using OPA/rego. The goal is to focus on things that are easy to extract based on Kubernetes resources definitions, e.g. number of replicas for a [K8S Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)

For our blog post, we will be using an open source project: [gatekeeper-library](https://github.com/open-policy-agent/gatekeeper-library) which contains a good set of example constraints. Moreover, the project structure is quite helpful in the sense of providing an example of how you can go about managing the OPA constraints. Rego language which is used for creating OPA policies should be unit tested thoroughly and in  [src folder](https://github.com/open-policy-agent/gatekeeper-library/tree/master/src/general), you can find  pure rego files and unit tests. In [this folder](https://github.com/open-policy-agent/gatekeeper-library/tree/master/library/general), there are the templates that are created out of the rego files in the [src folder](https://github.com/open-policy-agent/gatekeeper-library/tree/master/src/general) and for each template there's an example constraint together with some target data that would result in both positive and negative results for the constraint. Rego based policies can get quite complex, so in our view it's a must to have rego unit tests which cover both **happy & unhappy** paths.

As mentioned earlier, there might be certain constraints which you don't want to directly enforce(MUST ve NICE-TO-HAVE) e.g. on a dev cluster you might not want to enforce **>1 replicas** or maybe before you enforce a specific constraint you might want to give platform users enough grace period to take the necessary precautions. You control whether a constraint blocks a K8S resource from getting admitted is via ```spec.enforcementAction``` property. By default, ```enforcementAction``` is set to ```deny```. In our example, we will install all constraints with  ```enforcementAction: dryrun``` property, this will make sure that we don't directly impact any workload running on K8S clusters. The details of *dryrun* enforcement action are explained [here](https://open-policy-agent.github.io/gatekeeper/website/docs/violations/#dry-run-enforcement-action).


### Prometheus Metrics



- https://open-policy-agent.github.io/gatekeeper/website/docs/audit#audit-using-kinds-specified-in-the-constraints-only



### Action

#### Required tools/things

- Kubectl
- Ytt
- Kustomize
- A working K8S cluster
- Helm

#### Git submodule

Run ```git submodule update --init --recursive``` to get Gatekeeper dependency

#### Install OPA/Gatekeeper

We've used [Tanzu Mission Control(TMC)](https://tanzu.vmware.com/mission-control) to provision a Kubernetes cluster and TMC automatically gives us a cluster with Gatekeeper on it. If your cluster does not come with Gatekeeper preinstalled, you can use install it as explained [here](https://open-policy-agent.github.io/gatekeeper/website/docs/install/). If you are familiar with helm, the easiest way to install it is :
```bash
helm repo add gatekeeper https://open-policy-agent.github.io/gatekeeper/charts
helm install gatekeeper/gatekeeper --generate-name
```

#### Install Gatekeeper example constraints


Use script:```./apply_gatekeeper_constraints``` to create example constraints from gatekeeper-library.
Bear in mind that this script uses ytt to inject ```  enforcementAction: dryrun``` in order to not enforce any actions.


#### Install exporter 

```kubectl -n opa-exporter apply -f expoter-k8s-resources```

#### Install kube-stack

```cd kube-stack && ./install.sh```

This helm chart comes with some configuration to set up a Grafana dashboard.

#### Log on to Grafana


!TODO add screenshot