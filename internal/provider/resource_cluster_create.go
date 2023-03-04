package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	k3dCmdUtil "github.com/k3d-io/k3d/v5/cmd/util"
	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/client"
	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/k3d/cluster"
	"github.com/nikhilsbhat/terraform-provider-k3d/pkg/utils"
	k3dClient "github.com/rancher/k3d/v5/pkg/client"
	"github.com/rancher/k3d/v5/pkg/config/types"
	"github.com/rancher/k3d/v5/pkg/config/v1alpha4"
	"github.com/rancher/k3d/v5/pkg/runtimes"
	types2 "github.com/rancher/k3d/v5/pkg/types"
	"sigs.k8s.io/yaml"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		DeleteContext: resourceClusterDelete,
		// UpdateContext: resourceClusterUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of the Cluster to be created",
			},
			"servers_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Count of servers",
			},
			"agents_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Count of agents in the cluster",
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("K3D_IMAGE", nil),
				ForceNew:    true,
				Description: "Image name to be used for creation of cluster, it would be used along with kubernetes_version",
			},
			"network": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Network to be associated with the cluster",
			},
			"subnetwork": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Define a subnet for the newly created container network",
				Computed:    false,
			},
			"cluster_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "superSecretToken to be used",
			},
			"volumes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Mount volumes into the nodes (Format: [SOURCE:]DEST[@NODEFILTER[;NODEFILTER...]]",
				Elem: &schema.Resource{
					Schema: resourceClusterVolumeSchema(),
				},
			},
			"ports": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Description: "Map ports from the node containers (via the serverlb) to the host " +
					"(Format: [HOST:][HOSTPORT:]CONTAINERPORT[/PROTOCOL][@NODEFILTER])",
				Elem: &schema.Resource{
					Schema: resourceClusterPortsConfig(),
				},
			},
			"env": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Environment variables to be added nodes.",
				Elem: &schema.Resource{
					Schema: resourceClusterEnvsAndLabelsSchema(),
				},
			},
			"registries": {
				Type:        schema.TypeSet,
				Description: "Define how registries should be created or used",
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: resourceClusterRegistriesSchema(),
				},
			},
			"host_aliases": {
				Type:        schema.TypeSet,
				ForceNew:    true,
				Optional:    true,
				Description: "/etc/hosts style entries to be injected into /etc/hosts in the node containers and in the NodeHosts section in CoreDNS.",
				Elem: &schema.Resource{
					Schema: resourceHostAliasesConfig(),
				},
			},
			"kube_api": {
				Description: "same as `--api-port myhost.my.domain:6445` (where the name would resolve to 127.0.0.1)",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Description: "Important for the `server` setting in the kubeconfig.",
							ForceNew:    true,
							Optional:    true,
							Type:        schema.TypeString,
						},
						"host_ip": {
							Description:  "Where the Kubernetes API will be listening on.",
							ForceNew:     true,
							Optional:     true,
							Type:         schema.TypeString,
							ValidateFunc: validation.IsIPAddress,
						},
						"host_port": {
							Description:  "Specify the Kubernetes API server port exposed on the LoadBalancer.",
							ForceNew:     true,
							Optional:     true,
							Type:         schema.TypeInt,
							ValidateFunc: validation.IsPortNumber,
						},
					},
				},
			},
			"k3d_options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "k3d runtime settings",
				Elem: &schema.Resource{
					Schema: resourceClusterK3dOptionsSchema(),
				},
			},
			"k3s_options": {
				Type:        schema.TypeSet,
				ForceNew:    true,
				Optional:    true,
				Computed:    false,
				Description: "Options passed on to K3s itself",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"extra_args": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    false,
							Description: "additional arguments passed to the `k3s server|agent` command; same as `--k3s-arg`",
							Elem: &schema.Resource{
								Schema: resourceClusterEnvsAndLabelsSchema(),
							},
						},
						"node_labels": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    false,
							Description: "same as `--k3s-node-label 'foo=bar@agent:1'` -> this results in a Kubernetes node label",
							Elem: &schema.Resource{
								Schema: resourceClusterEnvsAndLabelsSchema(),
							},
						},
					},
				},
			},
			"kube_config": {
				Type:        schema.TypeSet,
				ForceNew:    true,
				Optional:    true,
				Description: "Way to manage the kubeconfig generated after creating k3d clusters.",
				Computed:    false,
				Elem: &schema.Resource{
					Schema: resourceKubeconfigConfig(),
				},
			},
			"runtime": {
				Description: "Runtime options for k3d",
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: resourceClusterRuntimeSchema(),
				},
			},
			"config_yaml": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "",
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	if !d.IsNewResource() {
		return nil
	}
	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceName))
		id = resourceID
	}

	clusterName := utils.String(d.Get(utils.TerraformResourceName))
	k3dImage := utils.String(d.Get(utils.TerraformResourceImage))
	if len(k3dImage) == 0 {
		k3dImage = defaultConfig.GetK3dImage()
	}

	cfg := &v1alpha4.SimpleConfig{
		ObjectMeta: types.ObjectMeta{
			Name: clusterName,
		},
		Servers:      utils.Int(d.Get(utils.TerraformResourceServersCount)),
		Agents:       utils.Int(d.Get(utils.TerraformResourceAgentsCount)),
		Image:        k3dImage,
		Network:      utils.String(d.Get(utils.TerraformResourceNetwork)),
		Subnet:       utils.String(d.Get(utils.TerraformResourceSubnet)),
		ClusterToken: utils.String(d.Get(utils.TerraformResourceClusterToken)),
		ExposeAPI:    flattenKubeAPI(d.Get(utils.TerraformKubeAPI)),
		Volumes:      flattenVolumes(d.Get(utils.TerraformResourceVolumes)),
		Ports:        flattenPorts(d.Get(utils.TerraformResourcePorts)),
		Env:          flattenEnvVars(d.Get(utils.TerraformResourceEnv)),
		Registries:   flattenRegistries(d.Get(utils.TerraformResourceRegistries)),
		HostAliases:  flattenHostAlias(d.Get(utils.TerraformHostAlias)),
	}

	k3dOptions, err := flattenK3DOptions(d.Get(utils.TerraformResourceK3dOptions))
	if err != nil {
		return diag.Errorf("fetching %s errored with: %v", utils.TerraformResourceK3dOptions, err)
	}

	cfg.Options = v1alpha4.SimpleConfigOptions{
		K3dOptions:        k3dOptions,
		K3sOptions:        flattenK3SOptions(d.Get(utils.TerraformResourceK3sOptions)),
		KubeconfigOptions: flattenKubeConfig(d.Get(utils.TerraformResourceKubeConfig)),
		Runtime:           flattenRuntime(d.Get(utils.TerraformK3dRuntime)),
	}

	if err = cluster.CreateCluster(ctx, defaultConfig.K3DRuntime, cfg); err != nil {
		if delErr := cluster.CheckAndDeleteCluster(ctx, defaultConfig.K3DRuntime, clusterName); delErr != nil {
			return diag.Errorf("creation of cluster '%s' FAILED with: %v\n, also FAILED to rollback changes!: %v", clusterName, err, delErr)
		}

		return diag.Errorf("creating cluster '%s' errored with: %v", clusterName, err)
	}

	d.SetId(id)

	return resourceClusterRead(ctx, d, meta)
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterName := utils.String(d.Get(utils.TerraformResourceName))

	k3dCluster, err := k3dClient.ClusterGet(ctx, runtimes.SelectedRuntime, &types2.Cluster{Name: clusterName})
	if err != nil {
		return diag.Errorf("fetching cluster '%s' errored with %v", clusterName, err)
	}

	if err = d.Set(utils.TerraformResourceNetwork, k3dCluster.Network.Name); err != nil {
		return diag.Errorf("setting %s errored with %v", utils.TerraformResourceNetwork, err)
	}

	if err = d.Set(utils.TerraformResourceClusterToken, k3dCluster.Token); err != nil {
		return diag.Errorf("setting %s errored with %v", utils.TerraformResourceClusterToken, err)
	}

	yamlOUT, err := yaml.Marshal(k3dCluster)
	if err != nil {
		return diag.Errorf("marshalling to yaml errored with: %v", err)
	}

	if err = d.Set(utils.TerrFormConfigYAML, string(yamlOUT)); err != nil {
		return diag.Errorf("setting %s errored with %v", utils.TerrFormConfigYAML, err)
	}

	return nil
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*client.Config)

	id := d.Id()
	if len(id) == 0 {
		return diag.Errorf("resource with the specified ID not found")
	}

	clusterName := utils.String(d.Get(utils.TerraformResourceName))
	if err := cluster.CheckAndDeleteCluster(ctx, defaultConfig.K3DRuntime, clusterName); err != nil {
		return diag.Errorf("deleting cluster '%s' errored with %v", clusterName, err)
	}

	return nil
}

