package octal

import (
	"context"
	"fmt"
	"log"

	resource_component "github.com/dylanturn/terraform-provider-octal/internal/component"
	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Create(ctx context.Context, d *schema.ResourceData, meta interface{}, component resource_component.Component) diag.Diagnostics {
	var diags diag.Diagnostics

	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Create a component called %s that contains %v parts.|", &component, component.GetName(), len(component.RenderManifests()))
	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Render manifests the component called %s |", &component, component.GetName())
	for _, manifest := range component.RenderManifests() {
		if manifest != "" {
			manifestObject, err := util.ContentToObject(manifest)
			if err != nil {
				log.Println("Failed to marshall the following manifest:")
				log.Println(manifest)
				tflog.Error(ctx, err.Error())
			}
			log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Add an object called %s to the component called %s |", &component, manifestObject.GetName(), component.GetName())
			component.AddUnstructuredObject(*manifestObject)
		}
	}

	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Component part count for %s: %v |", &component, component.GetName(), len(component.GetUnstructuredObjects()))
	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Set the component information for %s |", &component, component.GetName())
	err := d.Set(component.GetName(), flattenComponentRecord(ctx, component))
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to flatten the component record for %s. Error Message: %s", component.GetName(), err.Error()))
		return diags
	}

	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Deploy %v parts for the component named %s |", &component, len(component.GetUnstructuredObjects()), component.GetName())
	for _, componentPart := range component.GetUnstructuredObjects() {
		log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Deploy a part called %s for the component %s |", &component, componentPart.GetName(), component.GetName())
		util.ResourceK8sManifestCreate(ctx, d, meta, &componentPart)
	}

	log.Printf("[octal].[Create]::[ResourceComponent]:[component]:[%p] Component part deployment for %s complete! |", &component, component.GetName())

	Read(ctx, d, meta, component)

	return diags
}

func Read(ctx context.Context, d *schema.ResourceData, meta interface{}, component resource_component.Component) diag.Diagnostics {
	var diags diag.Diagnostics
	componentState := d.Get(component.GetName()).([]interface{})[0].(map[string]interface{})
	componentParts := componentState["parts"].([]interface{})
	for _, part := range componentParts {
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
	log.Printf("[octal].[flattenComponentRecord]::[ResourceComponent]:[component]:[%p] Flatten the component state entry |", &component)
	flatComponentRecord := make([]map[string]interface{}, 1)
	flatParts := []interface{}{}

	for _, part := range component.GetUnstructuredObjects() {
		log.Printf("[octal].[flattenComponentRecord]::[ResourceComponent]:[component]:[%p] Flatten the state entry for %s/%s |", &component, part.GetKind(), part.GetName())
		flatParts = append(flatParts, flattenUnstructuredObject(&part))
	}

	log.Printf("[octal].[flattenComponentRecord]::[ResourceComponent]:[component]:[%p] Add all the flat parts to the flat component record |", &component)
	flatComponentRecord[0] = map[string]interface{}{
		"parts": flatParts,
	}

	log.Printf("[octal].[flattenComponentRecord]::[ResourceComponent]:[component]:[%p] Component state entry is now flat! |", &component)
	return flatComponentRecord
}

func flattenUnstructuredObject(object *unstructured.Unstructured) map[string]interface{} {
	return map[string]interface{}{
		"group":       object.GetObjectKind().GroupVersionKind().Group,
		"version":     object.GetObjectKind().GroupVersionKind().Version,
		"kind":        object.GetObjectKind().GroupVersionKind().Kind,
		"name":        object.GetName(),
		"namespace":   object.GetNamespace(),
		"labels":      object.GetLabels(),
		"annotations": object.GetAnnotations(),
	}
}
