package oauth2proxy

import (
	"github.com/fatih/structs"

	"github.com/paulczar/oauth2-proxy/pkg/apis/oauth2proxy/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	image = "a5huynh/oauth2_proxy:2.2-debian"
)

// Deployment creates a oauth2 proxy deployment
func Deployment(cr *v1alpha1.Proxy) *appsv1.Deployment {
	labels := map[string]string{
		"oauth2-proxy": name(cr),
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name(cr),
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Proxy",
				}),
			},
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cr.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: podSpec(cr),
			},
		},
	}
	//spew.Dump(deployment)
	return deployment
}

func name(cr *v1alpha1.Proxy) string {
	return "oauth2-proxy-" + cr.Name
}

func labels(cr *v1alpha1.Proxy) map[string]string {
	return map[string]string{
		"oauth2-proxy": name(cr),
	}
}

func deploymentSpec(cr *v1alpha1.Proxy) appsv1.DeploymentSpec {
	spec := appsv1.DeploymentSpec{
		Replicas: &cr.Spec.Size,
		Selector: &metav1.LabelSelector{
			MatchLabels: labels(cr),
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels(cr),
			},
			Spec: podSpec(cr),
		},
	}
	return spec
}

// createPod creates a oauth2 proxy deployment
func podSpec(cr *v1alpha1.Proxy) apiv1.PodSpec {
	args, envs := podArgsAndEnvs(&cr.Spec.Config)
	if cr.Spec.Image == "" {
		cr.Spec.Image = image
	}
	pod := &apiv1.PodSpec{
		Containers: []apiv1.Container{
			{
				Name:    "oauth2-proxy",
				Image:   cr.Spec.Image,
				Command: []string{"/bin/oauth2_proxy"},
				Args:    args,
				Env:     envs,
			},
		},
	}

	return *pod
}

func podArgsAndEnvs(config *v1alpha1.ProxyConfig) ([]string, []apiv1.EnvVar) {
	m := structs.New(config)
	var args = []string{}
	var envs = []apiv1.EnvVar{}

	for _, f := range m.Fields() {
		if f.Value().(string) != "" {
			if f.Tag("cli") != "" {
				t := f.Tag("cli") + "=" + f.Value().(string)
				args = append(args, t)
			} else if f.Tag("env") != "" {
				e := apiv1.EnvVar{
					Name:  f.Tag("env"),
					Value: f.Value().(string),
				}
				envs = append(envs, e)
			}
		}
	}
	return args, envs
}
