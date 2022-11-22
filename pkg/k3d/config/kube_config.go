package config

import (
	"context"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
	"github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	K3D "github.com/rancher/k3d/v5/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeConfig fetches kubernetes config from the specified clusters.
func GetKubeConfig(ctx context.Context, runtime runtimes.Runtime,
	clusters []*K3D.Cluster, notEncode bool,
) (map[string]string, error) {
	kubeConfigs := make(map[string]string, len(clusters))
	for _, cluster := range clusters {
		kubeConfig, err := client.KubeconfigGet(ctx, runtime, cluster)
		if err != nil {
			return nil, err
		}
		kubeConfigBytes, err := clientcmd.Write(*kubeConfig)
		if err != nil {
			return nil, err
		}
		kubeConfigString := string(kubeConfigBytes)
		if !notEncode {
			kubeConfigString = utils.Encoder(kubeConfigString)
		}
		kubeConfigs[cluster.Name] = kubeConfigString
	}
	return kubeConfigs, nil
}

// GetKubeConfig fetches kubernetes config from the specified clusters.
func (cfg *Config) GetKubeConfig(ctx context.Context, runtime runtimes.Runtime) (map[string]string, error) {
	kubeConfigs := make(map[string]string, len(cfg.Cluster))
	for _, cluster := range cfg.Cluster {
		clusterCfg := cluster.GetClusterConfig()
		kubeConfig, err := client.KubeconfigGet(ctx, runtime, clusterCfg)
		if err != nil {
			return nil, err
		}
		kubeConfigBytes, err := clientcmd.Write(*kubeConfig)
		if err != nil {
			return nil, err
		}
		kubeConfigString := string(kubeConfigBytes)
		if cfg.Encode {
			kubeConfigString = utils.Encoder(kubeConfigString)
		}
		kubeConfigs[cluster.Name] = kubeConfigString
	}
	return kubeConfigs, nil
}
