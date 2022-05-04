//go:build integration
// +build integration

package tests

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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

var (
	containersConsul = []string{"consul-server1", "consul-server2", "consul-server3"}
	responseRegex    = regexp.MustCompile(`{"name":"com\.newrelic\.consul".*`)
)

func TestMain(m *testing.M) {
	flag.Parse()
	result := m.Run()
	os.Exit(result)
}

func TestMetricsCollectionInCluterWithSSL(t *testing.T) {
	if !waitForConsulClusterUpAndRunning(20) {
		t.Fatal("tests cannot be executed")
	}
	defer dockerComposeDown()
	hostname := "consul-server1"
	envVars := []string{
		fmt.Sprintf("HOSTNAME=%s", hostname),
		fmt.Sprintf("ENABLE_SSL=true"),
		fmt.Sprintf("CA_BUNDLE_FILE=/consul/config/certs/consul-agent-ca.pem"),
	}
	response, stderr, err := dockerComposeRun(envVars, containerIntegration)
	fmt.Println(response)
	fmt.Println(stderr)
	assert.Nil(t, err)
	assert.NotEmpty(t, response)

	cleanIntegrationResponse := responseRegex.FindString(response)
	err = validateJSONSchema(schema, cleanIntegrationResponse)
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
		Hostname:     "localhost",
		Port:         "8500",
		EnableSSL:    true,
		CABundleFile: "certs/consul-agent-ca.pem",
	}
	apiConfig := arg.CreateAPIConfig(arg.Hostname)

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
