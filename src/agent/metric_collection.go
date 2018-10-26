package agent

import (
	"sync"

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

		if err := agent.CollectMetrics(gaugeMetrics, counterMetrics, timerMetrics); err != nil {
			log.Error("Error collecting core metrics for Agent '%s': %s", agent.entity.Metadata.Name, err.Error())
		}
	}
}
