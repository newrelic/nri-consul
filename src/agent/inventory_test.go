package agent

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/data/inventory"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-consul/src/args"
	"github.com/newrelic/nri-consul/src/testutils"
)

func TestCollectInventory(t *testing.T) {
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

	entity, err := i.Entity("test", "agent")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	agents := []*Agent{
		{
			client: client,
			entity: entity,
		},
	}

	mux.HandleFunc("/v1/agent/self", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"Config": {
				"Datacenter": "dev-uss",
				"NodeName": "vault-dev-uss-0",
				"NodeID": "70ccd111-96de-a058-7be7-81a00fa7ca17",
				"Revision": "39f93f011",
				"Server": false,
				"Version": "1.2.1"
			},
			"DebugConfig": {
				"ACLDatacenter": "dev-uss",
				"ACLEnableKeyListPolicy": false,
				"AutopilotMaxTrailingLogs": 250,
				"AutopilotRedundancyZoneTag": "",
				"Checks": [],
				"ClientAddrs": [
					"127.0.0.1"
				],
				"ConnectCAConfig": {}
			}
		}`)
	})

	expected := inventory.Items{
		"Config/Datacenter": inventory.Item{
			"value": "dev-uss",
		},
		"Config/NodeName": inventory.Item{
			"value": "vault-dev-uss-0",
		},
		"Config/NodeID": inventory.Item{
			"value": "70ccd111-96de-a058-7be7-81a00fa7ca17",
		},
		"Config/Revision": inventory.Item{
			"value": "39f93f011",
		},
		"Config/Server": inventory.Item{
			"value": false,
		},
		"Config/Version": inventory.Item{
			"value": "1.2.1",
		},
		"DebugConfig/ACLDatacenter": inventory.Item{
			"value": "dev-uss",
		},
		"DebugConfig/ACLEnableKeyListPolicy": inventory.Item{
			"value": false,
		},
		"DebugConfig/AutopilotMaxTrailingLogs": inventory.Item{
			"value": float64(250),
		},
		"DebugConfig/ClientAddrs": inventory.Item{
			"value": []interface{}{
				"127.0.0.1",
			},
		},
	}

	doneChan := make(chan bool)
	go func() {
		CollectInventory(agents)
		close(doneChan)
	}()

	select {
	case <-doneChan:
		out := agents[0].entity.Inventory.Items()
		if !reflect.DeepEqual(out, expected) {
			t.Errorf("Expected %+v got %+v", expected, out)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timed out")
	}
}
