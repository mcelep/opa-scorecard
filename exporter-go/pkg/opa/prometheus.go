package opa

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

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

func ExportViolations(constraints []Constraint) []prometheus.Metric {
	unique := make(map[string]bool)
	m := make([]prometheus.Metric, 0)
	for _, c := range constraints {
		for _, v := range c.Status.Violations {
			key := fmt.Sprintf("%v-%v-%v-%v-%v", c.Meta.Kind, c.Meta.Name, v.Name, v.Namespace, v.Message)
			if _, ok := unique[key]; ok {
				log.Printf("Found duplicate metrics: %v\n", key)
				continue
			}
			unique[key] = true
			metric := prometheus.MustNewConstMetric(ConstraintViolation, prometheus.GaugeValue, 1, c.Meta.Kind, c.Meta.Name, v.Kind, v.Name, v.Namespace, v.Message, v.EnforcementAction)
			m = append(m, metric)
		}
	}
	return m
}

func ExportConstraintInformation(constraints []Constraint) []prometheus.Metric {
	m := make([]prometheus.Metric, 0)
	for _, c := range constraints {
		metric := prometheus.MustNewConstMetric(ConstraintInformation, prometheus.GaugeValue, 1, c.Meta.Kind, c.Meta.Name, c.Spec.EnforcementAction, fmt.Sprintf("%f", c.Status.TotalViolations))
		m = append(m, metric)
	}
	return m
}
