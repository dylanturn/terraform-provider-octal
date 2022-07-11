package schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Metadata() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"group": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: false,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: false,
		},
		"kind": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: false,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: false,
		},
		"namespace": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: false,
		},
		"labels": {
			Type:     schema.TypeMap,
			Computed: true,
			Optional: false,
		},
		"annotations": {
			Type:     schema.TypeMap,
			Computed: true,
			Optional: false,
		},
	}
}
