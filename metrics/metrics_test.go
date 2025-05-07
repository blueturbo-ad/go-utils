package metrics

import (
	"testing"
)

func TestPrometheus(t *testing.T) {

	t.Run("TestPrometheus", func(t *testing.T) {
		err := GetInstance().Init("eventserver", nil)
		if err != nil {
			t.Fatalf("Failed to initialize metrics: %v", err)
		}
		GetInstance().SetGauge("test_gauge", 1.0)

	})

}
