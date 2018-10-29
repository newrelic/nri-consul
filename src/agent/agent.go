// Package agent handles Agent entity creation and inventory/metric collection
package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-consul/src/args"
	"github.com/newrelic/nri-consul/src/metrics"
)

// number of workers there can be per pool
const workerCount = 5

// Agent represents a Consul agent.
// It's comprised of the client connected to that agent
// and the Entity representing it.
type Agent struct {
	entity *integration.Entity
	client *api.Client
}

// CreateAgents creates an Agent structure for every Agent member of the LAN cluster
func CreateAgents(client *api.Client, i *integration.Integration, args *args.ArgumentList) ([]*Agent, error) {
	members, err := client.Agent().Members(false)
	if err != nil {
		log.Error("Error getting members: %s", err.Error())
		return nil, err
	}

	agents := make([]*Agent, 0, len(members))
	for _, member := range members {
		var agent Agent

		agent.entity, err = i.Entity(member.Name, "agent")
		if err != nil {
			log.Error("Error creating entity for Agent '%s': %s", member.Name, err.Error())
			continue
		}

		agent.client, err = api.NewClient(args.CreateAPIConfig(member.Name))
		if err != nil {
			log.Error("Error creating client for Agent '%s': %s", member.Name, err.Error())
			continue
		}

		agents = append(agents, &agent)
	}

	return agents, nil
}

func (a *Agent) processConfig(config map[string]interface{}, configPrefix string) {
	for key, value := range config {
		switch v := value.(type) {
		case map[string]interface{}:
			log.Debug("Not processing config param '%s' nested object", key)
		case string:
			if v != "" {
				a.setInventoryItem(configPrefix+"/"+key, "value", v)
			}
		case []interface{}:
			if len(v) > 0 {
				if stringVal, err := arrayToString(v); err != nil {
					log.Debug("Unable to store config param '%s': %s", key, err.Error())
				} else {
					a.setInventoryItem(configPrefix+"/"+key, "value", *stringVal)
				}
			}
		default:
			a.setInventoryItem(configPrefix+"/"+key, "value", v)
		}
	}
}

// setInventoryItem adds a wrapper around setting an inventory item
func (a *Agent) setInventoryItem(key, field string, value interface{}) {
	if err := a.entity.SetInventoryItem(key, field, value); err != nil {
		log.Debug("Error setting Inventory item '%s' on Agent '%s': %s", key, a.entity.Metadata.Name)
	}
}

// collectPeerCount counts the number of peers for an agent.
func (a *Agent) collectPeerCount(metricSet *metric.Set) error {
	log.Debug("Starting peer count collection for Agent %s", a.entity.Metadata.Name)

	peers, err := a.client.Status().Peers()
	if err != nil {
		return err
	}

	metrics.SetMetric(metricSet, "consul.peers", len(peers), metric.GAUGE)

	log.Debug("Finished peer count collection for Agent %s", a.entity.Metadata.Name)
	return nil
}

func (a *Agent) collectLatencyMetrics(metricSet *metric.Set) error {
	log.Debug("Starting latency metric collection for Agent %s", a.entity.Metadata.Name)

	nodes, _, err := a.client.Coordinate().Nodes(nil)
	if err != nil {
		return err
	}

	if len(nodes) == 1 {
		return errors.New("cluster only contains 1 node")
	}

	agentNode := findNode(a.entity.Metadata.Name, nodes)
	if agentNode == nil {
		return errors.New("could not find node for agent")
	}

	// calculate and popluate metrics
	calculateLatencyMetrics(metricSet, agentNode, nodes)

	log.Debug("Finished latency metric collection for Agent %s", a.entity.Metadata.Name)
	return nil
}

// CollectCoreMetrics collects metrics for an Agent
func (a *Agent) CollectCoreMetrics(metricSet *metric.Set, gaugeDefs, counterDefs []*metrics.MetricDefinition, timerDefs []*metrics.TimerDefinition) error {
	log.Debug("Starting core metric collection for Agent %s", a.entity.Metadata.Name)
	metricInfo, err := a.client.Agent().Metrics()
	if err != nil {
		return err
	}

	// collect gauges
	collectGaugeMetrics(metricSet, metricInfo.Gauges, gaugeDefs)

	// collect counters
	collectCounterMetrics(metricSet, metricInfo.Counters, counterDefs)

	// collect timers
	collectTimerMetrics(metricSet, metricInfo.Samples, timerDefs)

	log.Debug("Finished core metric collection for Agent %s", a.entity.Metadata.Name)
	return nil
}

func collectGaugeMetrics(metricSet *metric.Set, gauges []api.GaugeValue, defs []*metrics.MetricDefinition) {
	for _, def := range defs {
		found := false

		// Look through all gauges for metric
		for _, gauge := range gauges {
			// If found, record and break
			if def.APIKey == gauge.Name {
				value := gauge.Value

				// special case where we need to convert nanoseconds to milliseconds
				if def.APIKey == "consul.runtime.total_gc_pause_ns" {
					value /= 1000000
				}

				found = true
				metrics.SetMetric(metricSet, def.MetricName, value, def.SourceType)
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
				metrics.SetMetric(metricSet, def.MetricName, counter.Count, def.SourceType)
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
		metrics.SetMetric(metricSet, def.MetricName, value, def.SourceType)
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

// arrayToString converts an interface array to a comma delimited string if possible
func arrayToString(input []interface{}) (*string, error) {
	stringElements := make([]string, len(input))

	for i, elem := range input {
		elemString, ok := elem.(string)
		if !ok {
			return nil, fmt.Errorf("could not convert %v of type %T to string", elem, elem)
		}

		stringElements[i] = elemString
	}

	outString := strings.Join(stringElements, ",")

	return &outString, nil
}
