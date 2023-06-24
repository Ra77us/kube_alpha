/*
Copyright 2023.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1 "kubiki.amocna/operator/api/v1"
)

// HephaestusDeploymentReconciler reconciles a HephaestusDeployment object
type HephaestusDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operator.kubiki.amocna,resources=hephaestusdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.kubiki.amocna,resources=hephaestusdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.kubiki.amocna,resources=hephaestusdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HephaestusDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *HephaestusDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var hephaestusDeployment operatorv1.HephaestusDeployment
	if err := r.Get(ctx, req.NamespacedName, &hephaestusDeployment); err != nil {
		log.Log.Error(err, "unable to fetch Hesphaestus Deployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Log.Info("Reconciling Test Hesphaestus Deployment", "Hesphaestus Deployment", hephaestusDeployment)

	//volume

	//volumeDeployment := getVolumeDeployment(hephaestusDeployment)
	//if err := r.Create(ctx, &volumeDeployment); err != nil {
	//	log.Log.Error(err, "unable to create volume Deployment", volumeDeployment.Name)
	//	return ctrl.Result{}, err
	//}
	//log.Log.Info("Created Deployment", volumeDeployment.Name)

	//gui
	log.FromContext(ctx).Info("GUI Version is ", "HephaestusGuiVersion", hephaestusDeployment.Spec.HephaestusGuiVersion)

	if hephaestusDeployment.Spec.HephaestusGuiVersion == "" {
		log.Log.Info("GUI Version is not set")
	} else {
		log.Log.Info("GUI Version is set", "HephaestusGuiVersion", hephaestusDeployment.Spec.HephaestusGuiVersion)
	}

	guiDeployment := getGuiDeployment(hephaestusDeployment)
	if err := r.Create(ctx, &guiDeployment); err != nil {
		log.Log.Error(err, "unable to create Deployment", "Deployment.Namespace", guiDeployment.Namespace, "Deployment.Name", guiDeployment.Name)
		return ctrl.Result{}, err
	}
	log.Log.Info("Created Deployment", "Deployment.Namespace", guiDeployment.Namespace, "Deployment.Name", guiDeployment.Name)

	//metrics-adapter
	log.FromContext(ctx).Info("Metrics Adapter Image is ", "MetricsAdapterImage", hephaestusDeployment.Spec.MetricsAdapterImage)

	if hephaestusDeployment.Spec.MetricsAdapterImage == "" {
		log.Log.Info("Metrics Adapter Image is not set")
	} else {
		log.Log.Info("Metrics Adapter Image is set", "MetricsAdapterImage", hephaestusDeployment.Spec.MetricsAdapterImage)
	}

	metricsAdapterDeployment := getMetricsAdapterDeployment(hephaestusDeployment)
	if err := r.Create(ctx, &guiDeployment); err != nil {
		log.Log.Error(err, "unable to create metrics adapter Deployment", "Deployment.Namespace", metricsAdapterDeployment.Namespace, "Deployment.Name", metricsAdapterDeployment.Name)
		return ctrl.Result{}, err
	}
	log.Log.Info("Created Deployment", "Deployment.Namespace", metricsAdapterDeployment.Namespace, "Deployment.Name", metricsAdapterDeployment.Name)

	//execution-controller
	log.FromContext(ctx).Info("Execution Controller Image is ", "ExecutionControllerImage", hephaestusDeployment.Spec.ExecutionControllerImage)

	if hephaestusDeployment.Spec.ExecutionControllerImage == "" {
		log.Log.Info("Execution Controller Image is not set")
	} else {
		log.Log.Info("Execution Controller Image is set", "ExecutionControllerImage", hephaestusDeployment.Spec.ExecutionControllerImage)
	}

	executionControllerDeployment := getExecutionControllerDeployment(hephaestusDeployment)
	if err := r.Create(ctx, &guiDeployment); err != nil {
		log.Log.Error(err, "unable to create execution controller Deployment", "Deployment.Namespace", executionControllerDeployment.Namespace, "Deployment.Name", executionControllerDeployment.Name)
		return ctrl.Result{}, err
	}
	log.Log.Info("Created Deployment", "Deployment.Namespace", executionControllerDeployment.Namespace, "Deployment.Name", executionControllerDeployment.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HephaestusDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1.HephaestusDeployment{}).
		Complete(r)
}
