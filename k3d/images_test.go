package k3d_test

import (
	"terraform-provider-rancherk3d/k3d"
	"testing"

	"github.com/rancher/k3d/v4/pkg/runtimes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

//func TestK3Dimages_StoreImages(t *testing.T) {
//	images1 := k3d.K3Dimages{
//		Images:       []string{"basnik/terragen:v0.2.0"},
//		Cluster:      "k3s-default",
//		StoreTarBall: false,
//		Config:       k3d.K3dConfig{K3DRuntime: runtimes.Docker},
//	}
//
//	images2 := images1
//	images2.Cluster = "k3s-default-1"
//
//	tests := []struct {
//		name    string
//		images  k3d.K3Dimages
//		want    *k3d.StoredImages
//		wantErr bool
//	}{
//		{
//			name:   "should be able to load images to the specified cluster",
//			images: images1,
//			want: &k3d.StoredImages{
//				Cluster: images1.Cluster,
//				Images:  images1.Images,
//			},
//			wantErr: false,
//		},
//		{
//			name:    "should fail to load images as specified cluster is not present",
//			images:  images2,
//			want:    nil,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			client := &k3d.K3Dimages{
//				Images:       tt.images.Images,
//				Cluster:      tt.images.Cluster,
//				StoreTarBall: tt.images.StoreTarBall,
//				StoredImages: tt.images.StoredImages,
//				Config:       tt.images.Config,
//			}
//			got, err := client.StoreImages()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("StoreImages() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("StoreImages() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestStoreImages(t *testing.T) {
	t.Run("should be able to load images to cluster k3s-default", func(t *testing.T) {
		client := &k3d.K3Dimages{
			Images:       []string{"basnik/terragen:v0.2.0"},
			Cluster:      "k3s-default",
			StoreTarBall: false,
			Config:       k3d.K3dConfig{K3DRuntime: runtimes.Docker},
			Context:      context.TODO(),
		}
		got, err := client.StoreImages()
		assert.Nil(t, err)
		assert.NotNil(t, got)
	})

	t.Run("should be able to load images to cluster test-cluster", func(t *testing.T) {
		client := &k3d.K3Dimages{
			Images:       []string{"basnik/terragen:v0.2.0"},
			Cluster:      "test-cluster",
			StoreTarBall: false,
			Config:       k3d.K3dConfig{K3DRuntime: runtimes.Docker},
			Context:      context.Background(),
		}
		got, err := client.StoreImages()
		assert.Nil(t, err)
		assert.NotNil(t, got)
	})
}
