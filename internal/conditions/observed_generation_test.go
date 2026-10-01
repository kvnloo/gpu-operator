/**
# Copyright (c) NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package conditions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	nvidiav1 "github.com/NVIDIA/gpu-operator/api/nvidia/v1"
	nvidiav1alpha1 "github.com/NVIDIA/gpu-operator/api/nvidia/v1alpha1"
)

func assertObservedGeneration(t *testing.T, want int64, conditions []metav1.Condition) {
	t.Helper()
	require.Len(t, conditions, 2)
	for _, condition := range conditions {
		assert.Equal(t, want, condition.ObservedGeneration, "condition %q", condition.Type)
	}
}

func TestConditionUpdatersSetObservedGeneration(t *testing.T) {
	const generation int64 = 7
	ctx := context.Background()

	t.Run("ClusterPolicy", func(t *testing.T) {
		cr := newClusterPolicy("cluster-policy")
		cr.Generation = generation
		c := newClusterPolicyClient(t, cr)
		updater := NewClusterPolicyUpdater(c)

		require.NoError(t, updater.SetConditionsReady(ctx, cr, Reconciled, "ready"))
		require.NoError(t, updater.SetConditionsError(ctx, cr, ReconcileFailed, "failed"))

		got := &nvidiav1.ClusterPolicy{}
		require.NoError(t, c.Get(ctx, types.NamespacedName{Name: cr.Name}, got))
		assertObservedGeneration(t, generation, got.Status.Conditions)
	})

	t.Run("NVIDIADriver", func(t *testing.T) {
		cr := newNvDriver("gpu-driver")
		cr.Generation = generation
		c := newNvDriverClient(t, cr)
		updater := NewNvDriverUpdater(c)

		require.NoError(t, updater.SetConditionsReady(ctx, cr, Reconciled, "ready"))
		require.NoError(t, updater.SetConditionsError(ctx, cr, ReconcileFailed, "failed"))

		got := &nvidiav1alpha1.NVIDIADriver{}
		require.NoError(t, c.Get(ctx, types.NamespacedName{Name: cr.Name}, got))
		assertObservedGeneration(t, generation, got.Status.Conditions)
	})

	t.Run("GPUCluster", func(t *testing.T) {
		cr := newGPUCluster("gpu-cluster", nvidiav1alpha1.Ready)
		cr.Generation = generation
		c := newGPUClusterClient(t, cr)
		updater := NewGPUClusterUpdater(c)

		require.NoError(t, updater.SetConditionsReady(ctx, cr, Reconciled, "ready"))
		require.NoError(t, updater.SetConditionsError(ctx, cr, ReconcileFailed, "failed"))

		got := getGPUCluster(t, c, cr.Name)
		assertObservedGeneration(t, generation, got.Status.Conditions)
	})
}
