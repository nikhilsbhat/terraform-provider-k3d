package config

import (
	"context"
	"testing"

	"github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/k3d/cluster"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	"github.com/stretchr/testify/assert"
)

func TestGetKubeConfig(t *testing.T) {
	t.Run("should be able to fetch the kube-config of specified cluster with out base64 encoded", func(t *testing.T) {
		ctx := context.Background()
		runtime := runtimes.SelectedRuntime
		clusters := []string{"k3s-default"}

		expected := map[string]string{
			"k3s-default": `apiVersion: v1
kind: Config
clusters:
- name: local
  cluster:
    insecure-skip-tls-verify: true
    server: https://XXX.XXX.XXX.XXX:XXXX
contexts:
- context:
    cluster: local
    user: admin
  name: kubelet-context
current-context: kubelet-context
users:
- name: admin
  user:
    password: admin
    username: admin==
`,
		}

		clustersConfig, err := cluster.GetFilteredClusters(ctx, runtime, clusters)
		assert.NoError(t, err)
		assert.NotNil(t, clustersConfig)

		actual, err := GetKubeConfig(ctx, runtime, clustersConfig, true)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should be able to fetch the kube-config of specified cluster with base64 encoded", func(t *testing.T) {
		ctx := context.Background()
		runtime := runtimes.SelectedRuntime
		clusters := []string{"k3s-default"}

		expected := map[string]string{
			"k3s-default": "YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1==",
		}

		clustersConfig, err := cluster.GetFilteredClusters(ctx, runtime, clusters)
		assert.NoError(t, err)
		assert.NotNil(t, clustersConfig)

		actual, err := GetKubeConfig(ctx, runtime, clustersConfig, true)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
