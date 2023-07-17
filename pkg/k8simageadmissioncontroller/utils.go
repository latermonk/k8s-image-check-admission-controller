package k8simageadmissioncontroller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/manifest"
	"github.com/regclient/regclient/types/platform"
	"github.com/regclient/regclient/types/ref"
	logging "github.com/sirupsen/logrus"
)

type Layer struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

func GetManifest(imageName string) (map[string]interface{}, error) {
	var plat platform.Platform
	var manifestMap map[string]interface{}

	ctx := context.Background()
	rc := regclient.New()
	ref, err := ref.New(imageName)
	if err != nil {
		logging.Errorf("Error %v :", err)
		return nil, err
	}

	m, err := rc.ManifestGet(ctx, ref)
	if err != nil {
		logging.Errorf("Cannot list all manifests %v :", err)
		return nil, err
	}

	plat, err = platform.Parse("linux/amd64")
	if err != nil {
		logging.Errorf("Cannot parse platform %v :", err)
		return nil, err
	}

	desc, err := manifest.GetPlatformDesc(m, &plat)
	if err != nil {
		logging.Errorf("Cannot retrieve manifest platform %v :", err)
		return nil, err
	}

	manifest, err := rc.ManifestGet(ctx, ref, regclient.WithManifestDesc(*desc))
	if err != nil {
		logging.Errorf("Cannot retrieve manifest %v :", err)
		return nil, err
	}

	r, err := manifest.MarshalJSON()
	if err != nil {
		logging.Errorf("Cannot marshall manifest %v :", err)
		return nil, err

	}

	err = json.Unmarshal(r, &manifestMap)
	if err != nil {
		logging.Errorf("Cannot unmarshal manifest %v :", err)
		return nil, err
	}

	return manifestMap, nil

}

func GetImageSize2(imageName string) (int64, error) {
	sizeTotal := 0
	manifestMap, err := GetManifest(imageName)
	if err != nil {
		logging.Errorf("Cannot retrieve manifest platform %v :", err)
	}

	if layers, ok := manifestMap["layers"].([]interface{}); ok {
		jsonString, _ := json.Marshal(layers)
		s := []Layer{}
		err := json.Unmarshal(jsonString, &s)
		if err != nil {
			logging.Errorf("Cannot unmarshal layers %v :", err)
		}
		fmt.Println(s)
		for _, layer := range s {
			sizeTotal += layer.Size
		}

	}
	return int64(sizeTotal), nil
}

func GetImageSize(imageName string) (int64, error) {

	sizeTotal := 0

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		panic(err)
	}

	// get image information
	imageName = "ibackchina2018/ubuntu-sshd:huge"
	args := filters.NewArgs(filters.Arg("reference", imageName))
	images, err := cli.ImageList(ctx, types.ImageListOptions{Filters: args})
	if err != nil {
		panic(err)
	}

	if len(images) != 1 {
		panic(fmt.Errorf("Can't find image %s", imageName))
	}
	//sizeTotal = float32(images[0].Size) / (1024 * 1024)

	sizeTotal = len(images)
	println(sizeTotal)

	return int64(sizeTotal), nil
}
