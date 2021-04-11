package rancherk3d

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getImagesStored(t *testing.T) {
	t.Run("should return the images to be stored in a required format", func(t *testing.T) {
		clusters := []string{"cluster1", "cluster2"}
		images := []string{"basnik/terragen:latest", "basnik/renderer:latest"}

		//expected := []map[string]interface{}{
		//	{
		//		"cluster": "cluster1",
		//		"tarball_stored": {
		//			"basnik/renderer:latest": "basnik/renderer:latest",
		//			"basnik/terragen:latest": "basnik/terragen:latest",
		//		},
		//	},
		//	{
		//		"cluster": "cluster2",
		//		"tarball_stored":
		//		"basnik/renderer:latest": "basnik/renderer:latest",
		//		"basnik/terragen:latest": "basnik/terragen:latest",
		//	},
		//}
		actual := getImagesStored(clusters, images)
		assert.Equal(t, "", actual)

	})
}

func getMapofinterface() map[string]interface{} {
	return map[string]interface{}{

	}
}
