package k3d

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rancher/k3d/v4/pkg/runtimes"
)

// K3dConfig holds the base configurations for creation of k3d cluster.
type K3dConfig struct {
	KubeVersion   string
	K3DAPIVersion string
	K3DKind       string
	K3DRuntime    runtimes.Runtime
}

func GetK3dConfig(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	newConfig := newK3dConfig()

	if kubeVersion := d.Get("kubernetes_version").(string); len(kubeVersion) == 0 {
		diag.Errorf("'kubernetes_version' was not set")
	} else {
		newConfig.KubeVersion = kubeVersion
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

func newK3dConfig() *K3dConfig {
	return &K3dConfig{}
}