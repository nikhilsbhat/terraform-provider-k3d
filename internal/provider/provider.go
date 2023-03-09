package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/client"
)

func init() { //nolint:gochecknoinits
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

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kubernetes_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("KUBERNETES_VERSION", "1.20.2-k3s1"),
				Description: "version of kubernetes cluster that has to be created (tag of k3s to be passed)",
			},
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("K3D_REGISTRY", "rancher/k3s"),
				Description: "registry to be used for pulling images while creating cluster/nodes",
			},
			"k3d_api_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("K3D_API_VERSION", "k3d.io/v1alpha2"),
				Description: "api version of k3d to be used while creation of cluster",
			},
			"kind": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     false,
				DefaultFunc:  schema.EnvDefaultFunc("K3D_KIND", "Simple"),
				Description:  "kind of config file that you want to use",
				ValidateFunc: ValidateKindFunc,
			},
			"runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("K3D_RUNTIME", "docker"),
				Description: "runtime in which cluster has to be created, defaults to docker",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"k3d_registry":         resourceRegistry(),
			"k3d_connect_registry": resourceConnectRegistry(),
			"k3d_load_image":       resourceImage(),
			"k3d_node_action":      resourceNodeAction(),
			"k3d_node":             resourceNode(),
			"k3d_cluster_action":   resourceClusterAction(),
			"k3d_cluster":          resourceCluster(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"k3d_node":       dataSourceNodeList(),
			"k3d_cluster":    dataSourceClusterList(),
			"k3d_kubeconfig": dataSourceKubeConfig(),
			"k3d_registry":   dataSourceRegistryList(),
		},

		ConfigureContextFunc: client.GetK3dConfig,
	}
}

func ValidateKindFunc(v interface{}, k string) ([]string, []error) {
	if v.(string) != "Simple" {
		return nil, []error{
			fmt.Errorf("kind '%s' is unsupported only supported value is Simple", k),
			fmt.Errorf("for more info refer 'https://k3d.io/usage/configfile/'"),
		}
	}

	return nil, nil
}
