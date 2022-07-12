package octal

import (
	"context"
	"fmt"

	"github.com/dylanturn/terraform-provider-octal/internal/octal/components"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Create(ctx context.Context, d *schema.ResourceData, meta interface{}, component components.OctalComponent) diag.Diagnostics {
	var diags diag.Diagnostics

	/*
		for _, manifest := range component.RenderManifests() {
			if manifest != "" {
				manifestObject, err := contentToObject(manifest)
				if err != nil {
					log.Println("Failed to marshall the following manifest:")
					log.Println(manifest)
					tflog.Error(ctx, err.Error())
				}
				component.AddUnstructuredObject(*manifestObject)
			}
		}

		err := d.Set(component.GetName(), flattenComponentRecord(ctx, component))
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Failed to flatten the component record for %s. Error Message: %s", component.GetName(), err.Error()))
			return diags
		}*/

	for _, componentPart := range component.GetObjects() {
		resourceK8sManifestCreate(ctx, d, meta, componentPart)
	}

	Read(ctx, d, meta, component)

	return diags
}

func Read(ctx context.Context, d *schema.ResourceData, meta interface{}, component components.OctalComponent) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, componentPart := range component.GetObjects() {
		resourceK8sManifestRead(ctx, d, meta, componentPart)
	}

	return diags
}

func Update(ctx context.Context, d *schema.ResourceData, meta interface{}, component components.OctalComponent) diag.Diagnostics {
	var diags diag.Diagnostics

	/*
		for _, manifest := range component.RenderManifests() {
			if manifest != "" {
				manifestObject, err := contentToObject(manifest)
				if err != nil {
					log.Println("Failed to marshall the following manifest:")
					log.Println(manifest)
					tflog.Error(ctx, err.Error())
				}
				component.AddUnstructuredObject(*manifestObject)
			}
		}*/

	// Update the component manifests
	err := d.Set(component.GetName(), component.GetFlat())
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Failed to flatten the component record for %s. Error Message: %s", component.GetName(), err.Error()))
		return diags
	}

	// Update each of the component objects
	for _, componentObject := range component.GetObjects() {
		object := componentObject.GetUnstructuredObject()
		resourceK8sManifestUpdate(ctx, d, meta, &object)
	}

	Read(ctx, d, meta, component)

	return diags
}

func Delete(ctx context.Context, d *schema.ResourceData, meta interface{}, component components.OctalComponent) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, componentObject := range component.GetObjects() {
		resourceK8sManifestDelete(ctx, d, meta, componentObject)
	}

	return diags
}
