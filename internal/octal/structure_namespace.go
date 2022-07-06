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

func loadNamespaceManifest(manifestPath string) (*v1.Namespace, error) {
	namespace := &v1.Namespace{}
	err := loadManifest(manifestPath).Decode(&namespace)
	return namespace, err
}

func createNamespace(ctx context.Context, meta interface{}, d *schema.ResourceData, namespace *v1.Namespace) diag.Diagnostics {
	var diags diag.Diagnostics

	nsObjMfst := namespace

	/*******************************\
	** Load Manifest Template      **
	\*******************************/
	// This loads the namespace from the manifest template into a v1.Namespace object
	/*nsObjMfst, nsObjMfstErr := loadNamespaceManifest(path)
	if nsObjMfstErr != nil {
		tflog.Error(ctx, nsObjMfstErr.Error())
		return diags
	}*/

	/*******************************\
	** Update Manifest MetaData    **
	\*******************************/
	// This applied the updates provided by the Terraform resource to the base Namespace Object
	updateMetadata(ctx, "namespace", false, &nsObjMfst.ObjectMeta, d)

	/*******************************\
	** Create Kubernetes Object    **
	\*******************************/
	// Here we create the object using a template object customized by the resource inputs.
	tflog.Info(ctx, fmt.Sprintf("[INFO] Creating new namespace: %#v", nsObjMfst.Name))
	// Get the K8s client
	client := meta.(*apiClient).clientset
	_, err := client.CoreV1().Namespaces().Create(ctx, nsObjMfst, metav1.CreateOptions{})
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}
	tflog.Info(ctx, "[INFO] created a resource")

	/*******************************\
	** Read New Object Back        **
	\*******************************/
	// Here we appear to read back the namespace state to make sure it's consistent?
	readNamespace(ctx, d, meta)

	return diags
}

func readNamespace(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	namespace, nsErr := getNamespace(ctx, d, meta)
	if nsErr != nil {
		return diags
	}

	d.Set("namespace", flattenMetadata("namespace", &namespace.ObjectMeta))
	return diags
}

func updateNamespace(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	namespace, err := getNamespace(ctx, d, meta)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}

	updateMetadata(ctx, "namespace", false, &namespace.ObjectMeta, d)

	client := meta.(*apiClient).clientset
	client.CoreV1().Namespaces().Update(ctx, namespace, metav1.UpdateOptions{})

	// Here we appear to read back the namespace state to make sure it's consistent?
	readNamespace(ctx, d, meta)
	return diags
}

func deleteNamespace(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	namespace, err := getNamespace(ctx, d, meta)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return diags
	}

	client := meta.(*apiClient).clientset
	client.CoreV1().Namespaces().Delete(ctx, namespace.Name, metav1.DeleteOptions{})

	return diags
}
