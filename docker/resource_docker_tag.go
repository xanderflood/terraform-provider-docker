package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDockerTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceDockerTagCreate,
		Read:   resourceDockerTagRead,
		Update: resourceDockerTagUpdate,
		Delete: resourceDockerTagDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The tagged repository name to pull images from. Must not contain a digest segment.",
			},
			"pull_triggers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"latest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The sha256 digest of the latest image in this tag.",
			},
			"full_image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"all": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of all images from this tag currently being maintained locally. This attribute is intended to internal resource management and should generally not be referenced externally. Use `latest` instead for all containers and services.",
			},
		},
	}
}

func resourceDockerTagCreate(d *schema.ResourceData, meta interface{}) error {
	authConfig := meta.(*ProviderConfig).AuthConfigs
	client := meta.(*ProviderConfig).DockerClient

	name := d.Get("name").(string)

	digest, err := getLatestImageDigestByName(name, authConfig)
	if err != nil {
		return err
	}

	err = pullImageByDigest(name, digest, client, authConfig)
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))
	d.Set("latest", digest)
	d.Set("full_image_name", name+"@"+digest)
	d.Set("all", []string{digest})
	return nil
}

func resourceDockerTagRead(d *schema.ResourceData, meta interface{}) error {
	authConfig := meta.(*ProviderConfig).AuthConfigs
	name := d.Get("name").(string)

	digest, err := getLatestImageDigestByName(name, authConfig)
	d.Set("latest", digest)
	d.Set("full_image_name", name+"@"+digest)
	d.Set("all", d.Get("all")) //TODO is this unnecessary?
	return err
}

func resourceDockerTagUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).DockerClient
	authConfig := meta.(*ProviderConfig).AuthConfigs

	name := d.Get("name").(string)
	digest := d.Get("latest").(string)

	err := pullImageByDigest(name, digest, client, authConfig)
	if err != nil {
		return err
	}

	managedDigests := interfaceSliceToStringSlice(d.Get("all").([]interface{}))
	err = removeImagesByDigest(name, managedDigests[1:], client)
	if err != nil {
		return err
	}

	d.Set("full_image_name", name+"@"+digest)
	d.Set("all", []string{"latest", managedDigests[0]})
	return resourceDockerTagRead(d, meta)
}

func resourceDockerTagDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).DockerClient

	name := d.Get("name").(string)

	managedDigests := interfaceSliceToStringSlice(d.Get("all").([]interface{}))
	err := removeImagesByDigest(name, managedDigests, client)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func pullImageByDigest(name string, digest string, client *client.Client, authConfig *AuthConfigs) error {
	name = name + "@" + digest
	err := pullImage(client, authConfig, name)
	if err != nil {
		return fmt.Errorf("Unable to pull image %s: %w", name, err)
	}
	return nil
}

func removeImagesByDigest(name string, managedDigests []string, client *client.Client) error {
	var data Data
	if err := fetchLocalImages(&data, client); err != nil {
		return err
	}

	for _, digest := range managedDigests {
		taggedName := name + "@" + digest
		err := forceRemoveImage(taggedName, data, client)
		if err != nil {
			return err
		}
	}

	return nil
}

func interfaceSliceToStringSlice(ifaces []interface{}) []string {
	strings := make([]string, len(ifaces))
	for i := range ifaces {
		strings[i] = ifaces[i].(string)
	}
	return strings
}
