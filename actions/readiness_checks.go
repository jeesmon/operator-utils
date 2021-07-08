/*
SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"fmt"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ConditionStatusSuccess = "True"
)

type ResourceNotReadyError struct {
	PartialObject client.Object
}

func (e *ResourceNotReadyError) Error() string {
	return fmt.Sprintf("%v/%v is not ready", e.PartialObject.GetNamespace(), e.PartialObject.GetName())
}

func IsResourceNotReadyError(err error) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
	case *ResourceNotReadyError:
		return true
	default:
		return false
	}
}

func IsDeploymentReady(resource *appsv1.Deployment) (bool, error) {
	if resource == nil {
		return false, nil
	}
	// A deployment has an array of conditions
	for _, condition := range resource.Status.Conditions {
		// One failure condition exists, if this exists, return the Reason
		if condition.Type == appsv1.DeploymentReplicaFailure {
			return false, errors.Errorf(condition.Reason)
			// A successful deployment will have the progressing condition type as true
		} else if condition.Type == appsv1.DeploymentProgressing && condition.Status != ConditionStatusSuccess {
			return false, nil
		}
	}
	return true, nil
}

func IsEndpointsReady(resource *corev1.Endpoints) (bool, error) {
	if resource == nil {
		return false, nil
	}

	for _, s := range resource.Subsets {
		if len(s.Addresses) > 0 {
			return true, nil
		}
	}

	return false, nil
}

func IsJobReady(resource *batchv1.Job) (bool, error) {
	if resource == nil {
		return false, nil
	}

	if resource.Status.Failed > 0 {
		return false, errors.New(fmt.Sprintf("Job Failed, check log for %v/%v", resource.Namespace, resource.Name))
	} else if resource.Status.Active > 0 || resource.Status.Succeeded == 0 {
		return false, nil
	} else if resource.Status.Succeeded > 0 {
		return true, nil
	}

	return false, nil
}

func IsServiceMeshControlPlaneReady(resource *unstructured.Unstructured) (bool, error) {
	if resource == nil {
		return false, nil
	}

	items, found, err := unstructured.NestedSlice(resource.Object, "status", "conditions")
	if !found || err != nil {
		return false, errors.Errorf("Status Conditions for ServiceMeshControlPlane is not found")
	}

	for _, item := range items {
		condition, ok := item.(map[string]interface{})
		if !ok {
			return false, errors.Errorf("Status Conditions for ServiceMeshControlPlane is not found")
		}

		if condition["type"] == "Ready" && condition["reason"] == "ComponentsReady" {
			return true, nil
		}
	}

	return false, nil
}

func IsServiceMeshMemberRollReady(resource *unstructured.Unstructured) (bool, error) {
	if resource == nil {
		return false, nil
	}

	items, found, err := unstructured.NestedSlice(resource.Object, "status", "conditions")
	if !found || err != nil {
		return false, errors.Errorf("Status Conditions for ServiceMeshMemberRoll is not found")
	}

	for _, item := range items {
		condition, ok := item.(map[string]interface{})
		if !ok {
			return false, errors.Errorf("Status Conditions for ServiceMeshMemberRoll is not found")
		}

		if condition["type"] == "Ready" && condition["status"] == ConditionStatusSuccess {
			return true, nil
		}
	}

	return false, nil
}

func IsServiceMeshMemberReady(resource *unstructured.Unstructured) (bool, error) {
	if resource == nil {
		return false, nil
	}

	items, found, err := unstructured.NestedSlice(resource.Object, "status", "conditions")
	if !found || err != nil {
		return false, errors.Errorf("Status Conditions for ServiceMeshMember is not found")
	}

	for _, item := range items {
		condition, ok := item.(map[string]interface{})
		if !ok {
			return false, errors.Errorf("Status Conditions for ServiceMeshMember is not found")
		}

		if condition["type"] == "Ready" && condition["status"] == ConditionStatusSuccess {
			return true, nil
		}
	}

	return false, nil
}
