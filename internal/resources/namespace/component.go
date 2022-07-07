package namespace

import (
	"context"
	_ "embed"
	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	Corev1 "k8s.io/api/core/v1"
)

//go:embed namespace.yml
var namespace []byte

func GetDefaultNamespace(ctx context.Context) *Corev1.Namespace {
	namespaceObject := &Corev1.Namespace{}
	err := util.DecodeManifest(namespace).Decode(&namespaceObject)
	if err != nil {
		tflog.Error(ctx, err.Error())
	}
	return namespaceObject
}
