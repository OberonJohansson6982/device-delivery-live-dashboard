package dashboard

import (
	"testing"
	"time"
)

func TestClassifyDeliverySignals(t *testing.T) {
	observed := time.Date(2026, 9, 24, 10, 30, 0, 0, time.FixedZone("CST", 8*60*60))
	tests := []struct {
		name   string
		input  Signal
		health string
	}{
		{"running build stays active", Signal{DeviceID: "ios-lab-07", Kind: "build", Status: "running"}, "active"},
		{"completed release is healthy", Signal{DeviceID: "android-lab-02", Kind: "release", Status: "released"}, "healthy"},
		{"diagnostic warning needs attention", Signal{DeviceID: "linux-rig-11", Kind: "diagnostic", Status: "warning"}, "attention"},
		{"failed build needs attention", Signal{DeviceID: "ios-lab-07", Kind: "build", Status: "failed"}, "attention"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event, err := Classify(test.input, observed)
			if err != nil {
				t.Fatal(err)
			}
			if event.Health != test.health {
				t.Fatalf("health = %q, want %q", event.Health, test.health)
			}
			if event.Observed.Location() != time.UTC {
				t.Fatalf("timestamp location = %v, want UTC", event.Observed.Location())
			}
		})
	}
}
