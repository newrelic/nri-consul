package agent

import (
	"math"
	"sort"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/serf/coordinate"
	"github.com/newrelic/infra-integrations-sdk/v3/data/metric"
	"github.com/newrelic/nri-consul/src/metrics"
)

func findNode(nodeName string, nodes []*api.CoordinateEntry) *api.CoordinateEntry {
	for _, node := range nodes {
		if node.Node == nodeName {
			return node
		}
	}
	return nil
}

func calculateLatencyMetrics(metricSet *metric.Set, node *api.CoordinateEntry, nodes []*api.CoordinateEntry) {
	latencies := make([]float64, 0)
	for _, other := range nodes {
		if other.Node == node.Node {
			continue
		}

		latencies = append(latencies, calcLatencyDist(node.Coord, other.Coord))
	}

	// Sort latencies
	sort.Float64s(latencies)

	// Set metrics
	metrics.SetMetric(metricSet, "net.agent.medianLatencyInMilliseconds", calcLatencyMedian(latencies), metric.GAUGE)
	metrics.SetMetric(metricSet, "net.agent.minLatencyInMilliseconds", latencies[0], metric.GAUGE)
	metrics.SetMetric(metricSet, "net.agent.maxLatencyInMilliseconds", latencies[len(latencies)-1], metric.GAUGE)
	metrics.SetMetric(metricSet, "net.agent.p25LatencyInMilliseconds", calcLatencyPercentile(latencies, 0.25), metric.GAUGE)
	metrics.SetMetric(metricSet, "net.agent.p75LatencyInMilliseconds", calcLatencyPercentile(latencies, 0.75), metric.GAUGE)
	metrics.SetMetric(metricSet, "net.agent.p90LatencyInMilliseconds", calcLatencyPercentile(latencies, 0.90), metric.GAUGE)
	metrics.SetMetric(metricSet, "net.agent.p95LatencyInMilliseconds", calcLatencyPercentile(latencies, 0.95), metric.GAUGE)
	metrics.SetMetric(metricSet, "net.agent.p99LatencyInMilliseconds", calcLatencyPercentile(latencies, 0.99), metric.GAUGE)
}

// calcLatencyDist calculates distance between two coordinates.
// Taken from Consul docs https://www.consul.io/docs/internals/coordinates.html
// In order to compute the latency between two nodes you must calculate
// the distance based on their coordinates.
func calcLatencyDist(a, b *coordinate.Coordinate) float64 {
	// Calculate the Euclidean distance plus the heights.
	sumsq := 0.0
	for i := 0; i < len(a.Vec); i++ {
		diff := a.Vec[i] - b.Vec[i]
		sumsq += diff * diff
	}
	rtt := math.Sqrt(sumsq) + a.Height + b.Height

	// Apply the adjustment components, guarding against negatives.
	adjusted := rtt + a.Adjustment + b.Adjustment
	if adjusted > 0.0 {
		rtt = adjusted
	}

	return rtt * 1000.0
}

// calcLatencyMedian is the median of a data set of latencies
func calcLatencyMedian(latencies []float64) float64 {
	numLatencies := len(latencies)
	halfIndex := numLatencies / 2

	if numLatencies%2 == 0 {
		return (latencies[halfIndex-1] + latencies[halfIndex]) / 2
	}

	return latencies[halfIndex]
}

func calcLatencyPercentile(latencies []float64, percent float64) float64 {
	numLatencies := float64(len(latencies))
	index := int(math.Ceil(numLatencies * percent))
	return latencies[index-1]
}
