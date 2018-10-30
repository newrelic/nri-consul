package cluster

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-consul/src/agent"
	"github.com/newrelic/nri-consul/src/args"
	"github.com/newrelic/nri-consul/src/testutils"
)

func TestNewCluster_NoLeader(t *testing.T) {
	i, err := integration.New("test", "1.0.0")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if _, err := NewCluster(nil, i); err == nil {
		t.Error("Expected error")
	}

}
func TestNewCluster_Normal(t *testing.T) {
	i, err := integration.New("test", "1.0.0")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	agent := &agent.Agent{}

	out, err := NewCluster(agent, i)
	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	if out.entity.Metadata.Name != "Cluster" {
		t.Fatalf("Entity was not named correctly %s", out.entity.Metadata.Name)
	} else if out.entity.Metadata.Namespace != "cluster" {
		t.Fatalf("Entity has wrong namespace %s", out.entity.Metadata.Namespace)
	}

	if out.leader != agent {
		t.Errorf("Agent was not set correctly expected %+v got %+v", agent, out.leader)
	}

}

func Test_Cluster_CollectMetrics_Full(t *testing.T) {
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

	clusterEntity, err := i.Entity("test", "cluster")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	agentEntity, err := i.Entity("leader", "agent")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	c := &Cluster{
		entity: clusterEntity,
		leader: agent.NewAgent(client, agentEntity),
	}

	setMetricMuxes(mux)

	expected := map[string]interface{}{
		"event_type":                          "ConsulClusterSample",
		"displayName":                         c.entity.Metadata.Name,
		"entityName":                          c.entity.Metadata.Namespace + ":" + c.entity.Metadata.Name,
		"leader":                              "leader",
		"raft.txns":                           float64(0),
		"raft.commitTimeAvgInMilliseconds":    float64(3),
		"raft.commitTimes":                    float64(0),
		"raft.commitTimeMedianInMilliseconds": float64(3),
		"raft.commitTimeMaxInMilliseconds":    float64(5),
		"catalog.registeredNodes":             float64(3),
		"catalog.criticalNodes":               float64(1),
		"catalog.upNodes":                     float64(1),
		"catalog.warningNodes":                float64(1),
		"catalog.passingNodes":                float64(1),
	}

	c.CollectMetrics()

	result := c.entity.Metrics[0].Metrics
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v got %+v", expected, result)
	}
}

func Test_Cluster_CollectMetrics_All_Endpoint_Fails(t *testing.T) {
	_, hostname, port, serverClose := testutils.SetupServer()
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

	clusterEntity, err := i.Entity("test", "cluster")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	agentEntity, err := i.Entity("leader", "agent")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	c := &Cluster{
		entity: clusterEntity,
		leader: agent.NewAgent(client, agentEntity),
	}

	expected := map[string]interface{}{
		"event_type":  "ConsulClusterSample",
		"displayName": c.entity.Metadata.Name,
		"entityName":  c.entity.Metadata.Namespace + ":" + c.entity.Metadata.Name,
		"leader":      "leader",
	}

	c.CollectMetrics()

	result := c.entity.Metrics[0].Metrics
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v got %+v", expected, result)
	}
}

