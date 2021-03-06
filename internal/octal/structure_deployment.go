package octal

import (
	"context"
	"fmt"

	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	Appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createDeployments(ctx context.Context, meta interface{}, d *schema.ResourceData, deployments map[string][]Appsv1.Deployment) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, deploymentList := range deployments {
		for _, deployment := range deploymentList {
			tflog.Info(ctx, deployment.String())
			util.ResourceK8sManifestCreate(ctx, d, meta, d.Get("namespace").(string), deployment.String())
			//createDeployment(ctx, meta, d, component, deployment)
			tflog.Info(ctx, fmt.Sprintf("%s Created!!", deployment.Kind))
		}
	}

	return diags
}

func createDeployment(ctx context.Context, meta interface{}, d *schema.ResourceData, component string, defaultDeployment Appsv1.Deployment) diag.Diagnostics {
	var diags diag.Diagnostics

	namespace := d.Get("namespace").(string)

	/*******************************\
	** Update Manifest MetaData    **
	\*******************************/
	// This applied the updates provided by the Terraform resource to the base Namespace Object
	updateMetadata(ctx, component, true, &defaultDeployment.ObjectMeta, d)

	/*******************************\
	** Create Kubernetes Object    **
	\*******************************/
	// Here we create the object using a template object customized by the resource inputs.
	tflog.Info(ctx, fmt.Sprintf("[INFO] Creating new %s: %#v", component, defaultDeployment.Name))
	// Get the K8s client
	client := meta.(*apiClient).clientset

	_, err := client.AppsV1().Deployments(namespace).Create(ctx, &defaultDeployment, metav1.CreateOptions{})
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}
	tflog.Info(ctx, "[INFO] created a resource")

	/*******************************\
	** Read New Object Back        **
	\*******************************/
	// Here we appear to read back the namespace state to make sure it's consistent?
	readDeployment(ctx, d, meta, component)

	return diags
}

func readDeployments(ctx context.Context, d *schema.ResourceData, meta interface{}, components []string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	for _, component := range components {
		readDeployment(ctx, d, meta, component)
	}

	return diags
}
func readDeployment(ctx context.Context, d *schema.ResourceData, meta interface{}, component string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	object, nsErr := getDeployment(ctx, d, meta)
	if nsErr != nil {
		return diags
	}

	d.Set(component, flattenMetadata(component, &object.ObjectMeta))

	updateMetadata(ctx, component, true, &object.ObjectMeta, d)

	return diags
}

func updateDeployments(ctx context.Context, d *schema.ResourceData, meta interface{}, components []string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	for _, component := range components {
		updateDeployment(ctx, d, meta, component)
	}

	return diags
}
func updateDeployment(ctx context.Context, d *schema.ResourceData, meta interface{}, component string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	namespace := d.Get("namespace").(string)

	object, err := getDeployment(ctx, d, meta)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}

	updateMetadata(ctx, component, true, &object.ObjectMeta, d)

	client := meta.(*apiClient).clientset
	client.AppsV1().Deployments(namespace).Update(ctx, object, metav1.UpdateOptions{})

	// Here we appear to read back the namespace state to make sure it's consistent?
	readNamespace(ctx, d, meta)
	return diags
}

func deleteDeployments(ctx context.Context, d *schema.ResourceData, meta interface{}, components []string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	for _, component := range components {
		deleteDeployment(ctx, d, meta, component)
	}

	return diags
}
func deleteDeployment(ctx context.Context, d *schema.ResourceData, meta interface{}, component string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	namespace := d.Get("namespace").(string)

	object, err := getDeployment(ctx, d, meta)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}

	client := meta.(*apiClient).clientset
	client.AppsV1().Deployments(namespace).Delete(ctx, object.Name, metav1.DeleteOptions{})

	return diags
}
