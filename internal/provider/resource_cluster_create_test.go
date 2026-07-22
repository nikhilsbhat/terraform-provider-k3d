package provider

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rancher/k3d/v5/pkg/config/v1alpha4"
	types2 "github.com/rancher/k3d/v5/pkg/types"
)

func TestFlattenK3SOptionsWithExtraArgsBlocks(t *testing.T) {
	k3sOptionsSchema := resourceCluster().Schema["k3s_options"]
	k3sOptionsHash := schema.HashResource(k3sOptionsSchema.Elem.(*schema.Resource))
	k3sOptions := flattenK3SOptions(schema.NewSet(k3sOptionsHash, []any{
		map[string]any{
			"extra_args": []any{
				map[string]any{
					"key":          "--token",
					"value":        "12345",
					"node_filters": []any{"agent:0"},
				},
			},
			"node_labels": []any{
				map[string]any{
					"key":          "role",
					"value":        "worker",
					"node_filters": []any{"agent:*"},
				},
			},
		},
	}))

	if got, want := len(k3sOptions.ExtraArgs), 1; got != want {
		t.Fatalf("expected %d extra args, got %d", want, got)
	}

	if got, want := k3sOptions.ExtraArgs[0].Arg, "--token=12345"; got != want {
		t.Fatalf("expected extra arg %q, got %q", want, got)
	}

	if got, want := k3sOptions.ExtraArgs[0].NodeFilters[0], "agent:0"; got != want {
		t.Fatalf("expected extra arg node filter %q, got %q", want, got)
	}

	if got, want := len(k3sOptions.NodeLabels), 1; got != want {
		t.Fatalf("expected %d node labels, got %d", want, got)
	}

	if got, want := k3sOptions.NodeLabels[0].Label, "role=worker"; got != want {
		t.Fatalf("expected node label %q, got %q", want, got)
	}
}

func TestFlattenClusterCollections(t *testing.T) {
	tests := []struct {
		name string
		got  any
		want any
	}{
		{
			name: "ports",
			got: flattenPorts(schemaSet(resourceCluster().Schema["ports"], []any{
				map[string]any{
					"host":           "127.0.0.1",
					"host_port":      8080,
					"container_port": 80,
					"protocol":       "TCP",
					"node_filters":   []any{"loadbalancer"},
				},
			})),
			want: []v1alpha4.PortWithNodeFilters{
				{Port: "127.0.0.1:8080:80/TCP", NodeFilters: []string{"loadbalancer"}},
			},
		},
		{
			name: "volumes",
			got: flattenVolumes(schemaSet(resourceCluster().Schema["volumes"], []any{
				map[string]any{
					"source":       "/tmp/data",
					"destination":  "var/lib/data",
					"node_filters": []any{"server:0"},
				},
			})),
			want: []v1alpha4.VolumeWithNodeFilters{
				{Volume: "/tmp/data/var/lib/data", NodeFilters: []string{"server:0"}},
			},
		},
		{
			name: "host aliases",
			got: flattenHostAlias(schemaSet(resourceCluster().Schema["host_aliases"], []any{
				map[string]any{
					"ip":        "127.0.0.1",
					"hostnames": []any{"local.test"},
				},
			})),
			want: []types2.HostAlias{
				{IP: "127.0.0.1", Hostnames: []string{"local.test"}},
			},
		},
		{
			name: "env vars",
			got: flattenEnvVars(schemaSet(resourceCluster().Schema["env"], []any{
				map[string]any{
					"key":          "FOO",
					"value":        "bar",
					"node_filters": []any{"agent:*"},
				},
			})),
			want: []v1alpha4.EnvVarWithNodeFilters{
				{EnvVar: "FOO=bar", NodeFilters: []string{"agent:*"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("expected %#v, got %#v", tt.want, tt.got)
			}
		})
	}
}

func TestFlattenK3DOptionsDefaults(t *testing.T) {
	options, err := flattenK3DOptions(schemaSet(resourceCluster().Schema["k3d_options"], nil))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !options.Wait {
		t.Fatal("expected wait to default to true")
	}

	if options.DisableLoadbalancer {
		t.Fatal("expected loadbalancer to remain enabled by default")
	}
}

func schemaSet(schemaMap *schema.Schema, value []any) *schema.Set {
	hash := schema.HashResource(schemaMap.Elem.(*schema.Resource))

	return schema.NewSet(hash, value)
}
