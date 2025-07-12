package composite

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/plasticgaming99/gomaze/gabagui"
	"github.com/plasticgaming99/gomaze/photon"
)

type EditGoph struct {
	EbitenScr *ebiten.Image
	graphical bool
	Photon    FlexPhoton
	GabaGUI   gabagui.GabaGUI
}

func NewEditor() EditGoph {
	return EditGoph{
		EbitenScr: ebiten.NewImage(1920, 1080),
		graphical: false,
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
		gph.GabaGUI.Draw(gph.EbitenScr)
	case false:
		gph.Photon.PhotonDrawer(gph.EbitenScr)
	}
}

func (gph *EditGoph) Ticker() {
	switch gph.graphical {
	case true:
		gph.GabaGUI.Tick()
	case false:
		gph.Photon.PhotonTicker()
	}
}

type FlexPhoton struct {
	PhotonStruct *photon.Editor
	ScreenX      int
	ScreenY      int
}

type Compositor struct {
	ebitenScreen *ebiten.Image
	EditScreen   *ebiten.Image
	MazeScreen   *ebiten.Image
	WindowX      int
	WindowY      int
}

func (ph *FlexPhoton) PhotonDrawer(scr *ebiten.Image) {
	ph.PhotonStruct.Draw(scr)
}

func (ph *FlexPhoton) PhotonTicker() {
	ph.PhotonStruct.Layout(ph.ScreenX, ph.ScreenY)
	ph.PhotonStruct.Update()
}

func (cp *Compositor) Draw(img *ebiten.Image, edit EditGoph) {
	edit.Drawer()
	img.DrawImage(edit.EbitenScr, nil)
}

func (cp *Compositor) Tick(edit EditGoph) {
	edit.Ticker()
}
