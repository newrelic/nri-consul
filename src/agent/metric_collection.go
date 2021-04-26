package agent

import (
	"sync"

	"github.com/newrelic/infra-integrations-sdk/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// CollectMetrics does a metric collect for a group of agents
func CollectMetrics(agents []*Agent) {
	var wg sync.WaitGroup
	agentChan := createMetricPool(&wg)

	for _, agent := range agents {
		agentChan <- agent
	}

	close(agentChan)

	wg.Wait()
}

func createMetricPool(wg *sync.WaitGroup) chan *Agent {
	agentChan := make(chan *Agent)
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go metricWorker(agentChan, wg)
	}

	return agentChan
}

func metricWorker(agentChan <-chan *Agent, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		agent, ok := <-agentChan
		if !ok {
			return
		}

		CollectMetricsFromOne(agent)

	}
}

// CollectMetricsFromOne does a metric collect for a single agent
func CollectMetricsFromOne(agent *Agent) {
	metricSet := agent.entity.NewMetricSet("ConsulAgentSample",
		attribute.Attribute{Key: "displayName", Value: agent.entity.Metadata.Name},
		attribute.Attribute{Key: "entityName", Value: agent.entity.Metadata.Namespace + ":" + agent.entity.Metadata.Name},
		attribute.Attribute{Key: "ip", Value: agent.ipAddr},
		attribute.Attribute{Key: "datacenter", Value: agent.datacenter},
	)

	// Collect core metrics
	if err := agent.CollectCoreMetrics(metricSet, gaugeMetrics, counterMetrics, timerMetrics); err != nil {
		log.Error("Error collecting core metrics for Agent '%s': %s", agent.entity.Metadata.Name, err.Error())
	}

	// Peer Count
	if err := agent.collectPeerCount(metricSet); err != nil {
		log.Error("Error collecting peer count for Agent '%s': %s", agent.entity.Metadata.Name, err.Error())
	}

	// Latency metrics
	if err := agent.collectLatencyMetrics(metricSet); err != nil {
		log.Error("Error collecting latency metrics for Agent '%s': %s", agent.entity.Metadata.Name, err.Error())
	}

}
