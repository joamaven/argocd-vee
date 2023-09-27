// Copyright 2019 ArgoCD Operator Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package argocd

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/argoproj-labs/argocd-operator/api/v1beta1"
	argoproj "github.com/argoproj-labs/argocd-operator/api/v1beta1"
	"github.com/argoproj-labs/argocd-operator/common"
	"github.com/argoproj-labs/argocd-operator/pkg/util"
)

// getArgoServerServiceType will return the server Service type for the ArgoCD.
func getArgoServerServiceType(cr *argoproj.ArgoCD) corev1.ServiceType {
	if len(cr.Spec.Server.Service.Type) > 0 {
		return cr.Spec.Server.Service.Type
	}
	return corev1.ServiceTypeClusterIP
}

// newService returns a new Service for the given ArgoCD instance.
func newService(cr *argoproj.ArgoCD) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    common.DefaultLabels(cr.Name, cr.Name, ""),
		},
	}
}

// newServiceWithName returns a new Service instance for the given ArgoCD using the given name.
func newServiceWithName(name string, component string, cr *argoproj.ArgoCD) *corev1.Service {
	svc := newService(cr)
	svc.ObjectMeta.Name = name

	lbls := svc.ObjectMeta.Labels
	lbls[common.AppK8sKeyName] = name
	lbls[common.AppK8sKeyComponent] = component
	svc.ObjectMeta.Labels = lbls

	return svc
}

// newServiceWithSuffix returns a new Service instance for the given ArgoCD using the given suffix.
func newServiceWithSuffix(suffix string, component string, cr *argoproj.ArgoCD) *corev1.Service {
	return newServiceWithName(fmt.Sprintf("%s-%s", cr.Name, suffix), component, cr)
}

// reconcileGrafanaService will ensure that the Service for Grafana is present.
func (r *ArgoCDReconciler) reconcileGrafanaService(cr *v1beta1.ArgoCD) error {
	svc := newServiceWithSuffix("grafana", "grafana", cr)
	if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
		if !cr.Spec.Grafana.Enabled {
			// Service exists but enabled flag has been set to false, delete the Service
			return r.Client.Delete(context.TODO(), svc)
		}
		return nil // Service found, do nothing
	}

	if !cr.Spec.Grafana.Enabled {
		return nil // Grafana not enabled, do nothing.
	}

	svc.Spec.Selector = map[string]string{
		common.AppK8sKeyName: util.NameWithSuffix(cr.Name, "grafana"),
	}

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       80,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(3000),
		},
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("creating service %s for Argo CD instance %s in namespace %s", svc.Name, cr.Name, cr.Namespace))
	return r.Client.Create(context.TODO(), svc)
}

// reconcileMetricsService will ensure that the Service for the Argo CD application controller metrics is present.
func (r *ArgoCDReconciler) reconcileMetricsService(cr *v1beta1.ArgoCD) error {
	svc := newServiceWithSuffix("metrics", "metrics", cr)
	if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
		// Service found, do nothing
		return nil
	}

	svc.Spec.Selector = map[string]string{
		common.AppK8sKeyName: util.NameWithSuffix(cr.Name, "application-controller"),
	}

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "metrics",
			Port:       8082,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(8082),
		},
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}
	return r.Client.Create(context.TODO(), svc)
}

// reconcileRedisHAAnnounceServices will ensure that the announce Services are present for Redis when running in HA mode.
func (r *ArgoCDReconciler) reconcileRedisHAAnnounceServices(cr *v1beta1.ArgoCD) error {
	for i := int32(0); i < common.ArgoCDDefaultRedisHAReplicas; i++ {
		svc := newServiceWithSuffix(fmt.Sprintf("redis-ha-announce-%d", i), "redis", cr)
		if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
			if !cr.Spec.HA.Enabled {
				return r.Client.Delete(context.TODO(), svc)
			}
			return nil // Service found, do nothing
		}

		if !cr.Spec.HA.Enabled {
			return nil //return as Ha is not enabled do nothing
		}

		svc.ObjectMeta.Annotations = map[string]string{
			common.ServiceAlphaK8sKeyTolerateUnreadyEndpoints: "true",
		}

		svc.Spec.PublishNotReadyAddresses = true

		svc.Spec.Selector = map[string]string{
			common.AppK8sKeyName:            util.NameWithSuffix(cr.Name, "redis-ha"),
			common.StatefulSetK8sKeyPodName: util.NameWithSuffix(cr.Name, fmt.Sprintf("redis-ha-server-%d", i)),
		}

		svc.Spec.Ports = []corev1.ServicePort{
			{
				Name:       "server",
				Port:       common.ArgoCDDefaultRedisPort,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromString("redis"),
			}, {
				Name:       "sentinel",
				Port:       common.ArgoCDDefaultRedisSentinelPort,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromString("sentinel"),
			},
		}

		if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
			return err
		}

		if err := r.Client.Create(context.TODO(), svc); err != nil {
			return err
		}
	}
	return nil
}

