package observability

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

var _ caddyfile.Unmarshaler = (*Observer)(nil)

func (m *Observer) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume directive name

	for d.Next() {
		switch d.Val() {
		case "service_id":
			if !d.NextArg() {
				return d.ArgErr()
			}
			serviceID := d.Val()

			m.ServiceID = serviceID
		case "process_id":
			if !d.NextArg() {
				return d.ArgErr()
			}

			m.ProcessID = d.Val()
		case "rule_id":
			if !d.NextArg() {
				return d.ArgErr()
			}

			ruleID := d.Val()

			m.RuleID = ruleID
		default:
			return d.ArgErr()
		}
	}

	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Observer
	err := m.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
