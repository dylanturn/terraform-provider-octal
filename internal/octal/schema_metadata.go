package octal

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func metadataSchema(componentName, componentType string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uid": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("The unique in time and space value for this %s. More info: http://kubernetes.io/docs/user-guide/identifiers#uids", componentType),
			Computed:    true,
			Optional:    false,
		},
		"resource_version": {
			Type:        schema.TypeString,
			Description: fmt.Sprintf("An opaque value that represents the internal version of this %s that can be used by clients to determine when %s has changed. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency", componentType, componentName),
			Computed:    true,
			Optional:    false,
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     componentName,
			Description: "A name that will be given to the deployment",
		},
		"labels": {
			Type: schema.TypeMap,
			//Computed:    true,
			Optional:    true,
			Default:     map[string]string{},
			Description: "Additional labels to add to the deployment",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"annotations": {
			Type: schema.TypeMap,
			//Computed:    true,
			Optional:    true,
			Default:     map[string]string{},
			Description: "Additional annotations to add to the deployment",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}
