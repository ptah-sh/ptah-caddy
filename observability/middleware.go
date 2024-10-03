package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

var _ caddyhttp.MiddlewareHandler = (*Observer)(nil)

func (m *Observer) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	startTime := time.Now()

	serverName := serverNameFromContext(r.Context())

	m.metrics.requestsInFlight.WithLabelValues(m.ServiceID, m.ProcessID, serverName, m.RuleID).Add(1)
	defer m.metrics.requestsInFlight.WithLabelValues(m.ServiceID, m.ProcessID, serverName, m.RuleID).Add(-1)

	recorder := newResponseRecorder(w)

	err := next.ServeHTTP(recorder, r)
	status := strconv.Itoa(recorder.Status())
	if err == nil {
		m.metrics.requestsCount.WithLabelValues(m.ServiceID, m.ProcessID, serverName, m.RuleID, status).Add(1)
	} else {
		// 500 status as the Caddy middleware failed itself
		m.metrics.requestsCount.WithLabelValues(m.ServiceID, m.ProcessID, serverName, m.RuleID, "500").Add(1)
	}

	if !recorder.firstByte.IsZero() {
		m.metrics.requestsTtfb.WithLabelValues(m.ServiceID, m.ProcessID, serverName, m.RuleID, status).Observe(time.Since(recorder.firstByte).Seconds())
	}

	m.metrics.requestsDuration.WithLabelValues(m.ServiceID, m.ProcessID, serverName, m.RuleID, status).Observe(time.Since(startTime).Seconds())

	return err
}

type ResponseRecorder struct {
	caddyhttp.ResponseRecorder

	firstByte time.Time
}

func newResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseRecorder: caddyhttp.NewResponseRecorder(w, nil, nil),
	}
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	if r.firstByte.IsZero() {
		r.firstByte = time.Now()
	}

	r.ResponseRecorder.WriteHeader(statusCode)
}
