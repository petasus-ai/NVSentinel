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

// Excluded when controller-gen loads the package: the CRD for this type comes from
// protoc-gen-crd, and controller-gen cannot build a schema from protobuf structs.
//go:build !controllergen

package v1alpha1

import (
	nvsentinelv1 "github.com/nvidia/nvsentinel/data-models/pkg/apis/nvsentinel/v1"
)

// MaintenanceRequest and its list type are defined in
// github.com/nvidia/nvsentinel/data-models/pkg/apis/nvsentinel/v1, so that
// systems outside this repository can import them. These aliases keep the
// lifecycle-manager names working.
type (
	MaintenanceRequest     = nvsentinelv1.MaintenanceRequest
	MaintenanceRequestList = nvsentinelv1.MaintenanceRequestList
)

// The nvsentinel.dgxc.nvidia.com/v1 registration lives in the same package. It
// registers ExternalRemediationRequest too, because both kinds share the group
// and version. It is independent of the SchemeBuilder for
// nvsentinel.nvidia.com (defined in groupversion_info.go), so the two groups
// can be registered independently.
var (
	MRGroupVersion  = nvsentinelv1.GroupVersion
	MRSchemeBuilder = nvsentinelv1.SchemeBuilder
	AddMRToScheme   = nvsentinelv1.AddToScheme
)
