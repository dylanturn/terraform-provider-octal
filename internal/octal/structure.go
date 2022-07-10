package octal

import (
	"context"
	"fmt"
	"log"
	"strings"

	resource_component "github.com/dylanturn/terraform-provider-octal/internal/component"
	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Create(ctx context.Context, d *schema.ResourceData, meta interface{}, component resource_component.Component) diag.Diagnostics {
	var diags diag.Diagnostics

	/* Here we have one manifest...
	 ** 2022-07-10T10:34:32.593-0500 [INFO]  provider.terraform-provider-octal: Component Part Count: 1: @module=octal tf_resource_type=octal_cert_manager tf_rpc=ApplyResourceChange tf_provider_addr=registry.terraform.io/hashicorp/octal tf_req_id=4855047f-a13c-626c-23d7-799a3cc9540f @caller=/Users/dylanturnbull/tmp/terraform-provider-octal/internal/octal/structure.go:19 timestamp=2022-07-10T10:34:32.593-0500
	 */
	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Create a component called %s that contains %v parts.|", &component, component.GetName(), len(component.RenderManifests()))

	// First render all the manifests, then turn them into objects and put them in a list
	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Render manifests the component called %s |", &component, component.GetName())
	for _, manifest := range component.RenderManifests() {
		manifestObject, err := util.ContentToObject(manifest)
		if err != nil {
			tflog.Error(ctx, err.Error())
		}
		/* We can see that all the objects are getting added
		 ** 2022-07-10T10:34:32.593-0500 [INFO]  provider.terraform-provider-octal: Add object: cert-manager-schema: @caller=/Users/dylanturnbull/tmp/terraform-provider-octal/internal/octal/structure.go:27 tf_resource_type=octal_cert_manager @module=octal tf_provider_addr=registry.terraform.io/hashicorp/octal tf_req_id=4855047f-a13c-626c-23d7-799a3cc9540f tf_rpc=ApplyResourceChange timestamp=2022-07-10T10:34:32.593-0500
		 */
		log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Add an object called %s to the component called %s |", &component, manifestObject.GetName(), component.GetName())
		component.AddUnstructuredObject(*manifestObject)

	}

	/* WHY IS THIS RETURNING 0??
	 ** 2022-07-10T10:34:32.590-0500 [INFO]  provider.terraform-provider-octal: !!!### Component part count: 0: tf_rpc=ApplyResourceChange @caller=/Users/dylanturnbull/tmp/terraform-provider-octal/internal/octal/structure.go:31 @module=octal tf_provider_addr=registry.terraform.io/hashicorp/octal tf_req_id=4855047f-a13c-626c-23d7-799a3cc9540f tf_resource_type=octal_cert_manager timestamp=2022-07-10T10:34:32.590-0500
	 */
	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Component part count for %s: %v |", &component, component.GetName(), len(component.GetUnstructuredObjects()))

	// Second, update the state data for the component.
	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Set the component information for %s |", &component, component.GetName())
	d.Set(component.GetName(), flattenComponentRecord(ctx, component))

	// Now create each of the parts in Kubernetes.
	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Deploy %v parts for the component named %s |", &component, len(component.GetUnstructuredObjects()), component.GetName())
	for _, componentPart := range component.GetUnstructuredObjects() {
		log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Deploy a part called %s for the component %s |", &component, componentPart.GetName(), component.GetName())
		util.ResourceK8sManifestCreate(ctx, d, meta, &componentPart)
	}

	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Component part deployment for %s complete! |", &component, component.GetName())

	return diags
}

func Read(ctx context.Context, d *schema.ResourceData, meta interface{}, component resource_component.Component) diag.Diagnostics {
	var diags diag.Diagnostics

	// Presumably this is where we read the runtime configuration/status of each component part.
	for _, part := range d.Get(component.GetName()).([]interface{})[0].(map[string]interface{})["parts"].(map[string]interface{}) {
		util.ResourceK8sManifestRead(ctx, d, meta, part.(map[string]interface{}))
	}

	return diags
}

func Update(ctx context.Context, d *schema.ResourceData, meta interface{}, component resource_component.Component) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func Delete(ctx context.Context, d *schema.ResourceData, meta interface{}, component resource_component.Component) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func flattenComponentRecord(ctx context.Context, component resource_component.Component) []map[string]interface{} {
	flatComponentRecord := make([]map[string]interface{}, 1)

	flatParts := map[string]interface{}{}

	for _, part := range component.GetUnstructuredObjects() {
		tflog.Info(ctx, fmt.Sprintf("##################### %s #####################", part))
		flatParts[buildId(&part)] = flattenUnstructuredObject(&part)
	}

	flatComponentRecord[0] = map[string]interface{}{
		"parts": flatParts,
	}

	return flatComponentRecord
}

func flattenUnstructuredObject(object *unstructured.Unstructured) map[string]interface{} {
	return map[string]interface{}{
		"group":            object.GetObjectKind().GroupVersionKind().Group,
		"version":          object.GetObjectKind().GroupVersionKind().Version,
		"kind":             object.GetObjectKind().GroupVersionKind().Kind,
		"uid":              object.GetUID(),
		"resource_version": object.GetResourceVersion(),
		"name":             object.GetName(),
		"labels":           object.GetLabels(),
		"annotations":      object.GetAnnotations(),
	}
}

const idSeparator = "::"

func buildId(object *unstructured.Unstructured) string {
	return strings.Join(
		[]string{
			object.GetNamespace(),
			object.GroupVersionKind().GroupVersion().String(),
			object.GroupVersionKind().Kind,
			object.GetName(),
		},
		idSeparator,
	)
}
