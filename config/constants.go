package config

import "image/color"

var BackgroundColor = color.RGBA{
	0, 0, 0, 255,
}

const (
	Scale        = 3.5
	ScreenWidth  = 336
	ScreenHeight = 475
	// ScreenHeight = 475
	// game screen viewport
	ViewPortHeight = ScreenWidth
	ViewPortWidth  = ScreenWidth
	WalkSpeed      = 3
	RunSpeed       = 5
)
