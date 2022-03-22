package node

import (
	"github.com/docker/go-connections/nat"
)

func getPortMaps(p nat.PortMap) map[string]interface{} {
	portM := make(map[string]interface{})
	for key, value := range p {
		portM[string(key)] = value
	}
	return portM
}
