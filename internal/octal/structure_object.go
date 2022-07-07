package octal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	Appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type KubeObject interface {
	metav1.Object
	runtime.Object
}

func createKubernetesObjectNew[T KubeObject](ctx context.Context, meta interface{}, d *schema.ResourceData, kind string, namespace string, deployment KubeObject, opts metav1.CreateOptions) (result rest.Result) {

	tflog.Info(ctx, fmt.Sprintf("The Deployment: %s", deployment.GetName()))
	tflog.Info(ctx, fmt.Sprintf("The namespace: %s", deployment.GetNamespace()))
	tflog.Info(ctx, fmt.Sprintf("The kind: %s", deployment.GetObjectKind().GroupVersionKind().Kind))

	/*client := meta.(*apiClient).clientset.RESTClient()
	result = client.Post().
		Namespace(deployment.GetNamespace()).
		Resource(deployment.GetObjectKind().GroupVersionKind().Kind).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(deployment).
		Do(ctx)
	return result*/
	return
}

func createKubernetesObject(ctx context.Context, meta interface{}, d *schema.ResourceData, kind string, namespace string, deployment KubeObject, opts metav1.CreateOptions) (result rest.Result) {

	tflog.Info(ctx, fmt.Sprintf("The Deployment: %s", deployment.GetName()))
	tflog.Info(ctx, fmt.Sprintf("The namespace: %s", deployment.GetNamespace()))
	tflog.Info(ctx, fmt.Sprintf("The kind: %s", deployment.GetObjectKind().GroupVersionKind().Kind))

	/*client := meta.(*apiClient).clientset.RESTClient()
	result = client.Post().
		Namespace(deployment.GetNamespace()).
		Resource(deployment.GetObjectKind().GroupVersionKind().Kind).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(deployment).
		Do(ctx)
	return result*/
	return
}

func createObject(ctx context.Context, meta interface{}, d *schema.ResourceData, component string, defaultDeployment Appsv1.Deployment) diag.Diagnostics {
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
	readDeployment(ctx, d, meta, component)

	return diags
}
