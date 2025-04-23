/*
Copyright 2025.

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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	monitorv1alpha1 "github.com/you/http-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HTTPMonitorReconciler reconciles a HTTPMonitor object
type HTTPMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Helper function to trigger deployment restart
func triggerDeploymentRestart(c client.Client, namespace, name string) error {
	// Generate a timestamped annotation
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]string{
						"redeploy-timestamp": time.Now().Format(time.RFC3339),
					},
				},
			},
		},
	}

	// Encode it into JSON
	patchBytes, err := json.Marshal(patchData)
	if err != nil {
		return fmt.Errorf("failed to marshal patch data: %v", err)
	}

	// Apply patch to the deployment
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	return c.Patch(context.TODO(), deploy, client.RawPatch(types.StrategicMergePatchType, patchBytes))
}

var (
	lastRestartTime time.Time
	errorTimestamps []time.Time
)

// +kubebuilder:rbac:groups=monitor.mydomain.com,resources=httpmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitor.mydomain.com,resources=httpmonitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitor.mydomain.com,resources=httpmonitors/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;patch;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HTTPMonitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *HTTPMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Printf("Reconciling HTTPMonitor: %s\n", req.NamespacedName)

	// 1. Fetch the HTTPMonitor instance
	var monitor monitorv1alpha1.HTTPMonitor
	if err := r.Get(ctx, req.NamespacedName, &monitor); err != nil {
		log.Printf("Failed to get HTTPMonitor: %v", err)
		return ctrl.Result{}, nil
	}

	// 2. Parse CRD spec
	failureThreshold := monitor.Spec.FailureThreshold
	if failureThreshold == 0 {
		failureThreshold = 5
	}
	windowMinutes := monitor.Spec.WindowMinutes
	if windowMinutes == 0 {
		windowMinutes = 5
	}
	logFilePath := monitor.Spec.LogFilePath
	if logFilePath == "" {
		logFilePath = "/var/log/nginx/access.log"
	}
	deploymentTarget := monitor.Spec.DeploymentTarget

	// 3. Read log file
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("Failed to open log file %s: %v", logFilePath, err)
		return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
	}
	defer file.Close()

	now := time.Now()
	windowStart := now.Add(-time.Duration(windowMinutes) * time.Minute)

	// 4. Parse log lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 9 {
			status := parts[8]
			if strings.HasPrefix(status, "5") {
				errorTimestamps = append(errorTimestamps, now)
				log.Printf("[FAILURE] %s", line)
			} else if strings.HasPrefix(status, "2") {
				log.Printf("[SUCCESS] %s", line)
			}
		}
	}

	// 5. Prune old timestamps
	var recentErrors []time.Time
	for _, t := range errorTimestamps {
		if t.After(windowStart) {
			recentErrors = append(recentErrors, t)
		}
	}
	errorTimestamps = recentErrors
	log.Printf("500 errors in the last %d minutes: %d\n", windowMinutes, len(errorTimestamps))

	// 6. Trigger restart if threshold exceeded
	if len(errorTimestamps) >= failureThreshold && time.Since(lastRestartTime) > 10*time.Minute {
		err := triggerDeploymentRestart(r.Client, req.Namespace, deploymentTarget)
		if err != nil {
			log.Printf("Failed to restart deployment: %v", err)
		} else {
			log.Printf("Restart triggered for deployment " + deploymentTarget)
			lastRestartTime = time.Now()
		}
	}

	return ctrl.Result{RequeueAfter: 15 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitorv1alpha1.HTTPMonitor{}).
		Complete(r)
}
