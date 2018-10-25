package agent

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
)

// MetricDefinition represents a all the definition to collect
// a metric from the API and send to Infrastructure
type MetricDefinition struct {
	APIKey     string
	MetricName string
	SourceType metric.SourceType
}

var gaugeMetrics = []*MetricDefinition{
	{
		APIKey:     "consul.runtime.num_goroutines",
		MetricName: "runtime.goroutines",
		SourceType: metric.GAUGE,
	},
	{
		APIKey:     "consul.runtime.alloc_bytes",
		MetricName: "runtime.allocationsInBytes",
		SourceType: metric.GAUGE,
	},
	{
		APIKey:     "consul.runtime.heap_objects",
		MetricName: "runtime.heapObjects",
		SourceType: metric.GAUGE,
	},
	{
		APIKey:     "consul.runtime.sys_bytes",
		MetricName: "runtime.virtualAddressSpaceInBytes",
		SourceType: metric.GAUGE,
	},
	{
		APIKey:     "consul.runtime.malloc_count",
		MetricName: "runtime.allocations",
		SourceType: metric.GAUGE,
	},
	{
		APIKey:     "consul.runtime.free_count",
		MetricName: "runtime.frees",
		SourceType: metric.GAUGE,
	},
	{
		APIKey:     "consul.runtime.total_gc_pause_ns",
		MetricName: "runtime.gcPauseInMilliseconds",
		SourceType: metric.GAUGE,
	},
	{
		APIKey:     "consul.runtime.total_gc_runs",
		MetricName: "runtime.gcCycles",
		SourceType: metric.GAUGE,
	},
}

var counterMetrics = []*MetricDefinition{
	{
		APIKey:     "consul.client.rpc",
		MetricName: "client.rpcLoad",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.client.rpc.exceeded",
		MetricName: "client.rpcRateLimited",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.client.rpc.failed",
		MetricName: "client.rpcFailed",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.acl.cache_hit",
		MetricName: "agent.aclCacheHit",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.acl.cache_miss",
		MetricName: "agent.aclCacheMiss",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.dns.stale_queries",
		MetricName: "agent.staleQueries",
		SourceType: metric.RATE,
	},
}

// TODO timer metrics
