/*
Copyright 2019 Rancher Labs.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"

	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes/typed/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type ServiceAccountHandler func(string, *v1.ServiceAccount) (*v1.ServiceAccount, error)

type ServiceAccountController interface {
	ServiceAccountClient

	OnChange(ctx context.Context, name string, sync ServiceAccountHandler)
	OnRemove(ctx context.Context, name string, sync ServiceAccountHandler)
	Enqueue(namespace, name string)

	Cache() ServiceAccountCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type ServiceAccountClient interface {
	Create(*v1.ServiceAccount) (*v1.ServiceAccount, error)
	Update(*v1.ServiceAccount) (*v1.ServiceAccount, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.ServiceAccount, error)
	List(namespace string, opts metav1.ListOptions) (*v1.ServiceAccountList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ServiceAccount, err error)
}

type ServiceAccountCache interface {
	Get(namespace, name string) (*v1.ServiceAccount, error)
	List(namespace string, selector labels.Selector) ([]*v1.ServiceAccount, error)

	AddIndexer(indexName string, indexer ServiceAccountIndexer)
	GetByIndex(indexName, key string) ([]*v1.ServiceAccount, error)
}

type ServiceAccountIndexer func(obj *v1.ServiceAccount) ([]string, error)

type serviceAccountController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.ServiceAccountsGetter
	informer          informers.ServiceAccountInformer
	gvk               schema.GroupVersionKind
}

func NewServiceAccountController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.ServiceAccountsGetter, informer informers.ServiceAccountInformer) ServiceAccountController {
	return &serviceAccountController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromServiceAccountHandlerToHandler(sync ServiceAccountHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.ServiceAccount
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.ServiceAccount))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *serviceAccountController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.ServiceAccount))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateServiceAccountOnChange(updater generic.Updater, handler ServiceAccountHandler) ServiceAccountHandler {
	return func(key string, obj *v1.ServiceAccount) (*v1.ServiceAccount, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1.ServiceAccount)
			}
		}

		return copyObj, err
	}
}

func (c *serviceAccountController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *serviceAccountController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *serviceAccountController) OnChange(ctx context.Context, name string, sync ServiceAccountHandler) {
	c.AddGenericHandler(ctx, name, FromServiceAccountHandlerToHandler(sync))
}

func (c *serviceAccountController) OnRemove(ctx context.Context, name string, sync ServiceAccountHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromServiceAccountHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *serviceAccountController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *serviceAccountController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *serviceAccountController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *serviceAccountController) Cache() ServiceAccountCache {
	return &serviceAccountCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *serviceAccountController) Create(obj *v1.ServiceAccount) (*v1.ServiceAccount, error) {
	return c.clientGetter.ServiceAccounts(obj.Namespace).Create(obj)
}

func (c *serviceAccountController) Update(obj *v1.ServiceAccount) (*v1.ServiceAccount, error) {
	return c.clientGetter.ServiceAccounts(obj.Namespace).Update(obj)
}

func (c *serviceAccountController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.ServiceAccounts(namespace).Delete(name, options)
}

func (c *serviceAccountController) Get(namespace, name string, options metav1.GetOptions) (*v1.ServiceAccount, error) {
	return c.clientGetter.ServiceAccounts(namespace).Get(name, options)
}

func (c *serviceAccountController) List(namespace string, opts metav1.ListOptions) (*v1.ServiceAccountList, error) {
	return c.clientGetter.ServiceAccounts(namespace).List(opts)
}

func (c *serviceAccountController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.ServiceAccounts(namespace).Watch(opts)
}

func (c *serviceAccountController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.ServiceAccount, err error) {
	return c.clientGetter.ServiceAccounts(namespace).Patch(name, pt, data, subresources...)
}

type serviceAccountCache struct {
	lister  listers.ServiceAccountLister
	indexer cache.Indexer
}

func (c *serviceAccountCache) Get(namespace, name string) (*v1.ServiceAccount, error) {
	return c.lister.ServiceAccounts(namespace).Get(name)
}

func (c *serviceAccountCache) List(namespace string, selector labels.Selector) ([]*v1.ServiceAccount, error) {
	return c.lister.ServiceAccounts(namespace).List(selector)
}

func (c *serviceAccountCache) AddIndexer(indexName string, indexer ServiceAccountIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.ServiceAccount))
		},
	}))
}

func (c *serviceAccountCache) GetByIndex(indexName, key string) (result []*v1.ServiceAccount, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.ServiceAccount))
	}
	return result, nil
}
