package agent

import (
	"sync"

	"github.com/newrelic/infra-integrations-sdk/log"
)

// CollectInventory collects inventory data for each Agent entity
func CollectInventory(agents []*Agent) {
	var wg sync.WaitGroup
	agentChan := createInventoryPool(&wg)

	for _, agent := range agents {
		agentChan <- agent
	}

	close(agentChan)

	wg.Wait()
}

func createInventoryPool(wg *sync.WaitGroup) chan *Agent {
	agentChan := make(chan *Agent)
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go inventoryWorker(agentChan, wg)
	}

	return agentChan
}

func inventoryWorker(agentChan <-chan *Agent, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		agent, ok := <-agentChan
		if !ok {
			return
		}

		selfData, err := agent.client.Agent().Self()
		if err != nil {
			log.Error("Error retrieving self configuration data for Agent '%s': %s", agent.entity.Metadata.Name, err.Error())
			continue
		}

		// Config data
		if configData, ok := selfData["Config"]; ok {
			agent.processConfig(configData, "Config")
		}

		// Debug config data
		if debugConfig, ok := selfData["DebugConfig"]; ok {
			agent.processConfig(debugConfig, "DebugConfig")
		}
	}
}
