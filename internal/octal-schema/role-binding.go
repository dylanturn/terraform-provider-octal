package octal_schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RoleBindingSchema() *schema.Resource {
	componentSpec := Metadata()

	componentSpec["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: false,
		Computed: true,
	}
	componentSpec["component"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: false,
		Computed: true,
	}

	return &schema.Resource{
		Schema: componentSpec,
	}
}
