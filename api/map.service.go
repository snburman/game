package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/config"
	"github.com/snburman/game/models"
	"github.com/snburman/game/objects"
)

type MapService struct {
	api           *API
	primaryMap    models.Map[[]models.Image]
	currentMap    models.Map[[]models.Image]
	currentImages []models.Image
}

func NewMapService(api *API, g objects.IGame) *MapService {
	ms := &MapService{
		api:           api,
		primaryMap:    models.Map[[]models.Image]{},
		currentMap:    models.Map[[]models.Image]{},
		currentImages: []models.Image{},
	}
	_, err := ms.GetPrimaryMap()
	if err != nil {
		panic(err)
	}
	ms.SetCurrentMap(g, ms.primaryMap)
	return ms
}

func (ms *MapService) PrimaryMap() models.Map[[]models.Image] {
	return ms.primaryMap
}

// GetPrimaryMap makes a get request to server for primary map
func (ms *MapService) GetPrimaryMap() (models.Map[[]models.Image], error) {
	_map := models.Map[[]models.Image]{}
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
func (ms *MapService) GetMapByID(id string) (models.Map[[]models.Image], error) {
	_map := models.Map[[]models.Image]{}

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

func (ms *MapService) CurrentMap() models.Map[[]models.Image] {
	return ms.currentMap
}

func (ms *MapService) SetCurrentMap(g objects.IGame, _map models.Map[[]models.Image]) {
	// set map
	ms.currentMap = _map

	// set images
	imgs := ms.ImagesFromMap(_map)
	ms.currentImages = imgs

	// set objects
	objs, player := objects.ObjectersFromImages(imgs)
	g.Objects().SetAll(objs)
	if player != nil {
		g.SetPlayer(player)
	}
}

func (ms *MapService) CurrentImages() []models.Image {
	return ms.currentImages
}

// ImagesFromMap creates ebiten images from a Map, sorting by tiles, other, then portals
func (ms *MapService) ImagesFromMap(_map models.Map[[]models.Image]) []models.Image {
	imagesUnsorted := []models.Image{}

	for key, image := range _map.Data {
		_map.Data[key].Image = models.ImageFromPixelData(image)
	}
	imagesUnsorted = append(imagesUnsorted, _map.Data...)

	// tiles first then non-tiles
	var allImages []models.Image
	var nonTiles []models.Image
	for _, img := range imagesUnsorted {
		if img.AssetType == models.Tile {
			allImages = append(allImages, img)
		} else {
			nonTiles = append(nonTiles, img)
		}
	}
	allImages = append(allImages, nonTiles...)

	// add portals as images
	for _, portal := range _map.Portals {
		// blank ebiten image for collision detection with portal
		am := models.Image{
			ID:        portal.MapID,
			UserID:    _map.UserID,
			Name:      "portal-" + portal.MapID,
			AssetType: models.MapPortal,
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
