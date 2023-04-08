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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kubikiv1 "kubiki/operator/api/v1"
)

type HephaestusDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *HephaestusDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var hephaestusDeployment kubikiv1.HephaestusDeployment
	if err := r.Get(ctx, req.NamespacedName, &hephaestusDeployment); err != nil {
		log.Log.Error(err, "unable to fetch Hesphaestus Deployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Log.Info("Reconciling Test Hesphaestus Deployment", "Hesphaestus Deployment", hephaestusDeployment)
	log.FromContext(ctx).Info("GUI Version is ", "GuiVersion", hephaestusDeployment.Spec.GuiVersion)

	if hephaestusDeployment.Spec.GuiVersion == "" {
		log.Log.Info("GUI Version is not set")
	} else {
		log.Log.Info("GUI Version is set", "GuiVersion", hephaestusDeployment.Spec.GuiVersion)
	}

	deployment := getGuiDeployment(hephaestusDeployment)
	if err := r.Create(ctx, &deployment); err != nil {
		log.Log.Error(err, "unable to create Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		return ctrl.Result{}, err
	}
	log.Log.Info("Created Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)

	return ctrl.Result{}, nil
}

func (r *HephaestusDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubikiv1.HephaestusDeployment{}).
		Complete(r)
}
