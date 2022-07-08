package octal

import (
	"context"

	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createDeployments(ctx context.Context, meta interface{}, d *schema.ResourceData, deployments []string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, deployment := range deployments {
		util.ResourceK8sManifestCreate(ctx, d, meta, deployment)
	}

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
