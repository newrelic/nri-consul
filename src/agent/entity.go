// Package agent handles Agent entity creation and inventory/metric collection
package agent

import (
	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-consul/src/args"
)

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
