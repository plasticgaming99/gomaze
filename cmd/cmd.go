package cmd

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	rd "github.com/plasticgaming99/gomaze/_lib/randoms"
	"github.com/plasticgaming99/gomaze/composite"
	"github.com/plasticgaming99/gomaze/gridsys"
	"github.com/plasticgaming99/gomaze/maze"
)

var (
	MazeMap   []maze.Maze
	DefGopher []gridsys.Vec2
	DefAxis   []int
	gridEd    = gridsys.New()
	comp      = composite.NewCompositor()
	edgp      = composite.NewEditor()
)

func init() {
	gridEd.SizeMult = 10
	edgp.Gridsys.InitializeSpace()
}

type Gaem struct {
	Maze  maze.Maze
	count int
}

var mzmap int

func init2() {
	MazeMap = append(MazeMap, maze.Maze{
		SizeX: 7,
		SizeY: 7,
		Map: [][]int{
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0},
			{0, 2, 0, 0, 0, 0, 0},
			{0, 1, 0, 0, 0, 0, 0},
			{0, 1, 0, 0, 0, 0, 0},
			{0, 1, 0, 0, 0, 0, 0},
		},
		Gopher: []maze.Gopher{
			{X: 1, Y: 6, Angle: 3},
		},
	})
	DefGopher = append(DefGopher, gridsys.Vec2{X: 1, Y: 6})
	DefAxis = append(DefAxis, 3)

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
			{0, 0, 0, 0, 0, 0, 0, 0, 1, 2},
		},
		Gopher: []maze.Gopher{
			{X: 0, Y: 1},
		},
	})
	DefGopher = append(DefGopher, gridsys.Vec2{X: 0, Y: 1})
	DefAxis = append(DefAxis, 0)
}

// init maps
func init() {
	init2()
}

var toggler = false

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

	if rd.RepeatingKeyPressed(ebiten.KeyR) {
		MazeMap[mzmap].Gopher[0].X = DefGopher[mzmap].X
		MazeMap[mzmap].Gopher[0].Y = DefGopher[mzmap].Y
		MazeMap[mzmap].Gopher[0].Angle = DefAxis[mzmap]

		gridsys.PointerX, gridsys.PointerY = 1, 1
	}

	if rd.RepeatingKeyPressed(ebiten.KeySpace) {
		fmt.Println("w")
		edgp.Gridsys.InterpretTick(&MazeMap[mzmap])
		fmt.Println(gridsys.Interpret)
		toggler = !toggler
		if toggler {
			gridsys.Interpret = true
		} else {
			gridsys.Interpret = false
			//MazeMap[mzmap].Gopher[0]
		}
	}

	if rd.RepeatingKeyPressed(ebiten.KeyEnter) {
		mzmap++
		MazeMap[mzmap].Gopher[0].X = DefGopher[mzmap].X
		MazeMap[mzmap].Gopher[0].Y = DefGopher[mzmap].Y
		MazeMap[mzmap].Gopher[0].Angle = DefAxis[mzmap]
	}

	//gridEd.Tick()
	edgp.Gridsys.Tick(&MazeMap[mzmap])
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
	ebiten.SetTPS(60)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// maze := gomaze.Maze{}
	mzmap = 0
	gaem := Gaem{}
	ebiten.RunGame(&gaem)
}
