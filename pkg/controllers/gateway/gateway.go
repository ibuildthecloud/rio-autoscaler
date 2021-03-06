package gateway

import (
	"context"
	"fmt"
	"sync"

	"github.com/rancher/rio-autoscaler/types"
	"github.com/rancher/wrangler/pkg/kv"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

var (
	EndpointChanMap = sync.Map{}
)

func Register(ctx context.Context, rContext *types.Context) error {
	logrus.Info("Starting gateway endpoint controller")
	e := endpointHandler{}
	rContext.Core.Core().V1().Endpoints().OnChange(ctx, "gateway-endpoint-watcher", e.Sync)

	return nil
}

type endpointHandler struct{}

func (e endpointHandler) Sync(key string, obj *corev1.Endpoints) (*corev1.Endpoints, error) {
	if obj == nil {
		namespace, name := kv.Split(key, "/")
		EndpointChanMap.Delete(fmt.Sprintf("%s.%s", name, namespace))
		return nil, nil
	}

	// todo: add a filter only for scale-to-zero services so that we don't have to keep a channel for every endpoint
	if obj != nil && obj.DeletionTimestamp == nil {
		ch := make(chan struct{}, 0)
		EndpointChanMap.LoadOrStore(fmt.Sprintf("%s.%s", obj.Name, obj.Namespace), ch)
		if isEndpointReady(obj) {
			o, ok := EndpointChanMap.Load(fmt.Sprintf("%s.%s", obj.Name, obj.Namespace))
			if ok {
				c := o.(chan struct{})
				close(c)
				EndpointChanMap.Delete(fmt.Sprintf("%s.%s", obj.Name, obj.Namespace))
			}
		}
	}
	return obj, nil
}

func isEndpointReady(obj *corev1.Endpoints) bool {
	ready := true
	if len(obj.Subsets) == 0 {
		ready = false
	}
	for _, subnet := range obj.Subsets {
		if len(subnet.NotReadyAddresses) > 0 {
			ready = false
		}
	}
	return ready
}
