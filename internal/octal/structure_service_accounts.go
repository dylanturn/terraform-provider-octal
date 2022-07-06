package octal

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func loadServiceAccountManifest(manifestPath string) (*v1.ServiceAccount, error) {
	serviceAccount := &v1.ServiceAccount{}
	err := loadManifest(manifestPath).Decode(&serviceAccount)
	return serviceAccount, err
}

func createServiceAccount(ctx context.Context, meta interface{}, d *schema.ResourceData, path string) diag.Diagnostics {
	var diags diag.Diagnostics

	/*******************************\
	** Load Manifest Template      **
	\*******************************/
	// This loads the namespace from the manifest template into a v1.Namespace object
	objManifest, objManifestErr := loadServiceAccountManifest(path)
	if objManifestErr != nil {
		tflog.Error(ctx, objManifestErr.Error())
		return diags
	}

	/*******************************\
	** Update Manifest MetaData    **
	\*******************************/
	// This applied the updates provided by the Terraform resource to the base Namespace Object
	updateMetadata(ctx, "ServiceAccount", false, &objManifest.ObjectMeta, d)

	/*******************************\
	** Create Kubernetes Object    **
	\*******************************/
	// Here we create the object using a template object customized by the resource inputs.
	tflog.Info(ctx, fmt.Sprintf("[INFO] Creating new service-account: %#v", objManifest.Name))
	// Get the K8s client
	client := meta.(*apiClient).clientset
	namespace := d.Get("namespace").([]map[string]interface{})[0]["name"].(string)
	_, err := client.CoreV1().ServiceAccounts(namespace).Create(ctx, objManifest, metav1.CreateOptions{})
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}
	tflog.Info(ctx, "[INFO] created a resource")

	/*******************************\
	** Read New Object Back        **
	\*******************************/
	// Here we appear to read back the namespace state to make sure it's consistent?
	readServiceAccount(ctx, d, meta)

	return diags
}

func readServiceAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	namespace, nsErr := getServiceAccount(ctx, d, meta)
	if nsErr != nil {
		return diags
	}

	d.Set("service_account", flattenMetadata("ServiceAccount", &namespace.ObjectMeta))
	return diags
}

func updateServiceAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	serviceAccount, err := getServiceAccount(ctx, d, meta)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}

	updateMetadata(ctx, "namespace", true, &serviceAccount.ObjectMeta, d)

	client := meta.(*apiClient).clientset
	namespace := d.Get("namespace").([]map[string]interface{})[0]["name"].(string)
	client.CoreV1().ServiceAccounts(namespace).Update(ctx, serviceAccount, metav1.UpdateOptions{})

	// Here we appear to read back the namespace state to make sure it's consistent?
	readServiceAccount(ctx, d, meta)
	return diags
}

func deleteServiceAccount(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	serviceAccount, err := getServiceAccount(ctx, d, meta)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}

	client := meta.(*apiClient).clientset
	namespace := d.Get("namespace").([]map[string]interface{})[0]["name"].(string)
	client.CoreV1().ServiceAccounts(namespace).Delete(ctx, serviceAccount.Name, metav1.DeleteOptions{})

	return diags
}
