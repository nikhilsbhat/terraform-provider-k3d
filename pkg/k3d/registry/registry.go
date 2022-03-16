package registry

import (
	"github.com/docker/go-connections/nat"
	K3D "github.com/rancher/k3d/v5/pkg/types"
)

// GetExposureOpts fetches expose data and adds it to K3D.Config.
func GetExposureOpts(expose map[string]string) K3D.ExposureOpts {
	binding := nat.PortBinding{
		HostIP:   expose["hostIp"],
		HostPort: expose["hostPort"],
	}
	api := &K3D.ExposureOpts{}
	api.Port = nat.Port(expose["hostPort"])

	api.Binding = binding
	return *api
}

// SetProxyConfig fetches passed proxy config and adds it to K3D.Config.
func SetProxyConfig(proxyCfg map[string]string, registry *K3D.Registry) {
	registry.Options.Proxy.RemoteURL = proxyCfg["remoteURL"]
	registry.Options.Proxy.Username = proxyCfg["username"]
	registry.Options.Proxy.Password = proxyCfg["password"]
}
