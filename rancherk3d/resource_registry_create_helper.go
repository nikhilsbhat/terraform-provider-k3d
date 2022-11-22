package rancherk3d

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utils2 "github.com/nikhilsbhat/terraform-provider-rancherk3d/pkg/utils"
)

func validateProxy(value map[string]string) bool {
	if len(value["remoteURL"]) == 0 || len(value["username"]) == 0 || len(value["password"]) == 0 {
		return false
	}
	return true
}

func validateAndSetProxy(d *schema.ResourceData, proxy map[string]string) map[string]string {
	if utils2.Bool(d.Get(utils2.TerraformUseProxy)) {
		if !validateProxy(proxy) {
			return map[string]string{}
		}
		fmt.Printf("proxy config validation failed, config cannot be empty, dropping proxy config")
	}
	return nil
}

func validateExpose(value map[string]string) bool {
	if len(value["hostIp"]) == 0 || len(value["hostPort"]) == 0 {
		return false
	}
	return true
}

func validateAndSetExpose(expose map[string]string) map[string]string {
	if !validateExpose(expose) {
		return map[string]string{
			"hostIp":   "0.0.0.0",
			"hostPort": "5200",
		}
	}
	return expose
}

func validateAndSetHost(d *schema.ResourceData) string {
	if len(utils2.String(d.Get(utils2.TerraformResourceHost))) == 0 {
		return utils2.String(d.Get(utils2.TerraformResourceName))
	}
	return utils2.String(d.Get(utils2.TerraformResourceHost))
}

func getMetadata(d *schema.ResourceData) map[string]string {
	metadata := make(map[string]string)
	if host := utils2.String(d.Get(utils2.TerraformResourceHost)); len(host) == 0 {
		metadata["host"] = utils2.String(d.Get(utils2.TerraformResourceName))
		return metadata
	}
	metadata["host"] = utils2.String(d.Get(utils2.TerraformResourceHost))
	return metadata
}
