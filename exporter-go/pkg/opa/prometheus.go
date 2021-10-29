package opa

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricSet struct {
	sync.RWMutex
	Metrics map[string]prometheus.Metric
}

var (
	namespace = "opa_scorecard"

	Up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last OPA scorecard query successful.",
		nil, nil,
	)
	ConstraintViolation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "constraint_violations"),
		"OPA violations for all constraints",
		[]string{"kind", "name", "violating_kind", "violating_name", "violating_namespace", "violation_msg", "violation_enforcement"}, nil,
	)
	ConstraintInformation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "constraint_information"),
		"Some general information of all constraints",
		[]string{"kind", "name", "enforcementAction", "totalViolations"}, nil,
	)
)

func ExportViolations(constraints []Constraint, metrics map[string]prometheus.Metric) {
	for _, c := range constraints {
		for _, v := range c.Status.Violations {
			// Abuse Go map key uniqueness to ensure we avoid duplicate metrics:
			// Since the entire content of a metric consists of strings, we can concatenate all properties of the metric to make a unique map key
			// If we encounter a duplicate metric, the previous instance (with the same map key) will simply be overwritten
			metrics[c.Meta.Kind+c.Meta.Name+v.Kind+v.Name+v.Namespace+v.Message+v.EnforcementAction] = prometheus.MustNewConstMetric(
				ConstraintViolation, prometheus.GaugeValue, 1, c.Meta.Kind, c.Meta.Name, v.Kind, v.Name, v.Namespace, v.Message, v.EnforcementAction)
		}
	}
}

func ExportConstraintInformation(constraints []Constraint, metrics map[string]prometheus.Metric) {
	for _, c := range constraints {
		metrics[c.Meta.Kind+c.Meta.Name+c.Spec.EnforcementAction+fmt.Sprintf("%f", c.Status.TotalViolations)] = prometheus.MustNewConstMetric(
			ConstraintInformation, prometheus.GaugeValue, 1, c.Meta.Kind, c.Meta.Name, c.Spec.EnforcementAction, fmt.Sprintf("%f", c.Status.TotalViolations))
	}
}
