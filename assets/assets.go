package assets

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/snburman/game/config"
)

type AssetType string

type Assets struct {
	Images []Image `json:"images"`
}

func Load() *Assets {
	//TODO: acquire ID from JS.globals
	id := "6794a98e48815ec0dd9c19d0"
	// Make get request
	client := http.Client{}
	req, err := http.NewRequest("GET", config.Env().SERVER_URL+"/game/wasm/map/"+id, nil)
	if err != nil {
		log.Println("error during new request")
		panic(err)
	}
	req.Header.Add("CLIENT_ID", config.Env().CLIENT_ID)
	req.Header.Add("CLIENT_SECRET", config.Env().CLIENT_SECRET)
	res, err := client.Do(req)

	if err != nil {
		log.Println("error from response")
		panic(err)
	}
	defer res.Body.Close()
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var _map Map[[]Image]
	err = json.Unmarshal(bts, &_map)
	if err != nil {
		panic(err)
	}

	var assets Assets = Assets{}

	for key, image := range _map.Data {
		img, err := imageFromPixelData(image)
		if err != nil {
			panic(err)
		}
		_map.Data[key].Image = img
	}
	assets.Images = append(assets.Images, _map.Data...)

	return &assets
}

type Pixel struct {
	X int `json:"x"`
	Y int `json:"y"`
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
	A int `json:"a"`
}

type PixelData = [][]Pixel

type PlayerAsset[T any] struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Data   T      `json:"data"`
}
