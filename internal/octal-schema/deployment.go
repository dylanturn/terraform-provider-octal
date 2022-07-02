package octal_schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DeploymentSchema() map[string]*schema.Schema {

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
	componentSpec["image_tag"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The image tag used by the deployment",
	}
	componentSpec["image_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The image name used by the deployment",
	}
	componentSpec["image_repository"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The image repository to use when pulling images",
	}
	componentSpec["image_pull_policy"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Determines when the image should be pulled prior to starting the container. `Always`: Always pull the image. | `IfNotPresent`: Only pull the image if it does not already exist on the node. | `Never`: Never pull the image",
	}
	return componentSpec
}
