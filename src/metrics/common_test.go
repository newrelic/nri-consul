package metrics

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/v3/data/metric"
	"github.com/newrelic/infra-integrations-sdk/v3/integration"
)

func TestSetMetric(t *testing.T) {

	i, err := integration.New("test", "1.0.0")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	entity, err := i.Entity("test", "agent")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	set := entity.NewMetricSet("ConsulTestSample")

	SetMetric(set, "test", float64(42), metric.GAUGE)

	value, ok := set.Metrics["test"]
	if !ok {
		t.Error("Metric was not added")
	} else if value != float64(42) {
		t.Error("Value was not set correctly for metric")
	}
}
