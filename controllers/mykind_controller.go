/*
Copyright 2021.

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
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mygroupv1 "mydomain/myproject/api/v1"
)

// MyKindReconciler reconciles a MyKind object
type MyKindReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var (
	podOwner = ".metadata.controller"
	apiGVStr = mygroupv1.GroupVersion.String()
)

//+kubebuilder:rbac:groups=mygroup.mydomain,resources=mykinds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mygroup.mydomain,resources=mykinds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mygroup.mydomain,resources=mykinds/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyKind object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *MyKindReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mykind", req.NamespacedName)

	// your logic here
	var MyRsrc mygroupv1.MyKind
	if err := r.Get(ctx, req.NamespacedName, &MyRsrc); err != nil {
		log.Error(err, "unable to fetch MyKind")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	var childPods corev1.PodList
	if err := r.List(ctx, &childPods, client.InNamespace(req.Namespace), client.MatchingFields{podOwner: req.Name}); err != nil {
		log.Error(err, "unable to list child Jobs")
		return ctrl.Result{}, err
	}
	diff := MyRsrc.Spec.NrPods - len(childPods.Items)
	if diff == 0 {
		log.V(1).Info("Nothing to be done")
	} else if diff > 0 {
		log.V(1).Info("Create some pods")
		for i := 1; i <= diff; i++ {
			name := fmt.Sprintf("%s-pod-%v", MyRsrc.Name, i)
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      make(map[string]string),
					Annotations: make(map[string]string),
					Name:        name,
					Namespace:   MyRsrc.Namespace,
				},
				Spec: *MyRsrc.Spec.PodTemplate.Spec.DeepCopy(),
			}
			if err := ctrl.SetControllerReference(&MyRsrc, pod, r.Scheme); err != nil {
				log.Error(err, "unable to construct pod from template")
				return ctrl.Result{}, err
			}
			if err := r.Create(ctx, pod); err != nil {
				log.Error(err, "unable to create Pods for MyKind", "pod", pod)
				return ctrl.Result{}, err
			}
			log.V(1).Info("created Pods for MyKind", "pod", pod)
		}
	} else {
		log.V(1).Info("Delete some pods")
		deletions := 0
		for _, childPod := range childPods.Items {
			if deletions >= MyRsrc.Spec.NrPods {
				break
			}
			if err := r.Delete(ctx, &childPod, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				log.Error(err, "unable to delete pod", "pod", childPod)
				return ctrl.Result{}, err
			}
			deletions += 1
		}
	}
	MyRsrc.Status.MyStatus = "Reconciled"
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyKindReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Pod{}, podOwner, func(rawObj client.Object) []string {
		// grab the pod object, extract the owner...
		pod := rawObj.(*corev1.Pod)
		owner := metav1.GetControllerOf(pod)
		if owner == nil {
			return nil
		}
		// ...make sure it's MyKind...
		if owner.APIVersion != apiGVStr || owner.Kind != "MyKind" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&mygroupv1.MyKind{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
