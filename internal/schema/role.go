package schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RoleSchema() *schema.Resource {
	componentSpec := Metadata()

	componentSpec["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    false,
		Computed:    true,
		Description: "The name of this deployment",
	}
	componentSpec["component"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    false,
		Computed:    true,
		Description: "The name of this deployment",
	}

	return &schema.Resource{
		Schema: componentSpec,
	}
}
