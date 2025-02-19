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
	api            *API
	conn           *Conn
	player         *objects.Player
	onlinePlayers  map[string]*objects.Player
	primaryMap     models.Map[[]models.Image]
	primaryObjects []objects.Objecter
	currentMap     models.Map[[]models.Image]
	currentObjects []objects.Objecter
	portalMaps     map[string]models.Map[[]models.Image]
	portalObjects  map[string][]objects.Objecter
}

func NewMapService(api *API) *MapService {
	ms := &MapService{
		api:            api,
		onlinePlayers:  map[string]*objects.Player{},
		primaryMap:     models.Map[[]models.Image]{},
		primaryObjects: []objects.Objecter{},
		currentMap:     models.Map[[]models.Image]{},
		currentObjects: []objects.Objecter{},
		portalMaps:     map[string]models.Map[[]models.Image]{},
		portalObjects:  map[string][]objects.Objecter{},
	}
	err := ms.GetPrimaryMap()
	if err != nil {
		log.Println("error getting primary map")
		panic(err)
	}
	ms.conn, err = NewConn(config.Env().WS_SERVER_URL+"/game/ws", ms)
	if err != nil {
		log.Println("error creating websocket connection")
		panic(err)
	}
	dispatch := NewDispatch(ms.conn, LoadNewOnlinePlayer, PlayerUpdate{
		UserID: ms.api.UserID(),
		MapID:  ms.primaryMap.ID.Hex(),
		Dir:    int(ms.player.Direction()),
		Pos:    *ms.player.Position(),
	})
	dispatch.MarshalAndPublish()

	return ms
}

func (ms *MapService) PrimaryMap() models.Map[[]models.Image] {
	return ms.primaryMap
}

// GetPrimaryMap makes a get request to server for primary map
func (ms *MapService) GetPrimaryMap() error {
	userID := ms.api.UserID()
	// get map
	_map := models.Map[[]models.Image]{}
	path := config.Env().SERVER_URL + "/game/wasm/map/primary/" + userID
	res := ms.api.Request(http.MethodGet, path)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}
	err := json.Unmarshal(res.Body, &_map)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// set map
	ms.SetCurrentMap(_map)
	ms.primaryMap = _map
	ms.player.Username = _map.UserName
	ms.primaryObjects = ms.currentObjects
	return nil
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
		log.Println(err.Error())
		return _map, err
	}
	return _map, nil
}

// GetPortalMaps makes a get request to server for all portal maps by ID
func (ms *MapService) GetPortalMaps(_map models.Map[[]models.Image]) error {
	// get updated map with portals
	updatedMap, err := ms.GetMapByID(_map.ID.Hex())
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// get all portal maps
	var ids []string
	for _, p := range updatedMap.Portals {
		ids = append(ids, p.MapID)
	}
	path := config.Env().SERVER_URL + "/game/wasm/map/ids" + "?ids="
	for i, id := range ids {
		if i == 0 {
			path += id
		} else {
			path += "&ids=" + id
		}
	}
	var _maps []models.Map[[]models.Image]
	res := ms.api.Request(http.MethodGet, path)
	if res.Error != nil {
		log.Println(res.Error.Error())
		return res.Error
	}
	err = json.Unmarshal(res.Body, &_maps)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// parse and set portal maps
	for _, _map := range _maps {
		ms.portalMaps[_map.ID.Hex()] = _map
		// ignore player object, already set
		objs, _ := objects.ObjectersFromImages(ms.ImagesFromMap(_map), ms.api.userID)
		// set portal objects
		ms.portalObjects[_map.ID.Hex()] = objs
	}
	return nil
}

func (ms *MapService) CurrentMap() models.Map[[]models.Image] {
	return ms.currentMap
}

// SetCurrentMap sets the current map and extracts object
// into ms.Player() and ms.CurrentObjects()
//
// it then fetches all portal maps in a go routine
func (ms *MapService) SetCurrentMap(_map models.Map[[]models.Image]) {
	// set map
	ms.currentMap = _map

	// extract images
	imgs := ms.ImagesFromMap(_map)

	// set objects
	objs, player := objects.ObjectersFromImages(imgs, ms.api.userID)
	ms.currentObjects = objs
	if ms.player == nil && player != nil {
		ms.player = player
	}
	if ms.player == nil && player == nil {
		ms.player = objects.NewDefaultPlayer(ms.api.userID, _map.Entrance.X, _map.Entrance.Y)
	}
	ms.player.SetPosition(objects.Position{
		X: _map.Entrance.X,
		Y: _map.Entrance.Y,
	})
	go ms.GetPortalMaps(_map)
}

