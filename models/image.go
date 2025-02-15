package models

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Tile        AssetType = "tile"
	Object      AssetType = "object"
	PlayerUp    AssetType = "player_up"
	PlayerDown  AssetType = "player_down"
	PlayerLeft  AssetType = "player_left"
	PlayerRight AssetType = "player_right"
	MapPortal   AssetType = "portal"
)

type Image struct {
	ID            string    `json:"_id"`
	UserID        string    `json:"user_id"`
	X             int       `json:"x"`
	Y             int       `json:"y"`
	Name          string    `json:"name"`
	AssetType     AssetType `json:"asset_type"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	Data          PixelData `json:"data"`
	*ebiten.Image `json:"-"`
}

type AssetType string

type Pixel struct {
	X int `json:"x"`
	Y int `json:"y"`
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
	A int `json:"a"`
}

type PixelData = [][]Pixel

func ImageFromPixelData(img Image) *ebiten.Image {
	// craete rectangle using width and height
	rect := image.NewRGBA(image.Rect(0, 0, img.Width, img.Height))
	// fill image with squares with rgba values
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			rect.Set(x, y, color.RGBA{
				R: uint8(img.Data[y][x].R),
				G: uint8(img.Data[y][x].G),
				B: uint8(img.Data[y][x].B),
				A: uint8(img.Data[y][x].A),
			})
		}
	}
	eImg := ebiten.NewImageFromImage(rect)
	return eImg
}
