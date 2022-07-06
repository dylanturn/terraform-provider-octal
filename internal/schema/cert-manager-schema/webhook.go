package cert_manager_schema

import (
	octal_schema "github.com/dylanturn/terraform-provider-octal/internal/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func WebhoookSchema() *schema.Resource {

	webhookSpec := *octal_schema.ComponentSchema()

	webhookSpec["mutating_webhook"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     octal_schema.MutatingWebhookConfiguration(),
	}

	webhookSpec["validating_webhook"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     octal_schema.ValidatingWebhookConfiguration(),
	}

	return &schema.Resource{
		Schema: webhookSpec,
	}
}
