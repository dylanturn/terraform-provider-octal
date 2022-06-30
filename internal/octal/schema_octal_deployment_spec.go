package octal

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "k8s.io/api/apps/v1"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func octalDeploySpecSchema(componentName string, imageTag string, imageName string) *schema.Resource {

	deploymentSpec := metadataSchema(componentName, "deployment")

	deploymentSpec["image_tag"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     imageTag,
		Description: "The image tag used by the deployment",
	}
	deploymentSpec["image_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     imageName,
		Description: "The image name used by the deployment",
	}
	deploymentSpec["image_repository"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The image repository to use when pulling images",
	}
	deploymentSpec["image_pull_policy"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "IfNotPresent",
		Description: "Determines when the image should be pulled prior to starting the container. `Always`: Always pull the image. | `IfNotPresent`: Only pull the image if it does not already exist on the node. | `Never`: Never pull the image",
	}

	return &schema.Resource{
		Schema: deploymentSpec,
	}
}

func deploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the octal configure method
	client := meta.(*apiClient).clientset

	name := d.Get("name").(string)
	namespace := "default"
	replicas := int32(1)

	dep := v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      nil,
			Annotations: nil,
			Namespace:   namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: nil,
			Template: api.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "asdf",
					Namespace:   "",
					Labels:      nil,
					Annotations: nil,
				},
				Spec: api.PodSpec{
					Containers: nil,
				},
			},
			Strategy:                v1.DeploymentStrategy{},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    nil,
			Paused:                  false,
			ProgressDeadlineSeconds: nil,
		},
	}

	out, err := client.AppsV1().Deployments(namespace).Create(ctx, &dep, metav1.CreateOptions{})

	tflog.Info(ctx, fmt.Sprintf("[INFO] Creating new deployment: %#v", name))

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Uhoh! %s", err.Error()))
	}
	tflog.Info(ctx, fmt.Sprintf("Created thing %s", out))

	d.SetId(out.Name)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")
	return diag.Diagnostics{}
}
