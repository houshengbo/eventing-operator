// +build e2e

/*
Copyright 2020 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"testing"

	"knative.dev/eventing-operator/test/common"
	"knative.dev/eventing-operator/test"
	"knative.dev/eventing-operator/test/resources"
	"knative.dev/pkg/test/logstream"
)

// TestKnativeEventingDeployment verifies the KnativeEventing creation, deployment recreation, and KnativeEventing deletion.
func TestKnativeEventingDeployment(t *testing.T) {
	cancel := logstream.Start(t)
	defer cancel()
	clients := common.Setup(t)

	names := test.ResourceNames{
		KnativeEventing: test.EventingOperatorName,
		Namespace:       test.EventingOperatorNamespace,
	}

	test.CleanupOnInterrupt(func() { test.TearDown(clients, names) })
	defer test.TearDown(clients, names)

	// Create a KnativeEventing
	if _, err := resources.CreateKnativeEventing(clients.KnativeEventing(), names); err != nil {
		t.Fatalf("KnativeService %q failed to create: %v", names.KnativeEventing, err)
	}

	// Test if KnativeEventing can reach the READY status
	t.Run("create", func(t *testing.T) {
		common.KnativeEventingVerify(t, clients, names)
	})

	// Delete the deployments one by one to see if they will be recreated.
	t.Run("restore", func(t *testing.T) {
		common.KnativeEventingVerify(t, clients, names)
		common.DeploymentRecreation(t, clients, names)
	})

	// Delete the KnativeEventing to see if all resources will be removed
	t.Run("delete", func(t *testing.T) {
		common.KnativeEventingVerify(t, clients, names)
		common.KnativeEventingDelete(t, clients, names)
	})
}
