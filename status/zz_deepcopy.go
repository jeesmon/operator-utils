//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package status

import (
	v1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonStatusSpec) DeepCopyInto(out *CommonStatusSpec) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RelatedObjects != nil {
		in, out := &in.RelatedObjects, &out.RelatedObjects
		*out = make([]corev1.ObjectReference, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonStatusSpec.
func (in *CommonStatusSpec) DeepCopy() *CommonStatusSpec {
	if in == nil {
		return nil
	}
	out := new(CommonStatusSpec)
	in.DeepCopyInto(out)
	return out
}
