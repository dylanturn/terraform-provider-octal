package schema

import (
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewComponentSchema() *map[string]*schema.Schema {

	componentSchema := map[string]*schema.Schema{
		"parts": {
			Type:     schema.TypeList,
			Optional: false,
			Computed: true,
			Elem: &schema.Resource{
				Schema: Metadata(),
			},
		},
	}
	return &componentSchema
}
