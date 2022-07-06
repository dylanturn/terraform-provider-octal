package schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ComponentSchema() *map[string]*schema.Schema {

	componentSchema := DeploymentSchema()

	componentSchema["service"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     Service(),
	}
	componentSchema["service_account"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     ServiceAccount(),
	}
	componentSchema["roles"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     RoleSchema(),
	}
	componentSchema["role_bindings"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     RoleBindingSchema(),
	}
	componentSchema["cluster_roles"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     ClusterRole(),
	}
	componentSchema["cluster_role_bindings"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Computed: true,
		Elem:     ClusterRoleBinding(),
	}

	return &componentSchema
}
