package cainjector

import (
	"context"
	"embed"
	"github.com/dylanturn/terraform-provider-octal/internal/octal"
	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	Appsv1 "k8s.io/api/apps/v1"
)

//go:embed deployment.yml
var deployment []byte

//go:embed service-account.yml
var serviceAccount []byte

//go:embed roles/*
var roles embed.FS

//go:embed role-bindings/*
var roleBindings embed.FS

//go:embed cluster-roles/*
var clusterRoles embed.FS

//go:embed cluster-role-bindings/*
var clusterRoleBindings embed.FS

func GetDefaultDeployment(ctx context.Context) *Appsv1.Deployment {
	deploymentObject := &Appsv1.Deployment{}
	err := util.DecodeManifest(deployment).Decode(&deploymentObject)
	if err != nil {
		tflog.Error(ctx, err.Error())
	}
	return deploymentObject
}

type Component octal.Component
type ResourceComponent octal.ResourceComponent

func GetComponent() octal.Component {

	roles.ReadDir(".")

	cainjector := octal.ResourceComponent{
		Name:                              "cainjector",
		DeploymentManifests:               []string{string(deployment)},
		ServiceAccountManifests:           []string{string(serviceAccount)},
		ServiceManifests:                  []string{},
		RoleManifests:                     nil,
		RoleBindingManifests:              nil,
		ClusterRolesManifests:             nil,
		ClusterRoleBindingsManifests:      nil,
		CustomResourceDefinitionManifests: nil,
	}

	return cainjector
}
