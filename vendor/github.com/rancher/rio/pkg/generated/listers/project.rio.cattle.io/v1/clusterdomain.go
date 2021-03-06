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
	v1 "github.com/rancher/rio/pkg/apis/project.rio.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ClusterDomainLister helps list ClusterDomains.
type ClusterDomainLister interface {
	// List lists all ClusterDomains in the indexer.
	List(selector labels.Selector) (ret []*v1.ClusterDomain, err error)
	// ClusterDomains returns an object that can list and get ClusterDomains.
	ClusterDomains(namespace string) ClusterDomainNamespaceLister
	ClusterDomainListerExpansion
}

// clusterDomainLister implements the ClusterDomainLister interface.
type clusterDomainLister struct {
	indexer cache.Indexer
}

// NewClusterDomainLister returns a new ClusterDomainLister.
func NewClusterDomainLister(indexer cache.Indexer) ClusterDomainLister {
	return &clusterDomainLister{indexer: indexer}
}

// List lists all ClusterDomains in the indexer.
func (s *clusterDomainLister) List(selector labels.Selector) (ret []*v1.ClusterDomain, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ClusterDomain))
	})
	return ret, err
}

// ClusterDomains returns an object that can list and get ClusterDomains.
func (s *clusterDomainLister) ClusterDomains(namespace string) ClusterDomainNamespaceLister {
	return clusterDomainNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ClusterDomainNamespaceLister helps list and get ClusterDomains.
type ClusterDomainNamespaceLister interface {
	// List lists all ClusterDomains in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.ClusterDomain, err error)
	// Get retrieves the ClusterDomain from the indexer for a given namespace and name.
	Get(name string) (*v1.ClusterDomain, error)
	ClusterDomainNamespaceListerExpansion
}

// clusterDomainNamespaceLister implements the ClusterDomainNamespaceLister
// interface.
type clusterDomainNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ClusterDomains in the indexer for a given namespace.
func (s clusterDomainNamespaceLister) List(selector labels.Selector) (ret []*v1.ClusterDomain, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ClusterDomain))
	})
	return ret, err
}

// Get retrieves the ClusterDomain from the indexer for a given namespace and name.
func (s clusterDomainNamespaceLister) Get(name string) (*v1.ClusterDomain, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("clusterdomain"), name)
	}
	return obj.(*v1.ClusterDomain), nil
}
