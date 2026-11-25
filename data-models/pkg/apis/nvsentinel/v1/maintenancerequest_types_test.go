// Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	protos "github.com/nvidia/nvsentinel/data-models/pkg/protos"
)

func newMR() *MaintenanceRequest {
	return &MaintenanceRequest{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "nvsentinel.dgxc.nvidia.com/v1",
			Kind:       "MaintenanceRequest",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:       "maintenance-node-01",
			Labels:     map[string]string{"team": "hardware"},
			Finalizers: []string{"nvsentinel.dgxc.nvidia.com/maintenance-request-cleanup"},
		},
		Spec: &protos.MaintenanceRequestSpec{
			HealthEvent: &protos.HealthEvent{
				Version:           1,
				Agent:             "lifecycle-manager",
				ComponentClass:    "Node",
				CheckName:         "PlannedMaintenance",
				NodeName:          "node-01.example-cluster.internal",
				RecommendedAction: protos.RecommendedAction_NONE,
				Message:           "Planned hardware maintenance",
				Metadata:          map[string]string{"maintenanceRequestRequesterAgent": "my-maintenance-tool"},
			},
			StartTime: timestamppb.New(time.Date(2026, 8, 15, 10, 0, 0, 0, time.UTC)),
		},
		Status: &protos.MaintenanceRequestStatus{
			Conditions: []*protos.Condition{
				{
					Type:               "HealthEventEmitted",
					Status:             "True",
					LastTransitionTime: timestamppb.New(time.Date(2026, 8, 15, 10, 0, 5, 0, time.UTC)),
					Reason:             "Emitted",
				},
			},
		},
	}
}

func TestMaintenanceRequest_DeepCopyObject_NilSafe(t *testing.T) {
	t.Parallel()

	var nilMR *MaintenanceRequest

	assert.Nil(t, nilMR.DeepCopyObject())
	assert.Nil(t, nilMR.DeepCopy())
}

func TestMaintenanceRequest_DeepCopy_FieldFidelity(t *testing.T) {
	t.Parallel()

	original := newMR()
	copied := original.DeepCopy()

	require.NotNil(t, copied)
	assert.Equal(t, original.TypeMeta, copied.TypeMeta)
	assert.Equal(t, original.Name, copied.Name)
	assert.Equal(t, original.Labels, copied.Labels)
	assert.Equal(t, original.Finalizers, copied.Finalizers)
	require.NotNil(t, copied.Spec)
	require.NotNil(t, copied.Spec.HealthEvent)
	assert.Equal(t, original.Spec.HealthEvent.Version, copied.Spec.HealthEvent.Version)
	assert.Equal(t, original.Spec.HealthEvent.NodeName, copied.Spec.HealthEvent.NodeName)
	assert.Equal(t, original.Spec.HealthEvent.CheckName, copied.Spec.HealthEvent.CheckName)
	assert.Equal(t, original.Spec.HealthEvent.Message, copied.Spec.HealthEvent.Message)
	assert.Equal(t, original.Spec.HealthEvent.Metadata, copied.Spec.HealthEvent.Metadata)
	require.NotNil(t, copied.Spec.StartTime)
	assert.Equal(t, original.Spec.StartTime.AsTime(), copied.Spec.StartTime.AsTime())
	require.NotNil(t, copied.Status)
	require.Len(t, copied.Status.Conditions, 1)
	assert.Equal(t, original.Status.Conditions[0].Type, copied.Status.Conditions[0].Type)
	assert.Equal(t, original.Status.Conditions[0].Status, copied.Status.Conditions[0].Status)
	assert.Equal(t, original.Status.Conditions[0].Reason, copied.Status.Conditions[0].Reason)
}

// TestMaintenanceRequest_DeepCopy_IsIndependent verifies that mutating the deep
// copy does not affect the original.
func TestMaintenanceRequest_DeepCopy_IsIndependent(t *testing.T) {
	t.Parallel()

	original := newMR()
	copied := original.DeepCopy()
	require.NotNil(t, copied)

	copied.Labels["mutated"] = "yes"
	assert.NotContains(t, original.Labels, "mutated", "ObjectMeta must be independently deep-copied")

	copied.Spec.HealthEvent.Message = "mutated"
	assert.Equal(t, "Planned hardware maintenance", original.Spec.HealthEvent.Message,
		"Spec.HealthEvent must be independently proto-cloned")

	copied.Spec.HealthEvent.Metadata["mutated"] = "yes"
	assert.NotContains(t, original.Spec.HealthEvent.Metadata, "mutated",
		"Spec.HealthEvent.Metadata must be independently proto-cloned")

	copied.Status.Conditions[0].Reason = "mutated"
	assert.Equal(t, "Emitted", original.Status.Conditions[0].Reason,
		"Status.Conditions must be independently proto-cloned")
}

