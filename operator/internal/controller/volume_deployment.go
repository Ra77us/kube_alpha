package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getVolumeDeployment(hephaestusDeployment operatorv1.HephaestusDeployment) appsv1.Deployment {
	one := int32(1)
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: hephaestusDeployment.Name + "-volume-deployment",
		},
		TypeMeta: metav1.TypeMeta{Kind: "PersistentVolume"},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: "hephaestus-manual",
			//AccessModes: []v1.PersistentVolumeAccessMode{
			//	"ReadWriteOnce",
			//},
			//Capacity: "25Mi",
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/mnt/hephaestus-gui-pv",
				},
			},
		},
	}
}
