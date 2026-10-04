package dashboard

import (
	"errors"
	"strings"
	"time"
)

type Signal struct {
	DeviceID string `json:"device_id"`
	Kind     string `json:"kind"`
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
}

type DeviceEvent struct {
	DeviceID  string    `json:"device_id"`
	Kind      string    `json:"kind"`
	Status    string    `json:"status"`
	Health    string    `json:"health"`
	Message   string    `json:"message,omitempty"`
	Observed  time.Time `json:"observed_at"`
	EventName string    `json:"event_name"`
}

func Classify(signal Signal, observed time.Time) (DeviceEvent, error) {
	signal.DeviceID = strings.TrimSpace(signal.DeviceID)
	if signal.DeviceID == "" {
		return DeviceEvent{}, errors.New("device_id is required")
	}

	allowed := map[string]map[string]bool{
		"build":      {"queued": true, "running": true, "succeeded": true, "failed": true},
		"release":    {"deploying": true, "released": true, "rolled_back": true},
		"diagnostic": {"ok": true, "warning": true, "error": true},
	}
	statuses, ok := allowed[signal.Kind]
	if !ok || !statuses[signal.Status] {
		return DeviceEvent{}, errors.New("kind and status do not form a known device signal")
	}

	health := "active"
	switch signal.Status {
	case "succeeded", "released", "ok":
		health = "healthy"
	case "failed", "rolled_back", "warning", "error":
		health = "attention"
	}

	return DeviceEvent{
		DeviceID:  signal.DeviceID,
		Kind:      signal.Kind,
		Status:    signal.Status,
		Health:    health,
		Message:   signal.Message,
		Observed:  observed.UTC(),
		EventName: "developer.device.status",
	}, nil
}
