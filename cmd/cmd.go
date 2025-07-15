package cmd

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	rd "github.com/plasticgaming99/gomaze/_lib/randoms"
	"github.com/plasticgaming99/gomaze/composite"
	"github.com/plasticgaming99/gomaze/gridsys"
	"github.com/plasticgaming99/gomaze/maze"
)

var (
	MazeMap []maze.Maze
	gridEd  = gridsys.New()
	comp    = composite.NewCompositor()
	edgp    = composite.NewEditor()
)

func init() {
	gridEd.SizeMult = 10
}

type Gaem struct {
	Maze maze.Maze
}

var mzmap int

// init maps
func init() {
	MazeMap = append(MazeMap, maze.Maze{
		SizeX: 7,
		SizeY: 7,
		Map: [][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 1, 1, 1, 1, 1, 1},
			{0, 1, 0, 0, 0, 0, 0},
			{0, 1, 0, 0, 0, 0, 0},
			{0, 1, 0, 0, 0, 0, 0},
		},
		Gopher: []maze.Gopher{
			{X: 1, Y: 6},
		},
	})

	MazeMap = append(MazeMap, maze.Maze{
		SizeX: 10,
		SizeY: 10,
		Map: [][]int{
			{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			{1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 1, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
		},
		Gopher: []maze.Gopher{
			{X: 0, Y: 1},
		},
	})
}

func (g *Gaem) Update() error {
	/*if rd.RepeatingKeyPressed(ebiten.KeyUp) {
		MazeMap[mzmap].Gopher[0].Up(&MazeMap[mzmap])
	} else if rd.RepeatingKeyPressed(ebiten.KeyDown) {
		MazeMap[mzmap].Gopher[0].Down(&MazeMap[mzmap])
	} else if rd.RepeatingKeyPressed(ebiten.KeyLeft) {
		MazeMap[mzmap].Gopher[0].Left(&MazeMap[mzmap])
	} else if rd.RepeatingKeyPressed(ebiten.KeyRight) {
		MazeMap[mzmap].Gopher[0].Right(&MazeMap[mzmap])
	} else*/
	// separator
	if rd.RepeatingKeyPressed(ebiten.KeyE) {
		MazeMap[mzmap].Gopher[0].Rotate(1)
	} else if rd.RepeatingKeyPressed(ebiten.KeyW) {
		MazeMap[mzmap].Gopher[0].Walk(&MazeMap[mzmap], 1)
	} else if rd.RepeatingKeyPressed(ebiten.KeyQ) {
		MazeMap[mzmap].Gopher[0].Rotate(-1)
	}
	//gridEd.Tick()
	edgp.Gridsys.Tick()
	return nil
}

// draw function
func (g *Gaem) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 50, 50})
	comp.Draw(screen, edgp, &MazeMap[mzmap])
	//gridEd.Draw(screen, v2, siz)
}

func (g *Gaem) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenWidth, screenHeight := ebiten.WindowSize()
	return screenWidth, screenHeight
}

func Exec() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// maze := gomaze.Maze{}
	mzmap = 1
	gaem := Gaem{}
	ebiten.RunGame(&gaem)
}
