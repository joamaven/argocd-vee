// Copyright 2021 ArgoCD Operator Developers
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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	argoprojv1a1 "github.com/argoproj-labs/argocd-operator/api/v1alpha1"
	"github.com/argoproj-labs/argocd-operator/common"
)

// newServiceWithName returns a new Service instance for the given ArgoCD using the given name.
func newServiceWithName(name string, component string, cr *argoprojv1a1.ArgoCD) *corev1.Service {
	svc := newService(cr)
	svc.ObjectMeta.Name = name

	lbls := svc.ObjectMeta.Labels
	lbls[common.ArgoCDKeyName] = name
	lbls[common.ArgoCDKeyComponent] = component
	svc.ObjectMeta.Labels = lbls

	return svc
}

// newService returns a new Service for the given ArgoCD instance.
func newService(cr *argoprojv1a1.ArgoCD) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labelsForCluster(cr),
		},
	}
}

// NewServiceWithSuffix returns a new Service instance for the given ArgoCD using the given suffix.
func NewServiceWithSuffix(suffix string, component string, cr *argoprojv1a1.ArgoCD) *corev1.Service {
	return newServiceWithName(fmt.Sprintf("%s-%s", cr.Name, suffix), component, cr)
}

// EnsureAutoTLSAnnotation applies TLS annotations
func EnsureAutoTLSAnnotation(cr *argoprojv1a1.ArgoCD, svc *corev1.Service) bool {
	autoTLSAnnotationName := ""
	if cr.Spec.Repo.AutoTLS == "openshift" {
		autoTLSAnnotationName = "service.beta.openshift.io/serving-cert-secret-name"
	}
	if autoTLSAnnotationName != "" {
		if svc.Annotations == nil {
			svc.Annotations = make(map[string]string)
		}
		val, ok := svc.Annotations[autoTLSAnnotationName]
		if !ok || val != common.ArgoCDRepoServerTLSSecretName {
			logr.Info(fmt.Sprintf("requesting AutoTLS on service %s", svc.ObjectMeta.Name))
			svc.Annotations[autoTLSAnnotationName] = common.ArgoCDRepoServerTLSSecretName
			return true
		}
	}

	return false
}

// GetArgoServerServiceType will return the server Service type for the ArgoCD.
func GetArgoServerServiceType(cr *argoprojv1a1.ArgoCD) corev1.ServiceType {
	if len(cr.Spec.Server.Service.Type) > 0 {
		return cr.Spec.Server.Service.Type
	}
	return corev1.ServiceTypeClusterIP
}
