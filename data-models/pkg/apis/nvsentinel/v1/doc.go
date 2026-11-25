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

// Package v1 contains the Go types of the ExternalRemediationRequest and
// MaintenanceRequest custom resources in the nvsentinel.dgxc.nvidia.com/v1 API
// group.
//
// The package is for external systems that read or write these resources with
// a typed Kubernetes client. Register the types with AddToScheme. The module
// that contains this package does not require any other NVSentinel module, so
// a module outside the NVSentinel repository can import it. The NVSentinel
// components use the same types.
//
// The protobuf messages in data-models/protobufs own the spec and status
// shapes: external_remediation.proto and maintenance_request.proto. The CRDs
// are generated from these messages with protoc-gen-crd. The types in this
// package add only the Kubernetes object machinery: TypeMeta, ObjectMeta,
// DeepCopy, and a JSON codec that encodes spec and status with protojson.
//
// The status conditions are protobuf messages. To use them with the
// controller-runtime condition helpers, such as meta.SetStatusCondition,
// convert them with github.com/nvidia/nvsentinel/data-models/pkg/condition.
package v1
