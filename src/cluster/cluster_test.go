package cluster

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/nri-consul/src/agent"
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
