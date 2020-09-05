package template

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Templater knows how to template Kubernetes resources.
type Templater interface {
	Template(ctx context.Context, obj metav1.Object) error
}

// NewLabelTemplater returns a new templater that will template with labels.
func NewLabelTemplater(templates map[string]string) Templater {
	return labeltemplater{templates: templates}
}

type labeltemplater struct {
	templates map[string]string
}

func (l labeltemplater) Template(_ context.Context, obj metav1.Object) error {
	labels := obj.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	for k, v := range l.templates {
		labels[k] = v
	}

	obj.SetLabels(labels)
	return nil
}

// DummyTemplater is a templater that doesn't do anything.
var DummyTemplater Templater = dummyMaker(0)

type dummyMaker int

func (dummyMaker) Template(_ context.Context, _ metav1.Object) error { return nil }
