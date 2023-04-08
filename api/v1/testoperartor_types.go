package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HephaestusdeploymentSpec struct {
	PodImage string `json:"podImage,omitempty"`
}

type HephaestusdeploymentStatus struct {
}

type Hephaesdeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HephaestusdeploymentSpec   `json:"spec,omitempty"`
	Status HephaestusdeploymentStatus `json:"status,omitempty"`
}

type Hephaestusdeployment struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Hephaesdeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Hephaesdeployment{}, &Hephaestusdeployment{})
}
