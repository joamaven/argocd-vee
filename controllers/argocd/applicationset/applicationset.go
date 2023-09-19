package applicationset

import (
	"github.com/argoproj-labs/argocd-operator/api/v1beta1"
	"github.com/argoproj-labs/argocd-operator/common"
	"github.com/argoproj-labs/argocd-operator/pkg/util"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplicationSetReconciler struct {
	Client   client.Client
	Scheme   *runtime.Scheme
	Instance *v1beta1.ArgoCD
	Logger   logr.Logger
}

var (
	resourceName   string
	resourceLabels map[string]string
)

func (asr *ApplicationSetReconciler) Reconcile() error {

	asr.Logger = ctrl.Log.WithName(ArgoCDApplicationSetControllerComponent).WithValues("instance", asr.Instance.Name, "instance-namespace", asr.Instance.Namespace)

	resourceName = util.GenerateUniqueResourceName(asr.Instance.Name, asr.Instance.Namespace, ArgoCDApplicationSetControllerComponent)
	resourceLabels = common.DefaultLabels(resourceName, asr.Instance.Name, ArgoCDApplicationSetControllerComponent)

	if err := asr.reconcileServiceAccount(); err != nil {
		asr.Logger.Info("reconciling applicationSet serviceaccount")
		return err
	}

	if err := asr.reconcileRole(); err != nil {
		asr.Logger.Info("reconciling applicationSet role")
		return err
	}

	if err := asr.reconcileRoleBinding(); err != nil {
		asr.Logger.Info("reconciling applicationSet roleBinding")
		return err
	}

	if err := asr.reconcileDeployment(); err != nil {
		asr.Logger.Info("reconciling applicationSet deployment")
		return err
	}

	if err := asr.reconcileService(); err != nil {
		asr.Logger.Info("reconciling applicationSet deployment")
		return err
	}

	return nil
}

func (asr *ApplicationSetReconciler) DeleteResources() error {

	var deletionError error = nil

	if err := asr.DeleteService(resourceName, asr.Instance.Namespace); err != nil {
		asr.Logger.Error(err, "DeleteResources: failed to delete service")
		deletionError = err
	}

	if err := asr.DeleteDeployment(resourceName, asr.Instance.Namespace); err != nil {
		asr.Logger.Error(err, "DeleteResources: failed to delete deployment")
		deletionError = err
	}

	if err := asr.DeleteRoleBinding(resourceName, asr.Instance.Namespace); err != nil {
		asr.Logger.Error(err, "DeleteResources: failed to delete roleBinding")
		deletionError = err
	}

	if err := asr.DeleteRole(resourceName, asr.Instance.Namespace); err != nil {
		asr.Logger.Error(err, "DeleteResources: failed to delete role")
		deletionError = err
	}

	if err := asr.DeleteServiceAccount(resourceName, asr.Instance.Namespace); err != nil {
		asr.Logger.Error(err, "DeleteResources: failed to delete serviceaccount")
		deletionError = err
	}

	return deletionError
}
