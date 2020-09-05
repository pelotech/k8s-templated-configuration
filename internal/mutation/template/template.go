package template

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	secretsAnnotation    = "template.pelo.tech/inject-secrets"
	configMapsAnnotation = "template.pelo.tech/inject-configmaps"
	volumesAnnotation    = "template.pelo.tech/into-volumes"
)

// Templater knows how to template Kubernetes pods.
type Templater interface {
	Template(ctx context.Context, pod *corev1.Pod) error
}

// NewTemplater returns a new templater that will template with labels.
func NewTemplater(config map[string]string) Templater {
	return templater{config: config}
}

// Templater knows how to template Kubernetes pods.
type templater struct {
	config map[string]string
}

// Template knows how to template Kubernetes pods.
func (t templater) Template(_ context.Context, pod *corev1.Pod) error {
	annotations := pod.GetAnnotations()
	if annotations == nil {
		return nil
	}
	secretsannotationvalue := annotations[secretsAnnotation]
	configmapsannotationvalue := annotations[configMapsAnnotation]
	volumesannotationvalue := annotations[volumesAnnotation]
	if "" == (secretsannotationvalue+configmapsannotationvalue) ||
		"" == volumesannotationvalue {
		return nil
	}
	secrets := strings.Split(secretsannotationvalue, ",")
	configmaps := strings.Split(configmapsannotationvalue, ",")
	volumes := strings.Split(volumesannotationvalue, ",")

	lenvolumemounts := len(volumes)
	lenenvfrom := len(secrets) + len(configmaps)

	volumemounts := make([]corev1.VolumeMount, lenvolumemounts)
	for _, volume := range volumes {
		volumemounts = append(volumemounts, corev1.VolumeMount{
			Name:      volume,
			MountPath: "/" + volume,
		})
	}

	envfrom := make([]corev1.EnvFromSource, lenenvfrom)
	for _, secret := range secrets {
		envfrom = append(envfrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: secret,
				},
			},
		})
	}
	for _, configmap := range configmaps {
		envfrom = append(envfrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configmap,
				},
			},
		})
	}

	initcontainer := corev1.Container{
		Name:         "templated-config",
		Image:        "pelotech/envsub", //TODO: parameterize
		VolumeMounts: volumemounts,
		EnvFrom:      envfrom,
		Args:         volumes,
	}

	pod.Spec.InitContainers = append(pod.Spec.InitContainers, initcontainer)

	return nil
}
