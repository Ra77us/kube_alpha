package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "kubiki.amocna/operator/api/v1"
)

func getGuiDeployment(hephaestusDeployment operatorv1.HephaestusDeployment) appsv1.Deployment {
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
							Image: "hephaestusmetrics/gui:" + hephaestusDeployment.Spec.HephaestusGuiVersion,
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
						},
					},
				},
			},
		},
	}
}
