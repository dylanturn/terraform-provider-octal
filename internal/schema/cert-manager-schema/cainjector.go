package cert_manager_schema

import (
	octal_schema "github.com/dylanturn/terraform-provider-octal/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CaiInjectorSchema() *schema.Resource {

	// componentSpec := *octal_schema.ComponentSchema()
	componentSpec := *octal_schema.NewComponentSchema()

	return &schema.Resource{
		Schema: componentSpec,
	}
}