func flattenPorts(ports interface{}) []v1alpha4.PortWithNodeFilters {
	k3dPorts := make([]v1alpha4.PortWithNodeFilters, 0)

	for _, port := range ports.(*schema.Set).List() {
		p := port.(map[string]interface{})
		k3dPorts = append(k3dPorts, v1alpha4.PortWithNodeFilters{
			Port:        getPortMappings(p),
			NodeFilters: utils.GetSlice(p["node_filters"].([]interface{})),
		})
	}

	return k3dPorts
}

func getPortMappings(p map[string]interface{}) string {
	return fmt.Sprintf("%s:%d:%d/%s", p["host"].(string), p["host_port"].(int), p["container_port"].(int), p["protocol"].(string))
}

func flattenVolumes(volumes interface{}) []v1alpha4.VolumeWithNodeFilters {
	k3dVolumes := make([]v1alpha4.VolumeWithNodeFilters, 0)

	for _, volume := range volumes.(*schema.Set).List() {
		v := volume.(map[string]interface{})
		k3dVolumes = append(k3dVolumes, v1alpha4.VolumeWithNodeFilters{
			Volume:      fmt.Sprintf("%s/%s", v["source"].(string), v["destination"].(string)),
			NodeFilters: utils.GetSlice(v["node_filters"].([]interface{})),
		})
	}

	return k3dVolumes
}

