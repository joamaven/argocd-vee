package redis

import (
	"fmt"

	"github.com/argoproj-labs/argocd-operator/common"
	"github.com/argoproj-labs/argocd-operator/pkg/cluster"
	"github.com/argoproj-labs/argocd-operator/pkg/workloads"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	haproxyCfgKey             = "haproxy.cfg"
	haproxyScriptKey          = "haproxy_init.sh"
	initScriptKey             = "init.sh"
	redisConfKey              = "redis.conf"
	sentinelConfKey           = "sentinel.Conf"
	livenessScriptKey         = "redis_liveness.sh"
	readinessScriptKey        = "redis_readiness.sh"
	sentinelLivenessScriptKey = "sentinel_liveness.sh"
)

// reconcileHAConfigMap will ensure that the Redis HA ConfigMap is present for the given ArgoCD instance
func (rr *RedisReconciler) reconcileHAConfigMap() error {
	cmRequest := workloads.ConfigMapRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.ArgoCDRedisHAConfigMapName,
			Namespace:   rr.Instance.Namespace,
			Labels:      resourceLabels,
			Annotations: rr.Instance.Annotations,
		},
		Data: map[string]string{
			haproxyCfgKey:    rr.GetHAProxyConfig(),
			haproxyScriptKey: rr.GetHAProxyScript(),
			initScriptKey:    rr.GetInitScript(),
			redisConfKey:     rr.GetConf(),
			sentinelConfKey:  rr.GetSentinelConf(),
		},
	}

	desiredCM, err := workloads.RequestConfigMap(cmRequest)
	if err != nil {
		rr.Logger.V(1).Info("reconcileHAConfigMap: one or more mutations could not be applied")
		return errors.Wrap(err, fmt.Sprintf("reconcileHAConfigMap: failed to request configMap %s in namespace %s", desiredCM.Name, desiredCM.Namespace))
	}

	namespace, err := cluster.GetNamespace(rr.Instance.Namespace, rr.Client)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("reconcileHAConfigMap: failed to retrieve namespace %s", rr.Instance.Namespace))
	}
	if namespace.DeletionTimestamp != nil {
		if err := rr.deleteConfigMap(desiredCM.Name, desiredCM.Namespace); err != nil {
			return errors.Wrap(err, fmt.Sprintf("reconcileHAConfigMap: failed to delete configMap %s in namespace %s", desiredCM.Name, desiredCM.Namespace))
		}
	}

	_, err = workloads.GetConfigMap(desiredCM.Name, desiredCM.Namespace, rr.Client)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrap(err, fmt.Sprintf("reconcileHAConfigMap: failed to retrieve configMap %s in namespace %s", desiredCM.Name, desiredCM.Namespace))
		}

		if err = controllerutil.SetControllerReference(rr.Instance, desiredCM, rr.Scheme); err != nil {
			rr.Logger.Error(err, "reconcileHAConfigMap: failed to set owner reference for configMap", "name", desiredCM.Name, "namespace", desiredCM.Namespace)
		}

		if err = workloads.CreateConfigMap(desiredCM, rr.Client); err != nil {
			return errors.Wrap(err, fmt.Sprintf("reconcileHAConfigMap: failed to create configMap %s in namespace %s", desiredCM.Name, desiredCM.Namespace))
		}
		rr.Logger.V(0).Info("reconcileHAConfigMap: configMap created", "name", desiredCM.Name, "namespace", desiredCM.Namespace)
		return nil
	}

	return nil
}

// reconcileHAHealthConfigMap will ensure that the Redis HA Health ConfigMap is present for the given ArgoCD.
func (rr *RedisReconciler) reconcileHAHealthConfigMap() error {
	cmRequest := workloads.ConfigMapRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:        common.ArgoCDRedisHAHealthConfigMapName,
			Namespace:   rr.Instance.Namespace,
			Labels:      resourceLabels,
			Annotations: rr.Instance.Annotations,
		},
		Data: map[string]string{
			livenessScriptKey:         rr.GetLivenessScript(),
			readinessScriptKey:        rr.GetReadinessScript(),
			sentinelLivenessScriptKey: rr.GetSentinelLivenessScript(),
		},
	}

	desiredCM, err := workloads.RequestConfigMap(cmRequest)
	if err != nil {
		rr.Logger.Error(err, "reconcileHAHealthConfigMap: failed to request configMap", "name", desiredCM.Name, "namespace", desiredCM.Namespace)
		rr.Logger.V(1).Info("reconcileHAHealthConfigMap: one or more mutations could not be applied")
		return err
	}

	namespace, err := cluster.GetNamespace(rr.Instance.Namespace, rr.Client)
	if err != nil {
		rr.Logger.Error(err, "reconcileHAHealthConfigMap: failed to retrieve namespace", "name", rr.Instance.Namespace)
		return err
	}
	if namespace.DeletionTimestamp != nil {
		if err := rr.deleteConfigMap(desiredCM.Name, desiredCM.Namespace); err != nil {
			return errors.Wrap(err, fmt.Sprintf("reconcileHAHealthConfigMap: failed to delete configMap %s in namespace %s", desiredCM.Name, desiredCM.Namespace))
		}
	}

	_, err = workloads.GetConfigMap(desiredCM.Name, desiredCM.Namespace, rr.Client)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrap(err, fmt.Sprintf("reconcileHAHealthConfigMap: failed to retrieve configMap %s in namespace %s", desiredCM.Name, desiredCM.Namespace))
		}

		if err = controllerutil.SetControllerReference(rr.Instance, desiredCM, rr.Scheme); err != nil {
			rr.Logger.Error(err, "reconcileHAHealthConfigMap: failed to set owner reference for configMap", "name", desiredCM.Name, "namespace", desiredCM.Namespace)
		}

		if err = workloads.CreateConfigMap(desiredCM, rr.Client); err != nil {
			return errors.Wrap(err, fmt.Sprintf("reconcileHAHealthConfigMap: failed to create configMap %s in namespace %s", desiredCM.Name, desiredCM.Namespace))
		}
		rr.Logger.V(0).Info("reconcileHAHealthConfigMap: configMap created", "name", desiredCM.Name, "namespace", desiredCM.Namespace)
		return nil
	}

	return nil
}

func (rr *RedisReconciler) deleteConfigMap(name, namespace string) error {
	if err := workloads.DeleteConfigMap(name, namespace, rr.Client); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return errors.Wrap(err, fmt.Sprintf("deleteConfigMap: failed to delete configMap %s in namespace %s", name, namespace))
	}
	rr.Logger.V(0).Info("DeleteConfigMap: configMap deleted", "name", name, "namespace", namespace)
	return nil
}
