package datacenter

import (
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/nri-consul/src/metrics"
)

var counterMetrics = []*metrics.MetricDefinition{
	{
		APIKey:     "consul.memberlist.msg.suspect",
		MetricName: "cluster.suspects",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.serf.member.flap",
		MetricName: "cluster.flaps",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.raft.state.leader",
		MetricName: "raft.completedLeaderElections",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.raft.state.candidate",
		MetricName: "raft.initiatedLeaderElections",
		SourceType: metric.RATE,
	},
	{
		APIKey:     "consul.raft.apply",
		MetricName: "raft.txns",
		SourceType: metric.RATE,
	},
}

var timerMetrics = []*metrics.TimerDefinition{
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.commitTime",
			MetricName: "raft.commitTimeAvgInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Average,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.commitTime",
			MetricName: "raft.commitTimes",
			SourceType: metric.RATE,
		},
		Operation: metrics.Count,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.commitTime",
			MetricName: "raft.commitTimeMedianInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Median,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.commitTime",
			MetricName: "raft.commitTimeMaxInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Max,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.leader.dispatchLog",
			MetricName: "raft.logDispatchAvgInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Average,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.leader.dispatchLog",
			MetricName: "raft.logDispatches",
			SourceType: metric.RATE,
		},
		Operation: metrics.Count,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.leader.dispatchLog",
			MetricName: "raft.logDispatchMedianInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Median,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.leader.dispatchLog",
			MetricName: "raft.logDispatchMaxInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Max,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.leader.lastContact",
			MetricName: "raft.lastContactAvgInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Average,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.leader.lastContact",
			MetricName: "raft.lastContacts",
			SourceType: metric.RATE,
		},
		Operation: metrics.Count,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.leader.lastContact",
			MetricName: "raft.lastContactMedianInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Median,
	},
	{
		MetricDefinition: metrics.MetricDefinition{
			APIKey:     "consul.raft.leader.lastContact",
			MetricName: "raft.lastContactMaxInMilliseconds",
			SourceType: metric.GAUGE,
		},
		Operation: metrics.Max,
	},
}
