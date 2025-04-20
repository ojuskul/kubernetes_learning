/*
Copyright 2025.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HTTPMonitorSpec defines the desired state of HTTPMonitor
type HTTPMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of HTTPMonitor. Edit httpmonitor_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// HTTPMonitorStatus defines the observed state of HTTPMonitor
type HTTPMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HTTPMonitor is the Schema for the httpmonitors API
type HTTPMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HTTPMonitorSpec   `json:"spec,omitempty"`
	Status HTTPMonitorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HTTPMonitorList contains a list of HTTPMonitor
type HTTPMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HTTPMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HTTPMonitor{}, &HTTPMonitorList{})
}
