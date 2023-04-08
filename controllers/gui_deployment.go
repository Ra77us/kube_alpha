package controllers

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubikiv1 "kubiki/operator/api/v1"
)

func getGuiDeployment(hephaestusDeployment kubikiv1.HephaestusDeployment) appsv1.Deployment {
	one := int32(1)
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-gui-deployment",
			Namespace: hephaestusDeployment.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &one,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": hephaestusDeployment.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": hephaestusDeployment.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  hephaestusDeployment.Name,
							Image: "hephaestusmetrics/gui:" + hephaestusDeployment.Spec.GuiVersion,
							Env: []corev1.EnvVar{
								{
									Name:  "prometheus.address",
									Value: hephaestusDeployment.Spec.PrometheusAddress,
								},
								{
									Name:  "saved.path",
									Value: "/../storage/metrics/savedMetrics/metrics.json",
								},
								{
									Name:  "config.path",
									Value: "/../storage/metrics/configMetrics/metrics.json",
								},
								{
									Name:  "logs.path",
									Value: "/../storage/logs",
								},
							},
							ImagePullPolicy: corev1.PullPolicy("Always"),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							//VolumeMounts: []corev1.VolumeMount{
							//	{
							//		Name:      "storage",
							//		MountPath: "storage",
							//	},
							//	{
							//		Name:      "config-volume",
							//		MountPath: "storage/metrics/configMetrics",
							//	},
							//},
						},
					},
					//Volumes: []corev1.Volume{
					//	{
					//		Name: "storage",
					//		VolumeSource: corev1.VolumeSource{
					//			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					//				ClaimName: "hephaestus-gui-pvc",
					//			},
					//		},
					//	},
					//	{
					//		Name: "config-map",
					//		VolumeSource: corev1.VolumeSource{
					//			ConfigMap: &corev1.ConfigMapVolumeSource{
					//				LocalObjectReference: corev1.LocalObjectReference{
					//					Name: "gui-configmap",
					//				},
					//				Items: []corev1.KeyToPath{
					//					{
					//						Key:  "metrics.json",
					//						Path: "metrics.json",
					//					},
					//				},
					//			},
					//		},
					//	},
					//},
				},
			},
		},
	}
}
