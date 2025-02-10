package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/config"
)

type MapService struct {
	api           *API
	primaryMap    assets.Map[[]assets.Image]
	currentMap    assets.Map[[]assets.Image]
	currentImages []assets.Image
}

func NewMapService(api *API) *MapService {
	ms := &MapService{
		api:           api,
		primaryMap:    assets.Map[[]assets.Image]{},
		currentMap:    assets.Map[[]assets.Image]{},
		currentImages: []assets.Image{},
	}
	_, err := ms.GetPrimaryMap()
	if err != nil {
		panic(err)
	}
	return ms
}

func (ms *MapService) PrimaryMap() assets.Map[[]assets.Image] {
	return ms.primaryMap
}

// GetPrimaryMap makes a get request to server for primary map
func (ms *MapService) GetPrimaryMap() (assets.Map[[]assets.Image], error) {
	_map := assets.Map[[]assets.Image]{}
	userID := ms.api.UserID()

	path := config.Env().SERVER_URL + "/game/wasm/map/primary/" + userID
	res := ms.api.Request(http.MethodGet, path)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return _map, res.Error
	}

	err := json.Unmarshal(res.Body, &_map)
	if err != nil {
		return _map, err
	}

	ms.primaryMap = _map
	return _map, nil
}

// GetMapByID makes a get request to server for map by id
func (ms *MapService) GetMapByID(id string) (assets.Map[[]assets.Image], error) {
	_map := assets.Map[[]assets.Image]{}

	path := config.Env().SERVER_URL + "/game/wasm/map?id=" + id + "&userID=" + ms.api.UserID()
	res := ms.api.Request(http.MethodGet, path)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return _map, res.Error
	}
	err := json.Unmarshal(res.Body, &_map)
	if err != nil {
		return _map, err
	}
	return _map, nil
}

func (ms *MapService) CurrentMap() assets.Map[[]assets.Image] {
	return ms.currentMap
}

func (ms *MapService) SetCurrentMap(_map assets.Map[[]assets.Image]) {
	ms.currentMap = _map
	ms.currentImages = ms.ImagesFromMap(_map)
}

func (ms *MapService) CurrentImages() []assets.Image {
	return ms.currentImages
}

// ImagesFromMap creates ebiten images from a Map, sorting by tiles, other, then portals
func (ms *MapService) ImagesFromMap(_map assets.Map[[]assets.Image]) []assets.Image {
	imagesUnsorted := []assets.Image{}

	for key, image := range _map.Data {
		img, err := assets.ImageFromPixelData(image)
		if err != nil {
			panic(err)
		}
		_map.Data[key].Image = img
	}
	imagesUnsorted = append(imagesUnsorted, _map.Data...)

	// tiles first then non-tiles
	var allImages []assets.Image
	var nonTiles []assets.Image
	for _, img := range imagesUnsorted {
		if img.AssetType == assets.Tile {
			allImages = append(allImages, img)
		} else {
			nonTiles = append(nonTiles, img)
		}
	}
	allImages = append(allImages, nonTiles...)

	// add portals as images
	for _, portal := range _map.Portals {
		// blank ebiten image for collision detection with portal
		am := assets.Image{
			ID:        portal.MapID,
			UserID:    _map.UserID,
			Name:      "portal-" + portal.MapID,
			AssetType: assets.MapPortal,
			Width:     16,
			Height:    16,
			X:         portal.X,
			Y:         portal.Y,
		}
		_img := ebiten.NewImage(am.Width, am.Height)
		am.Image = _img
		allImages = append(allImages, am)
	}

	return allImages
}
