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
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	rd "github.com/plasticgaming99/gomaze/_lib/randoms"
	"github.com/plasticgaming99/gomaze/assets/fonts"
	gridassets "github.com/plasticgaming99/gomaze/gridsys/assets"
	"github.com/plasticgaming99/gomaze/gridsys/gridlocale"
)

// i'm implementing
var (
	PointerX int
	PointerY int
	Move     bool
)

// let us set up internal image
var (
	StartBlockImg                 *ebiten.Image
	StartBlockCapImg              *ebiten.Image
	BracketBlockImg               *ebiten.Image
	BracketBlockEndImg            *ebiten.Image
	ForBlockHorizontalImg         *ebiten.Image
	IfBlockHorizontalFirstImg     *ebiten.Image
	IfBlockHorizontalEmptyImg     *ebiten.Image
	IfBlockHorizontalExtentionImg *ebiten.Image
	BlockBlueImg                  *ebiten.Image
	NormalGridImg                 *ebiten.Image
	PointerImg                    *ebiten.Image
)

// fonts
var (
	misakiGothic2ndSrc *text.GoTextFaceSource
)

// mainly init image assets
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

	sBC := bytes.NewReader(gridassets.StartBlockCap)
	StartBlockCapImg, _, err = ebitenutil.NewImageFromReader(sBC)
	handleErr(err)

	brB := bytes.NewReader(gridassets.BracketBlock)
	BracketBlockImg, _, err = ebitenutil.NewImageFromReader(brB)
	handleErr(err)

	brBE := bytes.NewReader(gridassets.BracketEnd)
	BracketBlockEndImg, _, err = ebitenutil.NewImageFromReader(brBE)
	handleErr(err)

	fBH := bytes.NewReader(gridassets.ForBlockHorz)
	ForBlockHorizontalImg, _, err = ebitenutil.NewImageFromReader(fBH)
	handleErr(err)

	iBHF := bytes.NewReader(gridassets.IfBlockHorzFirst)
	IfBlockHorizontalFirstImg, _, err = ebitenutil.NewImageFromReader(iBHF)
	handleErr(err)

	iBHEm := bytes.NewReader(gridassets.IfBlockHorzEmp)
	IfBlockHorizontalEmptyImg, _, err = ebitenutil.NewImageFromReader(iBHEm)
	handleErr(err)

	iBHEx := bytes.NewReader(gridassets.IfBlockHorzExt)
	IfBlockHorizontalExtentionImg, _, err = ebitenutil.NewImageFromReader(iBHEx)
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

// init fonts
func init() {
	var err error
	misakiGothic2ndSrc, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MisakiGothic2ndFont))
	if err != nil {
		log.Fatal(err)
	}
}

type BlockKind int

const (
	StartBlock = BlockKind(iota)
	IfBlock
	ForInfBlock
	ForRangeBlock
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

	Lower *int64

	// only for if block
	Mid     *int64
	Boolean *BooleanKind
}

type Gridsys struct {
	SizeMult float64 // multiplier
	tX       float64 // translate X
	tY       float64 // translate Y

	HeadBlocks []int64
	Blocks     map[int64]*CodeBlock

	displayFont *text.GoTextFace
}

// init new gridsys with default val
func New() *Gridsys {
	return &Gridsys{
		SizeMult: 10,
		tX:       0,
		tY:       0,
		Blocks:   make(map[int64]*CodeBlock),
	}
}

// Get upper block
/*func (gsys *Gridsys) GetUpper(i int64) *int64 {
	i64 := gsys.Blocks[i].Upper
	return i64
}*/

func (gsys *Gridsys) GetLower(i int64) *int64 {
	i64 := gsys.Blocks[i].Lower
	return i64
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
	gsys.HeadBlocks = append(gsys.HeadBlocks, 0)
	gsys.Blocks[0] = &CodeBlock{
		Kind: StartBlock,
		X:    1,
		Y:    1,
	}
	// test
	one := int64(1)
	gsys.Blocks[0].Lower = &one
	gsys.Blocks[1] = &CodeBlock{
		Kind: ForInfBlock,
		X:    1,
		Y:    2,
	}

	two := int64(2)
	gsys.Blocks[1].Lower = &two
	gsys.Blocks[2] = &CodeBlock{
		Kind: WalkBlock,
		X:    1,
		Y:    3,
	}
}

