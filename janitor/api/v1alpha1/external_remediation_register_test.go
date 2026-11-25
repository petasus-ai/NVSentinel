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

package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
)

// The tests below use the janitor alias names to prove that the aliases for the
// types defined in data-models/pkg/apis/nvsentinel/v1 keep working.

func TestAddNVSentinelToScheme_RegistersTypes(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	err := AddNVSentinelToScheme(scheme)
	require.NoError(t, err)

	gvk := NVSentinelGroupVersion.WithKind("ExternalRemediationRequest")
	obj, err := scheme.New(gvk)
	require.NoError(t, err)
	assert.IsType(t, &ExternalRemediationRequest{}, obj)

	listGVK := NVSentinelGroupVersion.WithKind("ExternalRemediationRequestList")
	listObj, err := scheme.New(listGVK)
	require.NoError(t, err)
	assert.IsType(t, &ExternalRemediationRequestList{}, listObj)
}

func TestAddNVSentinelToScheme_DoesNotConflictWithJanitorScheme(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	require.NoError(t, AddToScheme(scheme))           // janitor.dgxc.nvidia.com
	require.NoError(t, AddNVSentinelToScheme(scheme)) // nvsentinel.dgxc.nvidia.com

	// Both groups should resolve independently.
	_, err := scheme.New(GroupVersion.WithKind("RebootNode"))
	require.NoError(t, err, "RebootNode (janitor.dgxc.nvidia.com) must still resolve")

	_, err = scheme.New(NVSentinelGroupVersion.WithKind("ExternalRemediationRequest"))
	require.NoError(t, err, "ExternalRemediationRequest (nvsentinel.dgxc.nvidia.com) must resolve")
}
