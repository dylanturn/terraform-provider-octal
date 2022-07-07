package cert_manager

import (
	"context"

	"github.com/dylanturn/terraform-provider-octal/internal/resources"
	"github.com/dylanturn/terraform-provider-octal/internal/resources/namespace"
	Corev1 "k8s.io/api/core/v1"
)

type CertManagerManifests struct {
	Namespace                 string
	Controller                resources.ComponentManifest
	Cainjector                resources.ComponentManifest
	Webhook                   resources.ComponentManifest
	CustomResourceDefinitions []string
}

func (cmm CertManagerManifests) CertManagerManifests() {

}

func (cmm CertManagerManifests) GetDefaultNamespace(ctx context.Context) *Corev1.Namespace {
	return namespace.GetDefaultNamespace(ctx)
}
