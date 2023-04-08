package controllers

import (
	"context"

	"github.com/go-logr/logr"
	kapps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	appsv1 "tutorial.kubebuilder.io/project/api/v1"
)

// FooReconciler reconciles a Foo object
type FooReconciler struct {
	client.Client
	Log logr.Logger
	Scheme *runtime.Scheme
}

// RBAC permissions to monitor foo custom resources
//+kubebuilder:rbac:groups=tutorial.my.domain,resources=foos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tutorial.my.domain,resources=foos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tutorial.my.domain,resources=foos/finalizers,verbs=update

// RBAC permissions to monitor pods
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FooReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	/* */
	log := r.Log.WithValues("simpleDeployment", req.NamespacedName)

	var simpleDeployment aplhav1.Foo
	if err := r.Get(ctx, req.NamespacedName, &simpleDeployment); err != nil {
		log.Error(err, "unable to fetch SimpleDeployment")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// +kubebuilder:docs-gen:collapse=Begin the Reconcile

	/*
		Build the deployment that we want to see exist within the cluster
	*/

	deployment := &kapps.Deployment{}

	// Set the information you care about
	deployment.Spec.Replicas = simpleDeployment.Spec.Replicas

	/*
		Set the controller reference, specifying that this `Deployment` is controlled by the `SimpleDeployment` being reconciled.
		This will allow for the `SimpleDeployment` to be reconciled when changes to the `Deployment` are noticed.
	*/
	if err := controllerutil.SetControllerReference(simpleDeployment, deployment, r.scheme); err != nil {
		return ctrl.Result{}, err
	}

	/*
		Manage your `Deployment`.
		- Create it if it doesn't exist.
		- Update it if it is configured incorrectly.
	*/
	foundDeployment := &kapps.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		log.V(1).Info("Creating Deployment", "deployment", deployment.Name)
		err = r.Create(ctx, deployment)
	} else if err == nil {
		if foundDeployment.Spec.Replicas != deployment.Spec.Replicas {
			foundDeployment.Spec.Replicas = deployment.Spec.Replicas
			log.V(1).Info("Updating Deployment", "deployment", deployment.Name)
			err = r.Update(ctx, foundDeployment)
		}
	}

	return ctrl.Result{}, err
}

/*
Finally, we add this reconciler to the manager, so that it gets started
when the manager is started.
Since we create dependency `Deployments` during the reconcile, we can specify that the controller `Owns` `Deployments`.
This will tell the manager that if a `Deployment`, or its status, is updated, then the `SimpleDeployment` in its ownerRef field should be reconciled.
*/

// SetupWithManager sets up the controller with the Manager.
func (r *SimpleDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aplhav1.Foo{}).
		Owns(&kapps.Deployment{}).
		Complete(r)
}
