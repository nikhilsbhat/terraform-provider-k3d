package rancherk3d

import (
	"context"
	"testing"

	"github.com/rancher/k3d/v4/pkg/runtimes"
	"github.com/stretchr/testify/assert"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func Test_getImagesStored(t *testing.T) {
	t.Run("should return the images to be stored in a required format", func(t *testing.T) {
		cluster := "cluster1"
		images := []string{"basnik/terragen:latest", "basnik/renderer:latest"}

		expected := []map[string]interface{}{
			{
				"cluster": "cluster1",
				"tarball_stored": map[string]string{
					"basnik/renderer:latest": "basnik/renderer:latest",
					"basnik/terragen:latest": "basnik/terragen:latest",
				},
			},
		}
		actual, err := getImagesToBeStored(context.Background(), runtimes.SelectedRuntime, images, cluster, false)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

// func Test_resourceImage(t *testing.T) {
//	ri := acctest.RandInt()
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t) },
//		ProviderFactories: testAccProviderFactories,
//		CheckDestroy:      testAccCheckImageDestroy(),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccMachineImage_basic(ri),
//				Check:  testAccCheckMachineImageExists,
//			},
//		},
//	})
// }
//
// func getMapofinterface() map[string]interface{} {
//	return map[string]interface{}{
//
//	}
// }
//
// func testAccCheckImageDestroy() resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		for _, rs := range s.RootModule().Resources {
//			if rs.Type != "keycloak_role" {
//				continue
//			}
//
//			id := rs.Primary.ID
//			realm := rs.Primary.Attributes["realm_id"]
//
//			role, _ := keycloakClient.GetRole(realm, id)
//			if role != nil {
//				return fmt.Errorf("role with id %s still exists", id)
//			}
//		}
//
//		return nil
//	}
// }
//
// func testAccMachineImage_basic(rInt int) string {
//
//	identity_domain := os.Getenv("OPC_IDENTITY_DOMAIN")
//
//	testAccMachineImageBasic := `
// resource "rancherk3d_load_image" "loading" {
//  images       = ["basnik/terragen:latest", "basnik/renderer:latest"]
//  clusters     = ["k3s-default",]
//  keep_tarball = true
// }`
//
//	return fmt.Sprintf(testAccMachineImageBasic, identity_domain, rInt)
// }