func setMetricMuxes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/agent/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"Timestamp": "2018-10-26 14:17:50 +0000 UTC",
			"Gauges": [],
			"Points": [],
			"Counters": [
				{
					"Name": "consul.raft.apply",
					"Count": 2,
					"Rate": 0.2,
					"Sum": 2,
					"Min": 1,
					"Max": 1,
					"Mean": 1,
					"Stddev": 0,
					"Labels": {}
				}
			],
			"Samples": [
				{
					"Name": "consul.raft.commitTime",
					"Count": 1,
					"Rate": 0.16021829843521118,
					"Sum": 5,
					"Min": 1,
					"Max": 5,
					"Mean": 3,
					"Stddev": 0,
					"Labels": {}
				}
			]
		}`)
	})

	mux.HandleFunc("/v1/catalog/nodes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"ID": "1234",
				"Node": "consul-dev-0",
				"Address": "10.0.0.140",
				"Datacenter": "dev",
				"TaggedAddresses": {
					"lan": "10.0.0.140",
					"wan": "10.0.0.140"
				},
				"Meta": {
					"consul-network-segment": ""
				},
				"CreateIndex": 8,
				"ModifyIndex": 9
			},
			{
				"ID": "5678",
				"Node": "consul-dev-1",
				"Address": "10.0.0.142",
				"Datacenter": "dev",
				"TaggedAddresses": {
					"lan": "10.0.0.142",
					"wan": "10.0.0.142"
				},
				"Meta": {
					"consul-network-segment": ""
				},
				"CreateIndex": 65,
				"ModifyIndex": 67
			},
			{
				"ID": "91011",
				"Node": "vault-dev-1",
				"Address": "10.0.0.55",
				"Datacenter": "dev",
				"TaggedAddresses": {
					"lan": "10.0.0.55",
					"wan": "10.0.0.55"
				},
				"Meta": {
					"consul-network-segment": ""
				},
				"CreateIndex": 6,
				"ModifyIndex": 7
			}
		]`)
	})

	mux.HandleFunc("/v1/catalog/services", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"consul": [],
			"vault": [
				"active"
			]
		}`)
	})

	mux.HandleFunc("/v1/health/service/vault", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"Node": {
					"ID": "91011",
					"Node": "vault-dev-0",
					"Address": "10.0.0.55",
					"Datacenter": "dev",
					"TaggedAddresses": {
						"lan": "10.0.0.55",
						"wan": "10.0.0.55"
					},
					"Meta": {
						"consul-network-segment": ""
					},
					"CreateIndex": 20183,
					"ModifyIndex": 20185
				},
				"Service": {
					"ID": "vault:vault-dev-0.consul.localnet:8200",
					"Service": "vault",
					"Tags": [
						"active"
					],
					"Address": "vault-dev-0.consul.localnet",
					"Meta": null,
					"Port": 8200,
					"EnableTagOverride": false,
					"ProxyDestination": "",
					"Connect": {
						"Native": false,
						"Proxy": null
					},
					"CreateIndex": 20259,
					"ModifyIndex": 20328
				},
				"Checks": [
					{
						"Node": "vault-dev-0",
						"CheckID": "serfHealth",
						"Name": "Serf Health Status",
						"Status": "warning",
						"Notes": "",
						"Output": "Agent alive and reachable",
						"ServiceID": "",
						"ServiceName": "",
						"ServiceTags": [],
						"Definition": {},
						"CreateIndex": 20183,
						"ModifyIndex": 20183
					},
					{
						"Node": "vault-dev-0",
						"CheckID": "vault:vault-dev-0.consul.localnet:8200:vault-sealed-check",
						"Name": "Vault Sealed Status",
						"Status": "passing",
						"Notes": "Vault service is healthy when Vault is in an unsealed status and can become an active Vault server",
						"Output": "Vault Unsealed",
						"ServiceID": "vault:vault-dev-0.consul.localnet:8200",
						"ServiceName": "vault",
						"ServiceTags": [
							"active"
						],
						"Definition": {},
						"CreateIndex": 20260,
						"ModifyIndex": 20329
					}
				]
			}
		]`)
	})

	mux.HandleFunc("/v1/health/service/consul", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"Node": {
					"ID": "fbfe7e9b-5d30-284b-cc05-d2d5cc43688d",
					"Node": "consul-dev-0",
					"Address": "10.0.0.140",
					"Datacenter": "dev",
					"TaggedAddresses": {
						"lan": "10.0.0.140",
						"wan": "10.0.0.140"
					},
					"Meta": {
						"consul-network-segment": ""
					},
					"CreateIndex": 8,
					"ModifyIndex": 9
				},
				"Service": {
					"ID": "consul",
					"Service": "consul",
					"Tags": [],
					"Address": "",
					"Meta": null,
					"Port": 8300,
					"EnableTagOverride": false,
					"ProxyDestination": "",
					"Connect": {
						"Native": false,
						"Proxy": null
					},
					"CreateIndex": 8,
					"ModifyIndex": 8
				},
				"Checks": [
					{
						"Node": "consul-dev-0",
						"CheckID": "serfHealth",
						"Name": "Serf Health Status",
						"Status": "critical",
						"Notes": "",
						"Output": "Agent alive and reachable",
						"ServiceID": "",
						"ServiceName": "",
						"ServiceTags": [],
						"Definition": {},
						"CreateIndex": 8,
						"ModifyIndex": 8
					}
				]
			},
			{
				"Node": {
					"ID": "c7f88fba-f8d9-94a9-3627-523398acf7db",
					"Node": "consul-dev-1",
					"Address": "10.0.0.142",
					"Datacenter": "dev",
					"TaggedAddresses": {
						"lan": "10.0.0.142",
						"wan": "10.0.0.142"
					},
					"Meta": {
						"consul-network-segment": ""
					},
					"CreateIndex": 65,
					"ModifyIndex": 67
				},
				"Service": {
					"ID": "consul",
					"Service": "consul",
					"Tags": [],
					"Address": "",
					"Meta": null,
					"Port": 8300,
					"EnableTagOverride": false,
					"ProxyDestination": "",
					"Connect": {
						"Native": false,
						"Proxy": null
					},
					"CreateIndex": 65,
					"ModifyIndex": 65
				},
				"Checks": [
					{
						"Node": "consul-dev-1",
						"CheckID": "serfHealth",
						"Name": "Serf Health Status",
						"Status": "passing",
						"Notes": "",
						"Output": "Agent alive and reachable",
						"ServiceID": "",
						"ServiceName": "",
						"ServiceTags": [],
						"Definition": {},
						"CreateIndex": 65,
						"ModifyIndex": 65
					}
				]
			}
		]`)
	})
}
