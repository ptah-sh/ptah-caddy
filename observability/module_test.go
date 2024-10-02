// Based on https://github.com/caddyserver/caddy/blob/master/modules/caddyhttp/metrics.go

package observability

import (
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func TestObserver(t *testing.T) {
	observer := &Observer{}

	observer.Provision(caddy.Context{})
}

func TestObserver_UnmarshalCaddyfile(t *testing.T) {
	observer := &Observer{}

	d := caddyfile.NewTestDispenser(`ptah_observer service_id "123" process_id "456" rule_id "789"`)

	err := observer.UnmarshalCaddyfile(d)
	if err != nil {
		t.Fatalf("Error unmarshalling caddyfile: %v", err)
	}

	if observer.ServiceID != "123" {
		t.Fatalf("Expected service_id to be '123', got '%s'", observer.ServiceID)
	}

	if observer.ProcessID != "456" {
		t.Fatalf("Expected process_id to be '456', got '%s'", observer.ProcessID)
	}

	if observer.RuleID != "789" {
		t.Fatalf("Expected rule_id to be '789', got '%s'", observer.RuleID)
	}
}
