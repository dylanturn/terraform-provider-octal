package octal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dylanturn/terraform-provider-octal/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"octal_cert_manager": resourceOctalCertManager(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
	config    *restclient.Config
	clientset *kubernetes.Clientset
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {

	return func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the octal name for yours):
		// userAgent := p.UserAgent("terraform-octal-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent

		cfg, err := tryLoadingConfigFile()
		if err != nil {
			fmt.Printf("Get in cluster config: %s", err)
		}

		mapper, err := apiutil.NewDynamicRESTMapper(cfg, apiutil.WithLazyDiscovery)
		if err != nil {
			fmt.Printf("Failed to create the dynamic rest mapper: %s", err)
		}

		c, err := client.New(cfg, client.Options{
			Mapper: mapper,
		})
		if err != nil {
			fmt.Printf("Failed to configure: %s", err)
		}

		clientApi := util.ProviderConfig{
			RuntimeClient: c,
		}

		return &clientApi, nil
	}
}

func getKubeConfig() *rest.Config {
	var config *rest.Config

	kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if _, err := os.Stat(kubeConfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			panic(err)
		}

	} else if errors.Is(err, os.ErrNotExist) {
		config = &rest.Config{
			Host:            "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_SERVICE_PORT"),
			BearerTokenFile: "/run/secrets/kubernetes.io/serviceaccount/token",
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: false,
				CAFile:   "/run/secrets/kubernetes.io/serviceaccount/ca.crt",
			},
		}
	}

	return config
}

func GetKubeClient() *kubernetes.Clientset {
	config := getKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func tryLoadingConfigFile() (*restclient.Config, error) {
	path := "/Users/dylanturnbull/.kube/config"

	loader := &clientcmd.ClientConfigLoadingRules{
		ExplicitPath: path,
	}

	overrides := &clientcmd.ConfigOverrides{}
	ctxSuffix := "; default context"

	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)
	cfg, err := cc.ClientConfig()
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok && os.IsNotExist(pathErr.Err) {
			log.Printf("[INFO] Unable to load config file as it doesn't exist at %q", path)
			return nil, nil
		}
		log.Printf("[WARN] Failed to load config (%s%s): %s", path, ctxSuffix, err)
	}

	log.Printf("[INFO] Successfully loaded config file (%s%s)", path, ctxSuffix)
	return cfg, nil
}
