// Copyright 2020 ArgoCD Operator Developers
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

package common

const (
	// ArgoCDDefaultApplicationInstanceLabelKey is the default app name as a tracking label.
	ArgoCDDefaultApplicationInstanceLabelKey = "mycompany.com/appname"

	// ArgoCDDefaultArgoImage is the ArgoCD container image to use when not specified.
	ArgoCDDefaultArgoImage = "argoproj/argocd"

	// ArgoCDDefaultArgoServerOperationProcessors is the number of ArgoCD Server Operation Processors to use when not specified.
	ArgoCDDefaultArgoServerOperationProcessors = int32(10)

	// ArgoCDDefaultArgoServerStatusProcessors is the number of ArgoCD Server Status Processors to use when not specified.
	ArgoCDDefaultArgoServerStatusProcessors = int32(20)

	// ArgoCDDefaultArgoVersion is the ArgoCD container image tag to use when not specified.
	ArgoCDDefaultArgoVersion = "v1.4.1"

	// ArgoCDDefaultConfigManagementPlugins is the default configuration value for the config management plugins.
	ArgoCDDefaultConfigManagementPlugins = ""

	// ArgoCDDefaultDexConfig is the default dex configuration.
	ArgoCDDefaultDexConfig = ""

	// ArgoCDDefaultDexImage is the Dex container image to use when not specified.
	ArgoCDDefaultDexImage = "quay.io/dexidp/dex"

	// ArgoCDDefaultDexOAuthRedirectPath is the default path to use for the OAuth Redirect URI.
	ArgoCDDefaultDexOAuthRedirectPath = "/api/dex/callback"

	// ArgoCDDefaultDexServiceAccountName is the default Service Account name for the Dex server.
	ArgoCDDefaultDexServiceAccountName = "argocd-dex-server"

	// ArgoCDDefaultDexVersion is the Dex container image tag to use when not specified.
	ArgoCDDefaultDexVersion = "v2.21.0"

	// ArgoCDDefaultExportLocalCapicity is the default capacity to use for local export.
	ArgoCDDefaultExportLocalCapicity = "2Gi"

	// ArgoCDDefaultGATrackingID is the default Google Analytics tracking ID.
	ArgoCDDefaultGATrackingID = ""

	// ArgoCDDefaultGAAnonymizeUsers is the default value for anonymizing google analytics users.
	ArgoCDDefaultGAAnonymizeUsers = false

	// ArgoCDDefaultGrafanaAdminUsername is the Grafana admin username to use when not specified.
	ArgoCDDefaultGrafanaAdminUsername = "admin"

	// ArgoCDDefaultGrafanaAdminPasswordLength is the length of the generated default Grafana admin password.
	ArgoCDDefaultGrafanaAdminPasswordLength = 32

	// ArgoCDDefaultGrafanaAdminPasswordNumDigits is the number of digits to use for the generated default Grafana admin password.
	ArgoCDDefaultGrafanaAdminPasswordNumDigits = 5

	// ArgoCDDefaultGrafanaAdminPasswordNumSymbols is the number of symbols to use for the generated default Grafana admin password.
	ArgoCDDefaultGrafanaAdminPasswordNumSymbols = 5

	// ArgoCDDefaultGrafanaImage is the Grafana container image to use when not specified.
	ArgoCDDefaultGrafanaImage = "grafana/grafana"

	// ArgoCDDefaultGrafanaReplicas is the default Grafana replica count.
	ArgoCDDefaultGrafanaReplicas = int32(1)

	// ArgoCDDefaultGrafanaSecretKeyLength is the length of the generated default Grafana secret key.
	ArgoCDDefaultGrafanaSecretKeyLength = 20

	// ArgoCDDefaultGrafanaSecretKeyNumDigits is the number of digits to use for the generated default Grafana secret key.
	ArgoCDDefaultGrafanaSecretKeyNumDigits = 5

	// ArgoCDDefaultGrafanaSecretKeyNumSymbols is the number of symbols to use for the generated default Grafana secret key.
	ArgoCDDefaultGrafanaSecretKeyNumSymbols = 0

	// ArgoCDDefaultGrafanaConfigPath is the default Grafana configuration directory when not specified.
	ArgoCDDefaultGrafanaConfigPath = "/var/lib/grafana"

	// ArgoCDDefaultGrafanaVersion is the Grafana container image tag to use when not specified.
	ArgoCDDefaultGrafanaVersion = "6.6.1"

	// ArgoCDDefaultHelpChatURL is the default help chat URL.
	ArgoCDDefaultHelpChatURL = "https://mycorp.slack.com/argo-cd"

	// ArgoCDDefaultHelpChatText is the default help chat text.
	ArgoCDDefaultHelpChatText = "Chat now!"

	// ArgoCDDefaultIngressPath is the path to use for the Ingress when not specified.
	ArgoCDDefaultIngressPath = "/"

	// ArgoCDDefaultKustomizeBuildOptions is the default kustomize build options.
	ArgoCDDefaultKustomizeBuildOptions = ""

	// ArgoCDDefaultOIDCConfig is the default OIDC configuration.
	ArgoCDDefaultOIDCConfig = ""

	// ArgoCDDefaultPrometheusReplicas is the default Prometheus replica count.
	ArgoCDDefaultPrometheusReplicas = int32(1)

	// ArgoCDDefaultRBACPolicy is the default RBAC policy CSV data.
	ArgoCDDefaultRBACPolicy = ""

	// ArgoCDDefaultRBACDefaultPolicy is the default Argo CD RBAC policy.
	ArgoCDDefaultRBACDefaultPolicy = "role:readonly"

	// ArgoCDDefaultRBACScopes is the default Argo CD RBAC scopes.
	ArgoCDDefaultRBACScopes = "[groups]"

	// ArgoCDDefaultRedisImage is the Redis container image to use when not specified.
	ArgoCDDefaultRedisImage = "redis"

	// ArgoCDDefaultRedisVersion is the Redis container image tag to use when not specified.
	ArgoCDDefaultRedisVersion = "5.0.3"

	// ArgoCDDefaultResourceCustomizations is the default resource customizations.
	ArgoCDDefaultResourceCustomizations = ""

	// ArgoCDDefaultResourceExclusions is the default resource exlcusions.
	ArgoCDDefaultResourceExclusions = ""

	// ArgoCDDefaultRepositories is the default repositories.
	ArgoCDDefaultRepositories = ""

	// ArgoCDDefaultSSHKnownHosts is the default SSH Known hosts data.
	ArgoCDDefaultSSHKnownHosts = `bitbucket.org ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
gitlab.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=
gitlab.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAfuCHKVTjquxvt6CM6tdG4SLp1Btn/nOeHHE5UOzRdf
gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9
ssh.dev.azure.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H
vs-ssh.visualstudio.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7Hr1oTWqNqOlzGJOfGJ4NakVyIzf1rXYd4d7wo6jBlkLvCA4odBlL0mDUyZ0/QUfTTqeu+tm22gOsv+VrVTMk6vwRU75gY/y9ut5Mb3bR5BV58dKXyq9A9UeB5Cakehn5Zgm6x1mKoVyf+FFn26iYqXJRgzIZZcZ5V6hrE0Qg39kZm4az48o0AUbf6Sp4SLdvnuMa2sVNwHBboS7EJkm57XQPVU3/QpyNLHbWDdzwtrlS+ez30S3AdYhLKEOxAG8weOnyrtLJAUen9mTkol8oII1edf7mWWbWVf0nBmly21+nZcmCTISQBtdcyPaEno7fFQMDD26/s0lfKob4Kw8H	
`
)
