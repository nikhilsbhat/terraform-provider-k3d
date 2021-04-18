package rancherk3d

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/k3d"
	"github.com/nikhilsbhat/terraform-provider-rancherk3d/utils"
)

var (
	terraformResourceImageClusters = "clusters"
	terraformResourceImageCluster  = "cluster"
	terraformResourceImages        = "images"
	terraformResourceImagesStored  = "images_stored"
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
				Required:    true,
				Computed:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "list of images to be imported to the existing cluster",
			},
			"clusters": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				Description: "name of the existing cluster to which the images has to be imported to",
			},
			"keep_tarball": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				Description: "enable to keep the tarball of the loaded images locally",
			},
			"images_stored": {
				Type:     schema.TypeSet,
				Computed: true,
				ForceNew: false,
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

		images := getImages(d.Get(terraformResourceImages))
		keepTarball := d.Get(terraformResourceKeepTarball).(bool)
		cluster := d.Get(terraformResourceImageClusters).(string)

		imagesLoaded, err := uploadImagesToClusters(defaultConfig, ctx, keepTarball, images, cluster)
		if err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}
		_ = imagesLoaded
		d.SetId(id)
		return resourceLoadImageRead(ctx, d, meta)
	}

	log.Printf("resource %s already exists", d.Id())
	return nil
}

func resourceLoadImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	images := d.Get(terraformResourceImages)
	cluster := d.Get(terraformResourceImageClusters).(string)

	imagesStored := getImagesToBeStored(cluster, utils.GetSlice(images.([]interface{})))

	log.Printf("imagesToBeStored %v", imagesStored)
	if err := d.Set(terraformResourceImagesStored, imagesStored); err != nil {
		return diag.Errorf("oops setting 'images_stored' errored with : %v", err)
	}
	if err := d.Set(terraformResourceImages, images); err != nil {
		return diag.Errorf("oops setting 'images' errored with : %v", err)
	}
	if err := d.Set(terraformResourceImageClusters, cluster); err != nil {
		return diag.Errorf("oops setting 'clusters' errored with : %v", err)
	}
	return nil
}

func resourceLoadImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)
	_ = defaultConfig
	// could be properly implemented once k3d supports deleting loaded images from cluster.

	orderID := d.Id()
	if len(orderID) == 0 {
		return diag.Errorf("resource with the specified ID not found")
	}
	d.SetId("")
	return nil
}

func resourceLoadImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(*k3d.K3dConfig)
	fmt.Println(defaultConfig)

	log.Printf("uploading newer images to k3d clusters")
	if d.HasChange(terraformResourceImageClusters) || d.HasChange(terraformResourceImages) {

		updatedClusters, updatedImages := getUpdatedClusterAndImages(d)
		keepTarball := d.Get(terraformResourceKeepTarball).(bool)

		imagesLoaded, err := uploadImagesToClusters(defaultConfig, ctx, keepTarball, updatedImages, updatedClusters)
		if err != nil {
			return diag.Errorf("creation failed with error: %v", err)
		}
		_ = imagesLoaded
		if err := d.Set(terraformResourceImageClusters, updatedClusters); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", terraformResourceImageClusters, err)
		}
		if err := d.Set(terraformResourceImages, updatedImages); err != nil {
			return diag.Errorf("oops setting '%s' errored with : %v", terraformResourceImages, err)
		}
		return resourceLoadImageRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")
	return nil
}

func getImagesToBeStored(cluster string, images []string) []map[string]interface{} {
	imagesStored := make([]map[string]interface{}, 0)
	imageStored := make(map[string]interface{})
	imageStored[terraformResourceImageCluster] = cluster
	imageStored[terraformResourceTarballStored] = getImageStored(images)
	imagesStored = append(imagesStored, imageStored)
	return imagesStored
}

func getImageStored(images []string) map[string]string {
	imagesStored := make(map[string]string, len(images))
	for _, image := range images {
		imagesStored[image] = image
	}
	return imagesStored
}

func uploadImagesToClusters(defaultConfig *k3d.K3dConfig, ctx context.Context, keepArtifact bool, images []string, cluster string) ([]*k3d.StoredImages, error) {
	storedImages := make([]*k3d.StoredImages, 0)
	log.Printf("uploading images %v to k3d cluster %s", images, cluster)
	client := getImagesClient(defaultConfig, ctx, images, cluster, keepArtifact)
	imagesToStore, err := client.StoreImages()
	if err != nil {
		return nil, fmt.Errorf("oops an error occurred while storing images to cluster %s : %v", cluster, err)
	}
	storedImages = append(storedImages, imagesToStore)
	log.Printf("images %v were successfully uploaded to k3d cluster %s", images, cluster)

	return storedImages, nil
}

func getImagesClient(dfconfig *k3d.K3dConfig, ctx context.Context, images []string, cluster string, keepArtifact bool) *k3d.K3Dimages {
	imagesClient := k3d.NewK3dImages()
	imagesClient.Context = context.Background()
	imagesClient.Images = images
	imagesClient.Cluster = cluster
	imagesClient.Config.K3DRuntime = dfconfig.K3DRuntime
	imagesClient.StoreTarBall = keepArtifact
	return imagesClient
}

func getUpdatedClustersAndImages(d *schema.ResourceData) (clusters, images []string) {
	oldClusters, newClusters := d.GetChange(terraformResourceImageClusters)
	clusters = getClusters(oldClusters)
	if !cmp.Equal(oldClusters, newClusters) {
		clusters = getClusters(newClusters)
	}
	oldImages, newImages := d.GetChange(terraformResourceImages)
	images = getImages(oldImages)
	if !cmp.Equal(oldImages, newImages) {
		images = getImages(newImages)
	}
	return
}

func getUpdatedClusterAndImages(d *schema.ResourceData) (cluster string, images []string) {
	oldCluster, newCluster := d.GetChange(terraformResourceImageClusters)
	if !cmp.Equal(oldCluster, newCluster) {
		cluster = newCluster.(string)
	}
	oldImages, newImages := d.GetChange(terraformResourceImages)
	images = getImages(oldImages)
	if !cmp.Equal(oldImages, newImages) {
		images = getImages(newImages)
	}
	return
}

func getImages(images interface{}) []string {
	return utils.GetSlice(images.([]interface{}))
}

func getClusters(clusters interface{}) []string {
	return utils.GetSlice(clusters.([]interface{}))
}
