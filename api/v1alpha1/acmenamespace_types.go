/*


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

// AcmeNamespaceSpec defines the desired state of AcmeNamespace
type AcmeNamespaceSpec struct {

	// The name of the namespace
	NamespaceName string `json:"namespaceName"`

	// The username for the namespace admin
	AdminUsername string `json:"adminUsername"`
}

// AcmeNamespaceStatus defines the observed state of AcmeNamespace
type AcmeNamespaceStatus struct {

	// Tracks the phase of the AcmeNamespace
	// +optional
	// +kubebuilder:validation:Enum=CreationInProgress;Created
	Phase string `json:"phase"`
}

// +kubebuilder:resource:scope=Cluster
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AcmeNamespace is the Schema for the acmenamespaces API
type AcmeNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AcmeNamespaceSpec   `json:"spec,omitempty"`
	Status AcmeNamespaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AcmeNamespaceList contains a list of AcmeNamespace
type AcmeNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AcmeNamespace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AcmeNamespace{}, &AcmeNamespaceList{})
}
