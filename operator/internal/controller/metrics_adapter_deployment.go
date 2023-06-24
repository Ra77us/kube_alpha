package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getMetricsAdapterDeployment(hephaestusDeployment operatorv1.HephaestusDeployment) appsv1.Deployment {
	one := int32(1)
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hephaestusDeployment.Name + "-metricsadapter-deployment",
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
							Image: "hephaestusmetrics/metrics-adapter:" + hephaestusDeployment.Spec.MetricsAdapterVersion,
							Env: []corev1.EnvVar{
								{
									Name:  "backend",
									Value: "http://hephaestus-gui." + hephaestusDeployment.Namespace + ":8080",
								},
								{
									Name:  "kubernetes-management",
									Value: "http://execution-controller." + hephaestusDeployment.Namespace + ":8097",
								},
							},
							ImagePullPolicy: corev1.PullPolicy("Always"),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}
}
