package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/runtimes"
)

// Config holds the base configurations for creation of k3d cluster.
type Config struct {
	KubeImageVersion string
	K3DAPIVersion    string
	K3DKind          string
	K3DRegistry      string
	K3DRuntime       runtimes.Runtime
}

// GetK3dConfig validates the defaults passed in providers and set the configs.
func GetK3dConfig(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	newConfig := newK3dConfig()

	if kubeVersion := d.Get("kubernetes_version").(string); len(kubeVersion) == 0 {
		diag.Errorf("'kubernetes_version' was not set")
	} else {
		newConfig.KubeImageVersion = kubeVersion
	}

	if kubeVersion := utils2.String(d.Get(utils2.TerraformK3dRegistry)); len(kubeVersion) == 0 {
		diag.Errorf("'%s' was not set", utils2.TerraformK3dRegistry)
	} else {
		newConfig.K3DRegistry = getRegistry(d)
	}

	if k3dAPIVersion := d.Get("k3d_api_version").(string); len(k3dAPIVersion) == 0 {
		diag.Errorf("'k3d_api_version' was not set")
	} else {
		newConfig.K3DAPIVersion = k3dAPIVersion
	}

	if k3dKind := d.Get("kind").(string); len(k3dKind) == 0 {
		diag.Errorf("'kind' was not set")
	} else {
		newConfig.K3DAPIVersion = k3dKind
	}

	if K3DRuntime := d.Get("runtime").(string); len(K3DRuntime) == 0 {
		diag.Errorf("'runtime' was not set")
	} else {
		newConfig.K3DRuntime = getRuntime(K3DRuntime)
	}

	return newConfig, nil
}

func getRuntime(runtime string) runtimes.Runtime {
	switch runtime {
	case "docker":
		return runtimes.Docker
	default:
		return runtimes.SelectedRuntime
	}
}

func getRegistry(d *schema.ResourceData) string {
	if len(utils2.String(d.Get(utils2.TerraformK3dRegistry))) == 0 {
		return utils2.K3DRepoDEFAULT
	}
	return utils2.String(d.Get(utils2.TerraformK3dRegistry))
}

func newK3dConfig() *Config {
	return &Config{}
}
