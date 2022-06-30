package octal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func namespaceSchema(componentName string) *schema.Resource {
	return &schema.Resource{
		Schema: metadataSchema(componentName, "namespace"),
	}
}