// reconcileRedisHAMasterService will ensure that the "master" Service is present for Redis when running in HA mode.
func (r *ArgoCDReconciler) reconcileRedisHAMasterService(cr *v1beta1.ArgoCD) error {
	svc := newServiceWithSuffix("redis-ha", "redis", cr)
	if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
		if !cr.Spec.HA.Enabled {
			return r.Client.Delete(context.TODO(), svc)
		}
		return nil // Service found, do nothing
	}

	if !cr.Spec.HA.Enabled {
		return nil //return as Ha is not enabled do nothing
	}

	svc.Spec.Selector = map[string]string{
		common.AppK8sKeyName: util.NameWithSuffix(cr.Name, "redis-ha"),
	}

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "server",
			Port:       common.ArgoCDDefaultRedisPort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromString("redis"),
		}, {
			Name:       "sentinel",
			Port:       common.ArgoCDDefaultRedisSentinelPort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromString("sentinel"),
		},
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}
	return r.Client.Create(context.TODO(), svc)
}

// reconcileRedisHAProxyService will ensure that the HA Proxy Service is present for Redis when running in HA mode.
func (r *ArgoCDReconciler) reconcileRedisHAProxyService(cr *v1beta1.ArgoCD) error {
	svc := newServiceWithSuffix("redis-ha-haproxy", "redis", cr)
	if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {

		if !cr.Spec.HA.Enabled {
			return r.Client.Delete(context.TODO(), svc)
		}

		if ensureAutoTLSAnnotation(svc, common.ArgoCDRedisServerTLSSecretName, cr.Spec.Redis.WantsAutoTLS()) {
			return r.Client.Update(context.TODO(), svc)
		}
		return nil // Service found, do nothing
	}

	if !cr.Spec.HA.Enabled {
		return nil //return as Ha is not enabled do nothing
	}

	ensureAutoTLSAnnotation(svc, common.ArgoCDRedisServerTLSSecretName, cr.Spec.Redis.WantsAutoTLS())

	svc.Spec.Selector = map[string]string{
		common.AppK8sKeyName: util.NameWithSuffix(cr.Name, "redis-ha-haproxy"),
	}

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "haproxy",
			Port:       common.ArgoCDDefaultRedisPort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromString("redis"),
		},
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}
	return r.Client.Create(context.TODO(), svc)
}

// reconcileRedisHAServices will ensure that all required Services are present for Redis when running in HA mode.
func (r *ArgoCDReconciler) reconcileRedisHAServices(cr *v1beta1.ArgoCD) error {

	if err := r.reconcileRedisHAAnnounceServices(cr); err != nil {
		return err
	}

	if err := r.reconcileRedisHAMasterService(cr); err != nil {
		return err
	}

	if err := r.reconcileRedisHAProxyService(cr); err != nil {
		return err
	}
	return nil
}

// reconcileRedisService will ensure that the Service for Redis is present.
func (r *ArgoCDReconciler) reconcileRedisService(cr *v1beta1.ArgoCD) error {
	svc := newServiceWithSuffix("redis", "redis", cr)

	if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
		if ensureAutoTLSAnnotation(svc, common.ArgoCDRedisServerTLSSecretName, cr.Spec.Redis.WantsAutoTLS()) {
			return r.Client.Update(context.TODO(), svc)
		}
		if cr.Spec.HA.Enabled {
			return r.Client.Delete(context.TODO(), svc)
		}
		return nil // Service found, do nothing
	}

	if cr.Spec.HA.Enabled {
		return nil //return as Ha is enabled do nothing
	}

	ensureAutoTLSAnnotation(svc, common.ArgoCDRedisServerTLSSecretName, cr.Spec.Redis.WantsAutoTLS())

	svc.Spec.Selector = map[string]string{
		common.AppK8sKeyName: util.NameWithSuffix(cr.Name, "redis"),
	}

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "tcp-redis",
			Port:       common.ArgoCDDefaultRedisPort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(common.ArgoCDDefaultRedisPort),
		},
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}
	return r.Client.Create(context.TODO(), svc)
}

