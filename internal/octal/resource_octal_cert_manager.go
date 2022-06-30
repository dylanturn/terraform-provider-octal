package octal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceOctalCertManager() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOctalCertManagerCreate,
		ReadContext:   resourceOctalCertManagerRead,
		UpdateContext: resourceOctalCertManagerUpdate,
		DeleteContext: resourceOctalCertManagerDelete,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the namespace",
				Elem:        namespaceSchema("cert-manager"),
			},
			"controller": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the deployment",
				Elem: octalDeploySpecSchema(
					"cert-manager-cainjector",
					"v1.8.1",
					"jetstack/cert-manager-controller"),
			},
			"cainjector": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the deployment",
				Elem: octalDeploySpecSchema(
					"cert-manager-cainjector",
					"v1.8.1",
					"jetstack/cert-manager-cainjector"),
			},
			"webhook": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the deployment",

				Elem: octalDeploySpecSchema(
					"cert-manager-cainjector",
					"v1.8.1",
					"jetstack/cert-manager-webhook"),
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceOctalCertManagerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	d.SetId(resource.UniqueId())

	namespaceManifest := "/Users/dylanturnbull/tmp/terraform-provider-octal/internal/templates/cert-manager/namespace.yml"

	createNamespace(ctx, meta, d, namespaceManifest)
	//createDeployment(ctx, d, meta, deploymentManifest) <- That's next?

	resourceOctalCertManagerRead(ctx, d, meta)
	return diags
}

func resourceOctalCertManagerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	readNamespace(ctx, d, meta)
	//readDeployment(ctx, d, meta) <- That's next?
	return diags
}

func resourceOctalCertManagerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	updateNamespace(ctx, d, meta)
	//updateDeployment(ctx, d, meta) <- That's next?
	return diags
}

func resourceOctalCertManagerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	deleteNamespace(ctx, d, meta)
	//deleteDeployment(ctx, d, meta) <- That's next?
	return diags
}
