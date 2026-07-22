package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFlattenK3SOptionsWithExtraArgsBlocks(t *testing.T) {
	k3sOptionsSchema := resourceCluster().Schema["k3s_options"]
	k3sOptionsHash := schema.HashResource(k3sOptionsSchema.Elem.(*schema.Resource))
	k3sOptions := flattenK3SOptions(schema.NewSet(k3sOptionsHash, []interface{}{
		map[string]interface{}{
			"extra_args": []interface{}{
				map[string]interface{}{
					"key":          "--token",
					"value":        "12345",
					"node_filters": []interface{}{"agent:0"},
				},
			},
			"node_labels": []interface{}{
				map[string]interface{}{
					"key":          "role",
					"value":        "worker",
					"node_filters": []interface{}{"agent:*"},
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
