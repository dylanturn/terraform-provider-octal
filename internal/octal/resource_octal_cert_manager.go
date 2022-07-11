package octal

import (
	"context"
	"time"

	cainjector "github.com/dylanturn/terraform-provider-octal/internal/resources/cert-manager/cainjector"
	controller "github.com/dylanturn/terraform-provider-octal/internal/resources/cert-manager/controller"
	webhook "github.com/dylanturn/terraform-provider-octal/internal/resources/cert-manager/webhook"
	octal_schema "github.com/dylanturn/terraform-provider-octal/internal/schema"
	cert_manager_schema "github.com/dylanturn/terraform-provider-octal/internal/schema/cert-manager-schema"
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
				Elem:        cert_manager_schema.ControllerSchema(),
			},
			"cainjector": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the deployment",
				Elem:        cert_manager_schema.CaiInjectorSchema(),
			},
			"webhook": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Additional annotations to add to the deployment",
				Elem:        cert_manager_schema.WebhoookSchema(),
			},
			"custom_resources": {
				Type:        schema.TypeList,
				Optional:    false,
				Computed:    true,
				Description: "Additional annotations to add to the deployment",
				Elem:        octal_schema.CustomResourceDefinition(),
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

	Create(ctx, d, meta, webhook.GetComponent(d))
	Create(ctx, d, meta, cainjector.GetComponent(d))
	Create(ctx, d, meta, controller.GetComponent(d))

	resourceOctalCertManagerRead(ctx, d, meta)

	return diags
}

func resourceOctalCertManagerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	Read(ctx, d, meta, webhook.GetComponent(d))
	Read(ctx, d, meta, cainjector.GetComponent(d))
	Read(ctx, d, meta, controller.GetComponent(d))

	return diags
}

func resourceOctalCertManagerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	Update(ctx, d, meta, webhook.GetComponent(d))
	Update(ctx, d, meta, cainjector.GetComponent(d))
	Update(ctx, d, meta, controller.GetComponent(d))

	resourceOctalCertManagerRead(ctx, d, meta)

	return diags
}

func resourceOctalCertManagerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	Delete(ctx, d, meta, webhook.GetComponent(d))
	Delete(ctx, d, meta, cainjector.GetComponent(d))
	Delete(ctx, d, meta, controller.GetComponent(d))

	return diags
}