func (ms *MapService) CurrentObjects() []objects.Objecter {
	return ms.currentObjects
}

func (ms *MapService) Player() *objects.Player {
	return ms.player
}

func (ms *MapService) OnlinePlayers() map[string]*objects.Player {
	return ms.onlinePlayers
}

func (ms *MapService) LoadOnlinePlayers(imgs []models.Image) {
	if len(imgs) == 0 {
		log.Println("no online players image data")
		ms.RemoveOnlinePlayers()
		return
	}
	ms.onlinePlayers = objects.PlayersFromImages(imgs)
}

func (ms *MapService) LoadNewOnlinePlayer(imgs []models.Image) {
	var player *objects.Player
	if len(imgs) == 0 {
		log.Println("no new online player image data")
		player = objects.NewDefaultPlayer(ms.api.userID, ms.currentMap.Entrance.X, ms.currentMap.Entrance.Y)
	} else {
		playerSlice := objects.PlayersFromImages(imgs)
		p, ok := playerSlice[imgs[0].UserID]
		if !ok {
			log.Println("player not found")
			return
		}
		player = p
	}
	player.SetPosition(objects.Position{
		X: ms.currentMap.Entrance.X,
		Y: ms.currentMap.Entrance.Y,
	})
	ms.onlinePlayers[player.ID()] = player
}

func (ms *MapService) RemoveOnlinePlayers() {
	ms.onlinePlayers = map[string]*objects.Player{}
}

func (ms *MapService) RemoveOnlinePlayerByID(id string) {
	delete(ms.onlinePlayers, id)
}

func (ms *MapService) UpdateLocalPlayer(update PlayerUpdate) {
	var player *objects.Player
	if update.UserID == ms.api.userID {
		player = ms.player
	} else {
		p, ok := ms.onlinePlayers[update.UserID]
		if !ok {
			return
		}
		player = p
	}
	player.SetPosition(update.Pos)
	player.SetDirection(objects.Direction(update.Dir))
}

func (ms *MapService) DispatchUpdatePlayer(update *objects.Player) {
	dispatch := NewDispatch(ms.conn, UpdatePlayer, PlayerUpdate{
		UserID: update.ID(),
		MapID:  ms.currentMap.ID.Hex(),
		Dir:    int(update.Direction()),
		Pos:    *update.Position(),
	})
	dispatch.MarshalAndPublish()
}

func (ms *MapService) LoadMap(g objects.IGame, id string) error {
	// dispatch player update after map load
	defer g.DispatchUpdatePlayer()

	// check if map is current
	if id == ms.currentMap.ID.Hex() {
		g.Player().SetPosition(objects.Position{
			X: ms.currentMap.Entrance.X,
			Y: ms.currentMap.Entrance.Y,
		})
		return nil
	}

	// check if map is primary
	if id == ms.primaryMap.ID.Hex() {
		// remove all online players locally
		ms.RemoveOnlinePlayers()
		// set primary map
		ms.currentMap = ms.primaryMap
		ms.currentObjects = ms.primaryObjects
		g.Player().SetPosition(objects.Position{
			X: ms.currentMap.Entrance.X,
			Y: ms.currentMap.Entrance.Y,
		})
		go ms.GetPortalMaps(ms.currentMap)
		return nil
	}

	// check if map is portal
	portal, ok := ms.portalMaps[id]
	if ok {
		// remove all online players locally
		ms.RemoveOnlinePlayers()
		// set portal objects
		objs, ok := ms.portalObjects[id]
		if !ok {
			panic("portal map objects not found")
		}
		// set portal map
		ms.currentMap = portal
		ms.currentObjects = objs
		g.Player().SetPosition(objects.Position{
			X: ms.currentMap.Entrance.X,
			Y: ms.currentMap.Entrance.Y,
		})
		go ms.GetPortalMaps(ms.currentMap)
		return nil
	}

	// if no objects, fetch map
	_map, err := ms.GetMapByID(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	// remove all online players locally
	ms.RemoveOnlinePlayers()
	ms.SetCurrentMap(_map)

	return nil
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
