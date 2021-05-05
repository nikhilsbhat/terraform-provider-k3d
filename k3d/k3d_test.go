package k3d

import (
	"reflect"
	"testing"

	"github.com/rancher/k3d/v4/pkg/runtimes"
)

func Test_getRuntime(t *testing.T) {
	tests := []struct {
		name    string
		runtime string
		want    runtimes.Runtime
	}{
		{
			name:    "should return docker runtime",
			runtime: "docker",
			want:    runtimes.Docker,
		},
		{
			name:    "should return docker runtime",
			runtime: "containerd",
			want:    runtimes.SelectedRuntime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRuntime(tt.runtime); !reflect.DeepEqual(got, tt.want) { //nolint:scopelint
				t.Errorf("getRuntime() = %v, want %v", got, tt.want) //nolint:scopelint
			}
		})
	}
}
