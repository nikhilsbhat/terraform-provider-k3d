---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "k3d_load_image Resource - terraform-provider-k3d"
subcategory: ""
description: |-
  
---

# k3d_load_image (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster` (String) name of the existing cluster to which the images has to be imported to
- `images` (List of String) list of images to be imported to the existing cluster

### Optional

- `all` (Boolean) if enabled loads images to all available clusters
- `keep_tarball` (Boolean) enable to keep the tarball of the loaded images locally
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `images_stored` (List of Object) list of images loaded to the cluster (see [below for nested schema](#nestedatt--images_stored))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)


<a id="nestedatt--images_stored"></a>
### Nested Schema for `images_stored`

Read-Only:

- `cluster` (String)
- `images` (List of String)


