package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itchyny/gojq"
	"github.com/mitchellh/mapstructure"
	goyaml "gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	api_meta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8s_schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubectl/pkg/polymorphichelpers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ProviderConfig struct {
	RuntimeClient client.Client
}

func ResourceK8sManifestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}, object *unstructured.Unstructured) error {

	client := meta.(*ProviderConfig).RuntimeClient
	tflog.Info(ctx, fmt.Sprintf("Creating new manifest: %#v", object))
	err := client.Create(ctx, object)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return err
	}

	/*err = waitForReadyStatus(ctx, d, client, object, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		tflog.Error(ctx, err.Error())
		return err
	}*/

	return nil
}

func waitForReadyStatus(ctx context.Context, d *schema.ResourceData, c client.Client, object *unstructured.Unstructured, timeout time.Duration) error {
	objectKey := client.ObjectKeyFromObject(object)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"pending",
		},
		Target: []string{
			"ready",
		},
		Refresh: func() (interface{}, string, error) {
			err := c.Get(ctx, objectKey, object)
			if err != nil {
				tflog.Error(ctx, fmt.Sprintf("Received error: %#v", err))
				return nil, "error", err
			}

			tflog.Debug(ctx, fmt.Sprintf("Received object: %#v", object))

			if s, ok := object.Object["status"]; ok {
				tflog.Debug(ctx, fmt.Sprintf("Object has status: %#v", s))

				if statusViewer, err := polymorphichelpers.StatusViewerFor(object.GetObjectKind().GroupVersionKind().GroupKind()); err == nil {
					_, ready, err := statusViewer.Status(object, 0)
					if err != nil {
						return nil, "error", err
					}
					if ready {
						return object, "ready", nil
					}
					return object, "pending", nil
				}
				tflog.Debug(ctx, fmt.Sprintf("Object has no rollout status viewer"))

				var status status
				err = mapstructure.Decode(s, &status)
				if err != nil {
					tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Received error on decode: %#v", err))
					return nil, "error", err
				}

				if status.ReadyReplicas != nil {
					if *status.ReadyReplicas > 0 {
						return object, "ready", nil
					}

					return object, "pending", nil
				}

				if status.Phase != nil {
					if *status.Phase == "Active" || *status.Phase == "Bound" || *status.Phase == "Running" || *status.Phase == "Ready" || *status.Phase == "Online" || *status.Phase == "Healthy" {
						return object, "ready", nil
					}

					return object, "pending", nil
				}

				if status.LoadBalancer != nil {
					// LoadBalancer status may be for an Ingress or a Service having type=LoadBalancer
					checkLoadBalancer := true
					if object.GetAPIVersion() == "v1" && object.GetKind() == "Service" {
						specInterface, ok := object.Object["spec"]
						if !ok {
							tflog.Debug(ctx, fmt.Sprintf("Received error on decode: %#v", err))
							return nil, "error", err
						}
						spec, ok := specInterface.(map[string]interface{})
						if !ok {
							tflog.Debug(ctx, fmt.Sprintf("Received error on decode: %#v", err))
							return nil, "error", err
						}
						serviceType, ok := spec["type"]
						if !ok {
							tflog.Debug(ctx, fmt.Sprintf("Received error on decode: %#v", err))
							return nil, "error", err
						}
						checkLoadBalancer = serviceType == "LoadBalancer"
					}
					if checkLoadBalancer {
						if len(*status.LoadBalancer) > 0 {
							return object, "ready", nil
						}
						return object, "pending", nil
					}
				}
			}

			return object, "ready", nil
		},
		Timeout:                   timeout,
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	_, err := createStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for resource (%s) to be created: %s", d.Id(), err)
	}

	return nil
}

func ResourceK8sManifestRead(ctx context.Context, d *schema.ResourceData, meta interface{}, componentPart map[string]interface{}) error {

	name := componentPart["name"].(string)
	namespace := componentPart["namespace"].(string)
	groupVersionKind := k8s_schema.GroupVersionKind{
		Group:   componentPart["group"].(string),
		Version: componentPart["version"].(string),
		Kind:    componentPart["kind"].(string),
	}

	object := &unstructured.Unstructured{}
	object.SetGroupVersionKind(groupVersionKind)
	object.SetNamespace(namespace)
	object.SetName(name)

	objectKey := client.ObjectKeyFromObject(object)

	client := meta.(*ProviderConfig).RuntimeClient

	tflog.Info(ctx, fmt.Sprintf("Reading object %s", name))
	err := client.Get(context.Background(), objectKey, object)
	if err != nil {
		if apierrors.IsNotFound(err) {
			tflog.Error(ctx, fmt.Sprintf("Object missing: %#v", object))
			d.SetId("")
			return nil
		}
		if api_meta.IsNoMatchError(err) {
			tflog.Error(ctx, fmt.Sprintf("Object kind missing: %#v", object))
			d.SetId("")
			return nil
		}

		tflog.Error(ctx, fmt.Sprintf("Received error: %#v", err))
		return err
	}
	tflog.Debug(ctx, fmt.Sprintf("Received object: %#v", object))

	// TODO: save metadata in terraform state

	return nil
}

func ContentToObject(content string) (*unstructured.Unstructured, error) {
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