func flattenHostAlias(alias interface{}) []types2.HostAlias {
	k3dAlias := make([]types2.HostAlias, 0)

	for _, port := range alias.(*schema.Set).List() {
		a := port.(map[string]interface{})
		k3dAlias = append(k3dAlias, types2.HostAlias{
			IP:        a["ip"].(string),
			Hostnames: utils.GetSlice(a["hostnames"].([]interface{})),
		})
	}

	return k3dAlias
}

func flattenEnvVars(envs interface{}) []v1alpha4.EnvVarWithNodeFilters {
	k3dEnvs := make([]v1alpha4.EnvVarWithNodeFilters, 0)

	for _, port := range envs.(*schema.Set).List() {
		e := port.(map[string]interface{})
		k3dEnvs = append(k3dEnvs, v1alpha4.EnvVarWithNodeFilters{
			EnvVar:      fmt.Sprintf("%s=%s", e["key"].(string), e["value"].(string)),
			NodeFilters: utils.GetSlice(e["node_filters"].([]interface{})),
		})
	}

	return k3dEnvs
}

func flattenKubeAPI(api interface{}) v1alpha4.SimpleExposureOpts {
	var exposureOpts v1alpha4.SimpleExposureOpts

	if api.(*schema.Set).Len() == 0 {
		return exposureOpts
	}

	apiList := api.(*schema.Set).List()
	a := apiList[0].(map[string]interface{})

	hostPort, _ := k3dCmdUtil.GetFreePort()

	if a["host_port"].(int) == 0 {
		exposureOpts.HostPort = fmt.Sprintf("%d", hostPort)
	} else {
		exposureOpts.HostPort = fmt.Sprintf("%d", a["host_port"].(int))
	}

	exposureOpts.Host = a["host"].(string)
	exposureOpts.HostIP = a["host_ip"].(string)

	return exposureOpts
}

func flattenK3DOptions(k3d interface{}) (v1alpha4.SimpleConfigOptionsK3d, error) {
	k3dList := k3d.(*schema.Set).List()
	k := k3dList[0].(map[string]interface{})

	k3DOptions := v1alpha4.SimpleConfigOptionsK3d{
		Wait:                utils.Bool(k["wait"]),
		DisableLoadbalancer: utils.Bool(k["no_loadbalancer"]),
		DisableImageVolume:  utils.Bool(k["no_image_volume"]),
		NoRollback:          utils.Bool(k["no_rollback"]),
		Loadbalancer: v1alpha4.SimpleConfigOptionsK3dLoadbalancer{
			ConfigOverrides: utils.GetSlice(k["loadbalancer_config_overrides"].([]interface{})),
		},
	}

	if len(k["timeout"].(string)) != 0 {
		timeout, err := time.ParseDuration(k["timeout"].(string))
		if err != nil {
			return k3DOptions, err
		}
		k3DOptions.Timeout = timeout
	}

	return k3DOptions, nil
}

