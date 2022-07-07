package cainjector

import (
	"embed"

	resource_component "github.com/dylanturn/terraform-provider-octal/internal/component"
	"github.com/dylanturn/terraform-provider-octal/internal/util"
)

//go:embed deployment.yml
var deployment embed.FS

//go:embed service-account.yml
var serviceAccounts embed.FS

//go:embed roles/*
var roles embed.FS

//go:embed role-bindings/*
var roleBindings embed.FS

//go:embed cluster-roles/*
var clusterRoles embed.FS

//go:embed cluster-role-bindings/*
var clusterRoleBindings embed.FS

type Component resource_component.Component
type ResourceComponent resource_component.ResourceComponent

func GetComponent() resource_component.Component {

	cainjector := resource_component.ResourceComponent{
		Name:                              "cainjector",
		DeploymentManifests:               util.ReadEmbeddedFiles(deployment),
		ServiceAccountManifests:           util.ReadEmbeddedFiles(serviceAccounts),
		ServiceManifests:                  []string{},
		RoleManifests:                     util.ReadEmbeddedFiles(roles),
		RoleBindingManifests:              util.ReadEmbeddedFiles(roleBindings),
		ClusterRolesManifests:             util.ReadEmbeddedFiles(clusterRoles),
		ClusterRoleBindingsManifests:      util.ReadEmbeddedFiles(clusterRoleBindings),
		CustomResourceDefinitionManifests: nil,
	}

	return cainjector
}
