package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NoiseSpec defines the desired state of Noise
type NoiseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Count    int    `json:"count,omitempty"`
	Interval string `json:"interval,omitempty"`
}

// NoiseStatus defines the observed state of Noise
type NoiseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Noise is the Schema for the noises API
type Noise struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NoiseSpec   `json:"spec,omitempty"`
	Status NoiseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NoiseList contains a list of Noise
type NoiseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Noise `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Noise{}, &NoiseList{})
}