// just for settings
type Vec2 struct {
	X int
	Y int
}

type Vec2F struct {
	X float64
	Y float64
}

func (gsys *Gridsys) DrawBlockPart(ebitenScr *ebiten.Image, color color.RGBA, shadowColor color.RGBA, x float64, y float64, lengx float64) {
	vector.DrawFilledRect(ebitenScr, float32(x), float32(y), float32(lengx), 10*float32(gsys.SizeMult), color, false)
	vector.DrawFilledRect(ebitenScr, float32(x), float32(y)+float32(9*gsys.SizeMult), float32(lengx), 1*float32(gsys.SizeMult), shadowColor, false)
}

func (gsys *Gridsys) DrawBlueBlock(ebitenScr *ebiten.Image, str string, pos Vec2F, length float64) {
	op := &ebiten.DrawImageOptions{}
	top := &text.DrawOptions{}
	misakiGothic2ndFace := &text.GoTextFace{Source: misakiGothic2ndSrc, Size: 5 * gsys.SizeMult}

	op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
	op.GeoM.Translate(pos.X, pos.Y)
	ebitenScr.DrawImage(StartBlockImg, op)
	op.GeoM.Reset()
	op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
	op.GeoM.Translate(pos.X, pos.Y-(3*gsys.SizeMult))
	ebitenScr.DrawImage(StartBlockCapImg, op)
	gsys.DrawBlockPart(ebitenScr, color.RGBA{86, 147, 255, 255}, color.RGBA{224, 192, 0, 255}, pos.X+(gsys.SizeMult*float64(StartBlockImg.Bounds().Dx())), pos.Y, 100)
	top.GeoM.Translate(pos.X+gsys.SizeMult, pos.Y+gsys.SizeMult)
	text.Draw(ebitenScr, str, misakiGothic2ndFace, top)
}

func (gsys *Gridsys) IsTouched(codeblock *CodeBlock) {
	switch codeblock.Kind {
	case StartBlock:
	case IfBlock:
	case ForInfBlock:
	case ForRangeBlock:
	case WalkBlock:
	case TurnRightBlock:
	case TurnLeftBlock:
	case FlipBlock:
	}
}

