package agent

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-consul/src/args"
	"github.com/newrelic/nri-consul/src/testutils"
)

func TestCreateAgents(t *testing.T) {
	mux, hostname, port, serverClose := testutils.SetupServer()
	defer serverClose()

	arg := args.ArgumentList{
		Hostname:  hostname,
		Port:      port,
		EnableSSL: false,
	}

	client, err := api.NewClient(arg.CreateAPIConfig(arg.Hostname))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	i, err := integration.New("test", "1.0.0")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	mux.HandleFunc("/v1/agent/members", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"Name": "consul-0",
				"Addr": "10.0.0.1",
				"Port": 8301,
				"Tags": {
					"build": "1.2.1:39f93f01",
					"dc": "dev",
					"id": "c7f88fba-f8d9-94a9-3627-523398acf7db",
					"port": "8300",
					"raft_vsn": "3",
					"role": "consul",
					"segment": "",
					"vsn": "2",
					"vsn_max": "3",
					"vsn_min": "2",
					"wan_join_port": "8302"
				},
				"Status": 1,
				"ProtocolMin": 1,
				"ProtocolMax": 5,
				"ProtocolCur": 2,
				"DelegateMin": 2,
				"DelegateMax": 5,
				"DelegateCur": 4
			}
		]`)
	})

	out, err := CreateAgents(client, i, &arg)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if len(out) != 1 {
		t.Fatalf("Expected 1 agent got %d", len(out))
	}

	agent := out[0]
	if agent.entity.Metadata.Name != "consul-0" {
		t.Errorf("Expected Entity name 'consul-0' got %s", agent.entity.Metadata.Name)
	}
}
