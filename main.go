// ----------------------------------------------------------------------------
//
//	***     TERRAGEN GENERATED CODE    ***    TERRAGEN GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//	This file was auto generated by Terragen.
//	This autogenerated code has to be enhanced further to make it fully working terraform-provider.
//
//	Get more information on how terragen works.
//	https://github.com/nikhilsbhat/terragen
//
// ----------------------------------------------------------------------------
//
//nolint:gocritic
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/nikhilsbhat/terraform-provider-k3d/rancherk3d"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: rancherk3d.Provider,
	})
}
