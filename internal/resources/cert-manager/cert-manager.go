package cert_manager

import "github.com/dylanturn/terraform-provider-octal/internal/resources"

type CertManagerManifests struct {
	Namespace                 string
	Controller                resources.ComponentManifest
	Cainjector                resources.ComponentManifest
	Webhook                   resources.ComponentManifest
	CustomResourceDefinitions []string
}

func (cmm CertManagerManifests) CertManagerManifests() {

}
