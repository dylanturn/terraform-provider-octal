package octal

import (
	"context"
	"time"

	"github.com/dylanturn/terraform-provider-octal/internal/octal/components"
	cainjector "github.com/dylanturn/terraform-provider-octal/internal/octal/components/cert-manager/cainjector"
	controller "github.com/dylanturn/terraform-provider-octal/internal/octal/components/cert-manager/controller"
	webhook "github.com/dylanturn/terraform-provider-octal/internal/octal/components/cert-manager/webhook"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOctalCertManager() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOctalCertManagerCreate,
		ReadContext:   resourceOctalCertManagerRead,
		UpdateContext: resourceOctalCertManagerUpdate,
		DeleteContext: resourceOctalCertManagerDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "cert-manager",
				Description:  "A name that will be given to the deployment",
				ValidateFunc: validateName,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A name that will be given to the deployment",
				Default:     "1.8.2",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The namespace to deploy Project-Octal in",
			},
			"controller": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the deployment",
				Elem:        ResourceComponentSchema(),
			},
			"cainjector": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the deployment",
				Elem:        ResourceComponentSchema(),
			},
			"webhook": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the deployment",
				Elem:        ResourceComponentSchema(),
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

	Create(ctx, d, meta, components.GetComponent("webhook", d.Get("namespace").(string), d, webhook.GetManifestTemplates()))
	Create(ctx, d, meta, components.GetComponent("cainjector", d.Get("namespace").(string), d, cainjector.GetManifestTemplates()))
	Create(ctx, d, meta, components.GetComponent("webhook", d.Get("namespace").(string), d, controller.GetManifestTemplates()))

	resourceOctalCertManagerRead(ctx, d, meta)

	return diags
}

func resourceOctalCertManagerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	Read(ctx, d, meta, components.GetComponent("webhook", d.Get("namespace").(string), d, webhook.GetManifestTemplates()))
	Read(ctx, d, meta, components.GetComponent("cainjector", d.Get("namespace").(string), d, cainjector.GetManifestTemplates()))
	Read(ctx, d, meta, components.GetComponent("webhook", d.Get("namespace").(string), d, controller.GetManifestTemplates()))

	return diags
}

func resourceOctalCertManagerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	Update(ctx, d, meta, components.GetComponent("webhook", d.Get("namespace").(string), d, webhook.GetManifestTemplates()))
	Update(ctx, d, meta, components.GetComponent("cainjector", d.Get("namespace").(string), d, cainjector.GetManifestTemplates()))
	Update(ctx, d, meta, components.GetComponent("webhook", d.Get("namespace").(string), d, controller.GetManifestTemplates()))

	resourceOctalCertManagerRead(ctx, d, meta)

	return diags
}

func resourceOctalCertManagerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	Delete(ctx, d, meta, components.GetComponent("webhook", d.Get("namespace").(string), d, webhook.GetManifestTemplates()))
	Delete(ctx, d, meta, components.GetComponent("cainjector", d.Get("namespace").(string), d, cainjector.GetManifestTemplates()))
	Delete(ctx, d, meta, components.GetComponent("webhook", d.Get("namespace").(string), d, controller.GetManifestTemplates()))

	return diags
}