func (gsys *Gridsys) DrawBlock(ebitenScr *ebiten.Image, codeblock *CodeBlock, pos Vec2F, nestc int) {
	op := &ebiten.DrawImageOptions{}
	top := &text.DrawOptions{}
	misakiGothic2ndFace := &text.GoTextFace{Source: misakiGothic2ndSrc, Size: 5 * gsys.SizeMult}
	switch codeblock.Kind {
	case StartBlock:
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X, pos.Y)
		ebitenScr.DrawImage(StartBlockImg, op)
		op.GeoM.Reset()
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X, pos.Y-(3*gsys.SizeMult))
		ebitenScr.DrawImage(StartBlockCapImg, op)
		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, pos.X+(gsys.SizeMult*float64(StartBlockImg.Bounds().Dx())), pos.Y, 100)
		top.GeoM.Translate(pos.X+gsys.SizeMult, pos.Y+gsys.SizeMult)
		text.Draw(ebitenScr, gridlocale.Start, misakiGothic2ndFace, top)
	case IfBlock:
		forinflen := gsys.SizeMult * 10 * 3
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X, pos.Y)
		ebitenScr.DrawImage(BracketBlockImg, op)
		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y, forinflen)
		op.GeoM.Reset()
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y)
		ebitenScr.DrawImage(StartBlockImg, op)
		op.GeoM.Reset()
		top.GeoM.Translate(pos.X+(gsys.SizeMult*10), pos.Y+gsys.SizeMult)
		text.Draw(ebitenScr, gridlocale.ForInf, misakiGothic2ndFace, top)
		if nestc < 1 {
			nestc = 1
		}
		for i := 1; i < nestc+1; i++ {
			fmt.Println(i)
			op.GeoM.Reset()
			op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
			op.GeoM.Translate(pos.X, pos.Y+(10*gsys.SizeMult*float64(i)))
			ebitenScr.DrawImage(BracketBlockImg, op)
		}
		op.GeoM.Reset()
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X, pos.Y+(10*gsys.SizeMult*float64(nestc+1)))
		ebitenScr.DrawImage(BracketBlockEndImg, op)
		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y+(10*gsys.SizeMult*float64(nestc+1)), forinflen)
	case ForInfBlock:
		forinflen := gsys.SizeMult * 10 * 3
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X, pos.Y)
		ebitenScr.DrawImage(BracketBlockImg, op)
		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y, forinflen)
		op.GeoM.Reset()
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y)
		ebitenScr.DrawImage(StartBlockImg, op)
		op.GeoM.Reset()
		top.GeoM.Translate(pos.X+(gsys.SizeMult*10), pos.Y+gsys.SizeMult)
		text.Draw(ebitenScr, gridlocale.If, misakiGothic2ndFace, top)
		if nestc < 1 {
			nestc = 1
		}
		for i := 1; i < nestc+1; i++ {
			fmt.Println(i)
			op.GeoM.Reset()
			op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
			op.GeoM.Translate(pos.X, pos.Y+(10*gsys.SizeMult*float64(i)))
			ebitenScr.DrawImage(BracketBlockImg, op)
		}
		op.GeoM.Reset()
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X, pos.Y+(10*gsys.SizeMult*float64(nestc+1)))
		ebitenScr.DrawImage(BracketBlockEndImg, op)
		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y+(10*gsys.SizeMult*float64(nestc+1)), forinflen)
	case ForRangeBlock:
		forinflen := gsys.SizeMult * 10 * 3
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X, pos.Y)
		ebitenScr.DrawImage(BracketBlockImg, op)
		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y, forinflen)
		op.GeoM.Reset()
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y)
		ebitenScr.DrawImage(StartBlockImg, op)
		op.GeoM.Reset()
		top.GeoM.Translate(pos.X+(gsys.SizeMult*10), pos.Y+gsys.SizeMult)
		text.Draw(ebitenScr, gridlocale.ForInf, misakiGothic2ndFace, top)
		if nestc < 1 {
			nestc = 1
		}
		for i := 1; i < nestc+1; i++ {
			fmt.Println(i)
			op.GeoM.Reset()
			op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
			op.GeoM.Translate(pos.X, pos.Y+(10*gsys.SizeMult*float64(i)))
			ebitenScr.DrawImage(BracketBlockImg, op)
		}
		op.GeoM.Reset()
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(pos.X, pos.Y+(10*gsys.SizeMult*float64(nestc+1)))
		ebitenScr.DrawImage(BracketBlockEndImg, op)
		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, pos.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), pos.Y+(10*gsys.SizeMult*float64(nestc+1)), forinflen)
	case WalkBlock:
		gsys.DrawBlueBlock(ebitenScr, gridlocale.Walk, Vec2F{pos.X, pos.Y}, 50)
	case TurnRightBlock:
	case TurnLeftBlock:
	case FlipBlock:
	}
}

func (gsys *Gridsys) DrawBlockCluster(ebitenScr *ebiten.Image, startcodeblock int64, startpos Vec2F) {
	gsys.DrawBlock(ebitenScr, gsys.Blocks[startcodeblock], startpos, 2)
	if gsys.Blocks[startcodeblock].Lower != nil {
		startpos.Y += 10 * gsys.SizeMult
		gsys.DrawBlockCluster(ebitenScr, *gsys.Blocks[startcodeblock].Lower, startpos)
	}
}

func (gsys *Gridsys) DrawAllBlocks(ebitenScr *ebiten.Image, pos Vec2F) {
	for _, id := range gsys.HeadBlocks {
		gsys.DrawBlockCluster(ebitenScr, id, pos)
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
	// draw block
	gsys.DrawAllBlocks(ebitenScr, Vec2F{10 * gsys.SizeMult * gsys.tX, 10 * gsys.SizeMult * gsys.tY})

	// palette
	separate := 7
	edratio := 4

	vector.DrawFilledRect(ebitenScr, float32(wx/separate*edratio), 0, float32(wx/separate*separate-edratio), float32(wy), color.RGBA{150, 150, 150, 255}, false)
}