func TestMaintenanceRequest_DeepCopyObject_ReturnsRuntimeObject(t *testing.T) {
	t.Parallel()

	original := newMR()
	obj := original.DeepCopyObject()

	require.IsType(t, &MaintenanceRequest{}, obj)
	copied := obj.(*MaintenanceRequest)
	assert.Equal(t, original.Name, copied.Name)
}

// TestMaintenanceRequest_DeepCopyInto_NilSourceFields_ClearDirtyDestination
// verifies the documented MaintenanceRequest behaviour: a nil Spec or Status in
// the source sets the field in out to nil, so a reused out does not keep a
// stale value.
func TestMaintenanceRequest_DeepCopyInto_NilSourceFields_ClearDirtyDestination(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		nilSpec    bool
		nilStatus  bool
		wantSpec   bool
		wantStatus bool
	}{
		{name: "nil spec", nilSpec: true, wantStatus: true},
		{name: "nil status", nilStatus: true, wantSpec: true},
		{name: "nil spec and status", nilSpec: true, nilStatus: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			src := newMR()
			if tc.nilSpec {
				src.Spec = nil
			}

			if tc.nilStatus {
				src.Status = nil
			}

			dirty := newMR()
			dirty.Spec.HealthEvent.Message = "stale"
			dirty.Status.Conditions[0].Reason = "stale"

			src.DeepCopyInto(dirty)

			if tc.wantSpec {
				require.NotNil(t, dirty.Spec)
				assert.Equal(t, "Planned hardware maintenance", dirty.Spec.HealthEvent.Message)
			} else {
				assert.Nil(t, dirty.Spec, "a nil source Spec must clear the destination Spec")
			}

			if tc.wantStatus {
				require.NotNil(t, dirty.Status)
				assert.Equal(t, "Emitted", dirty.Status.Conditions[0].Reason)
			} else {
				assert.Nil(t, dirty.Status, "a nil source Status must clear the destination Status")
			}
		})
	}
}

func TestMaintenanceRequestList_DeepCopy_FieldFidelity(t *testing.T) {
	t.Parallel()

	second := newMR()
	second.Name = "maintenance-node-02"

	original := &MaintenanceRequestList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "nvsentinel.dgxc.nvidia.com/v1",
			Kind:       "MaintenanceRequestList",
		},
		ListMeta: metav1.ListMeta{ResourceVersion: "42"},
		Items:    []MaintenanceRequest{*newMR(), *second},
	}

	copied := original.DeepCopy()

	require.NotNil(t, copied)
	assert.Equal(t, original.TypeMeta, copied.TypeMeta)
	assert.Equal(t, original.ResourceVersion, copied.ResourceVersion)
	require.Len(t, copied.Items, 2)
	assert.Equal(t, original.Items[0].Name, copied.Items[0].Name)
	assert.Equal(t, original.Items[1].Name, copied.Items[1].Name)

	copied.Items[0].Spec.HealthEvent.Message = "mutated"
	assert.Equal(t, "Planned hardware maintenance", original.Items[0].Spec.HealthEvent.Message,
		"Items must be independently deep-copied")
}

func TestMaintenanceRequestList_DeepCopyInto_NilItems_ClearDirtyDestination(t *testing.T) {
	t.Parallel()

	src := &MaintenanceRequestList{}
	dirty := &MaintenanceRequestList{Items: []MaintenanceRequest{*newMR()}}

	src.DeepCopyInto(dirty)

	assert.Nil(t, dirty.Items, "nil source Items must clear the destination Items")
}

func TestMaintenanceRequestList_DeepCopy_NilSafe(t *testing.T) {
	t.Parallel()

	var nilList *MaintenanceRequestList

	assert.Nil(t, nilList.DeepCopyObject())
	assert.Nil(t, nilList.DeepCopy())
}
