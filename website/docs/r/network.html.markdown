---
layout: "docker"
page_title: "Docker: docker_network"
sidebar_current: "docs-docker-resource-network"
description: |-
  Manages a Docker Network.
---

# docker\_network

Manages a Docker Network. This can be used alongside
[docker\_container](/docs/providers/docker/r/container.html)
to create virtual networks within the docker environment.

## Example Usage

```hcl
# Create a new docker network
resource "docker_network" "private_network" {
  name = "my_network"
}

# Access it somewhere else with ${docker_network.private_network.name}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the Docker network.
* `labels` - (Optional, map of string/string key/value pairs) User-defined key/value metadata.
* `check_duplicate` - (Optional, boolean) Requests daemon to check for networks
  with same name.
* `driver` - (Optional, string) Name of the network driver to use. Defaults to
  `bridge` driver.
* `options` - (Optional, map of strings) Network specific options to be used by
  the drivers.
* `internal` - (Optional, boolean) Restrict external access to the network.
  Defaults to `false`.
* `attachable` - (Optional, boolean) Enable manual container attachment to the network.
  Defaults to `false`.
* `ingress` - (Optional, boolean) Create swarm routing-mesh network.
  Defaults to `false`.
* `ipv6` - (Optional, boolean) Enable IPv6 networking.
  Defaults to `false`.
* `ipam_driver` - (Optional, string) Driver used by the custom IP scheme of the
  network.
* `ipam_config` - (Optional, block) See [IPAM config](#ipam_config-1) below for
  details.

<a id="ipam_config-1"></a>
### IPAM config
Configuration of the custom IP scheme of the network.

The `ipam_config` block supports:

* `subnet` - (Optional, string)
* `ip_range` - (Optional, string)
* `gateway` - (Optional, string)
* `aux_address` - (Optional, map of string)

## Attributes Reference

The following attributes are exported in addition to the above configuration:

* `id` (string)
* `scope` (string)
