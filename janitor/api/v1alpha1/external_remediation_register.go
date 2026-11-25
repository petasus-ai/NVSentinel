// Copyright (c) 2025, NVIDIA CORPORATION.  All rights reserved.
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

// ExternalRemediationRequest and its list type are defined in
// github.com/nvidia/nvsentinel/data-models/pkg/apis/nvsentinel/v1, so that
// systems outside this repository can import them. These aliases keep the
// janitor names working.
type (
	ExternalRemediationRequest     = nvsentinelv1.ExternalRemediationRequest
	ExternalRemediationRequestList = nvsentinelv1.ExternalRemediationRequestList
)

// The nvsentinel.dgxc.nvidia.com/v1 registration lives in the same package. It
// registers MaintenanceRequest too, because both kinds share the group and
// version. It is independent of the SchemeBuilder for janitor.dgxc.nvidia.com
// (defined in groupversion_info.go), so callers register both groups.
var (
	NVSentinelGroupVersion  = nvsentinelv1.GroupVersion
	NVSentinelSchemeBuilder = nvsentinelv1.SchemeBuilder
	AddNVSentinelToScheme   = nvsentinelv1.AddToScheme
)
