package resource_component

import (
	"context"

	Appsv1 "k8s.io/api/apps/v1"
	Corev1 "k8s.io/api/core/v1"
	Rbacv1 "k8s.io/api/rbac/v1"

	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Component interface {
	GetName() string
	GetDefaultDeployments(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Appsv1.Deployment
	GetDefaultServices(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Corev1.Service
	GetDefaultServiceAccounts(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Corev1.ServiceAccount
	GetDefaultRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Rbacv1.Role
	GetDefaultRoleBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Rbacv1.RoleBinding
	GetDefaultClusterRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Rbacv1.ClusterRole
	GetDefaultClusterRoleBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Rbacv1.ClusterRoleBinding
}

type ResourceComponent struct {
	Name                              string
	DeploymentManifests               []string
	ServiceManifests                  []string
	ServiceAccountManifests           []string
	RoleManifests                     []string
	RoleBindingManifests              []string
	ClusterRolesManifests             []string
	ClusterRoleBindingsManifests      []string
	CustomResourceDefinitionManifests []string
}

func (component ResourceComponent) GetName() string {
	return component.Name
}

func (component ResourceComponent) GetDefaultDeployments(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Appsv1.Deployment {
	manifests := component.DeploymentManifests
	objects := make([]Appsv1.Deployment, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, err.Error())
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultServiceAccounts(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Corev1.ServiceAccount {
	manifests := component.ServiceAccountManifests
	objects := make([]Corev1.ServiceAccount, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, err.Error())
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultServices(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Corev1.Service {
	manifests := component.ServiceManifests
	objects := make([]Corev1.Service, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, err.Error())
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Rbacv1.Role {
	manifests := component.RoleManifests
	objects := make([]Rbacv1.Role, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, err.Error())
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultRoleBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Rbacv1.RoleBinding {
	manifests := component.RoleBindingManifests
	objects := make([]Rbacv1.RoleBinding, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, err.Error())
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultClusterRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Rbacv1.ClusterRole {
	manifests := component.ClusterRolesManifests
	objects := make([]Rbacv1.ClusterRole, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, err.Error())
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultClusterRoleBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Rbacv1.ClusterRoleBinding {
	manifests := component.ClusterRoleBindingsManifests
	objects := make([]Rbacv1.ClusterRoleBinding, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, err.Error())
		}
	}
	return &objects
}
