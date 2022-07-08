package webhook

import (
	"embed"

	resource_component "github.com/dylanturn/terraform-provider-octal/internal/component"
	"github.com/dylanturn/terraform-provider-octal/internal/util"
)

//go:embed deployments/*
var deployments embed.FS

//go:embed services/*
var services embed.FS

//go:embed service-accounts/*
var serviceAccounts embed.FS

//go:embed roles/*
var roles embed.FS

//go:embed role-bindings/*
var roleBindings embed.FS

//go:embed cluster-roles/*
var clusterRoles embed.FS

//go:embed cluster-role-bindings/*
var clusterRoleBindings embed.FS

//go:embed custom-resource-definitions/*
var customResourceDefinitionManifests embed.FS

//go:embed mutating-webhook-configurations/*
var mutatingWebhookConfigurations embed.FS

//go:embed validating-webhook-configurations/*
var validatingWebhookConfigurations embed.FS

type Component resource_component.Component
type ResourceComponent resource_component.ResourceComponent

func GetComponent() resource_component.Component {

	webhook := resource_component.ResourceComponent{
		Name:                                    "webhook",
		DeploymentManifests:                     util.ReadEmbeddedFiles(deployments),
		ServiceAccountManifests:                 util.ReadEmbeddedFiles(serviceAccounts),
		ServiceManifests:                        util.ReadEmbeddedFiles(services),
		RoleManifests:                           util.ReadEmbeddedFiles(roles),
		RoleBindingManifests:                    util.ReadEmbeddedFiles(roleBindings),
		ClusterRoleManifests:                    util.ReadEmbeddedFiles(clusterRoles),
		ClusterRoleBindingManifests:             util.ReadEmbeddedFiles(clusterRoleBindings),
		CustomResourceDefinitionManifests:       util.ReadEmbeddedFiles(customResourceDefinitionManifests),
		MutatingWebhookConfigurationManifests:   util.ReadEmbeddedFiles(mutatingWebhookConfigurations),
		ValidatingWebhookConfigurationManifests: util.ReadEmbeddedFiles(validatingWebhookConfigurations),
	}

	return webhook
}
