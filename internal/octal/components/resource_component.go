package components

import (
	"bytes"
	"log"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceTemplates struct {
	DeploymentManifests                     []string
	ServiceManifests                        []string
	ServiceAccountManifests                 []string
	RoleManifests                           []string
	RoleBindingManifests                    []string
	ClusterRoleManifests                    []string
	ClusterRoleBindingManifests             []string
	MutatingWebhookConfigurationManifests   []string
	ValidatingWebhookConfigurationManifests []string
	CustomResourceDefinitionManifests       []string
}

type ResourceComponent struct {
	name              string
	namespace         string
	config            map[string]interface{}
	componentObjects  []OctalComponentObject
	resourceTemplates ResourceTemplates
	resourceManifests []string
}

type OctalComponent interface {
	GetName() string
	GetNamespace() string
	GetConfig() map[string]string
	GetObjects() []OctalComponentObject
	GetFlat() []map[string]interface{}
}

func GetComponent(componentName, componentNamespace string, d *schema.ResourceData, resourceTemplates ResourceTemplates) OctalComponent {
	componentConfig := map[string]interface{}{}

	// Get the component's configuration from the resource block.
	component, exists := d.GetOk(componentName)
	if exists {
		if component != nil && len(component.([]interface{})) > 0 {
			if component.([]interface{}) != nil {
				if component.([]interface{})[0] != nil {
					componentConfig = component.([]interface{})[0].(map[string]interface{})
				}
			}
		}
	}

	// Gather all the components manifest templates
	manifestTemplates := []string{}
	manifestTemplates = append(manifestTemplates, resourceTemplates.CustomResourceDefinitionManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.RoleManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.ClusterRoleManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.ServiceAccountManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.RoleBindingManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.ClusterRoleBindingManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.MutatingWebhookConfigurationManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.ValidatingWebhookConfigurationManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.DeploymentManifests...)
	manifestTemplates = append(manifestTemplates, resourceTemplates.ServiceManifests...)

	// Render the components templates and construct the OctalComponent object.
	renderedManifests := []string{}
	componentObjects := []OctalComponentObject{}
	for _, manifest := range manifestTemplates {
		manifestTemplate, err := template.New(componentName).Parse(manifest)
		if err != nil {
			panic(err)
		}

		flatConfig := map[string]string{}
		flatConfig["namespace"] = componentNamespace
		for k, v := range componentConfig {
			flattenConfig(k, v, flatConfig)
		}

		var parsedManifest bytes.Buffer
		err = manifestTemplate.Execute(&parsedManifest, flatConfig)
		if err != nil {
			panic(err)
		}
		renderedManifests = append(renderedManifests, parsedManifest.String())

		stringManifest := parsedManifest.String()

		if stringManifest != "" {
			componentObject, err := objectFromManifest(stringManifest)
			if err != nil {
				log.Printf("Failed to get object from manifest: %s", err.Error())
			} else {
				componentObjects = append(componentObjects, componentObject)
			}
		}
	}

	managedComponent := ResourceComponent{
		name:              componentName,
		namespace:         componentNamespace,
		config:            componentConfig,
		componentObjects:  componentObjects,
		resourceTemplates: resourceTemplates,
		resourceManifests: renderedManifests,
	}

	return &managedComponent
}

func (component *ResourceComponent) GetName() string {
	return component.name
}

func (component *ResourceComponent) GetNamespace() string {
	return component.namespace
}

func (component *ResourceComponent) GetConfig() map[string]string {
	flatmap := map[string]string{}
	flatmap["namespace"] = component.namespace
	for k, v := range component.config {
		flattenConfig(k, v, flatmap)
	}

	return flatmap
}

func (component *ResourceComponent) GetObjects() []OctalComponentObject {
	return component.componentObjects
}

func (component *ResourceComponent) GetFlat() []map[string]interface{} {
	flatComponentRecord := make([]map[string]interface{}, 1)
	flatParts := []interface{}{}

	for _, part := range component.GetObjects() {
		flatParts = append(flatParts, part.GetFlat())
	}

	flatComponentRecord[0] = map[string]interface{}{
		"objects": flatParts,
	}

	return flatComponentRecord
}
