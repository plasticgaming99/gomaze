package composite

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/plasticgaming99/gomaze/gridsys"
)

type EditGoph struct {
	EbitenScr *ebiten.Image
	graphical bool
	Gridsys   gridsys.Gridsys
}

func NewEditor() EditGoph {
	return EditGoph{
		EbitenScr: ebiten.NewImage(ebiten.Monitor().Size()),
		graphical: true,
	}
}

func (gph *EditGoph) SetGraphical() {
	gph.graphical = true
}

func (gph *EditGoph) SetTextmode() {
	gph.graphical = false
}

func (gph *EditGoph) Drawer() {
	switch gph.graphical {
	case true:
		gph.Gridsys.Draw(gph.EbitenScr)
	case false:
		//gph.Photon.PhotonDrawer(gph.EbitenScr)
	}
}

func (gph *EditGoph) Ticker() {
	switch gph.graphical {
	case true:
		gph.Gridsys.Tick()
	case false:
		//gph.Photon.PhotonTicker()
	}
}

type Compositor struct {
	EditorImg *ebiten.Image
	MazeImg   *ebiten.Image
}

func NewCompositor() Compositor {
	return Compositor{
		EditorImg: ebiten.NewImage(1920, 1080),
		MazeImg:   ebiten.NewImage(1920, 1080),
	}
}

func (cp *Compositor) Draw(img *ebiten.Image, edit EditGoph) {
	edit.Drawer()
	img.DrawImage(edit.EbitenScr, nil)
}

func (cp *Compositor) Tick(edit EditGoph) {
	edit.Ticker()
}
