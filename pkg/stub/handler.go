package stub

import (
	"context"

	"github.com/paulczar/oauth2-proxy/pkg/apis/oauth2proxy/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
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

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Proxy:
		err := sdk.Create(newOauth2Proxy(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create busybox pod : %v", err)
			return err
		}
	}
	return nil
}

// todo make this a config map that writes a config file
var args = []string{
	"--cookie-secure=false",
	"--upstream='http://upstream:80'",
	"--http-address='0.0.0.0:4180'",
	"--redirect-url='http://example.com/oauth2/callback'",
	"--email-domain='example.com'",
	"--cookie-secret='cookie-secret'",
	"--client-id='client-id'",
	"--client-secret='client-secret'",
}

// todo make this create a deployment + configmap
func newOauth2Proxy(cr *v1alpha1.Proxy) *v1.Pod {
	labels := map[string]string{
		"app": "oauth2-proxy",
	}
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "oauth2-proxy",
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
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "oauth2-proxy",
					Image:   "a5huynh/oauth2_proxy:2.2-debian",
					Command: []string{"/bin/oauth2_proxy"},
					Args:    args,
				},
			},
		},
	}
}
