package agent

import (
	"sync"

	"github.com/newrelic/infra-integrations-sdk/log"
)

const inventoryWorkers = 5

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
	wg.Add(inventoryWorkers)
	for i := 0; i < inventoryWorkers; i++ {
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
				a.setInventoryItem(configPrefix+"/"+key, "value", v)
			}
		default:
			a.setInventoryItem(configPrefix+"/"+key, "value", v)
		}
	}
}

func (a *Agent) setInventoryItem(key, field string, value interface{}) {
	if err := a.entity.SetInventoryItem(key, field, value); err != nil {
		log.Debug("Error setting Inventory item '%s' on Agent '%s': %s", key, a.entity.Metadata.Name)
	}
}
