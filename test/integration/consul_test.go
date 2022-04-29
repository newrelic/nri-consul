//go:build integration
// +build integration

package tests

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-consul/src/args"
	"github.com/stretchr/testify/assert"
)

const (
	containerIntegration = "nri-consul"
	schema               = "consul-schema.json"
)

var containersConsul = []string{"consul-server1", "consul-server2", "consul-server3"}

func TestMain(m *testing.M) {
	flag.Parse()
	result := m.Run()
	os.Exit(result)
}

func TestSuccessConnection(t *testing.T) {
	if !waitForConsulClusterUpAndRunning(20) {
		t.Fatal("tests cannot be executed")
	}
	hostname := "consul-server1"
	envVars := []string{
		fmt.Sprintf("HOSTNAME=%s", hostname),
	}
	response, stderr, err := dockerComposeRun(envVars, containerIntegration)
	fmt.Println(response)
	fmt.Println(stderr)
	assert.Nil(t, err)
	assert.NotEmpty(t, response)
	err = validateJSONSchema(schema, response)
	assert.NoError(t, err)
}

func waitForConsulClusterUpAndRunning(maxTries int) bool {
	stdout, stderr, err := dockerComposeUp([]string{}, containersConsul)
	fmt.Println(stdout)
	fmt.Println(stderr)
	if err != nil {
		log.Fatal(err)
	}
	arg := args.ArgumentList{
		Hostname: "localhost",
		Port:     "8500",
	}
	apiConfig, err := arg.CreateAPIConfig(arg.Hostname)
	if err != nil {
		log.Error("Error creating HTTP API client, please check configuration: %s", err.Error())
		os.Exit(1)
	}

	// create client
	client, err := api.NewClient(apiConfig)
	if err != nil {
		log.Error("Error creating API client, please check configuration: %s", err.Error())
		os.Exit(1)
	}

	for ; maxTries > 0; maxTries-- {
		log.Info("try to establish de connection with the Consul cluster...")
		m, err := client.Agent().Metrics()
		if err != nil {
			log.Warn("Api not ready")
			time.Sleep(2 * time.Second)
		}
		if m != nil && len(m.Gauges) > 0 {
			nodes, _, err := client.Coordinate().Nodes(nil)
			if err != nil {
				log.Warn("Api not ready")
				time.Sleep(2 * time.Second)
				continue
			}

			if len(nodes) <= 1 {
				log.Warn("nodes ready of 3: %d", len(nodes))
				time.Sleep(2 * time.Second)
				continue
			}

			log.Info("consul cluster is up & running!")
			return true
		}
	}
	return true
}
