package permissions

import (
	"context"

	"github.com/argoproj-labs/argocd-operator/common"
	"github.com/argoproj-labs/argocd-operator/pkg/argoutil"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlClient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterRoleBindingRequest struct {
	Name              string
	InstanceName      string
	InstanceNamespace string
	Component         string
	Labels            map[string]string
	Annotations       map[string]string
	RoleRef           rbacv1.RoleRef
	Subjects          []rbacv1.Subject
}

// newClusterclusterRoleBinding returns a new clusterclusterRoleBinding instance.
func newClusterRoleBinding(name, instanceName, instanceNamespace, component string, labels, annotations map[string]string, roleRef rbacv1.RoleRef, subjects []rbacv1.Subject) *rbacv1.ClusterRoleBinding {
	crbName := argoutil.GenerateUniqueResourceName(instanceName, instanceNamespace, component)
	if name != "" {
		crbName = name
	}

	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        crbName,
			Labels:      argoutil.MergeMaps(common.DefaultLabels(crbName, instanceName, component), labels),
			Annotations: argoutil.MergeMaps(common.DefaultAnnotations(instanceName, instanceNamespace), annotations),
		},
		RoleRef:  roleRef,
		Subjects: subjects,
	}
}

// RequestClusterRoleBinding creates a ClusterRoleBinding object based on the provided ClusterRoleBindingRequest.
func RequestClusterRoleBinding(request ClusterRoleBindingRequest) *rbacv1.ClusterRoleBinding {
	return newClusterRoleBinding(request.Name, request.InstanceName, request.InstanceNamespace, request.Component, request.Labels, request.Annotations, request.RoleRef, request.Subjects)
}

// CreateClusterRoleBinding creates the specified ClusterRoleBinding using the provided client.
func CreateClusterRoleBinding(crb *rbacv1.ClusterRoleBinding, client ctrlClient.Client) error {
	return client.Create(context.TODO(), crb)
}

// GetClusterRoleBinding retrieves the ClusterRoleBinding with the given name using the provided client.
func GetClusterRoleBinding(name string, client ctrlClient.Client) (*rbacv1.ClusterRoleBinding, error) {
	existingCRB := &rbacv1.ClusterRoleBinding{}
	err := client.Get(context.TODO(), ctrlClient.ObjectKey{Name: name}, existingCRB)
	if err != nil {
		return nil, err
	}
	return existingCRB, nil
}

// ListClusterRoleBindings returns a list of ClusterRoleBinding objects using the provided client and list options.
func ListClusterRoleBindings(client ctrlClient.Client, listOptions []ctrlClient.ListOption) (*rbacv1.ClusterRoleBindingList, error) {
	existingCRBs := &rbacv1.ClusterRoleBindingList{}
	err := client.List(context.TODO(), existingCRBs, listOptions...)
	if err != nil {
		return nil, err
	}
	return existingCRBs, nil
}

// UpdateClusterRoleBinding updates the specified ClusterRoleBinding using the provided client.
func UpdateClusterRoleBinding(crb *rbacv1.ClusterRoleBinding, client ctrlClient.Client) error {
	_, err := GetClusterRoleBinding(crb.Name, client)
	if err != nil {
		return err
	}

	if err = client.Update(context.TODO(), crb); err != nil {
		return err
	}

	return nil
}

// DeleteClusterRoleBinding deletes the ClusterRoleBinding with the given name using the provided client.
// It ignores the "not found" error if the ClusterRoleBinding does not exist.
func DeleteClusterRoleBinding(name string, client ctrlClient.Client) error {
	existingCRB, err := GetClusterRoleBinding(name, client)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
		return nil
	}

	if err := client.Delete(context.TODO(), existingCRB); err != nil {
		return err
	}
	return nil
}
