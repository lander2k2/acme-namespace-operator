/*


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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tenancyv1alpha1 "github.com/lander2k2/acme-namespace-operator/api/v1alpha1"
)

const (
	statusCreated    = "Created"
	statusInProgress = "CreationInProgress"
)

// AcmeNamespaceReconciler reconciles a AcmeNamespace object
type AcmeNamespaceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tenancy.acme.com,resources=acmenamespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tenancy.acme.com,resources=acmenamespaces/status,verbs=get;update;patch

func (r *AcmeNamespaceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("acmenamespace", req.NamespacedName)

	var acmeNs tenancyv1alpha1.AcmeNamespace
	if err := r.Get(ctx, req.NamespacedName, &acmeNs); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("resource deleted")
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	nsName := acmeNs.Spec.NamespaceName
	adminUsername := acmeNs.Spec.AdminUsername

	switch acmeNs.Status.Phase {
	case statusCreated:
		// do nothing
		log.Info("AcmeNamespace child resources have been created")
	case statusInProgress:
		// TODO: query and create as needed
		log.Info("AcmeNamespace child resource creation in progress")
	default:
		log.Info("AcmeNamespace child resources not created")

		// set status to statusInProgress
		acmeNs.Status.Phase = statusInProgress
		if err := r.Status().Update(ctx, &acmeNs); err != nil {
			log.Error(err, "unable to update AcmeNamespace status")
		}

		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: nsName,
				Labels: map[string]string{
					"admin": adminUsername,
				},
			},
		}

		// set owner reference for the namespace
		err := ctrl.SetControllerReference(&acmeNs, ns, r.Scheme)
		if err != nil {
			log.Error(err, "unable to set owner reference on namespace")
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, ns); err != nil {
			log.Error(err, "unable to create namespace")
			return ctrl.Result{}, err
		}

		// create role
		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "namespace-admin",
				Namespace: nsName,
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{"*"},
					Resources: []string{"*"},
					Verbs:     []string{"*"},
				},
			},
		}

		// set owner reference for the role
		err = ctrl.SetControllerReference(&acmeNs, role, r.Scheme)
		if err != nil {
			log.Error(err, "unable to set owner reference on role")
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, role); err != nil {
			log.Error(err, "unable to create namespace-admin role")
			return ctrl.Result{}, err
		}

		// create role binding
		binding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "namespace-admin",
				Namespace: nsName,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     "namespace-admin",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "User",
					Name:      adminUsername,
					Namespace: nsName,
				},
			},
		}

		// set owner reference for the role binding
		err = ctrl.SetControllerReference(&acmeNs, binding, r.Scheme)
		if err != nil {
			log.Error(err, "unable to set reference on role binding")
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, binding); err != nil {
			log.Error(err, "unable to create namespace-admin role binding")
			return ctrl.Result{}, err
		}

		// set status to statusCreated
		acmeNs.Status.Phase = statusCreated
		if err := r.Status().Update(ctx, &acmeNs); err != nil {
			log.Error(err, "unable to update AcmeNamespace status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *AcmeNamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenancyv1alpha1.AcmeNamespace{}).
		Complete(r)
}
