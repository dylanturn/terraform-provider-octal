package octal

import (
	"context"

	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Create(ctx context.Context, d *schema.ResourceData, meta interface{}, component ResourceComponent) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, manifest := range component.RenderManifests() {
		util.ResourceK8sManifestCreate(ctx, d, meta, component, manifest)
	}

	return diags
}

func Read(ctx context.Context, d *schema.ResourceData, meta interface{}, component ResourceComponent) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func Update(ctx context.Context, d *schema.ResourceData, meta interface{}, component ResourceComponent) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func Delete(ctx context.Context, d *schema.ResourceData, meta interface{}, component ResourceComponent) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
