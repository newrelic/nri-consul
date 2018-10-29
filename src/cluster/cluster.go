// Package cluster handles grabbing of cluster level metrics via the leader agent
package cluster

import (
	"errors"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-consul/src/agent"
)

// Cluster represents the cluster
// Wraps the leader agent and cluster entity
type Cluster struct {
	entity *integration.Entity
	leader *agent.Agent
}

// NewCluster creates a new cluster wrapped around the leader Agent
func NewCluster(leader *agent.Agent, i *integration.Integration) (*Cluster, error) {
	if leader == nil {
		return nil, errors.New("leader must not be nil")
	}

	clusterEntity, err := i.Entity("Cluster", "cluster")
	if err != nil {
		return nil, err
	}

	return &Cluster{
		entity: clusterEntity,
		leader: leader,
	}, nil
}

// CollectMetrics collects all cluster level metrics
func (c *Cluster) CollectMetrics() {
	metricSet := c.entity.NewMetricSet("ConsulClusterSample",
		metric.Attribute{Key: "displayName", Value: c.entity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: c.entity.Metadata.Namespace + ":" + c.entity.Metadata.Name},
		metric.Attribute{Key: "leaderNode", Value: c.leader.Name()},
	)

	// collect leader agent metrics
	if err := c.leader.CollectCoreMetrics(metricSet, nil, counterMetrics, timerMetrics); err != nil {
		log.Error("Error collecting leader metrics for Cluster: %s", err.Error())
	}
}
