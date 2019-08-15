/*
Copyright 2019 The Kubernetes Authors.

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

package controllers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Cluster Reconciler", func() {

	It("Should create a Cluster", func() {
		instance := &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "foo-",
				Namespace:    "default",
			},
			Spec: clusterv1.ClusterSpec{},
		}

		// Create the Cluster object and expect the Reconcile and Deployment to be created
		Expect(k8sClient.Create(ctx, instance)).ToNot(HaveOccurred())
		defer k8sClient.Delete(ctx, instance)

		// Make sure the Cluster exists.
		Eventually(func() bool {
			key := client.ObjectKey{Namespace: instance.Namespace, Name: instance.Name}
			if err := k8sClient.Get(ctx, key, instance); err != nil {
				return false
			}
			return len(instance.Finalizers) > 0
		}, timeout).Should(BeTrue())
	})

	It("Should successfully patch a cluster object if the status diff is empty but the spec diff is not", func() {
		// Setup
		cluster := &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-cluster",
				Namespace: "default",
			},
		}
		key := client.ObjectKey{Name: "test-cluster", Namespace: "default"}
		Expect(k8sClient.Create(ctx, cluster)).To(BeNil())
		defer k8sClient.Delete(ctx, cluster)

		// Reconcile
		Eventually(func() bool {
			patch := client.MergeFrom(cluster.DeepCopy())
			cluster.Spec.InfrastructureRef = &v1.ObjectReference{Name: "test"}
			Expect(clusterReconciler.patchCluster(ctx, cluster, patch)).To(BeNil())
			return true
		}, timeout).Should(BeTrue())

		// Assertions
		Eventually(func() bool {
			instance := &clusterv1.Cluster{}
			Expect(k8sClient.Get(ctx, key, instance)).To(BeNil())
			Expect(instance.Spec.InfrastructureRef.Name).To(BeEquivalentTo("test"))
			return true
		}, timeout).Should(BeTrue())
	})

	It("Should successfully patch a cluster object if the spec diff is empty but the status diff is not", func() {
		// Setup
		cluster := &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-cluster",
				Namespace: "default",
			},
		}
		key := client.ObjectKey{Name: "test-cluster", Namespace: "default"}
		Expect(k8sClient.Create(ctx, cluster)).To(BeNil())
		defer k8sClient.Delete(ctx, cluster)

		// Reconcile
		Eventually(func() bool {
			patch := client.MergeFrom(cluster.DeepCopy())
			cluster.Status.InfrastructureReady = true
			Expect(clusterReconciler.patchCluster(ctx, cluster, patch)).To(BeNil())
			return true
		}, timeout).Should(BeTrue())

		// Assertions
		Eventually(func() bool {
			instance := &clusterv1.Cluster{}
			Expect(k8sClient.Get(ctx, key, instance)).To(BeNil())
			Expect(instance.Status.InfrastructureReady).To(BeTrue())
			return true
		}, timeout).Should(BeTrue())
	})
})
