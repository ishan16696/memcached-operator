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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	MemcachedPhaseCreating MemcachedPhase = "Creating"
	MemcachedPhaseRunning  MemcachedPhase = "Running"
	MemcachedPhaseFailed   MemcachedPhase = "Failed"
	MemcachedPhaseDeleting MemcachedPhase = "Deleting"
)

type MemcachedPhase string

// MemCachedSpec defines the desired state of MemCached
type MemCachedSpec struct {
	// +required
	Type string `json:"type,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +required
	Replicas int32 `json:"replicas,omitempty"`
}

// MemCachedStatus defines the observed state of MemCached
type MemCachedStatus struct {
	// ReadyReplicas is the count of ready replicas.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
	// Nodes are the names of the memcached pods
	Nodes []string `json:"nodes"`
	// Phase is the memcached phase
	Phase MemcachedPhase `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".status.readyReplicas",name="Ready",type="integer"
// +kubebuilder:printcolumn:JSONPath=".spec.replicas",name="Desired",type="integer"
// +kubebuilder:printcolumn:JSONPath=".status.phase",name="Status",type="string"

// MemCached is the Schema for the memcacheds API
type MemCached struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemCachedSpec   `json:"spec,omitempty"`
	Status MemCachedStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MemCachedList contains a list of MemCached
type MemCachedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MemCached `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MemCached{}, &MemCachedList{})
}
