package stub

import (
	"context"
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/paulczar/oauth2-proxy/pkg/apis/oauth2proxy/v1alpha1"
	"github.com/paulczar/oauth2-proxy/pkg/oauth2proxy"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
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
		o2p := o
		spec := o2p.Spec
		// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
		// All secondary resources must have the CR set as their OwnerReference for this to be the case
		if event.Deleted {
			logrus.Debugf("deleting %s", o2p.Name)
		}
		// construct deployment from request
		dep := oauth2proxy.Deployment(o2p)

		// attempt to create deployment
		err = sdk.Create(dep)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create oauth2proxy deployment : %v", err)
			return err
		} else if errors.IsAlreadyExists(err) {
			logrus.Debugf("%s already exists!", o2p.Name)
		}

		// Grab the deployment that either already existed, or was just created.
		err = sdk.Get(dep)
		if err != nil {
			return fmt.Errorf("failed to get deployment: %v", err)
		}

		// attempt to resize deployment to match request.
		rep := spec.Size
		logrus.Debugf("updating replica count %v->%v for %s", *dep.Spec.Replicas, rep, o2p.Name)
		if *dep.Spec.Replicas != rep {
			//logrus.Infof("updating replica count %v->%v for %s", *dep.Spec.Replicas, o.Spec.Replicas, o.Name)
			dep.Spec.Replicas = &rep
			err = sdk.Update(dep)
			if err != nil {
				return fmt.Errorf("failed to update deployment: %v", err)
			}
		}

	}
	return nil
}
