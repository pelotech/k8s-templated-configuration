package template_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pelotech/k8s-templated-configuration/internal/mutation/template"
)

func TestTemplaterTemplate(t *testing.T) {
	tests := map[string]struct {
		templates map[string]string
		obj       *corev1.Pod
		expObj    *corev1.Pod
	}{
		"Given a pod, with no annotations nothing should happen.": {
			templates: map[string]string{
				"test1": "value1",
				"test2": "value2",
			},
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
			expObj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
			},
		},

		// "Having a service, the labels should be mutated.": {
		// 	templates: map[string]string{
		// 		"test1": "value1",
		// 		"test2": "value2",
		// 	},
		// 	obj: &corev1.Service{
		// 		ObjectMeta: metav1.ObjectMeta{
		// 			Name: "test",
		// 		},
		// 	},
		// 	expObj: &corev1.Service{
		// 		ObjectMeta: metav1.ObjectMeta{
		// 			Name: "test",
		// 			Labels: map[string]string{
		// 				"test1": "value1",
		// 				"test2": "value2",
		// 			},
		// 		},
		// 	},
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			m := template.NewTemplater(test.templates)

			err := m.Template(context.TODO(), test.obj)
			require.NoError(err)

			assert.Equal(test.expObj, test.obj)
		})
	}
}
