package webhook

import (
	"net/http"
)

// routes wires the routes to handlers on a specific router.
func (h handler) routes(router *http.ServeMux) error {
	templatepod, err := h.templatePod()
	if err != nil {
		return err
	}
	router.Handle("/wh/mutating/template", templatepod)

	return nil
}
