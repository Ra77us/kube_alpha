package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TestoperartorSpec struct {
	PodImage string `json:"podImage,omitempty"`
}

type TestoperartorStatus struct {
}

type Testoperartor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestoperartorSpec   `json:"spec,omitempty"`
	Status TestoperartorStatus `json:"status,omitempty"`
}

type TestoperartorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Testoperartor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Testoperartor{}, &TestoperartorList{})
}
