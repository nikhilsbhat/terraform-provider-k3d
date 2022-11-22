package rancherk3d_test

// import (
//	"os"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// )
//
// var testAccProviderFactories map[string]func() (*schema.Provider, error)
// var testAccProviders map[string]*schema.Provider
// var testAccProvider *schema.Provider
//
// func init() {
//	testAccProvider = Provider()
//	testAccProviders = map[string]*schema.Provider{
//		"rancherk3d": testAccProvider,
//	}
//	testAccProviderFactories = testRancherProvider()
// }
//
// func TestProvider(t *testing.T) {
//	if err := Provider().InternalValidate(); err != nil {
//		t.Fatalf("err: %s", err)
//	}
// }
//
// func TestProvider_impl(t *testing.T) {
//	var _ *schema.Provider = Provider()
// }
//
// func testAccPreCheck(t *testing.T) {
//	if err := os.Getenv("K3D_RUNTIME"); err == "" {
//		t.Fatal("K3D_RUNTIME must be set for acceptance tests")
//	}
// }
//
// func testRancherProvider() map[string]func() (*schema.Provider, error) {
//	rancherProvider := map[string]func() (*schema.Provider, error){
//		"rancherk3d": func() (*schema.Provider, error) {
//			return Provider(), nil
//		},
//	}
//	return rancherProvider
// }
