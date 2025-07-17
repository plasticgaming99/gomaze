package composite

import (
	"fmt"
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
	gs := *gridsys.New()
	gs.SizeMult = 6
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
	wx, wy := ebiten.WindowSize()
	edit.Drawer()
	maze.DrawMaze(cp.MazeImg, mz)
	img.DrawImage(edit.EbitenScr, nil)

	op := &ebiten.DrawImageOptions{}

	// maze space size
	mx := float64(wx / separate * (separate - edratio))
	my := float64(wy)

	// tile ratio
	mtx := mx / float64(maze.MazeTileSize*mz.SizeX)
	mty := my / float64(maze.MazeTileSize*mz.SizeY)

	var (
		rescalebase float64
		zt          bool
	)
	if mtx < mty {
		rescalebase = mtx
		zt = true
	} else {
		rescalebase = mty
		zt = false
	}

	// just maze size
	msx := float64(maze.MazeTileSize*mz.SizeX) * rescalebase
	msy := float64(maze.MazeTileSize*mz.SizeY) * rescalebase

	//fmt.Println(mx, my, mtx, mty, rescalebase)

	op.GeoM.Scale(float64(rescalebase), float64(rescalebase))
	basex := float64((wx / separate) * (edratio))
	if zt {
		op.GeoM.Translate(basex, (my-msy)/2)
	} else {
		fmt.Println(my, mx/mtx)
		op.GeoM.Translate((mx-msx)/2+basex, 0)
	}
	img.DrawImage(cp.MazeImg, op)
	op.GeoM.Reset()
}

/*func (cp *Compositor) Tick(edit EditGoph) {
	edit.Ticker()
}*/
