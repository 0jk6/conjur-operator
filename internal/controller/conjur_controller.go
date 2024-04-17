/*
Copyright 2024 0jk6.

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
	"encoding/base64"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"conjur-operator.0jk6.github.io/api/v1alpha1"
	conjurapi "conjur-operator.0jk6.github.io/internal/utils"
)

// ConjurReconciler reconciles a Conjur object
type ConjurReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.0jk6.github.io,resources=conjurs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.0jk6.github.io,resources=conjurs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.0jk6.github.io,resources=conjurs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Conjur object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *ConjurReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	//let's try to pull our custom resource
	conjur := &v1alpha1.Conjur{}
	log.Log.Info("Reconciling custom resource", "name", req.Name)
	refreshInterval := 60

	err := r.Get(ctx, req.NamespacedName, conjur)

	if err != nil {
		log.Log.Error(err, "unable to fetch custom resource")

		//make sure that the refresh interval is not less than 60 seconds
		if conjur.Spec.RefreshInterval < 60 {
			refreshInterval = 60
		} else {
			refreshInterval = conjur.Spec.RefreshInterval
		}

		if errors.IsNotFound(err) {
			return ctrl.Result{RequeueAfter: time.Duration(refreshInterval) * time.Second}, nil
		}
		return ctrl.Result{RequeueAfter: time.Duration(refreshInterval) * time.Second}, client.IgnoreNotFound(err)
	} else {
		//if we found our custom resource, let's create a project
		for secretName, secretToPull := range conjur.Spec.Data {
			secret := &corev1.Secret{}

			//pull the required data from the applied kind.
			conjurAcct := conjur.Spec.ConjurAcct
			conjurHost := conjur.Spec.ConjurHost
			hostname := conjur.Spec.Hostname
			secretIdentifier := secretToPull.SecretIdentifier

			//pull the conjur api key from the applied secret
			apiKey := ""
			apiKeyFromSecret := conjur.Spec.ApiKeyFromSecret

			//first pull the data from "apiKeyFromSecret"
			if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: apiKeyFromSecret}, secret); err != nil {
				log.Log.Info("API Key Secret", apiKeyFromSecret, "not found")
			} else {
				log.Log.Info("pulling api key from existing secret", "secret", apiKeyFromSecret)
				decodedApiKey, _ := base64.StdEncoding.DecodeString(string(secret.Data["apikey"]))
				apiKey = strings.ReplaceAll(string(decodedApiKey), "\n", "")
			}

			clientObject := client.ObjectKey{Namespace: req.Namespace, Name: secretName}

			if err := r.Get(ctx, clientObject, secret); err != nil {
				//create the secret
				log.Log.Info("Secret", secretName, "not found")
				log.Log.Info("Creating secret", secretName, secretIdentifier)

				secretFromConjur := conjurapi.PullSecret(conjurHost, conjurAcct, hostname, secretIdentifier, apiKey)

				secret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      secretName,
						Namespace: req.Namespace,
					},
					StringData: map[string]string{
						"data": secretFromConjur,
					},
				}

				if err := r.Create(ctx, secret); err != nil {
					log.Log.Error(err, "unable to create secret")
				} else {
					log.Log.Info("Secret created", "name", secret.Name)
				}

			} else {
				//if the secret is found, sync it with the vault
				log.Log.Info("Secret found", "name", secretName)

				secretFromConjur := conjurapi.PullSecret(conjurHost, conjurAcct, hostname, secretIdentifier, apiKey)
				log.Log.Info("Secret syncing", "name", secretName)
				secret.Data["data"] = []byte(secretFromConjur)

				if err := r.Update(ctx, secret); err != nil {
					log.Log.Error(err, "unable to update secret")
				} else {
					log.Log.Info("Secret synced", "name", secretName)
				}

			}
		}
	}

	if conjur.Spec.RefreshInterval < 60 {
		refreshInterval = 60
	} else {
		refreshInterval = conjur.Spec.RefreshInterval
	}

	return ctrl.Result{RequeueAfter: time.Duration(refreshInterval) * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConjurReconciler) SetupWithManager(mgr ctrl.Manager) error {
	leaderElection := new(bool)
	*leaderElection = false

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Conjur{}).
		Complete(r)
}
