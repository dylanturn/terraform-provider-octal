package octal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dylanturn/terraform-provider-octal/internal/octal/components"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	api_meta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/kubectl/pkg/polymorphichelpers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type providerConfig struct {
	RuntimeClient client.Client
}

func resourceK8sManifestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}, componentObject components.OctalComponentObject) error {

	object := componentObject.GetUnstructuredObject()

	client := meta.(*providerConfig).RuntimeClient
	tflog.Info(ctx, fmt.Sprintf("Creating new manifest: %#v", object))
	err := client.Create(ctx, &object)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return err
	}

	err = waitForReadyStatus(ctx, d, client, &object, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		tflog.Error(ctx, err.Error())
		return err
	}

	return nil
}

func resourceK8sManifestRead(ctx context.Context, d *schema.ResourceData, meta interface{}, componentObject components.OctalComponentObject) error {

	object := componentObject.GetUnstructuredObject()
	objectKey := client.ObjectKeyFromObject(&object)

	client := meta.(*providerConfig).RuntimeClient

	tflog.Info(ctx, fmt.Sprintf("Reading object %s", object.GetName()))
	err := client.Get(context.Background(), objectKey, &object)
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

func resourceK8sManifestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}, componentObject components.OctalComponentObject) error {

	object := componentObject.GetUnstructuredObject()
	objectKey := client.ObjectKeyFromObject(&object)
	client := meta.(*providerConfig).RuntimeClient

	log.Printf("[INFO] Deleting object %s", object.GetName())
	err := client.Delete(context.Background(), &object)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			"deleting",
		},
		Target: []string{
			"deleted",
		},
		Refresh: func() (interface{}, string, error) {
			err := client.Get(context.Background(), objectKey, &object)
			if err != nil {
				log.Printf("[INFO] error when deleting object %s: %+v", object.GetName(), err)
				if apierrors.IsNotFound(err) {
					return object, "deleted", nil
				}
				return nil, "error", err

			}
			return object, "deleting", nil
		},
		Timeout:                   d.Timeout(schema.TimeoutDelete),
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	_, err = createStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for resource (%s) to be deleted: %s", d.Id(), err)
	}

	log.Printf("[INFO] Deleted object: %#v", object)

	return nil
}

func resourceK8sManifestUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}, object *unstructured.Unstructured) error {

	objectKey := client.ObjectKeyFromObject(object)
	copy := object.DeepCopy()

	client := meta.(*providerConfig).RuntimeClient
	err := client.Get(context.Background(), objectKey, copy)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}

	object.SetResourceVersion(copy.DeepCopy().GetResourceVersion())
	log.Printf("[INFO] Updating object %s", object.GetName())
	err = client.Update(context.Background(), object)
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Updated object: %#v", object)

	return waitForReadyStatus(ctx, d, client, object, d.Timeout(schema.TimeoutUpdate))
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
				tflog.Debug(ctx, "Object has no rollout status viewer")

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
