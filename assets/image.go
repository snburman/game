package assets

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/hajimehoshi/ebiten"
)

const (
	Tile        AssetType = "tile"
	Object      AssetType = "object"
	PlayerUp    AssetType = "player_up"
	PlayerDown  AssetType = "player_down"
	PlayerLeft  AssetType = "player_left"
	PlayerRight AssetType = "player_right"
)

type Image struct {
	ID        string    `json:"_id"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Name      string    `json:"name"`
	AssetType AssetType `json:"asset_type"`
	// Path          string      `json:"path"`
	Width  int `json:"width"`
	Height int `json:"height"`
	// Frames        []FrameSpec `json:"frames"`
	Data          PixelData `json:"data"`
	*ebiten.Image `json:"-"`
}

func PngBytesFromFile(file io.Reader) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ImageFromBytes(data []byte) (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

func imageFromPixelData(img Image) (*ebiten.Image, error) {
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
	eImg, err := ebiten.NewImageFromImage(rect, ebiten.FilterDefault)

	return eImg, err
}
