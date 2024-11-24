package assets

import (
	"github.com/hajimehoshi/ebiten"
)

type Image struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Name string `json:"name"`
	// Path          string      `json:"path"`
	Width  int `json:"width"`
	Height int `json:"height"`
	// Frames        []FrameSpec `json:"frames"`
	Data          PixelData `json:"data"`
	*ebiten.Image `json:"-"`
}

// func GenerateSprites(assets *Assets) {
// 	count := len(assets.Images.Sprites)
// 	log.Printf("Generating %d characters", count)
// 	for k, s := range assets.Images.Sprites {

// 		// Generate the sprite image
// 		img, err := imageFromBytes(s.Data)
// 		if err != nil {
// 			log.Printf("error decoding sprite: %s", err)
// 			continue
// 		}
// 		s.Image = img
// 		assets.Images.Sprites[k] = s
// 		log.Printf("generated sprite size: %d", len(s.Data))
// 	}
// }
