package resource_component

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ResourceComponent struct {
	Name                                    string
	Namespace                               string
	Config                                  map[string]interface{}
	UnstructuredObjects                     []unstructured.Unstructured
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
	GetConfig() map[string]string
	RenderManifests() []string
	GetUnstructuredObjects() []unstructured.Unstructured
	AddUnstructuredObject(unstructuredObjects unstructured.Unstructured)
}

func (component ResourceComponent) GetName() string {
	return component.Name
}

func (component ResourceComponent) GetNamespace() string {
	return component.Namespace
}

func (component ResourceComponent) GetConfig() map[string]string {
	flatmap := map[string]string{}
	for k, v := range component.Config {
		flattenConfig(k, v, flatmap)
	}

	return flatmap
}

func (component ResourceComponent) RenderManifests() []string {
	renderedManifests := []string{}
	log.Printf("[resource_component].[RenderManifests]::[ResourceComponent]:[component]:[%p] Rendering templates for component: %s |", &component, component.Name)
	for _, manifest := range component.DeploymentManifests {
		log.Printf("[resource_component].[RenderManifests]::[ResourceComponent]:[component]:[%p] Render component manifest |", &component)
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
		log.Printf("[resource_component].[RenderManifests]::[ResourceComponent]:[component]:[%p] Adding the rendered manifest to renderedManifests |", &component)
		renderedManifests = append(renderedManifests, parsedManifest.String())
	}
	log.Printf("[resource_component].[RenderManifests]::[ResourceComponent]:[component]:[%p] Finished templating %v manifests |", &component, len(renderedManifests))
	return renderedManifests
}

func (component ResourceComponent) AddUnstructuredObject(unstructuredObject unstructured.Unstructured) {
	/* Objects getting added here
	 ** 2022-07-10T10:34:32.593-0500 [INFO]  provider.terraform-provider-octal: 2022/07/10 10:34:32 ####!!!! Add unstructuredObject to component.UnstructuredObjects !!!!####: timestamp=2022-07-10T10:34:32.593-0500
	 */
	log.Printf("[resource_component].[AddUnstructuredObject]::[ResourceComponent]:[component]:[%p] Add unstructuredObject %s to component %s |", &component, unstructuredObject.GetName(), component.GetName())
	component.UnstructuredObjects = append(component.UnstructuredObjects, unstructuredObject)

	/* We can see the unstructured objects list is being built here...
	 ** 2022-07-10T10:34:32.593-0500 [INFO]  provider.terraform-provider-octal: 2022/07/10 10:34:32 Get one of []unstructured.Unstructured{unstructured.Unstructured{Object:map[string]interface {}{"apiVersion":"apps/v1", "kind":"Deployment", "metadata":map[string]interface {}{"labels":map[string]interface {}{"app.kubernetes.io/component":"controller", "app.kubernetes.io/created-by":"terraform", "app.kubernetes.io/instance":"", "app.kubernetes.io/managed-by":"terraform", "app.kubernetes.io/name":"cert-manager-schema", "app.kubernetes.io/part-of":"cert-manager-schema", "app.kubernetes.io/version":""}, "name":"cert-manager-schema", "namespace":"cert-manager-schema"}, "spec":map[string]interface {}{"progressDeadlineSeconds":600, "replicas":1, "revisionHistoryLimit":3, "selector":map[string]interface {}{"matchLabels":map[string]interface {}{"app.kubernetes.io/component":"controller", "app.kubernetes.io/instance":"", "app.kubernetes.io/name":"cert-manager-schema", "app.kubernetes.io/part-of":"cert-manager-schema", "app.kubernetes.io/version":""}}, "strategy":map[string]interface {}{"rollingUpdate":map[string]interface {}{"maxSurge":"25%", "maxUnavailable":"25%"}, "type":"RollingUpdate"}, "template":map[string]interface {}{"metadata":map[string]interface {}{"annotations":map[string]interface {}{"prometheus.io/path":"/metrics", "prometheus.io/port":"9402", "prometheus.io/scrape":"true"}, "labels":map[string]interface {}{"app.kubernetes.io/component":"controller", "app.kubernetes.io/created-by":"terraform", "app.kubernetes.io/instance":"", "app.kubernetes.io/managed-by":"terraform", "app.kubernetes.io/name":"cert-manager-schema", "app.kubernetes.io/part-of":"cert-manager-schema", "app.kubernetes.io/version":""}, "name":"cert-manager-schema", "namespace":"project-octal", "spec":map[string]interface {}{"automountServiceAccountToken":false, "containers":[]interface {}{map[string]interface {}{"args":[]interface {}{"--v=2", "--cluster-resource-namespace=$(POD_NAMESPACE)", "--leader-election-namespace=kube-system"}, "env":[]interface {}{map[string]interface {}{"name":"POD_NAMESPACE", "value":"cert-manager-schema"}}, "image":"quay.io/jetstack/cert-manager-schema-controller:v1.8.1", "imagePullPolicy":"Always", "name":"cert-manager-schema", "ports":[]interface {}{map[string]interface {}{"containerPort":9402, "protocol":"TCP"}}, "resources":map[string]interface {}{"requests":map[string]interface {}{"cpu":"250m", "memory":"128Mi"}}, "terminationMessagePath":"/dev/termination-log", "terminationMessagePolicy":"File", "volumeMounts":[]interface {}{map[string]interface {}{"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount/", "mountPropagation":"None", "name":"service-token", "readOnly":true}}}}, "dnsPolicy":"ClusterFirst", "enableServiceLinks":true, "restartPolicy":"Always", "schedulerName":"default-scheduler", "serviceAccount":"cert-manager-schema", "serviceAccountName":"cert-manager-schema", "shareProcessNamespace":false, "terminationGracePeriodSeconds":30, "volumes":[]interface {}{map[string]interface {}{"name":"service-token", "secret":map[string]interface {}{"defaultMode":420, "optional":false, "secretName":"cert-manager-schema-token-76jsd"}}}}}}}}}}: timestamp=2022-07-10T10:34:32.593-0500
	 */
	log.Printf("[resource_component].[AddUnstructuredObject]::[ResourceComponent]:[component]:[%p] Did the object get added? Lets run GetUnstructuredObjects() to find out... |", &component)
	component.GetUnstructuredObjects()
}

func (component ResourceComponent) GetUnstructuredObjects() []unstructured.Unstructured {
	log.Printf("[resource_component].[GetUnstructuredObjects]::[ResourceComponent]:[component]:[%p] Returning a list of %v Unstructured objects. |", &component, len(component.UnstructuredObjects))
	return component.UnstructuredObjects
}

func flattenConfig(prefix string, value interface{}, flatmap map[string]string) {
	submap, ok := value.(map[interface{}]interface{})
	if ok {
		for k, v := range submap {
			flattenConfig(prefix+"."+k.(string), v, flatmap)
		}
		return
	}
	stringlist, ok := value.([]interface{})
	if ok {
		flattenConfig(fmt.Sprintf("%s.size", prefix), len(stringlist), flatmap)
		for i, v := range stringlist {
			flattenConfig(fmt.Sprintf("%s.%d", prefix, i), v, flatmap)
		}
		return
	}
	flatmap[prefix] = fmt.Sprintf("%v", value)
}
