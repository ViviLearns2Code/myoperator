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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var mykindlog = logf.Log.WithName("mykind-resource")

func (r *MyKind) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-mygroup-mydomain-v1-mykind,mutating=true,failurePolicy=fail,sideEffects=None,groups=mygroup.mydomain,resources=mykinds,verbs=create;update,versions=v1,name=mmykind.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &MyKind{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MyKind) Default() {
	mykindlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	r.Status.MyStatus = "DefaultValue"
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-mygroup-mydomain-v1-mykind,mutating=false,failurePolicy=fail,sideEffects=None,groups=mygroup.mydomain,resources=mykinds,verbs=create;update,versions=v1,name=vmykind.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &MyKind{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MyKind) ValidateCreate() error {
	mykindlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateRsrc()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MyKind) ValidateUpdate(old runtime.Object) error {
	mykindlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateRsrc()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MyKind) ValidateDelete() error {
	mykindlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *MyKind) validateRsrc() error {
	var allErrs field.ErrorList
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-7 {
		// The k8s name length is 63 character
		// the controller appends 7 characters "name-pod-xx" for pod names
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 56 characters"))
	}
	if r.Spec.NrPods > 10 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("nrPods"), r.Spec.NrPods, "must not exceed 10"))
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "mygroup.mydomain", Kind: "MyKind"},
		r.Name, allErrs)
}
