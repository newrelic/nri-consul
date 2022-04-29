//go:build integration
// +build integration

package tests

import (
	"flag"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
	validateJSONSchema(schema, response)
}

func waitForConsulClusterUpAndRunning(maxTries int) bool {
	stdout, stderr, err := dockerComposeUp([]string{}, containersConsul)
	fmt.Println(stdout)
	fmt.Println(stderr)
	if err != nil {
		log.Fatal(err)
	}
	for ; maxTries > 0; maxTries-- {
		log.Info("try to establish de connection with the Consul cluster...")

	}
	return true
}
