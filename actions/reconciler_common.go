/*
SPDX-License-Identifier: Apache-2.0
*/

package actions

import (
	"context"
	"time"

	"github.com/jeesmon/operator-utils/status"
	conditions "github.com/openshift/custom-resource-status/conditions/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	RequeueDelay      = 60 * time.Minute
	RequeueDelayError = 5 * time.Second
)

type ResourceState interface {
	Read(context.Context, client.Object) error
	IsResourcesReady(client.Object) (bool, error)
}

func ManageError(client client.Client, ctx context.Context, instance client.Object, statusConditions *[]conditions.Condition, issue error) (reconcile.Result, error) {
	condition := conditions.Condition{
		Type:    conditions.ConditionAvailable,
		Status:  v1.ConditionFalse,
		Message: issue.Error(),
	}

	if IsResourceNotReadyError(issue) {
		condition.Reason = string(status.ReasonInitializing)
	} else {
		condition.Reason = string(status.ReasonFailing)
	}

	conditions.SetStatusCondition(statusConditions, condition)

	err := client.Status().Update(ctx, instance)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	return reconcile.Result{
		RequeueAfter: RequeueDelayError,
		Requeue:      true,
	}, nil
}

func ManageSuccess(client client.Client, ctx context.Context, instance client.Object, statusConditions *[]conditions.Condition, resourcesReady bool) (reconcile.Result, error) {
	condition := conditions.Condition{
		Type: conditions.ConditionAvailable,
	}

	// If resources are ready and we have not errored before now, we are in a reconciling phase
	if resourcesReady {
		condition.Status = v1.ConditionTrue
		condition.Reason = string(status.ReasonReconciling)
		condition.Message = "All resource are ready"
	} else {
		condition.Status = v1.ConditionFalse
		condition.Reason = string(status.ReasonInitializing)
		condition.Message = "One or more resources are not ready"
	}

	conditions.SetStatusCondition(statusConditions, condition)

	err := client.Status().Update(ctx, instance)
	if err != nil {
		log.Error(err, "unable to update status")
		return reconcile.Result{
			RequeueAfter: RequeueDelayError,
			Requeue:      true,
		}, nil
	}

	return reconcile.Result{RequeueAfter: RequeueDelay}, nil
}

func RunDesiredStateActions(client client.Client, scheme *runtime.Scheme, ctx context.Context, instance client.Object, conditions *[]conditions.Condition, currentState ResourceState, desiredState DesiredResourceState) (reconcile.Result, error) {
	// Run the actions to reach the desired state
	actionRunner := NewControllerActionRunner(ctx, client, scheme, instance)
	err := actionRunner.RunAll(desiredState)
	if err != nil {
		return ManageError(client, ctx, instance, conditions, err)
	}

	resourcesReady, err := currentState.IsResourcesReady(instance)
	if err != nil {
		return ManageError(client, ctx, instance, conditions, err)
	}

	return ManageSuccess(client, ctx, instance, conditions, resourcesReady)
}

func ReadCurrentState(client client.Client, ctx context.Context, instance client.Object, conditions *[]conditions.Condition, currentState ResourceState) (reconcile.Result, error) {
	err := currentState.Read(ctx, instance)
	if err != nil {
		return ManageError(client, ctx, instance, conditions, err)
	}

	return reconcile.Result{}, nil
}
