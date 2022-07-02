package cert_manager_schema

import (
	octal_schema "github.com/dylanturn/terraform-provider-octal/internal/octal-schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ControllerSchema() *schema.Resource {

	componentSpec := *octal_schema.ComponentSchema()

	return &schema.Resource{
		Schema: componentSpec,
	}
}
