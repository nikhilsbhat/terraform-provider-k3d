package rancherk3d

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"terraform-provider-rancherk3d/k3d"
	"terraform-provider-rancherk3d/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	terraformResourceClusters      = "clusters"
	terraformResourceCluster       = "cluster"
	terraformResourceImages        = "images"
	terraformResourceKeepTarball   = "keep_tarball"
	terraformResourceTarballStored = "tarball_stored"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadImageLoad,
		ReadContext:   resourceLoadImageRead,
		DeleteContext: resourceLoadImageDelete,
		UpdateContext: resourceLoadImageUpdate,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"images": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of images to be imported to cluster",
			},
			"clusters": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of clusters to which the images has to be imported to",
			},
			"keep_tarball": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				Description: "enable to keep the tarball of the loaded images locally",
			},
			"images_stored": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "cluster to which the below images are stored",
						},
						"tarball_stored": {
							Type:        schema.TypeMap,
							Computed:    true,
							Optional:    true,
							Description: "details of images and its tarball stored, if in case keep_tarball is enabled",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["cluster"].(string)))
					return utils.GetHash(buf.String())
				},
			},
		},
	}
}

func resourceLoadImageLoad(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)

	if d.IsNewResource() {
		id := d.Id()
		fmt.Print(defaultConfig)

		id, err := utils.GetRandomID()
		if err != nil {
			return diag.Errorf("errored while fetching randomID %v", err)
		}

		clusters := utils.GetSlice((d.Get(terraformResourceClusters)).([]interface{}))
		for _, cluster := range clusters {
			client := getImagesClient(defaultConfig, d, cluster)
			imagesStored, err := client.StoreImages()
			if err != nil {
				return diag.Errorf("oops an error occurred while storing images to cluster %s : %v", cluster, err)
			}
			log.Printf("image that would be stored are %v", imagesStored)
		}

		d.SetId(id)

		return resourceLoadImageUpdate(ctx, d, meta)
	}

	log.Printf("resource %s already exists", d.Id())
	return nil
}

func resourceLoadImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

func resourceLoadImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

func resourceLoadImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)

	fmt.Println(defaultConfig)
	log.Printf("updating loaded images data")
	if d.HasChange(terraformResourceClusters) || d.HasChange(terraformResourceImages) {
		images := d.Get(terraformResourceImages)
		clusters := d.Get(terraformResourceClusters)

		d.GetChange()
		log.Printf("images: %v", images)
		log.Printf("images: %v", clusters)
		imagesStored := getImagesStored(images, clusters.([]interface{}))

		log.Printf("%v", imagesStored)
		if err := d.Set("images_stored", nil); err != nil {
			return diag.Errorf("oops error occurred while setting 'images_stored' for update operation %v", err)
		}
		return nil
	}

	log.Printf("nothing to update skipping")
	return nil
}

func updateImagesStored(d *schema.ResourceData, ) diag.Diagnostics {
	images := d.Get(terraformResourceImages)
	clusters := d.Get(terraformResourceClusters)
	if err := d.Set("images_stored", imagesStored); err != nil {
		return diag.Errorf("oops error occurred while setting 'images_stored' for create operation %v", err)
	}
}

func getImagesStored(clusters, images []string) []map[string]interface{} {
	imagesStored := make([]map[string]interface{}, 0)
	for _, cluster := range clusters {
		imageStored := make(map[string]interface{})
		imageStored[terraformResourceCluster] = cluster
		imageStored[terraformResourceTarballStored] = getImageStored(images)
		imagesStored = append(imagesStored, imageStored)
	}
	return imagesStored
}

func getImageStored(images []string) map[string]string {
	imagesStored := make(map[string]string, len(images))
	for _, image := range images {
		imagesStored[image] = image
	}
	return imagesStored
}

func getImagesClient(dfconfig *k3d.K3dConfig, d *schema.ResourceData, cluster string) *k3d.K3Dimages {
	images := d.Get(terraformResourceImages)
	keepTarball := d.Get(terraformResourceKeepTarball)

	imagesClient := k3d.NewK3dImages()
	imagesClient.Images = utils.GetSlice(images.([]interface{}))
	imagesClient.Cluster = cluster
	imagesClient.K3DRuntime = dfconfig.K3DRuntime
	imagesClient.StoreTarBall = keepTarball.(bool)

	return imagesClient
}

// Below snippet is commented for a while until these are used
/*func getImagesStoredWithCluster(clusters, images []string) (storedImages map[string]*k3d.StoredImages) {
	for _, cluster := range clusters {
		storedImages[cluster] = getImagesStored(cluster, images)
	}
	return
}

func getImagesStored(cluster string, images []string) (imagesStored *k3d.StoredImages) {
	imagesStored = &k3d.StoredImages{
		Cluster: cluster,
		TarBall: getTarballData(images),
	}
	return
}

func getTarballData(images []string) (tarballData []*k3d.TarBallData) {
	for _, image := range images {
		tarballData = append(tarballData, &k3d.TarBallData{
			Image: image,
		})
	}
	return
}*/
