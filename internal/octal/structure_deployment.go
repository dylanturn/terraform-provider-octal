package octal

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	Appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDeployment(ctx context.Context, meta interface{}, d *schema.ResourceData, component string, defaultDeployment Appsv1.Deployment) diag.Diagnostics {
	var diags diag.Diagnostics

	namespace := d.Get("name").(string)

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
	ReadDeployment(ctx, d, meta, component)

	return diags
}

func ReadDeployment(ctx context.Context, d *schema.ResourceData, meta interface{}, component string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	object, nsErr := getDeployment(ctx, d, meta)

	if nsErr != nil {
		return diags
	}

	d.Set(component, flattenMetadata(component, &object.ObjectMeta))

	updateMetadata(ctx, component, true, &object.ObjectMeta, d)

	return diags
}

func UpdateDeployment(ctx context.Context, d *schema.ResourceData, meta interface{}, component string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	namespace := d.Get("name").(string)

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

func DeleteDeployment(ctx context.Context, d *schema.ResourceData, meta interface{}, component string) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	namespace := d.Get("namespace").([]interface{})[0].(map[string]interface{})["name"].(string)

	object, err := getDeployment(ctx, d, meta)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}

	client := meta.(*apiClient).clientset
	client.AppsV1().Deployments(namespace).Delete(ctx, object.Name, metav1.DeleteOptions{})

	return diags
}
