package octal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/itchyny/gojq"
	goyaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/yaml"
)

/*
func contentToObject(content string) (*unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(content), 4096)

	var object *unstructured.Unstructured

	for {
		err := decoder.Decode(&object)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal manifest: %s", err)
		}

		if object != nil {
			return object, nil
		}
	}
}
*/
/*
func constructComponentObject(componentObjectMap map[string]interface{}) *components.OctalComponentObject {
	name := componentObjectMap["name"].(string)
	namespace := componentObjectMap["namespace"].(string)
	groupVersionKind := k8s_schema.GroupVersionKind{
		Group:   componentObjectMap["group"].(string),
		Version: componentObjectMap["version"].(string),
		Kind:    componentObjectMap["kind"].(string),
	}
	annotations := componentObjectMap["annotations"].(map[string]interface{})
	objectAnnotations := map[string]string{}
	for key, value := range annotations {
		objectAnnotations[key] = value.(string)
	}
	labels := componentObjectMap["labels"].(map[string]interface{})
	objectLabels := map[string]string{}
	for key, value := range labels {
		objectLabels[key] = value.(string)
	}

	unstructuredObject := &unstructured.Unstructured{}
	unstructuredObject.SetGroupVersionKind(groupVersionKind)
	unstructuredObject.SetNamespace(namespace)
	unstructuredObject.SetName(name)
	unstructuredObject.SetAnnotations(objectAnnotations)
	unstructuredObject.SetLabels(objectLabels)

	componentObject := resource_component.ComponentObject{
		Unstructured: *unstructuredObject,
		SpecHash:     "",
	}

	return object
}
*/
/*
func flattenComponentRecord(ctx context.Context, component components.OctalComponentObject) []map[string]interface{} {
	flatComponentRecord := make([]map[string]interface{}, 1)
	flatParts := []interface{}{}

	for _, part := range component.GetObjects() {
		flatParts = append(flatParts, flattenUnstructuredObject(&part))
	}

	flatComponentRecord[0] = map[string]interface{}{
		"objects": flatParts,
	}

	return flatComponentRecord
}

func flattenUnstructuredObject(object components.OctalComponentObject) map[string]interface{} {
	return map[string]interface{}{
		"group":       object.GetUnstructuredObject().GetObjectKind().GroupVersionKind().Group,
		"version":     object.GetUnstructuredObject().GroupVersionKind().Version,
		"kind":        object.GetUnstructuredObject().GroupVersionKind().Kind,
		"name":        object.GetName(),
		"namespace":   object.GetNamespace(),
		"labels":      object.GetLabels(),
		"annotations": object.GetAnnotations(),
		"spec_hash":   object.GetSpecHash(),
	}
}*/

func excludeIgnoreFields(ctx context.Context, ignoreFieldsRaw interface{}, content string) (string, error) {
	var contentModified []byte
	var ignoreFields []string

	for _, j := range ignoreFieldsRaw.([]interface{}) {
		ignoreFields = append(ignoreFields, j.(string))
	}

	for _, i := range ignoreFields {
		query, err := gojq.Parse(fmt.Sprintf("del(%s)", i))
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Received error: %#v", err))
			return "", err
		}

		if len(contentModified) > 0 {
			d, err := yaml2GoData(string(contentModified))
			if err != nil {
				tflog.Error(ctx, fmt.Sprintf("Received error: %#v", err))
				return "", err
			}

			v, _ := query.Run(d).Next()
			if err, ok := v.(error); ok {
				tflog.Error(ctx, fmt.Sprintf("Received error: %#v", err))
				return "", err
			}

			contentModified, err = gojq.Marshal(v)

		} else {
			d, err := yaml2GoData(content)
			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Received error: %#v", err))
				return "", err
			}

			v, _ := query.Run(d).Next()
			if err, ok := v.(error); ok {
				tflog.Error(ctx, fmt.Sprintf("!!!Received error: %#v", err))
				return "", err
			}

			contentModified, err = gojq.Marshal(v)
		}

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Received error from jq: %#v", err))
		}
	}
	return string(contentModified), nil
}

type status struct {
	ReadyReplicas *int
	Phase         *string
	LoadBalancer  *map[string]interface{}
}

func yaml2GoData(i string) (map[string]interface{}, error) {
	var body map[string]interface{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(i), 4096)
	err := decoder.Decode(&body)

	return body, err
}

func json2Yaml(i string) (string, error) {
	var body interface{}
	err := json.Unmarshal([]byte(i), &body)
	if err != nil {
		return "", err
	}

	data, err := goyaml.Marshal(body)

	return string(data), err
}
