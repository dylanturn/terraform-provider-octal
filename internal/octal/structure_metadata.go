package octal

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func flattenMetadata(objectMeta *metav1.ObjectMeta) []map[string]interface{} {
	flatNamespaces := make([]map[string]interface{}, 1)
	flatNamespaces[0] = map[string]interface{}{
		"uid":              fmt.Sprintf("%s", objectMeta.GetUID()),
		"resource_version": objectMeta.GetResourceVersion(),
		"name":             objectMeta.GetName(),
		"labels":           objectMeta.GetLabels(),
		"annotations":      objectMeta.GetAnnotations(),
	}

	delete(flatNamespaces[0]["labels"].(map[string]string), "project-octal.io/cert-manager")
	delete(flatNamespaces[0]["labels"].(map[string]string), "kubernetes.io/metadata.name")

	return flatNamespaces
}

func updateMetadata(component string, metaData *metav1.ObjectMeta, d *schema.ResourceData) *metav1.ObjectMeta {

	namespaceConfig := d.Get(component).([]interface{})[0].(map[string]interface{})
	namespaceLabels := namespaceConfig["labels"].(map[string]interface{})
	namespaceAnnotations := namespaceConfig["annotations"].(map[string]interface{})

	metaData.Name = namespaceConfig["name"].(string)

	if len(namespaceLabels) > 0 {
		// Merge the Terraform'd labels with the labels that exist.
		for key, value := range namespaceLabels {
			metaData.Labels[key] = value.(string)
		}
	} else {
		metaData.Labels = map[string]string{}
	}
	if len(namespaceAnnotations) > 0 {
		for key, value := range namespaceAnnotations {
			metaData.Annotations[key] = value.(string)
		}
	} else {
		metaData.Annotations = map[string]string{}
	}

	// Add the label that contains the instance id
	metaData.Labels["project-octal.io/cert-manager"] = d.Id()

	return metaData
}
