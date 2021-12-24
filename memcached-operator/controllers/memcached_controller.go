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
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1 "memcached/api/v1"
)

const (
	RequeueTime = 30 * time.Second
)

// MemCachedReconciler reconciles a MemCached object
type MemCachedReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch

func (r *MemCachedReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("memcached", req.NamespacedName)

	log.Info("Reconciling memcached ...")
	memcache := &cachev1.MemCached{}
	if err := r.Get(ctx, req.NamespacedName, memcache); err != nil {
		if errors.IsNotFound(err) {
			// Return and don't requeue
			log.Info("Memcached resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Memcached.")
		return ctrl.Result{Requeue: true, RequeueAfter: RequeueTime}, err
	}

	deploy := &appsv1.Deployment{}
	// check if deployment already exits
	if err := r.Get(ctx, types.NamespacedName{Name: memcache.Name, Namespace: memcache.Namespace}, deploy); err != nil {
		if errors.IsNotFound(err) {
			deploy = createDesiredDeployment(memcache)
			log.Info("creating a new Deployment")
			if err := r.Create(ctx, deploy); err != nil {
				log.Info("failed to create a deployment: %v", err)
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else {
			log.Info("failed to get deployment")
			return ctrl.Result{Requeue: true, RequeueAfter: RequeueTime}, err
		}
	}

	log.Info("Succeddfully reconciled...")
	return ctrl.Result{}, nil
}
func (r *MemCachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1.MemCached{}).
		Complete(r)
}

func createDesiredDeployment(memcache *cachev1.MemCached) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcache.Name,
			Namespace: memcache.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &memcache.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": memcache.ObjectMeta.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      memcache.ObjectMeta.Name,
					Namespace: memcache.ObjectMeta.Namespace,
					Labels: map[string]string{
						"app": memcache.ObjectMeta.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "memcached:1.4.36-alpine",
						Name:    "memcached",
						Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 11211,
								Name:          "memcached",
							},
						},
					}},
				},
			},
		},
	}
}

// // serviceForMemcached function takes in a Memcached object and returns a Service for that object.
// func serviceForMemcached(m *cachev1.MemCached) *corev1.Service {
// 	ls := map[string]string{
// 		"app": m.ObjectMeta.Name,
// 	}
// 	ser := &corev1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      m.Name,
// 			Namespace: m.Namespace,
// 		},
// 		Spec: corev1.ServiceSpec{
// 			Selector: ls,
// 			Ports: []corev1.ServicePort{
// 				{
// 					Port: 11211,
// 					Name: m.Name,
// 				},
// 			},
// 		},
// 	}

// 	return ser
// }
