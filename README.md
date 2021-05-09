# Terraform Provider For K3d


[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/terraform-provider-rancherk3d)](https://goreportcard.com/report/github.com/nikhilsbhat/terraform-provider-rancherk3d)  [![shields](https://img.shields.io/badge/license-mit-brightgreen)](https://github.com/nikhilsbhat/terraform-provider-rancherk3d/blob/master/LICENSE) [![shields](https://godoc.org/github.com/nikhilsbhat/terraform-provider-rancherk3d?status.svg)](https://godoc.org/github.com/nikhilsbhat/terraform-provider-rancherk3d)

[terraform](https://www.terraform.io/) provider for [k3d](https://k3d.io/), which helps in performing all operation that k3d does.

## Requirements

* Terraform v0.13.x [`tested`]
* Go 1.16
* Docker

## Features supported by the provider at the moment.

| component    | list/fetch | start/stop   | create   | delete    | load    |
| :----------: | :--------: | :----------: | :------: | :-------: | ------: |
|  `cluster`   | yes        | yes          | no       |  no       |  no     |
|  `node`      | yes        | yes          | yes      |  yes      |  no     | 
|  `registry`  | yes        | yes          | yes      |  yes      |  no     |
| `kubeconfig` | yes        | no           | no       |  yes      |  no     |
|    `image`   | no         | no           | no       |  no       |  yes    |

## Documentation

* Examples on the provider can be found in [examples](https://github.com/nikhilsbhat/terraform-provider-rancherk3d/tree/master/examples). <br><br> 
* Document that can help on how the [data_source](https://www.terraform.io/docs/language/data-sources/index.html) and [resource](https://www.terraform.io/docs/language/resources/syntax.html) could be used is [here](https://github.com/nikhilsbhat/terraform-provider-rancherk3d/tree/master/doc/README.md).

## TODO

* [ ] Support for cluster creation.
* [ ] Support for configuring registry with config-file.
* [ ] Terraform module.