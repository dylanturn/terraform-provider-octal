package octal

import (
	"context"
	octal_schema "github.com/dylanturn/terraform-provider-octal/internal/octal-schema"
	cert_manager_schema "github.com/dylanturn/terraform-provider-octal/internal/octal-schema/cert-manager-schema"
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
			"name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "cert-manager-schema",
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

type certManagerManifests struct {
	namespace                string
	controller               deploymentManifests
	cainjector               deploymentManifests
	webhook                  deploymentManifests
	customResourceDefinition []string
}
type deploymentManifests struct {
	deployment          string
	service             string
	serviceAccount      string
	role                string
	roleBinding         string
	clusterRoles        []string
	clusterRoleBindings []string
}

func resourceOctalCertManagerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags = diag.Diagnostics{}

	d.SetId(resource.UniqueId())

	namespaceManifest := "resources/cert-manager-schema/namespace.yml"
	serviceAccountManifest := "resources/cert-manager-schema/controller/service-account.yml"

	createNamespace(ctx, meta, d, namespaceManifest)
	createServiceAccount(ctx, meta, d, serviceAccountManifest)

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
