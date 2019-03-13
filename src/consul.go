package main

import (
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-consul/src/agent"
	"github.com/newrelic/nri-consul/src/args"
	"github.com/newrelic/nri-consul/src/datacenter"
)

const (
	integrationName    = "com.newrelic.consul"
	integrationVersion = "0.1.1"
)

func main() {
	var args args.ArgumentList
	// Create Integration
	i, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// Setup logging with verbose
	log.SetupLogging(args.Verbose)

	// create client
	client, err := api.NewClient(args.CreateAPIConfig(args.Hostname))
	if err != nil {
		log.Error("Error creating API client, please check configuration: %s", err.Error())
		os.Exit(1)
	}

	// Create the list of agents in LAN pool
	agents, leader, err := agent.CreateAgents(client, i, &args)
	if err != nil {
		log.Error("Error creating Agent entities: %s", err.Error())
		os.Exit(1)
	}

	dc, err := datacenter.NewDatacenter(leader, i)
	if err != nil {
		log.Error("Error creating Datacenter entity: %s", err.Error())
		if !args.TolerateLeaderConnectionIssues {
			os.Exit(1)
		}
	}

	// Collect inventory for agents
	if args.HasInventory() {
		agent.CollectInventory(agents)
	}

	// Collect metrics for Agents and cluster
	if args.HasMetrics() {
		agent.CollectMetrics(agents)
		if dc != nil {
			dc.CollectMetrics()
		}
	}

	if err = i.Publish(); err != nil {
		log.Error(err.Error())
	}
}
