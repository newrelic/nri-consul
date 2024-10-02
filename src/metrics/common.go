package metrics

import (
	"github.com/newrelic/infra-integrations-sdk/v3/data/metric"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
)

// SetMetric is a wrappper around metric.Set.SetMetric with error logging
func SetMetric(metricSet *metric.Set, name string, value interface{}, sourceType metric.SourceType) {
	if err := metricSet.SetMetric(name, value, sourceType); err != nil {
		log.Error("Error setting metric %s: %s", name, err.Error())
	}
}
