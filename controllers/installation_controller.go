/*
Copyright 2020 cs3238-tsuzu.

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
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/xerrors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ghappv1alpha1 "github.com/cs3238-tsuzu/ghapp-controller/api/v1alpha1"
	"github.com/cs3238-tsuzu/ghapp-controller/pkg/installations"
)

// InstallationReconciler reconciles a Installation object
type InstallationReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	RefreshBefore time.Duration
}

// +kubebuilder:rbac:groups=ghapp.tsuzu.dev,resources=installations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ghapp.tsuzu.dev,resources=installations/status,verbs=get;update;patch

func (r *InstallationReconciler) Reconcile(req ctrl.Request) (res ctrl.Result, err error) {
	ctx := context.Background()
	log := r.Log.WithValues("installation", req.NamespacedName)

	ins := &ghappv1alpha1.Installation{}
	if err := r.Get(ctx, req.NamespacedName, ins); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	defer func() {
		if err != nil {
			ins.Status.Ready = false
			ins.Status.Message = err.Error()

			log.Error(err, "failed to update status")
			if err := r.Client.Status().Update(ctx, ins); err != nil {
				log.Error(err, "failed to update status", "status", ins.Status.Message)

				res.Requeue = true
			}
		}
	}()

	insclient := installations.Client{
		Client:        r.Client,
		Installation:  ins,
		RefreshBefore: r.RefreshBefore,
	}

	status, generated, expiredAt, err := insclient.Run(ctx)

	if err != nil {
		return ctrl.Result{
			Requeue: true,
		}, xerrors.Errorf("failed to check the current secret: %w", err)
	}

	switch status {
	case installations.Undesirable:
		if err := r.Update(ctx, generated); err != nil {
			return ctrl.Result{
				Requeue: true,
			}, xerrors.Errorf("failed to update status: %w", err)
		}
	case installations.NotExisting:
		if err := r.Create(ctx, generated); err != nil {
			return ctrl.Result{
				Requeue: true,
			}, xerrors.Errorf("failed to update status: %w", err)
		}
	case installations.Desired:
		// do nothing
	}

	ins.Status.Ready = true
	ins.Status.Secret = generated.Name
	ins.Status.Secret = ins.Name

	if err := r.Client.Status().Update(ctx, ins); err != nil {
		return ctrl.Result{
			Requeue: true,
		}, xerrors.Errorf("failed to update status: %w", err)
	}

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: expiredAt.Add(-r.RefreshBefore).Sub(time.Now()),
	}, nil
}

func (r *InstallationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ghappv1alpha1.Installation{}).
		Complete(r)
}
