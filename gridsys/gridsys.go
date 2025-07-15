// gridsys is graphical editor
package gridsys

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	StartBlockImg *ebiten.Image
	BlockBlueImg  *ebiten.Image
	NormalGridImg *ebiten.Image
	PointerImg    *ebiten.Image
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
	StartBlockImg, _, err = ebitenutil.NewImageFromReader(sB)
	handleErr(err)
	bB := bytes.NewReader(gridassets.BlockBlue)
	BlockBlueImg, _, err = ebitenutil.NewImageFromReader(bB)
	handleErr(err)
	nG := bytes.NewReader(gridassets.NormalGrid)
	NormalGridImg, _, err = ebitenutil.NewImageFromReader(nG)
	handleErr(err)
	pT := bytes.NewReader(gridassets.Pointer)
	PointerImg, _, err = ebitenutil.NewImageFromReader(pT)
	handleErr(err)
}

type BlockKind int

const (
	StartBlock = BlockKind(iota)
	IfBlock
	WalkBlock
	TurnRightBlock
	TurnLeftBlock
	FlipBlock
)

type BooleanKind int

const (
	FrontIsWall = BooleanKind(iota)
	BackIsWall
	LeftIsWall
)

type CodeBlock struct {
	Kind BlockKind

	X int // It's a grid!!
	Y int // Grid too!

	upper int64
	lower int64

	// only for if block
	mid     int64
	boolean BooleanKind
}

type Gridsys struct {
	SizeMult float64 // multiplier
	tX       float64 // translate X
	tY       float64 // translate Y

	HeadBlocks []int64
	Blocks     map[int64]CodeBlock
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
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		gsys.tX = gsys.tX + 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		gsys.tX = gsys.tX - 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		gsys.tY = gsys.tY + 0.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		gsys.tY = gsys.tY - 0.05
	}
}

// reset all and add one head block
func (gsys *Gridsys) InitializeSpace() {
	gsys.HeadBlocks = append(gsys.HeadBlocks, 1)
	gsys.Blocks[1] = CodeBlock{
		Kind: StartBlock,
	}
}

// just for settings
type Vec2 struct {
	X int
	Y int
}

func (gsys *Gridsys) DrawBlock(ebitenScr *ebiten.Image, codeblock *CodeBlock, pos Vec2) {
	switch codeblock.Kind {
	case StartBlock:

	case IfBlock:
	case WalkBlock:
	case TurnRightBlock:
	case TurnLeftBlock:
	case FlipBlock:
	}
}

func (gsys *Gridsys) DrawAllBlocks(ebitenScr *ebiten.Image, pos Vec2) {
	for _, id := range gsys.HeadBlocks {
		fmt.Println(id)
	}
}

func (gsys *Gridsys) Draw(ebitenScr *ebiten.Image, pos Vec2, size Vec2) {
	// GRID
	wx, wy := size.X, size.Y
	baseX, baseY := math.Mod(gsys.tX, 1), math.Mod(gsys.tY, 1)
	var relTileX, relTileY int
	if 0 <= baseX {
		relTileX = int(math.Floor(gsys.tX))
	} else {
		relTileX = int(math.Floor(gsys.tX))
	}
	if 0 <= baseY {
		relTileY = int(math.Floor(gsys.tY))
	} else {
		relTileY = int(math.Floor(gsys.tY))
	}
	dc := false
	op := &ebiten.DrawImageOptions{}
	var mx, my float64
	rs := 10 * int(gsys.SizeMult)
	for y := -2; y < wy/rs*2; y++ {
		for x := -2; x < wx/rs+2; x++ {
			op.GeoM.Scale(float64(gsys.SizeMult), float64(gsys.SizeMult))
			op.GeoM.Translate(
				(float64(x*10)*gsys.SizeMult)+baseX*gsys.SizeMult*10,
				(float64(y*10)*gsys.SizeMult)+baseY*gsys.SizeMult*10,
			)
			ebitenScr.DrawImage(NormalGridImg, op)
			op.GeoM.Reset()
			// draw if pointer is available
			if x == PointerX+relTileX && y == PointerY+relTileY {
				var adjx, adjy int
				if 0 > relTileX {
					adjx += rs
				}
				if 0 > relTileY {
					adjy += rs
				}
				op := &ebiten.DrawImageOptions{}
				//op.GeoM.Scale(float64(gsys.SizeMult), float64(gsys.SizeMult))
				mx = (float64(x*10) * gsys.SizeMult) + (baseX * (gsys.SizeMult * 10)) + float64(adjx)
				my = (float64(y*10) * gsys.SizeMult) + (baseY * (gsys.SizeMult * 10)) + float64(adjy)
				//fmt.Println(math.Floor(mx), math.Floor(my))
				//op.GeoM.Translate(mx, my)
				//ebitenScr.DrawImage(Pointer, op)
				op.GeoM.Reset()
				dc = true
			}
			//fmt.Println("x", x, "relTileX", relTileX, "PointerX", PointerX, "=>", x-relTileX)
		}
	}
	if dc {
		var adjx, adjy int
		if 0 > relTileX {
			adjx += 10 * int(gsys.SizeMult)
		}
		if 0 > relTileY {
			adjy += 10 * int(gsys.SizeMult)
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(gsys.SizeMult), float64(gsys.SizeMult))
		//mx := (float64(baseX*10) * gsys.SizeMult) + (baseX * (gsys.SizeMult * 10)) + float64(adjx)
		//my := (float64(baseY*10) * gsys.SizeMult) + (baseY * (gsys.SizeMult * 10)) + float64(adjy)
		op.GeoM.Translate(mx, my)
		ebitenScr.DrawImage(PointerImg, op)
		op.GeoM.Reset()
	}

	// palette
	separate := 7
	edratio := 4
	vector.DrawFilledRect(ebitenScr, float32(wx/separate*edratio), 0, float32(wx/separate*separate-edratio), float32(wy), color.RGBA{150, 150, 150, 255}, false)
}
