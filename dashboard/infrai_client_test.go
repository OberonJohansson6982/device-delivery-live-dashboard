package dashboard

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestClientDecodesBusinessErrorBeforeStatus(t *testing.T) {
	client := NewInfraiClient("https://example.test", "test-key", &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 400, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"ok":false,"data":null,"error":{"code":"BAD_SIGNAL","message":"signal rejected"},"metadata":null}`))}, nil
	})})

	err := client.Publish(context.Background(), "device-delivery", DeviceEvent{DeviceID: "lab-1", EventName: "developer.device.status"}, "event-1")
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != 400 || apiErr.Code != "BAD_SIGNAL" {
		t.Fatalf("error = %#v", apiErr)
	}
}

func TestDeleteChannel(t *testing.T) {
	client := NewInfraiClient("https://example.test", "test-key", &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodDelete {
			t.Fatalf("method = %s, want DELETE", req.Method)
		}
		if req.URL.Path != "/v1/realtime/channel/delete/device-delivery" {
			t.Fatalf("path = %s", req.URL.Path)
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"ok":true,"data":{},"error":null,"metadata":null}`))}, nil
	})})

	if err := client.DeleteChannel(context.Background(), "device-delivery"); err != nil {
		t.Fatalf("DeleteChannel() error = %v", err)
	}
}
