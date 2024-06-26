/*
Copyright 2024 0jk6.

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

// custom secret data
type SecretToPull struct {
	SecretIdentifier string `json:"secretIdentifier"`
}

// ConjurSpec defines the desired state of Conjur
type ConjurSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	RefreshInterval  int                     `json:"refreshInterval"`
	ApiKeyFromSecret string                  `json:"apiKeyFromSecret"`
	ConjurHost       string                  `json:"conjurHost"`
	ConjurAcct       string                  `json:"conjurAcct"`
	Hostname         string                  `json:"hostname"`
	Data             map[string]SecretToPull `json:"data"`
}

// ConjurStatus defines the observed state of Conjur
type ConjurStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Conjur is the Schema for the conjurs API
type Conjur struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConjurSpec   `json:"spec,omitempty"`
	Status ConjurStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConjurList contains a list of Conjur
type ConjurList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Conjur `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Conjur{}, &ConjurList{})
}
