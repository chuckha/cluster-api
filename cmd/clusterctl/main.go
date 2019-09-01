/*
Copyright 2018 The Kubernetes Authors.

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

package main

import (
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/klogr"
	clusterv1 "github.com/chuckha/cluster-api/api/v1alpha2"
	"github.com/chuckha/cluster-api/cmd/clusterctl/cmd"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func main() {
	log.SetLogger(klogr.New())
	clusterv1.AddToScheme(scheme.Scheme)
	cmd.Execute()
}
