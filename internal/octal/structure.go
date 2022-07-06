package octal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	Appsv1 "k8s.io/api/apps/v1"
	Corev1 "k8s.io/api/core/v1"
	Rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
	"log"
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

func (component ResourceComponent) GetDefaultDeployment(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Appsv1.Deployment {
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

func (component ResourceComponent) GetDefaultServiceAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Corev1.ServiceAccount {
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

func (component ResourceComponent) GetDefaultService(ctx context.Context, d *schema.ResourceData, meta interface{}) *[]Corev1.Service {
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

func readFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func loadManifest(manifestPath string) *k8Yaml.YAMLOrJSONDecoder {
	return k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(readFile(manifestPath)), 1000)
}

func octalListOptions(resourceId string) metav1.ListOptions {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"project-octal.io/cert-manager-schema": resourceId}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	return listOptions
}

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func getNamespace(ctx context.Context, d *schema.ResourceData, meta interface{}) (*Corev1.Namespace, error) {
	client := meta.(*apiClient).clientset
	namespaces, err := client.CoreV1().Namespaces().List(ctx, octalListOptions(d.Id()))
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get objects! %s", err.Error()))
		return nil, err
	}
	if len(namespaces.Items) > 1 {
		return nil, errors.New(fmt.Sprintf("Found more than one object with the same id! Objects Found: %s", len(namespaces.Items)))
	}
	if len(namespaces.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("Couldn't find object with the id! Objects Found: %s", len(namespaces.Items)))
	}

	updateMetadata(ctx, "namespace", false, &namespaces.Items[0].ObjectMeta, d)

	return &namespaces.Items[0], nil
}

func getDeployment(ctx context.Context, d *schema.ResourceData, meta interface{}) (*Appsv1.Deployment, error) {
	namespace := d.Get("namespace").([]interface{})[0].(map[string]interface{})["name"].(string)
	client := meta.(*apiClient).clientset
	namespaces, err := client.AppsV1().Deployments(namespace).List(ctx, octalListOptions(d.Id()))
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get objects! %s", err.Error()))
		return nil, err
	}
	if len(namespaces.Items) > 1 {
		return nil, errors.New(fmt.Sprintf("Found more than one object with the same id! Objects Found: %s", len(namespaces.Items)))
	}
	if len(namespaces.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("Couldn't find object with the id! Objects Found: %s", len(namespaces.Items)))
	}

	return &namespaces.Items[0], nil
}

func getService(ctx context.Context, d *schema.ResourceData, meta interface{}) (*Corev1.Service, error) {

	namespace := d.Get("namespace").([]interface{})[0].(map[string]interface{})["name"].(string)

	client := meta.(*apiClient).clientset
	objects, err := client.CoreV1().Services(namespace).List(ctx, octalListOptions(d.Id()))
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get objects! %s", err.Error()))
		return nil, err
	}
	if len(objects.Items) > 1 {
		return nil, errors.New(fmt.Sprintf("Found more than one object with the same id! Objects Found: %s", len(objects.Items)))
	}
	if len(objects.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("Couldn't find object with the id! Objects Found: %s", len(objects.Items)))
	}
	return &objects.Items[0], nil
}

func getServiceAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) (*Corev1.ServiceAccount, error) {
	client := meta.(*apiClient).clientset

	namespace := d.Get("name").(string)

	serviceAccounts, err := client.CoreV1().ServiceAccounts(namespace).List(ctx, octalListOptions(d.Id()))

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get objects! %s", err.Error()))
		return nil, err
	}
	if len(serviceAccounts.Items) > 1 {
		return nil, errors.New(fmt.Sprintf("Found more than one object with the same id! Objects Found: %s", len(serviceAccounts.Items)))
	}
	if len(serviceAccounts.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("Couldn't find object with the id! Objects Found: %s", len(serviceAccounts.Items)))
	}

	return &serviceAccounts.Items[0], nil
}

func getRole(ctx context.Context, d *schema.ResourceData, meta interface{}) (*Rbacv1.Role, error) {

	namespace := d.Get("namespace").([]interface{})[0].(map[string]interface{})["name"].(string)

	client := meta.(*apiClient).clientset
	objects, err := client.RbacV1().Roles(namespace).List(ctx, octalListOptions(d.Id()))
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get objects! %s", err.Error()))
		return nil, err
	}
	if len(objects.Items) > 1 {
		return nil, errors.New(fmt.Sprintf("Found more than one object with the same id! Objects Found: %s", len(objects.Items)))
	}
	if len(objects.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("Couldn't find object with the id! Objects Found: %s", len(objects.Items)))
	}
	return &objects.Items[0], nil
}

func getRoleBinding(ctx context.Context, d *schema.ResourceData, meta interface{}) (*Rbacv1.RoleBinding, error) {

	namespace := d.Get("namespace").([]interface{})[0].(map[string]interface{})["name"].(string)

	client := meta.(*apiClient).clientset
	objects, err := client.RbacV1().RoleBindings(namespace).List(ctx, octalListOptions(d.Id()))
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get objects! %s", err.Error()))
		return nil, err
	}
	if len(objects.Items) > 1 {
		return nil, errors.New(fmt.Sprintf("Found more than one object with the same id! Objects Found: %s", len(objects.Items)))
	}
	if len(objects.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("Couldn't find object with the id! Objects Found: %s", len(objects.Items)))
	}
	return &objects.Items[0], nil
}

func getClusterRole(ctx context.Context, d *schema.ResourceData, meta interface{}) (*Rbacv1.ClusterRole, error) {
	client := meta.(*apiClient).clientset
	objects, err := client.RbacV1().ClusterRoles().List(ctx, octalListOptions(d.Id()))
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get objects! %s", err.Error()))
		return nil, err
	}
	if len(objects.Items) > 1 {
		return nil, errors.New(fmt.Sprintf("Found more than one object with the same id! Objects Found: %s", len(objects.Items)))
	}
	if len(objects.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("Couldn't find object with the id! Objects Found: %s", len(objects.Items)))
	}
	return &objects.Items[0], nil
}

func getClusterRoleBinding(ctx context.Context, d *schema.ResourceData, meta interface{}) (*Rbacv1.ClusterRoleBinding, error) {
	client := meta.(*apiClient).clientset
	objects, err := client.RbacV1().ClusterRoleBindings().List(ctx, octalListOptions(d.Id()))
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to get objects! %s", err.Error()))
		return nil, err
	}
	if len(objects.Items) > 1 {
		return nil, errors.New(fmt.Sprintf("Found more than one object with the same id! Objects Found: %s", len(objects.Items)))
	}
	if len(objects.Items) < 1 {
		return nil, errors.New(fmt.Sprintf("Couldn't find object with the id! Objects Found: %s", len(objects.Items)))
	}
	return &objects.Items[0], nil
}
