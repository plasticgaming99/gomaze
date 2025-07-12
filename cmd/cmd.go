package cmd

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/plasticgaming99/gomaze/composite"
	"github.com/plasticgaming99/gomaze/maze"
)

var (
	MazeMap []maze.Maze
	Editor  composite.EditGoph
	comp    composite.Compositor
)

func repeatingKeyPressed(key ebiten.Key) bool {
	var (
		delay    = ebiten.TPS() / 2
		interval = ebiten.TPS() / 18
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
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

// init compositor
func init() {
	Editor = composite.NewEditor()
}

func (g *Gaem) Update() error {
	if repeatingKeyPressed(ebiten.KeyUp) {
		MazeMap[mzmap].Gopher[0].Up(&MazeMap[mzmap])
	} else if repeatingKeyPressed(ebiten.KeyDown) {
		MazeMap[mzmap].Gopher[0].Down(&MazeMap[mzmap])
	} else if repeatingKeyPressed(ebiten.KeyLeft) {
		MazeMap[mzmap].Gopher[0].Left(&MazeMap[mzmap])
	} else if repeatingKeyPressed(ebiten.KeyRight) {
		MazeMap[mzmap].Gopher[0].Right(&MazeMap[mzmap])
	} else
	// separator
	if repeatingKeyPressed(ebiten.KeyE) {
		MazeMap[mzmap].Gopher[0].Rotate(1)
	} else if repeatingKeyPressed(ebiten.KeyW) {
		MazeMap[mzmap].Gopher[0].Walk(&MazeMap[mzmap], 1)
	} else if repeatingKeyPressed(ebiten.KeyQ) {
		MazeMap[mzmap].Gopher[0].Rotate(-1)
	}
	return nil
}

func (g *Gaem) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 0})
	maze.DrawMaze(screen, &MazeMap[mzmap])
	comp.Draw(screen, Editor)
}

func (g *Gaem) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenWidth, screenHeight := ebiten.WindowSize()
	return screenWidth, screenHeight
}

func Exec() {
	ebiten.SetWindowResizable(true)
	//maze := gomaze.Maze{}
	mzmap = 1
	gaem := Gaem{}
	ebiten.RunGame(&gaem)
}
