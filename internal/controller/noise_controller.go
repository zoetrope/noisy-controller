package controller

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	noisyv1 "github.com/zoetrope/noizy-controller/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NoiseReconciler reconciles a Noise object
type NoiseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=noisy.zoetrope.github.io,resources=noises,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=noisy.zoetrope.github.io,resources=noises/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=noisy.zoetrope.github.io,resources=noises/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Noise object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *NoiseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var noise noisyv1.Noise
	err := r.Get(ctx, req.NamespacedName, &noise)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !noise.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	err = r.reconcileConfigMaps(ctx, &noise)
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(noise.Spec.Interval) == 0 {
		return ctrl.Result{}, nil
	}

	interval, err := time.ParseDuration(noise.Spec.Interval)
	if err != nil {
		return ctrl.Result{}, err
	}
	logger.Info("requeue after", "interval", interval.String())
	return ctrl.Result{
		RequeueAfter: interval,
	}, nil
}

func (r *NoiseReconciler) reconcileConfigMaps(ctx context.Context, noise *noisyv1.Noise) error {
	logger := log.FromContext(ctx)

	for i := 0; i < noise.Spec.Count; i++ {
		name := fmt.Sprintf("configmap-%s-%d", noise.Name, i)

		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: noise.Namespace,
				Name:      name,
			},
		}
		op, err := ctrl.CreateOrUpdate(ctx, r.Client, cm, func() error {
			cm.Data = map[string]string{
				"value": randomString(16),
			}
			return ctrl.SetControllerReference(noise, cm, r.Scheme)
		})

		if err != nil {
			return err
		}

		logger.Info("update configmap successfully", "configmap", name, "op", op)
	}
	return nil
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// SetupWithManager sets up the controller with the Manager.
func (r *NoiseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&noisyv1.Noise{}).
		//Owns(&corev1.ConfigMap{}).
		Complete(r)
}
