package agent

import (
	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-consul/src/metrics"
)



// CollectMetrics collects metrics for an Agent
func (a *Agent) CollectMetrics(gaugeDefs, counterDefs []*metrics.MetricDefinition, timerDefs []*metrics.TimerDefinition) error {
	log.Debug("Starting Metric collection for Agent %s", a.entity.Metadata.Name)
	metricInfo, err := a.client.Agent().Metrics()
	if err != nil {
		return err
	}

	metricSet := a.entity.NewMetricSet("ConsulAgentSample",
		metric.Attribute{Key: "displayName", Value: a.entity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: a.entity.Metadata.Namespace + ":" + a.entity.Metadata.Name},
	)

	// collect gauges
	collectGaugeMetrics(metricSet, metricInfo.Gauges, gaugeDefs)

	// collect counters
	collectCounterMetrics(metricSet, metricInfo.Counters, counterDefs)

	// collect timers
	collectTimerMetrics(metricSet, metricInfo.Samples, timerDefs)

	log.Debug("Finished Metric collection for Agent %s", a.entity.Metadata.Name)
	return nil
}

func collectGaugeMetrics(metricSet *metric.Set, gauges []api.GaugeValue, defs []*metrics.MetricDefinition) {
	for _, def := range defs {
		found := false

		// Look through all gauges for metric
		for _, gague := range gauges {
			// If found, record and break
			if def.APIKey == gague.Name {
				found = true
				if err := metricSet.SetMetric(def.MetricName, gague.Value, def.SourceType); err != nil {
					log.Error("Error setting metric %s: %s", def.MetricName, err.Error())
				}

				break
			}
		}

		if !found {
			log.Debug("Did not find metric '%s' matching API key '%s'", def.MetricName, def.APIKey)
		}
	}
}

func collectCounterMetrics(metricSet *metric.Set, counters []api.SampledValue, defs []*metrics.MetricDefinition) {
	for _, def := range defs {
		found := false

		// Look through all counters for metric
		for _, counter := range counters {
			// If found, record and break
			if def.APIKey == counter.Name {
				found = true
				if err := metricSet.SetMetric(def.MetricName, counter.Count, def.SourceType); err != nil {
					log.Error("Error setting metric %s: %s", def.MetricName, err.Error())
				}

				break
			}
		}

		if !found {
			log.Debug("Did not find metric '%s' matching API key '%s'", def.MetricName, def.APIKey)
		}
	}
}

func collectTimerMetrics(metricSet *metric.Set, timers []api.SampledValue, defs []*metrics.TimerDefinition) {
	lookup := make(map[string]*api.SampledValue)

	for _, def := range defs {

		// Check if the timer is cached, if not search for it.
		sample, ok := lookup[def.APIKey]
		if !ok {
			for _, timer := range timers {
				if def.APIKey == timer.Name {
					sample = &timer
					lookup[def.APIKey] = sample
					break
				}
			}
		}

		if sample == nil {
			log.Debug("Did not find metric '%s' matching API key '%s'", def.MetricName, def.APIKey)
			continue
		}

		// Calculate/collect statistical sample
		value := calculateStatValue(def.Operation, sample)
		if err := metricSet.SetMetric(def.MetricName, value, def.SourceType); err != nil {
			log.Error("Error setting metric %s: %s", def.MetricName, err.Error())
		}
	}
}

func calculateStatValue(operation metrics.StatOperation, sample *api.SampledValue) float64 {
	var value float64
	switch operation {
	case metrics.Average:
		value = sample.Mean
	case metrics.Max:
		value = sample.Max
	case metrics.Median:
		value = (sample.Min + sample.Max) / 2.0
	case metrics.Count:
		value = float64(sample.Count)
	}

	return value
}