func flattenK3SOptions(k3s interface{}) v1alpha4.SimpleConfigOptionsK3s {
	var k3sOptions v1alpha4.SimpleConfigOptionsK3s

	if k3s.(*schema.Set).Len() == 0 {
		return k3sOptions
	}

	k3sList := k3s.(*schema.Set).List()
	k := k3sList[0].(map[string]interface{})
	k3sOptions.ExtraArgs = flattenExtraArgs(k["extra_args"].([]interface{}))
	k3sOptions.NodeLabels = flattenNodeLabels(k["node_labels"].([]interface{}))

	return k3sOptions
}

func flattenExtraArgs(extraArgs interface{}) []v1alpha4.K3sArgWithNodeFilters {
	k3sExtraArgs := make([]v1alpha4.K3sArgWithNodeFilters, 0)

	for _, port := range extraArgs.(*schema.Set).List() {
		e := port.(map[string]interface{})
		k3sExtraArgs = append(k3sExtraArgs, v1alpha4.K3sArgWithNodeFilters{
			Arg:         fmt.Sprintf("--%s=%s", e["key"].(string), e["value"].(string)),
			NodeFilters: utils.GetSlice(e["node_filters"].([]interface{})),
		})
	}

	return k3sExtraArgs
}

func flattenNodeLabels(nodeLabels interface{}) []v1alpha4.LabelWithNodeFilters {
	k3sNodeLabels := make([]v1alpha4.LabelWithNodeFilters, 0)

	for _, port := range nodeLabels.(*schema.Set).List() {
		e := port.(map[string]interface{})
		k3sNodeLabels = append(k3sNodeLabels, v1alpha4.LabelWithNodeFilters{
			Label:       fmt.Sprintf("%s=%s", e["key"].(string), e["value"].(string)),
			NodeFilters: utils.GetSlice(e["node_filters"].([]interface{})),
		})
	}

	return k3sNodeLabels
}

func flattenKubeConfig(cfg interface{}) v1alpha4.SimpleConfigOptionsKubeconfig {
	var kubeConfig v1alpha4.SimpleConfigOptionsKubeconfig

	if cfg.(*schema.Set).Len() == 0 {
		return kubeConfig
	}

	cfgList := cfg.(*schema.Set).List()
	c := cfgList[0].(map[string]interface{})

	kubeConfig.SwitchCurrentContext = c["switch_context"].(bool)
	kubeConfig.UpdateDefaultKubeconfig = c["update_default"].(bool)

	return kubeConfig
}

func flattenRuntime(run interface{}) v1alpha4.SimpleConfigOptionsRuntime {
	var runtime v1alpha4.SimpleConfigOptionsRuntime

	if run.(*schema.Set).Len() == 0 {
		return runtime
	}

	runList := run.(*schema.Set).List()
	k := runList[0].(map[string]interface{})

	runtime.GPURequest = k["gpu_request"].(string)
	runtime.ServersMemory = k["servers_memory"].(string)
	runtime.AgentsMemory = k["agents_memory"].(string)
	runtime.HostPidMode = k["host_pid_mode"].(bool)
	runtime.Labels = flattenNodeLabels(k["labels"].([]interface{}))

	return runtime
}

func flattenRegistries(reg interface{}) v1alpha4.SimpleConfigRegistries {
	regs := reg.(*schema.Set).List()
	if len(regs) == 0 || regs[0] == nil {
		return v1alpha4.SimpleConfigRegistries{}
	}

	r := regs[0].(map[string]interface{})

	create := r["create"].(bool)
	if !create {
		return v1alpha4.SimpleConfigRegistries{
			Use: utils.GetSlice(r["use"].([]interface{})),
		}
	}

	return v1alpha4.SimpleConfigRegistries{}
}
