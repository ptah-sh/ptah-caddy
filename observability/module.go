package observability

import (
	"context"
	"fmt"

	caddy "github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type httpMetrics struct {
	requestsInFlight *prometheus.GaugeVec
	requestsCount    *prometheus.CounterVec
	requestsTtfb     *prometheus.HistogramVec
	requestsDuration *prometheus.HistogramVec
}

func init() {
	caddy.RegisterModule(Observer{})

	// httpcaddyfile.RegisterGlobalOption("ptah_node_id", parseCaddyfile)

	httpcaddyfile.RegisterHandlerDirective("ptah_observer", parseCaddyfile)
}

type Observer struct {
	ServiceID string `json:"service_id"`
	ProcessID string `json:"process_id"`
	RuleID    string `json:"rule_id"`

	metrics httpMetrics
}

func (Observer) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.ptah_observer",
		New: func() caddy.Module {
			return new(Observer)
		},
	}
}

func (m *Observer) Provision(ctx caddy.Context) error {
	if m.ServiceID == "" {
		return fmt.Errorf("service_id is required")
	}

	if m.ProcessID == "" {
		return fmt.Errorf("process_id is required")
	}

	if m.RuleID == "" {
		return fmt.Errorf("rule_id is required")
	}

	namespace := "ptah"
	subsystem := "caddy_http"

	labels := []string{"server_name", "service_id", "process_id", "rule_id"}

	m.metrics.requestsInFlight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "requests_in_flight",
		Help:        "Counter of HTTP(S) requests in flight.",
		ConstLabels: prometheus.Labels{},
	}, labels)

	m.metrics.requestsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "requests_count",
		Help:        "Counter of HTTP(S) requests made.",
		ConstLabels: prometheus.Labels{},
	}, append(labels, "status_code"))

	m.metrics.requestsTtfb = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "requests_ttfb",
		Help:        "Histogram of HTTP(S) requests time to first byte.",
		ConstLabels: prometheus.Labels{},
	}, labels)

	m.metrics.requestsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   namespace,
		Subsystem:   subsystem,
		Name:        "requests_duration",
		Help:        "Histogram of HTTP(S) requests duration.",
		ConstLabels: prometheus.Labels{},
	}, labels)

	return nil
}

var _ caddy.Provisioner = (*Observer)(nil)

// serverNameFromContext extracts the current server name from the context.
// Returns "UNKNOWN" if none is available (should probably never happen).
func serverNameFromContext(ctx context.Context) string {
	srv, ok := ctx.Value(caddyhttp.ServerCtxKey).(*caddyhttp.Server)
	if !ok || srv == nil || srv.Name() == "" {
		return "UNKNOWN"
	}
	return srv.Name()
}
