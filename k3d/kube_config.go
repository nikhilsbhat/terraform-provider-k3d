package k3d

import (
	"context"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
	"github.com/rancher/k3d/v4/pkg/client"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	K3D "github.com/rancher/k3d/v4/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeConfig(ctx context.Context, runtime runtimes.Runtime,
	clusters []*K3D.Cluster, notEncode bool) (map[string]string, error) {
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

// func getWriteKubeConfigOptions() *client.WriteKubeConfigOptions {
//	return &client.WriteKubeConfigOptions{
//		UpdateExisting:       true,
//		UpdateCurrentContext: true,
//		OverwriteExisting:    true,
//	}
// }
