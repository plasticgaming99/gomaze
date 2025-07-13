// gridsys is graphical editor
package gridsys

import (
	"bytes"
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	rd "github.com/plasticgaming99/gomaze/_lib/randoms"
	gridassets "github.com/plasticgaming99/gomaze/gridsys/assets"
)

// i'm implementing
var (
	PointerX int
	PointerY int
	Move     bool
)

// let us set up internal image
var (
	StartBlock *ebiten.Image
	BlockBlue  *ebiten.Image
	NormalGrid *ebiten.Image
	Pointer    *ebiten.Image
)

// let us initialize
func init() {
	var err error
	handleErr := func(err error) {
		if err != nil {
			log.Fatal("error initializing grid assets")
		}
	}
	sB := bytes.NewReader(gridassets.StartBlock)
	StartBlock, _, err = ebitenutil.NewImageFromReader(sB)
	handleErr(err)
	bB := bytes.NewReader(gridassets.BlockBlue)
	BlockBlue, _, err = ebitenutil.NewImageFromReader(bB)
	handleErr(err)
	nG := bytes.NewReader(gridassets.NormalGrid)
	NormalGrid, _, err = ebitenutil.NewImageFromReader(nG)
	handleErr(err)
	pT := bytes.NewReader(gridassets.Pointer)
	Pointer, _, err = ebitenutil.NewImageFromReader(pT)
	handleErr(err)
}

type CodeBlock struct {
	X int // It's a grid!!
	Y int // Grid too!

	mouseMovementX int // Relative maybe
	mouseMovementY int

	upper int64
	lower int64
}

type Gridsys struct {
	SizeMult float64 // multiplier
	tX       float64 // translate X
	tY       float64 // translate Y

	HeadBlock []int64
	Blocks    map[int64]CodeBlock
}

// init new gridsys with default val
func New() *Gridsys {
	return &Gridsys{
		SizeMult: 10,
		tX:       0,
		tY:       0,
	}
}

// Get upper block
func (gsys *Gridsys) GetUpper(i int64) *int64 {
	i64 := gsys.Blocks[i].upper
	return &i64
}

func (gsys *Gridsys) GetLower(i int64) *int64 {
	i64 := gsys.Blocks[i].lower
	return &i64
}

func (gsys *Gridsys) Tick() {
	//x, y := ebiten.CursorPosition()
	switch {
	case rd.RepeatingKeyPressed(ebiten.KeyLeft):
		PointerX--
	case rd.RepeatingKeyPressed(ebiten.KeyRight):
		PointerX++
	case rd.RepeatingKeyPressed(ebiten.KeyUp):
		PointerY--
	case rd.RepeatingKeyPressed(ebiten.KeyDown):
		PointerY++
	}

	//tpsC := float64(ebiten.TPS()) / 100
	if ebiten.IsKeyPressed(ebiten.KeyNumpad4) {
		gsys.tX = gsys.tX + 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyNumpad6) {
		gsys.tX = gsys.tX - 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyNumpad8) {
		gsys.tY = gsys.tY + 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyNumpad2) {
		gsys.tY = gsys.tY - 0.05
	}
}

func (gsys *Gridsys) Draw(ebitenScr *ebiten.Image) {
	// GRID
	wx, wy := ebiten.WindowSize()
	baseX, baseY := math.Mod(gsys.tX, 1), math.Mod(gsys.tY, 1)
	var relTileX, relTileY int
	if 0 <= baseX {
		relTileX = int(math.Floor(gsys.tX))
	} else {
		relTileX = int(math.Floor(-gsys.tX))
	}
	if 0 <= baseY {
		relTileY = int(math.Floor(gsys.tY))
	} else {
		relTileY = int(math.Floor(-gsys.tY))
	}
	op := &ebiten.DrawImageOptions{}
	for y := -1; y < wy/100+2; y++ {
		for x := -1; x < wx/100+2; x++ {
			op.GeoM.Scale(float64(gsys.SizeMult), float64(gsys.SizeMult))
			op.GeoM.Translate(
				(float64(x*10)*gsys.SizeMult)+baseX*gsys.SizeMult*10,
				(float64(y*10)*gsys.SizeMult)+baseY*gsys.SizeMult*10,
			)
			ebitenScr.DrawImage(NormalGrid, op)
			op.GeoM.Reset()
			// draw if pointer is available
			if x == PointerX+relTileX && y == PointerY+relTileY {
				fmt.Println(
					(float64(x*10)*gsys.SizeMult)+baseX*gsys.SizeMult*10,
					(float64(y*10)*gsys.SizeMult)+baseY*gsys.SizeMult*10,
					PointerX,
					PointerY,
				)
				op.GeoM.Scale(float64(gsys.SizeMult), float64(gsys.SizeMult))
				op.GeoM.Translate(
					(float64(x*10)*gsys.SizeMult)+baseX*(gsys.SizeMult*10),
					(float64(y*10)*gsys.SizeMult)+baseY*(gsys.SizeMult*10),
				)
				ebitenScr.DrawImage(Pointer, op)
			}
			//fmt.Println("x", x, "relTileX", relTileX, "PointerX", PointerX, "=>", x-relTileX)

			op.GeoM.Reset()
		}
	}
}
