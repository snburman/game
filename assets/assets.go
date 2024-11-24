package assets

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type AssetType string

type Asset struct {
	Path string
	Data []byte
}

func NewAsset(t AssetType, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	switch t {
	case PNG:
		return pngBytesFromFile(file)
	default:
		return nil, errors.New("unsupported asset type")
	}
}

type Assets struct {
	Images Images `json:"images"`
}

func (a *Assets) Sprite(name string) Image {
	return a.Images.Sprites[name]
}

type Images struct {
	Sprites map[string]Image `json:"sprites"`
}

type FrameSpec struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

func Load() *Assets {
	// Make get request
	res, err := http.Get("http://localhost:9191/game/player/assets")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	bts, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var images []Image
	err = json.Unmarshal(bts, &images)
	if err != nil {
		panic(err)
	}

	var assets Assets = Assets{}

	for key, image := range images {
		img, err := imageFromPixelData(image)
		if err != nil {
			panic(err)
		}
		images[key].Image = img
	}

	assets.Images.Sprites = make(map[string]Image)
	for _, image := range images {
		assets.Images.Sprites[image.Name] = image
	}

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
