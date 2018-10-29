package metrics

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// MetricDefinition represents a all the definition to collect
// a metric from the API and send to Infrastructure
type MetricDefinition struct {
	APIKey     string
	MetricName string
	SourceType metric.SourceType
}

// StatOperation represents a statistical operation for Timer Metrics
type StatOperation int

const (

	// Average represents the average of a Timer metric
	Average StatOperation = iota

	// Max represents the max of a Timer metric
	Max

	// Median represents the median of a Timer metric
	Median

	// Count represents the count of a Timer metric
	Count
)

// TimerDefinition represents a Timer metric and it's statistical
// operation from the timer data set
type TimerDefinition struct {
	MetricDefinition
	Operation StatOperation
}
