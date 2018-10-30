// Package cluster handles grabbing of cluster level metrics via the leader agent
package cluster

import (
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-consul/src/agent"
	"github.com/newrelic/nri-consul/src/metrics"
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
		metric.Attribute{Key: "leader", Value: c.leader.Name()},
	)

	// collect leader agent metrics
	if err := c.leader.CollectCoreMetrics(metricSet, nil, counterMetrics, timerMetrics); err != nil {
		log.Error("Error collecting leader metrics for Cluster: %s", err.Error())
	}

	// collect node count
	if err := c.setNodeCountMetric(metricSet); err != nil {
		log.Error("Error collecting node count: %s", err.Error())
	}

	// collect node health counts
	if err := c.collectStatusCounts(metricSet); err != nil {
		log.Error("Error getting node health counts: %s", err.Error())
	}
}

func (c *Cluster) setNodeCountMetric(metricSet *metric.Set) error {
	nodes, _, err := c.leader.Client.Catalog().Nodes(nil)
	if err != nil {
		return err
	}

	metrics.SetMetric(metricSet, "catalog.registeredNodes", len(nodes), metric.GAUGE)
	return nil
}

// collectStatusCounts aggregates health status across services on a node.
func (c *Cluster) collectStatusCounts(metricSet *metric.Set) error {
	services, _, err := c.leader.Client.Catalog().Services(nil)
	if err != nil {
		return err
	}

	// keeps track of counts
	nodeCounts := map[string]int{
		"critical": 0,
		"up":       0,
		"warning":  0,
		"passing":  0,
	}

	// for each service look at the nodes that host it and count health
	for service := range services {
		entries, _, err := c.leader.Client.Health().Service(service, "", false, nil)
		if err != nil {
			log.Error("Error getting nodes for service: %s", service)
			return err
		}

		for _, entry := range entries {
			switch entry.Checks.AggregatedStatus() {
			case api.HealthCritical:
				nodeCounts["critical"]++
			case api.HealthWarning:
				nodeCounts["warning"]++
			case api.HealthPassing:
				nodeCounts["up"]++
				nodeCounts["passing"]++
			}
		}
	}

	for status, count := range nodeCounts {
		metrics.SetMetric(metricSet, fmt.Sprintf("catalog.%sNodes", status), count, metric.GAUGE)
	}

	return nil
}
