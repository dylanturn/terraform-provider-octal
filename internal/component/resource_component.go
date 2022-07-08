package resource_component

import (
	"context"
	"fmt"

	AdmissionV1 "k8s.io/api/admissionregistration/v1"
	AppsV1 "k8s.io/api/apps/v1"
	CoreV1 "k8s.io/api/core/v1"
	RbacV1 "k8s.io/api/rbac/v1"

	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Component interface {
	GetName() string
	GetDefaultDeployments(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]AppsV1.Deployment
	GetDefaultServices(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]CoreV1.Service
	GetDefaultServiceAccounts(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]CoreV1.ServiceAccount
	GetDefaultRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]RbacV1.Role
	GetDefaultRoleBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]RbacV1.RoleBinding
	GetDefaultClusterRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]RbacV1.ClusterRole
	GetDefaultClusterRoleBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]RbacV1.ClusterRoleBinding
	GetDefaultMutatingWebhookConfigurations(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]AdmissionV1.MutatingWebhookConfiguration
	GetDefaultValidatingWebhookConfigurations(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]AdmissionV1.ValidatingWebhookConfiguration
}

type ResourceComponent struct {
	Name                                    string
	DeploymentManifests                     []string
	ServiceManifests                        []string
	ServiceAccountManifests                 []string
	RoleManifests                           []string
	RoleBindingManifests                    []string
	ClusterRoleManifests                    []string
	ClusterRoleBindingManifests             []string
	MutatingWebhookConfigurationManifests   []string
	ValidatingWebhookConfigurationManifests []string
	CustomResourceDefinitionManifests       []string
}

func (component ResourceComponent) GetName() string {
	return component.Name
}

func (component ResourceComponent) GetDefaultDeployments(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]AppsV1.Deployment {
	manifests := component.DeploymentManifests
	objects := make([]AppsV1.Deployment, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {

			tflog.Error(ctx, fmt.Sprintf("Failed to decode deployment for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultServiceAccounts(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]CoreV1.ServiceAccount {
	manifests := component.ServiceAccountManifests
	objects := make([]CoreV1.ServiceAccount, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode service-account for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultServices(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]CoreV1.Service {
	manifests := component.ServiceManifests
	objects := make([]CoreV1.Service, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode service for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]RbacV1.Role {
	manifests := component.RoleManifests
	objects := make([]RbacV1.Role, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode role for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultRoleBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]RbacV1.RoleBinding {
	manifests := component.RoleBindingManifests
	objects := make([]RbacV1.RoleBinding, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode role-binding for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultClusterRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]RbacV1.ClusterRole {
	manifests := component.ClusterRoleManifests
	objects := make([]RbacV1.ClusterRole, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode cluster-role for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultClusterRoleBindings(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]RbacV1.ClusterRoleBinding {
	manifests := component.ClusterRoleBindingManifests
	objects := make([]RbacV1.ClusterRoleBinding, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode cluster-role-binding for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultMutatingWebhookConfigurations(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]AdmissionV1.MutatingWebhookConfiguration {
	manifests := component.MutatingWebhookConfigurationManifests
	objects := make([]AdmissionV1.MutatingWebhookConfiguration, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode MutatingWebhookConfiguration for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}

func (component ResourceComponent) GetDefaultValidatingWebhookConfigurations(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]AdmissionV1.ValidatingWebhookConfiguration {
	manifests := component.ValidatingWebhookConfigurationManifests
	objects := make([]AdmissionV1.ValidatingWebhookConfiguration, len(manifests))

	for index, manifest := range manifests {
		err := util.DecodeManifest([]byte(manifest)).Decode(&objects[index])
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to decode ValidatingWebhookConfigurationManifests for %s. Error: %s", component.Name, err.Error()))
		}
	}
	return &objects
}
