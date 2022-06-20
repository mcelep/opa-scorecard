# opa-exporter

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0.4](https://img.shields.io/badge/AppVersion-0.0.4-informational?style=flat-square)

Prometheus exporter for OPA Gatekeeper.

## Get the Helm repository

```shell
helm repo add opa-exporter https://mcelep.github.io/opa-scorecard
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Installing the chart

To install the chart with the release name `my-release`:

```shell
helm install my-release opa-exporter/opa-exporter
```

## Uninstalling the chart

To uninstall the `my-release` release:

```shell
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Pod affinity |
| image.pullSecrets | list | `[]` | List of image pull secrets |
| image.repository | string | `"mcelep/opa_scorecard_exporter"` | Image repository and name |
| image.tag | string | `""` | Overrides the image tag whose default is the chart `appVersion` |
| nodeSelector | object | `{}` | Pod node selector |
| podAnnotations | object | `{}` | Pod annotations |
| podSecurityContext | object | `{}` | Pod security context |
| rbac.create | bool | `true` | Whether to create Cluster Role and Cluster Role Binding |
| rbac.extraClusterRoleRules | list | `[]` | Extra ClusterRole rules |
| rbac.useExistingRole | string | `nil` | Use an existing ClusterRole/Role |
| replicaCount | int | `1` | Count of Pod replicas |
| resources | object | `{}` | Resources for the Agent container |
| securityContext | object | `{}` | Security context for the Agent container |
| service.port | int | `80` | Service port |
| service.type | string | `"ClusterIP"` | Service type |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Whether to create the Service Account used by the Pod |
| serviceAccount.name | string | `""` | If not set and `create` is `true`, a name is generated using the fullname template |
| serviceMonitor.enabled | bool | `true` | Wherter to install `ServiceMonitor` or not |
| tolerations | list | `[]` | Pod tolerations |
