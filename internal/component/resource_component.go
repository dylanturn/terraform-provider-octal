package resource_component

import (
	"bytes"
	"text/template"
)

type ResourceComponent struct {
	Name                                    string
	Namespace                               string
	Config                                  map[string]interface{}
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

type Component interface {
	GetName() string
	GetNamespace() string
	GetConfig() map[string]interface{}
	RenderManifests() []string
}

func (component ResourceComponent) GetName() string {
	return component.Name
}

func (component ResourceComponent) GetNamespace() string {
	return component.Namespace
}

func (component ResourceComponent) GetConfig() map[string]interface{} {
	return component.Config
}

func (component ResourceComponent) RenderManifests() []string {
	renderedManifests := []string{}

	for _, manifest := range component.DeploymentManifests {
		// "You have a task named \"{{ .Name}}\" with description: \"{{ .Description}}\""
		manifestTemplate, err := template.New(component.Name).Parse(manifest)
		if err != nil {
			panic(err)
		}

		var parsedManifest bytes.Buffer
		err = manifestTemplate.Execute(&parsedManifest, component.GetConfig())
		if err != nil {
			panic(err)
		}

		renderedManifests = append(renderedManifests, parsedManifest.String())
	}
	return renderedManifests
}
