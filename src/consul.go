package main

import (
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-consul/src/args"
)

const (
	integrationName    = "com.newrelic.consul"
	integrationVersion = "0.1.0"
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
	client, err := api.NewClient(args.CreateAPIConfig())
	if err != nil {
		log.Error("Error creating API client, please check configuration: %s", err.Error())
		os.Exit(1)
	}

	members, err := client.Agent().Members(false)
	if err != nil {
		log.Error("Error getting members: %s", err.Error())
		os.Exit(1)
	}

	for _, member := range members {
		log.Info(member.Name)
	}

	if err = i.Publish(); err != nil {
		log.Error(err.Error())
	}
}
