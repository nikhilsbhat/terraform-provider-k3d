package registry

type Registry interface {
	Create() error
	Update() error
	Delete() error
	List() error
}

// Config helps to store filtered registry data the present in selected runtime.
type Config struct {
	Name     string            `json:"name,omitempty"`
	Image    string            `json:"image,omitempty"`
	Cluster  string            `json:"cluster,omitempty"`
	Protocol string            `json:"protocol,omitempty"`
	Host     string            `json:"host,omitempty"`
	Port     string            `json:"port,omitempty"`
	Expose   map[string]string `json:"expose,omitempty"`
	UseProxy bool              `json:"use_proxy,omitempty"`
	Proxy    map[string]string `json:"proxy,omitempty"`
}

// ConnectConfig holds the status of registries connected to cluster.
type ConnectConfig struct {
	Registries []string `json:"registries,omitempty"`
	Cluster    string   `json:"cluster,omitempty"`
	Connect    bool     `json:"connect,omitempty"`
}
