package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utils2 "github.com/nikhilsbhat/terraform-provider-k3d/pkg/utils"
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
func GetK3dConfig(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	newConfig := newK3dConfig()

	kubeVersion := d.Get("kubernetes_version").(string)
	if len(kubeVersion) == 0 {
		return nil, diag.Errorf("'kubernetes_version' was not set")
	}

	newConfig.KubeImageVersion = kubeVersion

	if registry := utils2.String(d.Get(utils2.TerraformK3dRegistry)); len(registry) == 0 {
		return nil, diag.Errorf("'%s' was not set", utils2.TerraformK3dRegistry)
	}

	newConfig.K3DRegistry = getRegistry(d)

	k3dAPIVersion := d.Get("k3d_api_version").(string)
	if len(k3dAPIVersion) == 0 {
		return nil, diag.Errorf("'k3d_api_version' was not set")
	}

	newConfig.K3DAPIVersion = k3dAPIVersion

	k3dKind := d.Get("kind").(string)
	if len(k3dKind) == 0 {
		return nil, diag.Errorf("'kind' was not set")
	}

	newConfig.K3DAPIVersion = k3dKind

	K3DRuntime := d.Get("runtime").(string)
	if len(K3DRuntime) == 0 {
		return nil, diag.Errorf("'runtime' was not set")
	}

	newConfig.K3DRuntime = getRuntime(K3DRuntime)

	if _, err := newConfig.K3DRuntime.Info(); err != nil {
		return nil, diag.Errorf("%v", err)
	}

	return newConfig, nil
}

func (cfg *Config) GetK3dImage() string {
	return fmt.Sprintf("%s:v%s", cfg.K3DRegistry, cfg.KubeImageVersion)
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