// ensureAutoTLSAnnotation ensures that the service svc has the desired state
// of the auto TLS annotation set, which is either set (when enabled is true)
// or unset (when enabled is false).
//
// Returns true when annotations have been updated, otherwise returns false.
//
// When this method returns true, the svc resource will need to be updated on
// the cluster.
func ensureAutoTLSAnnotation(svc *corev1.Service, secretName string, enabled bool) bool {
	var autoTLSAnnotationName, autoTLSAnnotationValue string

	// We currently only support OpenShift for automatic TLS
	if IsRouteAPIAvailable() {
		autoTLSAnnotationName = common.ServiceBetaOpenshiftKeyCertSecret
		if svc.Annotations == nil {
			svc.Annotations = make(map[string]string)
		}
		autoTLSAnnotationValue = secretName
	}

	if autoTLSAnnotationName != "" {
		val, ok := svc.Annotations[autoTLSAnnotationName]
		if enabled {
			if !ok || val != secretName {
				log.Info(fmt.Sprintf("requesting AutoTLS on service %s", svc.ObjectMeta.Name))
				svc.Annotations[autoTLSAnnotationName] = autoTLSAnnotationValue
				return true
			}
		} else {
			if ok {
				log.Info(fmt.Sprintf("removing AutoTLS from service %s", svc.ObjectMeta.Name))
				delete(svc.Annotations, autoTLSAnnotationName)
				return true
			}
		}
	}

	return false
}

// reconcileRepoService will ensure that the Service for the Argo CD repo server is present.
func (r *ArgoCDReconciler) reconcileRepoService(cr *v1beta1.ArgoCD) error {
	svc := newServiceWithSuffix("repo-server", "repo-server", cr)

	if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
		if ensureAutoTLSAnnotation(svc, common.ArgoCDRepoServerTLSSecretName, cr.Spec.Repo.WantsAutoTLS()) {
			return r.Client.Update(context.TODO(), svc)
		}
		return nil // Service found, do nothing
	}

	ensureAutoTLSAnnotation(svc, common.ArgoCDRepoServerTLSSecretName, cr.Spec.Repo.WantsAutoTLS())

	svc.Spec.Selector = map[string]string{
		common.AppK8sKeyName: util.NameWithSuffix(cr.Name, "repo-server"),
	}

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "server",
			Port:       common.ArgoCDDefaultRepoServerPort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(common.ArgoCDDefaultRepoServerPort),
		}, {
			Name:       "metrics",
			Port:       common.ArgoCDDefaultRepoMetricsPort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(common.ArgoCDDefaultRepoMetricsPort),
		},
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}
	return r.Client.Create(context.TODO(), svc)
}

// reconcileServerMetricsService will ensure that the Service for the Argo CD server metrics is present.
func (r *ArgoCDReconciler) reconcileServerMetricsService(cr *v1beta1.ArgoCD) error {
	svc := newServiceWithSuffix("server-metrics", "server", cr)
	if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
		return nil // Service found, do nothing
	}

	svc.Spec.Selector = map[string]string{
		common.AppK8sKeyName: util.NameWithSuffix(cr.Name, "server"),
	}

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "metrics",
			Port:       8083,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(8083),
		},
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}
	return r.Client.Create(context.TODO(), svc)
}

// reconcileServerService will ensure that the Service is present for the Argo CD server component.
func (r *ArgoCDReconciler) reconcileServerService(cr *v1beta1.ArgoCD) error {
	svc := newServiceWithSuffix("server", "server", cr)
	if util.IsObjectFound(r.Client, cr.Namespace, svc.Name, svc) {
		if ensureAutoTLSAnnotation(svc, common.ArgoCDServerTLSSecretName, cr.Spec.Server.WantsAutoTLS()) {
			return r.Client.Update(context.TODO(), svc)
		}
		return nil // Service found, do nothing
	}

	ensureAutoTLSAnnotation(svc, common.ArgoCDServerTLSSecretName, cr.Spec.Server.WantsAutoTLS())

	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       80,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(8080),
		}, {
			Name:       "https",
			Port:       443,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(8080),
		},
	}

	svc.Spec.Selector = map[string]string{
		common.AppK8sKeyName: util.NameWithSuffix(cr.Name, "server"),
	}

	svc.Spec.Type = getArgoServerServiceType(cr)

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		return err
	}
	return r.Client.Create(context.TODO(), svc)
}

// reconcileServices will ensure that all Services are present for the given ArgoCD.
func (r *ArgoCDReconciler) reconcileServices(cr *v1beta1.ArgoCD) error {

	if err := r.reconcileDexService(cr); err != nil {
		log.Error(err, "error reconciling dex service")
	}

	err := r.reconcileGrafanaService(cr)
	if err != nil {
		return err
	}

	err = r.reconcileMetricsService(cr)
	if err != nil {
		return err
	}

	err = r.reconcileRedisHAServices(cr)
	if err != nil {
		return err
	}

	err = r.reconcileRedisService(cr)
	if err != nil {
		return err
	}

	err = r.reconcileRepoService(cr)
	if err != nil {
		return err
	}

	err = r.reconcileServerMetricsService(cr)
	if err != nil {
		return err
	}

	err = r.reconcileServerService(cr)
	if err != nil {
		return err
	}
	return nil
}
