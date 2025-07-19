package maze

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MazeBlock = iota
	MazeFree
	MazeGoal
)

var (
	MazeImg   = []*ebiten.Image{}
	GopherImg *ebiten.Image
)

const (
	MazeTileSize int = 64
)

func init() {
	MazeImg = make([]*ebiten.Image, 3, 3)
	MazeImg[MazeFree] = ebiten.NewImage(MazeTileSize, MazeTileSize)
	MazeImg[MazeFree].Fill(color.White)
	MazeImg[MazeBlock] = ebiten.NewImage(MazeTileSize, MazeTileSize)
	MazeImg[MazeBlock].Fill(color.RGBA{100, 100, 100, 255})
	MazeImg[MazeGoal] = ebiten.NewImage(MazeTileSize, MazeTileSize)
	MazeImg[MazeGoal].Fill(color.RGBA{0, 255, 0, 255})

	GiG := ebiten.NewImage(MazeTileSize/4, MazeTileSize)
	GiG.Fill(color.RGBA{0, 170, 170, 255})
	GopherImg = ebiten.NewImage(MazeTileSize, MazeTileSize)
	GopherImg.Fill(color.RGBA{0, 100, 100, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((MazeTileSize/4)*3), 0)
	GopherImg.DrawImage(GiG, op)
}

type Gopher struct {
	X     int
	Y     int
	Angle int
}

func (goph *Gopher) Walk(mz *Maze, i int) {
	for range i {
		switch goph.Angle {
		case 0:
			goph.Right(mz)
		case 1:
			goph.Down(mz)
		case 2:
			goph.Left(mz)
		case 3:
			goph.Up(mz)
		}
	}
}

func (goph *Gopher) Up(mz *Maze) {
	if goph.Y == 0 {
		return
	}
	if len(mz.Map) < goph.Y || len(mz.Map) == 0 {
		return
	}
	if len(mz.Map[goph.Y]) <= goph.X || len(mz.Map[goph.Y]) == 0 {
		return
	}
	if mz.Map[goph.Y-1][goph.X] == MazeBlock {
		return
	}
	goph.Y--
}

func (goph *Gopher) Down(mz *Maze) {
	if goph.Y == len(mz.Map)-1 {
		return
	}
	if len(mz.Map) < goph.Y || len(mz.Map) == 0 {
		return
	}
	if len(mz.Map[goph.Y]) <= goph.X || len(mz.Map[goph.Y]) == 0 {
		return
	}
	if mz.Map[goph.Y+1][goph.X] == MazeBlock {
		return
	}
	goph.Y++
}

func (goph *Gopher) Left(mz *Maze) {
	if goph.X == 0 {
		return
	}
	if len(mz.Map) <= goph.Y || len(mz.Map) < 0 {
		return
	}
	if len(mz.Map[goph.Y]) <= goph.X || len(mz.Map[goph.Y]) == 0 {
		return
	}
	if mz.Map[goph.Y][goph.X-1] == MazeBlock {
		return
	}
	goph.X--
}

func (goph *Gopher) Right(mz *Maze) {
	if len(mz.Map) <= goph.Y || len(mz.Map) < 0 {
		return
	}
	if len(mz.Map[goph.Y]) <= goph.X+1 || len(mz.Map[goph.Y]) < 0 {
		return
	}
	if mz.Map[goph.Y][goph.X+1] == MazeBlock {
		return
	}
	goph.X++
}

func (goph *Gopher) Rotate(i int) {
	fmt.Println(i)
	goph.Angle = (goph.Angle + i) % 4
	if goph.Angle < 0 {
		goph.Angle = 4 + goph.Angle
	}
}

// integer between 0-3
func (goph *Gopher) SetAngle(i int) {
	if i != 0 && i < 4 {
		return
	}
	goph.Angle = i
}

type Maze struct {
	SizeX  int
	SizeY  int
	Map    [][]int // order x y
	Gopher []Gopher
}

func (maze *Maze) FillFree() {
	for x, ct := range maze.Map {
		for y, _ := range ct {
			maze.Map[x][y] = MazeFree
		}
	}
}

func (maze *Maze) FillBlock() {
	for x, ct := range maze.Map {
		for y, _ := range ct {
			maze.Map[x][y] = MazeBlock
		}
	}
}

func DrawMaze(ebitenScr *ebiten.Image, mazeMap *Maze) {
	for y, c := range mazeMap.Map {
		for x, ct := range c {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(MazeTileSize*x), float64(MazeTileSize*y))
			ebitenScr.DrawImage(MazeImg[ct], op)
			if mazeMap.Gopher[0].X == x &&
				mazeMap.Gopher[0].Y == y {
				op2 := &ebiten.DrawImageOptions{}
				switch mazeMap.Gopher[0].Angle {
				case 0:
					ebitenScr.DrawImage(GopherImg, op)
				case 1:
					op2.GeoM.Rotate(float64(math.Pi / 2))
					op2.GeoM.Translate(float64(MazeTileSize*(x+1)), float64(MazeTileSize*y))
					ebitenScr.DrawImage(GopherImg, op2)
				case 2:
					op2.GeoM.Rotate(float64(math.Pi))
					op2.GeoM.Translate(float64(MazeTileSize*(x+1)), float64(MazeTileSize*(y+1)))
					ebitenScr.DrawImage(GopherImg, op2)
				case 3:
					op2.GeoM.Rotate(float64(3 * math.Pi / 2))
					op2.GeoM.Translate(float64(MazeTileSize*x), float64(MazeTileSize*(y+1)))
					ebitenScr.DrawImage(GopherImg, op2)
				default:
					fmt.Println("unexpected")
				}
			}
		}
	}
}
