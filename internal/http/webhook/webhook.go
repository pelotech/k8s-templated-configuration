package webhook

import (
	"fmt"
	"net/http"

	"github.com/pelotech/k8s-templated-configuration/internal/log"
	"github.com/pelotech/k8s-templated-configuration/internal/mutation/template"
)

// Config is the handler configuration.
type Config struct {
	MetricsRecorder MetricsRecorder
	Templater       template.Templater
	Logger          log.Logger
}

func (c *Config) defaults() error {
	if c.Templater == nil {
		return fmt.Errorf("templater is required")
	}

	if c.MetricsRecorder == nil {
		c.MetricsRecorder = dummyMetricsRecorder
	}

	if c.Logger == nil {
		c.Logger = log.Dummy
	}

	return nil
}

type handler struct {
	templater template.Templater
	handler   http.Handler
	metrics   MetricsRecorder
	logger    log.Logger
}

// New returns a new webhook handler.
func New(config Config) (http.Handler, error) {
	err := config.defaults()
	if err != nil {
		return nil, fmt.Errorf("handler configuration is not valid: %w", err)
	}

	mux := http.NewServeMux()

	h := handler{
		handler:   mux,
		templater: config.Templater,
		metrics:   config.MetricsRecorder,
		logger:    config.Logger.WithKV(log.KV{"service": "webhook-handler"}),
	}

	// Register all the routes with our router.
	err = h.routes(mux)
	if err != nil {
		return nil, fmt.Errorf("could not register routes on handler: %w", err)
	}

	// Register root handler middleware.
	h.handler = h.measuredHandler(h.handler) // Add metrics middleware.

	return h, nil
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}
