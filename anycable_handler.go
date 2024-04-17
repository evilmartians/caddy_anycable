package caddy_anycable

import (
	"fmt"
	"github.com/anycable/anycable-go/cli"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"log/slog"
	"net/http"
)

func init() {
	caddy.RegisterModule(AnyCableHandler{})
	httpcaddyfile.RegisterHandlerDirective("anycable", parseCaddyfile)
}

type AnyCableHandler struct {
	logger   *slog.Logger
	options  []string
	anycable *cli.Embedded
	handler  http.Handler
}

func (AnyCableHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.anycable",
		New: func() caddy.Module { return new(AnyCableHandler) },
	}
}

func (h AnyCableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	handler := h.handler

	if r.URL.Path == "/cable" {
		handler.ServeHTTP(w, r)

		return nil
	}

	return next.ServeHTTP(w, r)
}

func (h AnyCableHandler) runAnyCable() (*cli.Embedded, error) {
	argsWithProg := append([]string{"anycable-go"}, h.options...)
	c, err, _ := cli.NewConfigFromCLI(argsWithProg)
	if err != nil {
		return nil, err
	}

	opts := []cli.Option{
		cli.WithName("AnyCable"),
		cli.WithDefaultRPCController(),
		cli.WithDefaultBroker(),
		cli.WithDefaultSubscriber(),
		cli.WithDefaultBroadcaster(),
		cli.WithLogger(h.logger),
	}

	runner, err := cli.NewRunner(c, opts)

	if err != nil {
		return nil, err
	}

	anycable, err := runner.Embed()

	return anycable, err
}

func (h *AnyCableHandler) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			key := d.Val()
			if d.NextArg() {
				h.options = append(h.options, fmt.Sprintf("--%s=%s", key, d.Val()))
			} else {
				return d.Errf("expected 1 argument for %s but none provided", key)
			}
			if d.NextArg() {
				return d.Errf("expected only 1 argument for %s", key)
			}
		}
	}

	return nil
}

func (h *AnyCableHandler) Provision(ctx caddy.Context) error {
	anycable, err := h.runAnyCable()

	if err != nil {
		return err
	}

	h.handler, _ = anycable.WebSocketHandler()

	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var anyCable AnyCableHandler
	err := anyCable.UnmarshalCaddyfile(h.Dispenser)
	return anyCable, err
}

var (
	_ caddyhttp.MiddlewareHandler = (*AnyCableHandler)(nil)
	_ caddy.Provisioner           = (*AnyCableHandler)(nil)
	_ caddyfile.Unmarshaler       = (*AnyCableHandler)(nil)
)