package controller

import (
	"embed"

	"github.com/dylanturn/terraform-provider-octal/internal/octal/components"
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

// type Component components.OctalComponent
// type ResourceComponent components.ResourceComponent

func GetManifestTemplates() components.ResourceTemplates {
	componentManifests := components.ResourceTemplates{
		DeploymentManifests:                     components.ReadEmbeddedFiles(deployments),
		ServiceAccountManifests:                 components.ReadEmbeddedFiles(serviceAccounts),
		ServiceManifests:                        components.ReadEmbeddedFiles(services),
		RoleManifests:                           components.ReadEmbeddedFiles(roles),
		RoleBindingManifests:                    components.ReadEmbeddedFiles(roleBindings),
		ClusterRoleManifests:                    components.ReadEmbeddedFiles(clusterRoles),
		ClusterRoleBindingManifests:             components.ReadEmbeddedFiles(clusterRoleBindings),
		CustomResourceDefinitionManifests:       components.ReadEmbeddedFiles(customResourceDefinitionManifests),
		MutatingWebhookConfigurationManifests:   components.ReadEmbeddedFiles(mutatingWebhookConfigurations),
		ValidatingWebhookConfigurationManifests: components.ReadEmbeddedFiles(validatingWebhookConfigurations),
	}
	return componentManifests
}

/*
func GetComponent(d *schema.ResourceData) components.OctalComponent {

	componentName := "controller"
	componentConfig := map[string]interface{}{}

	// Get the component's configuration from the resource block.
	component, exists := d.GetOk(componentName)
	if exists {
		if component != nil && len(component.([]interface{})) > 0 {
			if component.([]interface{}) != nil {
				if component.([]interface{})[0] != nil {
					componentConfig = component.([]interface{})[0].(map[string]interface{})
				}
			}
		}
	}

	managedComponent := components.ResourceComponent{
		Name:             componentName,
		Namespace:        d.Get("namespace").(string),
		Config:           componentConfig,
		ComponentObjects: []components.OctalComponentObject{},
	}

	return &managedComponent
}
*/
