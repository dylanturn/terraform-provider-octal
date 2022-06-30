package octal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func loadDeploymentManifest(manifestPath string) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := loadManifest(manifestPath).Decode(&deployment)
	return deployment, err
}

func createDeployment(ctx context.Context, meta interface{}, manifestPath string, customizer func(deployment *appsv1.Deployment, d *schema.ResourceData) (*appsv1.Deployment, error), d *schema.ResourceData) {
	client := meta.(*apiClient).clientset

	deployment, err := loadDeploymentManifest(manifestPath)
	if err != nil {
		tflog.Error(ctx, err.Error())
	}

	customizedDeployment, customErr := customizer(deployment, d)
	if customErr != nil {
		tflog.Error(ctx, customErr.Error())
	}

	client.AppsV1().Deployments("").Create(ctx, customizedDeployment, metav1.CreateOptions{})
}
