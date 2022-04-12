package rancherk3d

//
//import (
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//	"strings"
//	"testing"
//)
//
//var testProviders map[string]func() (*schema.Provider, error)
//
//func TestResourceAudit(t *testing.T) {
//	t.Parallel()
//
//	node := acctest.RandomWithPrefix("tf-acc")
//
//	resource.Test(t, resource.TestCase{
//		ProviderFactories: testProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testDataSourceNodeList(node),
//				Check:  testResourceAudit_initialCheck(path),
//			},
//		},
//	})
//}
//
//func testDataSourceNodeList(node string) string {
//	return fmt.Sprintf(`
//data "rancherk3d_node_list" "k3s_default" {
//  cluster = "k3s-default"
//  nodes   = ["%s"]
//}
//`, node)
//}
//
//func testAccCheckKeycloakRoleDestroy() resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		for name, rs := range s.RootModule().Resources {
//			if rs.Type != "keycloak_role" || strings.HasPrefix(name, "data") {
//				continue
//			}
//
//			id := rs.Primary.ID
//			realm := rs.Primary.Attributes["realm_id"]
//
//			role, _ := keycloakClient.GetRole(testCtx, realm, id)
//			if role != nil {
//				return fmt.Errorf("%s with id %s still exists", name, id)
//			}
//		}
//
//		return nil
//	}
//}
