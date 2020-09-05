package webhook

import (
	"context"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"

	"github.com/pelotech/k8s-templated-configuration/internal/log"
	whhttp "github.com/slok/kubewebhook/pkg/http"
	mutatingwh "github.com/slok/kubewebhook/pkg/webhook/mutating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// templatePod sets up the webhook handler for adding templating initContainers using Kubewebhook library.
func (h handler) templatePod() (http.Handler, error) {
	mt := mutatingwh.MutatorFunc(func(ctx context.Context, obj metav1.Object) (bool, error) {
		pod, ok := obj.(*corev1.Pod)
		if !ok {
			h.logger.Warningf("received object is not an v1.Pod")
			return false, nil
		}
		err := h.templater.Template(ctx, pod)
		if err != nil {
			return false, fmt.Errorf("could not add templating initContainer to the resource: %w", err)
		}

		return false, nil
	})

	logger := h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "template"})
	wh, err := mutatingwh.NewWebhook(mutatingwh.WebhookConfig{Name: "template"}, mt, nil, h.metrics, logger)
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}
	whHandler, err := whhttp.HandlerFor(wh)
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}
