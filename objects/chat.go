package objects

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	uiImage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/snburman/game/config"
	"github.com/snburman/game/models"
	"golang.org/x/image/font/gofont/goregular"
)

type Chat struct {
	Object    *Object
	ui        *ebitenui.UI
	container *widget.Container
	textInput *widget.TextInput
	width     int
	height    int
}

func NewChat() *Chat {
	c := &Chat{
		width:  config.ScreenWidth,
		height: 37,
	}
	c.container = c.newContainer()
	c.ui = &ebitenui.UI{
		Container: c.container,
	}

	rect := image.Rect(0, 0, c.width, c.height)
	_img := ebiten.NewImageFromImage(image.NewRGBA(rect))
	_img.Fill(color.RGBA{255, 255, 255, 255})

	img := models.Image{
		Name:      "chat",
		AssetType: models.Object,
		Width:     c.width,
		Height:    c.height,
		X:         0,
		Y:         config.ScreenHeight - c.height,
		Image:     _img,
	}

	c.Object = NewObject(img, ObjectOptions{
		Position: Position{
			X: img.X,
			Y: img.Y,
		},
		Direction: Down,
	})

	return c
}

func (c *Chat) newContainer() *widget.Container {
	face, err := loadFont(16)
	if err != nil {
		panic(err)
	}

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(uiImage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255})),

		// the container will use a row layout to layout the textinput widgets
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			// widget.RowLayoutOpts.Spacing(20),
			// widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
		)),
	)

	c.textInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),

		// Set the keyboard type when opened on mobile devices.
		widget.TextInputOpts.MobileInputMode("text"),

		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     uiImage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: uiImage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(face),

		//Set the colors for the text and caret
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		//Set how much padding there is between the edge of the input and the text
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(10)),

		//Set the font and width of the caret
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 2),
		),

		//This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("Type to chat..."),

		//This is called when the user hits the "Enter" key.
		//There are other options that can configure this behavior
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
			c.textInput.SetText("")
		}),

		//This is called whenver there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)
		}),
	)

	rootContainer.AddChild(c.textInput)

	return rootContainer
}

func (c *Chat) Update(g IGame, tick uint) {
	c.ui.Update()

	for _, touch := range g.TouchManager().Touches() {
		if c.Object.IsPressed(touch.X, touch.Y) {
			c.textInput.Focus(true)
		}
	}
}

func (c *Chat) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(c.Object.Position().X), float64(c.Object.Position().Y))

	c.ui.Draw(c.Object.Image())
	screen.DrawImage(c.Object.Image(), opts)
}

func (c *Chat) Container() *widget.Container {
	return c.container
}

func (c *Chat) TextInput() *widget.TextInput {
	return c.textInput
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
