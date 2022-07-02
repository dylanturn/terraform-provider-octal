package octal_schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Metadata() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uid": {
			Type:        schema.TypeString,
			Description: "The unique in time and space value for component",
			Computed:    true,
			Optional:    false,
		},
		"resource_version": {
			Type:        schema.TypeString,
			Description: "An opaque value that represents the internal version of this component that can be used by clients to determine when it has changed",
			Computed:    true,
			Optional:    false,
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Additional labels to add to the deployment",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"annotations": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Additional annotations to add to the deployment",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}
