/*
SPDX-License-Identifier: Apache-2.0
*/

package autodetect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var serverGroupsAndResourcesMock func() ([]*metav1.APIGroup, []*metav1.APIResourceList, error)

type discoveryClientMock struct{}

func TestDetectCapabilities(t *testing.T) {
	config := DetectConfig{
		GroupVersionKinds: []schema.GroupVersionKind{
			{
				Group:   "group",
				Version: "version",
				Kind:    "kind1",
			},
			{
				Group:   "group",
				Version: "version",
				Kind:    "kind2",
			},
		},
	}

	bg := &Background{
		config: config,
		dc:     discoveryClientMock{},
	}

	serverGroupsAndResourcesMock = func() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
		apiList := []*metav1.APIResourceList{}
		for _, gvk := range config.GroupVersionKinds {
			gv := schema.GroupVersion{Group: gvk.Group, Version: gvk.Version}
			resources := []metav1.APIResource{
				{
					Kind: gvk.Kind,
				},
			}
			apiList = append(apiList, &metav1.APIResourceList{
				GroupVersion: gv.String(),
				APIResources: resources,
			})
		}
		return nil, apiList, nil
	}

	bg.DetectCapabilities()

	assert.Equal(t, IsResourceAvailable(schema.GroupVersionKind{
		Group:   "group",
		Version: "version",
		Kind:    "kind1",
	}), true)

	assert.Equal(t, IsResourceAvailable(schema.GroupVersionKind{
		Group:   "group",
		Version: "version",
		Kind:    "kind2",
	}), true)

	assert.NotEqual(t, IsResourceAvailable(schema.GroupVersionKind{
		Group:   "group",
		Version: "version",
		Kind:    "unknown",
	}), true)
}

func (dc discoveryClientMock) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	return nil, nil
}

func (dc discoveryClientMock) ServerResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (dc discoveryClientMock) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	return serverGroupsAndResourcesMock()
}

func (dc discoveryClientMock) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}

func (dc discoveryClientMock) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}
