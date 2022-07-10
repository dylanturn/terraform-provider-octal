package cert_manager_schema

import (
	octal_schema "github.com/dylanturn/terraform-provider-octal/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func WebhoookSchema() *schema.Resource {

	//componentSpec := *octal_schema.ComponentSchema()
	componentSpec := *octal_schema.NewComponentSchema()

	componentSpec["mutating_webhook"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     octal_schema.MutatingWebhookConfiguration(),
	}

	componentSpec["validating_webhook"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     octal_schema.ValidatingWebhookConfiguration(),
	}

	return &schema.Resource{
		Schema: componentSpec,
	}
}
