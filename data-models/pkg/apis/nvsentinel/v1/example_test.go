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

package v1_test

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nvsentinelv1 "github.com/nvidia/nvsentinel/data-models/pkg/apis/nvsentinel/v1"
	"github.com/nvidia/nvsentinel/data-models/pkg/condition"
	protos "github.com/nvidia/nvsentinel/data-models/pkg/protos"
)

// An external system reports the result of a repair. It changes only the
// ExternalRemediationComplete condition and keeps every other condition, because
// a status update replaces the full list.
func ExampleExternalRemediationRequest_writeResult() {
	// The ExternalRemediationRequest as the external system reads it, after
	// janitor released the node.
	extrr := &nvsentinelv1.ExternalRemediationRequest{
		ObjectMeta: metav1.ObjectMeta{Name: "extrr-he-7f0b3e2c"},
		Spec: &protos.ExternalRemediationRequestSpec{
			HealthEvent: &protos.HealthEvent{NodeName: "node-01"},
		},
		Status: &protos.ExternalRemediationRequestStatus{
			Conditions: []*protos.Condition{
				{
					Type:   nvsentinelv1.ConditionNVSentinelOwnershipReleased,
					Status: string(metav1.ConditionTrue),
					Reason: "ReleaseTaintApplied",
				},
				{
					Type:   nvsentinelv1.ConditionExternalRemediationComplete,
					Status: string(metav1.ConditionUnknown),
					Reason: "AwaitingExternalSystem",
				},
			},
		},
	}

	conditions := condition.ToMetav1Slice(extrr.Status.Conditions)
	meta.SetStatusCondition(&conditions, metav1.Condition{
		Type:    nvsentinelv1.ConditionExternalRemediationComplete,
		Status:  metav1.ConditionTrue,
		Reason:  "RepairSucceeded",
		Message: "The node is repaired.",
	})
	extrr.Status.Conditions = condition.FromMetav1Slice(conditions)

	// A controller-runtime client then writes the status:
	// err := c.Status().Update(ctx, extrr)

	for _, c := range extrr.Status.Conditions {
		fmt.Printf("%s=%s reason=%s\n", c.Type, c.Status, c.Reason)
	}

	// Output:
	// NVSentinelOwnershipReleased=True reason=ReleaseTaintApplied
	// ExternalRemediationComplete=True reason=RepairSucceeded
}

// A maintenance tool asks NVSentinel to cordon and drain a node. The JSON is
// what a client sends to the API server.
func ExampleMaintenanceRequest() {
	mr := &nvsentinelv1.MaintenanceRequest{
		TypeMeta: metav1.TypeMeta{
			APIVersion: nvsentinelv1.GroupVersion.String(),
			Kind:       nvsentinelv1.MaintenanceRequestKind,
		},
		ObjectMeta: metav1.ObjectMeta{Name: "maintenance-node-01"},
		Spec: &protos.MaintenanceRequestSpec{
			HealthEvent: &protos.HealthEvent{
				Version:           1,
				Agent:             "my-maintenance-tool",
				ComponentClass:    "Node",
				CheckName:         "PlannedMaintenance",
				NodeName:          "node-01",
				IsHealthy:         false,
				IsFatal:           false,
				RecommendedAction: protos.RecommendedAction_NONE,
				Message:           "Planned hardware maintenance",
			},
		},
	}

	// A controller-runtime client creates the object:
	// err := c.Create(ctx, mr)

	b, err := json.MarshalIndent(mr, "", "  ")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(string(b))

	// Output:
	// {
	//   "kind": "MaintenanceRequest",
	//   "apiVersion": "nvsentinel.dgxc.nvidia.com/v1",
	//   "metadata": {
	//     "name": "maintenance-node-01"
	//   },
	//   "spec": {
	//     "healthEvent": {
	//       "version": 1,
	//       "agent": "my-maintenance-tool",
	//       "componentClass": "Node",
	//       "checkName": "PlannedMaintenance",
	//       "message": "Planned hardware maintenance",
	//       "nodeName": "node-01"
	//     }
	//   }
	// }
}
