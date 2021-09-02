/*
SPDX-License-Identifier: Apache-2.0
*/

package status

import (
	conditions "github.com/openshift/custom-resource-status/conditions/v1"
	objectreferences "github.com/openshift/custom-resource-status/objectreferences/v1"
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/reference"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StatusReason string

var (
	ReasonReconciling  StatusReason = "Reconciling"
	ReasonFailing      StatusReason = "Failing"
	ReasonInitializing StatusReason = "Initializing"
)

// CommonStatusSpec defines the Common Status Spec
type CommonStatusSpec struct {
	// conditions describes the state of the operator's reconciliation functionality.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	// Conditions is a list of conditions related to operator reconciliation
	Conditions []conditions.Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`
	// RelatedObjects is a list of objects that are "interesting" or related to this operator
	//+operator-sdk:csv:customresourcedefinitions:type=status
	RelatedObjects []corev1.ObjectReference `json:"relatedObjects,omitempty"`
}

func UpdateStatusRelatedObjects(objects *[]corev1.ObjectReference, scheme *runtime.Scheme, resource client.Object) error {
	objectRef, err := reference.GetReference(scheme, resource)
	if err != nil {
		return err
	}
	objectreferences.SetObjectReference(objects, *objectRef)

	return nil
}
