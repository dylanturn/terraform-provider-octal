package cainjector

import (
	"embed"
	"github.com/dylanturn/terraform-provider-octal/internal/util"
	Appsv1 "k8s.io/api/apps/v1"
	Corev1 "k8s.io/api/core/v1"
	Rbacv1 "k8s.io/api/rbac/v1"
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

type Cainjector struct {
	Deployment          Appsv1.Deployment
	Service             Corev1.Service
	ServiceAccount      Corev1.ServiceAccount
	Role                []Rbacv1.Role
	RoleBinding         []Rbacv1.RoleBinding
	ClusterRoles        []Rbacv1.ClusterRole
	ClusterRoleBindings []Rbacv1.ClusterRoleBinding
}

func (Cainjector) GetDefaultDeployment() (Appsv1.Deployment, error) {
	deploymentObj := &Appsv1.Deployment{}
	err := util.DecodeManifest(deployment).Decode(&deploymentObj)
	return *deploymentObj, err
}
