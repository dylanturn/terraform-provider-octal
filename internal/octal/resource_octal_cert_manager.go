package octal

import (
	"context"
	"embed"
	cert_manager "github.com/dylanturn/terraform-provider-octal/internal/octal/resources/cert-manager"
	"github.com/dylanturn/terraform-provider-octal/internal/octal/resources/cert-manager/cainjector"
	octal_schema "github.com/dylanturn/terraform-provider-octal/internal/schema"
	cert_manager_schema "github.com/dylanturn/terraform-provider-octal/internal/schema/cert-manager-schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

//go:embed resources/cert-manager/cainjector/**
var certManagerCainjectorManifests embed.FS

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
				Type:        schema.TypeList,
				Optional:    false,
				Computed:    true,
				Description: "Additional annotations to add to the namespace",
				Elem:        octal_schema.NamespaceSchema(),
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

	//namespaceManifest := "resources/cert-manager-schema/namespace.yml"

	createNamespace(ctx, meta, d, cert_manager.GetDefaultNamespace(ctx))

	CreateDeployment(ctx, meta, d, "webhook", *cainjector.GetDefaultDeployment(ctx))
	CreateDeployment(ctx, meta, d, "cainjector", *cainjector.GetDefaultDeployment(ctx))
	CreateDeployment(ctx, meta, d, "controller", *cainjector.GetDefaultDeployment(ctx))

	resourceOctalCertManagerRead(ctx, d, meta)
	return diags
}

func resourceOctalCertManagerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	readNamespace(ctx, d, meta)
	readServiceAccount(ctx, d, meta)

	return diags
}

func resourceOctalCertManagerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	updateNamespace(ctx, d, meta)
	updateServiceAccount(ctx, d, meta)

	resourceOctalCertManagerRead(ctx, d, meta)
	return diags
}

func resourceOctalCertManagerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	deleteNamespace(ctx, d, meta)
	deleteServiceAccount(ctx, d, meta)

	return diags
}
