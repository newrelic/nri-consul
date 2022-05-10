package main_test

import (
	"fmt"
	"net/http"
	"testing"

	"time"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/nri-consul/src/args"
	"github.com/newrelic/nri-consul/src/testutils"
	"github.com/stretchr/testify/require"
)

func Test_ClientTimeoutIsHonored(t *testing.T) {
	mux, hostname, port, serverClose := testutils.SetupServer()
	defer serverClose()

	mux.HandleFunc("/v1/agent/members", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
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

	testCases := []struct {
		name          string
		timeout       string
		errorExpected bool
	}{
		{
			name:          "When the timeout is exceeded Then an error is returned",
			timeout:       "1s",
			errorExpected: true,
		},
		{
			name:          "When the timeout is not exceeded a correct response is retrieved",
			timeout:       "2s",
			errorExpected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			arg := args.ArgumentList{
				Hostname:  hostname,
				Port:      port,
				EnableSSL: false,
				Timeout:   tt.timeout,
			}
			apiConfig, err := arg.CreateAPIConfig(arg.Hostname)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err.Error())
			}
			client, err := api.NewClient(apiConfig)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err.Error())
			}

			agentMembers, err := client.Agent().Members(false)
			if tt.errorExpected {
				require.Error(t, err)
				return
			}
			require.NotEmpty(t, agentMembers)
		})
	}
}
