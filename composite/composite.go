package composite

import (
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/plasticgaming99/gomaze/gridsys"
	"github.com/plasticgaming99/gomaze/maze"
)

type EditGoph struct {
	EbitenScr *ebiten.Image
	graphical bool
	Gridsys   gridsys.Gridsys
}

func NewEditor() EditGoph {
	gs := gridsys.Gridsys{
		SizeMult: 8,
	}
	return EditGoph{
		// autorescale may work
		EbitenScr: ebiten.NewImage(1, 1),
		graphical: true,
		Gridsys:   gs,
	}
}

func (gph *EditGoph) SetGraphical() {
	gph.graphical = true
}

func (gph *EditGoph) SetTextmode() {
	gph.graphical = false
}

const (
	separate = 9
	edratio  = 4
)

func (gph *EditGoph) Drawer() {
	wx, wy := ebiten.WindowSize()
	wl := gridsys.Vec2{
		X: 0,
		Y: 0,
	}
	ws := gridsys.Vec2{
		X: (wx / separate) * edratio,
		Y: wy,
	}
	if ws.X != gph.EbitenScr.Bounds().Dx() ||
		ws.Y != gph.EbitenScr.Bounds().Dy() {
		gph.EbitenScr = ebiten.NewImage(ws.X, ws.Y)
		runtime.GC()
	}
	if gph.graphical {
		gph.Gridsys.Draw(gph.EbitenScr, wl, ws)
	} else {
		//gph.Photon.PhotonDrawer(gph.EbitenScr)
	}
}

func (gph *EditGoph) Ticker() {
	if gph.graphical {
		//gph.Gridsys.Tick()
	} else {
		//gph.Photon.PhotonTicker()
	}
}

type Compositor struct {
	EditorImg *ebiten.Image
	MazeImg   *ebiten.Image
}

func NewCompositor() Compositor {
	msx, msy := ebiten.Monitor().Size()
	edimg := ebiten.NewImage(msx/5*2, msy)
	return Compositor{
		EditorImg: edimg,
		MazeImg:   ebiten.NewImage(ebiten.Monitor().Size()),
	}
}

func (cp *Compositor) Draw(img *ebiten.Image, edit EditGoph, mz *maze.Maze) {
	wx, _ := ebiten.WindowSize()
	edit.Drawer()
	maze.DrawMaze(cp.MazeImg, mz)
	img.DrawImage(edit.EbitenScr, nil)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((wx/separate)*(edratio)), 0)
	img.DrawImage(cp.MazeImg, op)
	op.GeoM.Reset()
}

/*func (cp *Compositor) Tick(edit EditGoph) {
	edit.Ticker()
}*/
