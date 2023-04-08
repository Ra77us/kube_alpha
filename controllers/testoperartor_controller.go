package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kubikiv1 "kubiki.alpha/test2/api/v1"
)

type TestoperartorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *TestoperartorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var testOperator kubikiv1.Testoperartor
	if err := r.Get(ctx, req.NamespacedName, &testOperator); err != nil {
		log.Log.Error(err, "unable to fetch Test Operator")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Log.Info("Reconciling Test Operator", "Test Operator", testOperator)
	log.FromContext(ctx).Info("Pod Image is ", "PodImageName", testOperator.Spec.PodImage)
	if testOperator.Spec.PodImage == "" {
		log.Log.Info("Pod Image is not set")
	} else {
		log.Log.Info("Pod Image is set", "PodImageName", testOperator.Spec.PodImage)
	}
	one := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testOperator.Name + "-deployment",
			Namespace: testOperator.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &one,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": testOperator.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": testOperator.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  testOperator.Name,
							Image: testOperator.Spec.PodImage,
							Env: []corev1.EnvVar{
								{
									Name:  "prometheus.address",
									Value: "prometheus.monitoring:9090",
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
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "storage",
									MountPath: "storage",
								},
								{
									Name:      "config-volume",
									MountPath: "storage/metrics/configMetrics",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "hephaestus-gui-pvc",
								},
							},
						},
						{
							Name: "config-map",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "gui-configmap",
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "metrics.json",
											Path: "metrics.json",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if err := r.Create(ctx, deployment); err != nil {
		log.Log.Error(err, "unable to create Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		return ctrl.Result{}, err
	}
	log.Log.Info("Created Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)

	return ctrl.Result{}, nil
}

func (r *TestoperartorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubikiv1.Testoperartor{}).
		Complete(r)
}
