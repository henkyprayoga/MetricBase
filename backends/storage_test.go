package backends

import (
	"github.com/msiebuhr/MetricBase"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func generateTestStoreAndGet(backend MetricBase.Backend, t *testing.T) {
	// Start backend
	backend.Start()
	defer backend.Stop()

	// Load some data and read it back out
	addChan := make(chan MetricBase.Metric, 10)
	backend.AddMetrics(addChan)
	addChan <- *MetricBase.NewMetric("foo.bar", 3.14, 100)
	close(addChan)

	time.Sleep(time.Millisecond)

	// Read back list of metrics
	metricNames := GetMetricsAsList(backend)
	if len(metricNames) != 1 {
		t.Errorf("Expected to get one metric back, got %d", len(metricNames))
	} else if metricNames[0] != "foo.bar" {
		t.Errorf("Expected the metric name to be 'foo.bar', got '%v'", metricNames[0])
	}

	// Read back the data
	data := GetDataAsList(backend, "foo.bar", 0, 0)
	if len(data) != 1 {
		t.Fatalf("Expected to get one result, got %d", len(data))
	}

	if data[0].Value != 3.14 {
		t.Errorf("Expected data[0].Value=3.14, got '%f'.", data[0].Value)
	}
	if data[0].Time != 100 {
		t.Errorf("Expected data[0].Time=100, got '%d'.", data[0].Time)
	}
}

func TestLevelDbStoreAndGet(t *testing.T) {
	// Create tempdir (& remove afterwards)
	dir, err := ioutil.TempDir(".", "tmp_storage_test")
	if err != nil {
		t.Fatalf("Could not create tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	generateTestStoreAndGet(NewLevelDb(dir), t)
}

func TestMemoryStoreAndGet(t *testing.T) {
	// Create tempdir (& remove afterwards)
	generateTestStoreAndGet(NewMemoryBackend(), t)
}