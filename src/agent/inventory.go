package agent

import (
	"errors"
	"strings"
	"sync"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// CollectInventory collects inventory data for each Agent entity
func CollectInventory(agentEntities *[]integration.Entity) {

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

		// Special case to get Agents role
		if memberConfig, ok := selfData["Member"]; ok {
			role, err := getRole(memberConfig)
			if err != nil {
				log.Warn("Error finding role attribute for Agent '%s': %s", agent.entity.Metadata.Name, err)
			} else {
				agent.setInventoryItem("Member/Tags/role", "value", role)
			}
		}
	}
}

func getRole(memberConfig map[string]interface{}) (interface{}, error) {
	tagsRaw, ok := memberConfig["Tags"]
	if !ok {
		return nil, errors.New("could not find Tags structure in Member")
	}

	tagsData := tagsRaw.(map[string]interface{})
	role, ok := tagsData["role"]
	if !ok {
		return nil, errors.New("could not find role attribute in Tags")
	}

	return role, nil
}

func (a *Agent) processConfig(config map[string]interface{}, configPrefix string) {
	for key, value := range config {
		switch v := value.(type) {
		case map[string]interface{}:
			log.Debug("Not processing config param '%s' nested object", key)
		case []string:
			a.setInventoryItem(configPrefix+"/"+key, "value", strings.Join(v, ","))
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
