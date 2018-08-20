package stub

import (
	"context"

	"github.com/fatih/structs"

	"github.com/davecgh/go-spew/spew"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/paulczar/oauth2-proxy/pkg/apis/oauth2proxy/v1alpha1"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

var (
	DefaultTagName = "json"
)

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	var err error
	switch o := event.Object.(type) {
	case *v1alpha1.Proxy:
		err = sdk.Create(newOauth2ProxyDeployment(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create oauth2proxy deployment : %v", err)
			return err
		}
	}
	return nil
}

// todo make this create a deployment + configmap
func newOauth2ProxyDeployment(cr *v1alpha1.Proxy) *appsv1.Deployment {
	// process config struct into map[string]string
	var args = []string{}
	var envs = []apiv1.EnvVar{}
	m := structs.New(cr.Spec.Config)
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
	// todo should be a random string if name not set ?
	name := "oauth2-proxy-" + cr.Name
	labels := map[string]string{
		"oauth2-proxy": name,
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
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
			Replicas: &cr.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "oauth2-proxy",
							Image:   "a5huynh/oauth2_proxy:2.2-debian",
							Command: []string{"/bin/oauth2_proxy"},
							Args:    args,
							Env:     envs,
						},
					},
				},
			},
		},
	}
	spew.Dump(deployment)
	return deployment
}
