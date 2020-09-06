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
		container corev1.Container
		obj       *corev1.Pod
		expObj    *corev1.Pod
	}{
		"Given a pod, with no annotations nothing should happen.": {
			container: corev1.Container{},
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

		"Given a pod, with secrets and volumes annotations initContainer should be created.": {
			container: corev1.Container{},
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Annotations: map[string]string{
						"template.pelo.tech/inject-secrets": "secret1",
						"template.pelo.tech/into-volumes":   "volume1",
					},
				},
			},
			expObj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Annotations: map[string]string{
						"template.pelo.tech/inject-secrets": "secret1",
						"template.pelo.tech/into-volumes":   "volume1",
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "templated-config",
							Image: "pelotech/envtemplate",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "volume1",
									MountPath: "/volume1",
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "secret1",
										},
									},
								},
							},
							Args: []string{"volume1"},
						},
					},
				},
			},
		},
		"Given a pod, with configMaps and volumes annotations initContainer should be created.": {
			container: corev1.Container{
				Image: "pelotech/image",
			},
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Annotations: map[string]string{
						"template.pelo.tech/inject-configmaps": "config1",
						"template.pelo.tech/into-volumes":      "volume1",
					},
				},
			},
			expObj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Annotations: map[string]string{
						"template.pelo.tech/inject-configmaps": "config1",
						"template.pelo.tech/into-volumes":      "volume1",
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "templated-config",
							Image: "pelotech/image",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "volume1",
									MountPath: "/volume1",
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "config1",
										},
									},
								},
							},
							Args: []string{"volume1"},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			m, err := template.NewTemplater(test.container)
			require.NoError(err)

			err1 := m.Template(context.TODO(), test.obj)
			require.NoError(err1)

			assert.Equal(test.expObj, test.obj)
		})
	}
}
