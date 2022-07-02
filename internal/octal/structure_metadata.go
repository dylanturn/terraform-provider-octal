package octal

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func flattenMetadata(component string, objectMeta *metav1.ObjectMeta) []map[string]interface{} {
	flatMetadata := make([]map[string]interface{}, 1)
	flatMetadata[0] = map[string]interface{}{
		"uid":              fmt.Sprintf("%s", objectMeta.GetUID()),
		"resource_version": objectMeta.GetResourceVersion(),
		"name":             objectMeta.GetName(),
		"component":        component,
		"labels":           objectMeta.GetLabels(),
		"annotations":      objectMeta.GetAnnotations(),
	}

	// resourceName := d.Get("name").(string)
	// componentName := fmt.Sprintf("%s-%s", resourceName, component)
	//
	// // Add the component labels that get added to everything
	// flatMetadata[0]["labels"].(map[string]string)["project-octal.io/cert-manager-schema"] = d.Id()
	// flatMetadata[0]["labels"].(map[string]string)["app.kubernetes.io/instance"] = d.Id()
	// flatMetadata[0]["labels"].(map[string]string)["app.kubernetes.io/version"] = d.Get("version").(string)
	// flatMetadata[0]["labels"].(map[string]string)["app.kubernetes.io/name"] = componentName
	// flatMetadata[0]["labels"].(map[string]string)["app.kubernetes.io/component"] = component
	// flatMetadata[0]["labels"].(map[string]string)["app.kubernetes.io/part-of"] = resourceName
	// flatMetadata[0]["labels"].(map[string]string)["app.kubernetes.io/created-by"] = "terraform"
	// flatMetadata[0]["labels"].(map[string]string)["app.kubernetes.io/managed-by"] = "terraform"

	// delete(flatMetadata[0]["labels"].(map[string]string), "project-octal.io/cert-manager-schema")
	// delete(flatMetadata[0]["labels"].(map[string]string), "kubernetes.io/metadata.name")

	return flatMetadata
}

// This applied the updates provided by the Terraform resource to the base Namespace Object
// Adds the labels and annotations defined by the Terraform resource.
func updateMetadata(component string, namespaced bool, metaData *metav1.ObjectMeta, d *schema.ResourceData) {

	// Get the components configuration from the resource block.
	componentConfig := d.Get(component).([]interface{})[0].(map[string]interface{})

	/*******************************\
	** MetaData Object Namespace   **
	\*******************************/

	// Set the namespace of the component, if the component is a namespaced resource
	if namespaced {
		metaData.SetNamespace(componentConfig["namespace"].(string))
	}

	/*******************************\
	** MetaData Object Name        **
	\*******************************/

	// Get the name specified in the resource block.
	// This will be combined with the component type to produce the component name.
	resourceName := d.Get("name").(string)
	var componentName string
	if component == "namespace" {
		componentName = resourceName
	} else {
		componentName = fmt.Sprintf("%s-%s", resourceName, component)
	}

	// Set the name of the component
	metaData.SetName(componentName)

	/*******************************\
	** MetaData Object Labels      **
	\*******************************/

	// Get the labels specified in the resource block
	componentConfigLabels := componentConfig["labels"].(map[string]interface{})

	// Create an object that will hold the metadata objects labels
	componentLabels := map[string]string{}

	// Get the labels from the component config
	if len(componentConfigLabels) > 0 {
		for key, value := range componentConfigLabels {
			componentLabels[key] = value.(string)
		}
	}

	// Add the component labels that get added to everything
	componentLabels["project-octal.io/cert-manager-schema"] = d.Id()
	componentLabels["app.kubernetes.io/instance"] = d.Id()
	componentLabels["app.kubernetes.io/version"] = d.Get("version").(string)
	componentLabels["app.kubernetes.io/name"] = componentName
	componentLabels["app.kubernetes.io/component"] = component
	componentLabels["app.kubernetes.io/part-of"] = resourceName
	componentLabels["app.kubernetes.io/created-by"] = "terraform"
	componentLabels["app.kubernetes.io/managed-by"] = "terraform"

	// Apply the label update to the metadata object
	metaData.SetLabels(componentLabels)

	/*******************************\
	** MetaData Object Annotations **
	\*******************************/

	// Get the annotations specified in the resource block
	componentConfigAnnotations := componentConfig["annotations"].(map[string]interface{})

	// Create an object that will hold the metadata objects annotations
	componentAnnotations := map[string]string{}

	// Get the annotations from the component config
	if len(componentConfigAnnotations) > 0 {
		for key, value := range componentConfigAnnotations {
			componentAnnotations[key] = value.(string)
		}
	}

	// Apply the annotation update to the metadata object
	metaData.SetAnnotations(metaData.Annotations)
}
