// Package utils modifies some cluster api utils that have too many responsibilities
package utils

import (
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/storage/names"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

const (
	// TemplateSuffix is the object kind suffix used by infrastructure references associated
	// with MachineSet or MachineDeployments.
	TemplateSuffix = "Template"
)

type CloneTemplateInput struct {
	// TemplateRef is a reference to the template that needs to be cloned.
	// +required
	Template *unstructured.Unstructured

	// ClusterName is the cluster this object is linked to.
	// +required
	ClusterName string

	// OwnerRef is an optional OwnerReference to attach to the cloned object.
	// +optional
	OwnerRef *metav1.OwnerReference

	// Labels is an optional map of labels to be added to the object.
	// +optional
	Labels map[string]string
}

// CloneTemplate uses the client and the reference to create a new object from the template.
func CloneTemplate(in *CloneTemplateInput) (*unstructured.Unstructured, error) {
	template, found, err := unstructured.NestedMap(in.Template.Object, "spec", "template")
	if !found {
		return nil, errors.Errorf("missing Spec.Template on %v %q", in.Template.GroupVersionKind(), in.Template.GetName())
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve Spec.Template map on %v %q", in.Template.GroupVersionKind(), in.Template.GetName())
	}

	// Create the unstructured object from the template.
	to := &unstructured.Unstructured{Object: template}
	to.SetResourceVersion("")
	to.SetFinalizers(nil)
	to.SetUID("")
	to.SetSelfLink("")
	to.SetName(names.SimpleNameGenerator.GenerateName(in.Template.GetName() + "-"))
	to.SetNamespace(in.Template.GetNamespace())

	// Set labels.
	labels := to.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	for key, value := range in.Labels {
		labels[key] = value
	}
	labels[clusterv1.ClusterLabelName] = in.ClusterName
	to.SetLabels(labels)

	// Set the owner reference.
	if in.OwnerRef != nil {
		to.SetOwnerReferences([]metav1.OwnerReference{*in.OwnerRef})
	}

	// Set the object APIVersion.
	if to.GetAPIVersion() == "" {
		to.SetAPIVersion(in.Template.GetAPIVersion())
	}

	// Set the object Kind and strip the word "Template" if it's a suffix.
	if to.GetKind() == "" {
		to.SetKind(strings.TrimSuffix(in.Template.GetKind(), TemplateSuffix))
	}

	return to, nil
}

func GetObjectReference(from *unstructured.Unstructured) *corev1.ObjectReference {
	return &corev1.ObjectReference{
		APIVersion: from.GetAPIVersion(),
		Kind:       from.GetKind(),
		Name:       from.GetName(),
		Namespace:  from.GetNamespace(),
		UID:        from.GetUID(),
	}
}
