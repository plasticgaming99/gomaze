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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	BracketBlockEnd2Img           *ebiten.Image
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

	brBE2 := bytes.NewReader(gridassets.BracketEnd2)
	BracketBlockEnd2Img, _, err = ebitenutil.NewImageFromReader(brBE2)
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

// block position or some
type Vec2 struct {
	X int
	Y int
}

type CodeBlock struct {
	Kind BlockKind

	Pos *Vec2 // tile
	Vec *Vec2 // for dragging

	Length int // pls set

	// stroke
	dragged bool

	// for if and some
	Boolean *BooleanKind
	// if implements bool, YOU MUST SET
	BoolStart int
	BoolEnd   int
}

type Gridsys struct {
	SizeMult float64 // multiplier
	tX       float64 // translate X
	tY       float64 // translate Y

	Blocks map[Vec2]*CodeBlock

	strokes map[*Stroke]struct{}

	displayFont *text.GoTextFace
}

// init new gridsys with default val
func New() *Gridsys {
	return &Gridsys{
		SizeMult: 10,
		tX:       0,
		tY:       0,
		Blocks:   make(map[Vec2]*CodeBlock),
		strokes:  make(map[*Stroke]struct{}),
	}
}

type strokemouse struct{}

func (strkm *strokemouse) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (strkm *strokemouse) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

type Stroke struct {
	source strokemouse

	// offsetX and offsetY represents a relative value from the sprite's upper-left position to the cursor position.
	offsetX int
	offsetY int

	// sprite represents a sprite being dragged.
	codeblock *CodeBlock
}

func NewStroke(source strokemouse, codeblock *CodeBlock, blockpos *Vec2, gsys *Gridsys) *Stroke {
	codeblock.dragged = true
	x, y := source.Position()

	// ビューポートのオフセットを考慮して Vec を基にオフセットを計算
	offsetX := x - (codeblock.Vec.X + int(10*gsys.SizeMult*gsys.tX))
	offsetY := y - (codeblock.Vec.Y + int(10*gsys.SizeMult*gsys.tY))

	return &Stroke{
		source:    source,
		offsetX:   offsetX,
		offsetY:   offsetY,
		codeblock: codeblock,
	}
}

func (block *CodeBlock) In(x, y int, gsys *Gridsys) bool {
	if block.Pos == nil {
		return false
	}

	// ヒットボックスを拡大（例: 1.5倍）
	scale := 1.
	blockWidth := int(10 * gsys.SizeMult * 8 * scale)
	blockHeight := int(10 * gsys.SizeMult * scale)
	gridX := int(float64(block.Pos.X)*10*gsys.SizeMult + 10*gsys.SizeMult*gsys.tX - float64(blockWidth-int(10*gsys.SizeMult*8))/2)
	gridY := int(float64(block.Pos.Y)*10*gsys.SizeMult + 10*gsys.SizeMult*gsys.tY - float64(blockHeight-int(10*gsys.SizeMult))/2)

	return x >= gridX && x <= gridX+blockWidth && y >= gridY && y <= gridY+blockHeight
}

func (gsys *Gridsys) BlockAt(x, y int) *CodeBlock {
	// As the sprites are ordered from back to front,
	// search the clicked/touched sprite in reverse order.
	for _, cb := range gsys.Blocks {
		if cb.In(x, y, gsys) {
			return cb
		}
	}
	/*for i := len() - 1; i >= 0; i-- {
	s := g.gsprites[i]
	if s.In(x, y) {
		return s
	}
	}*/
	return nil
}

func (cb *CodeBlock) MoveTo(x, y int, sizemult float64, screensiz Vec2, gsys *Gridsys) {
	// ビューポートのオフセットを考慮
	offsetX := int(10 * sizemult * gsys.tX)
	offsetY := int(10 * sizemult * gsys.tY)
	cb.Vec.X = x - offsetX*2
	cb.Vec.Y = y - offsetY*2

	// 画面外に出ないように制限
	/*w, h := int(10*sizemult*8), int(10*sizemult)
	if cb.Vec.X < 0 {
		cb.Vec.X = 0
	}
	if cb.Vec.X > screensiz.X-w {
		cb.Vec.X = screensiz.X - w
	}
	if cb.Vec.Y < 0 {
		cb.Vec.Y = 0
	}
	if cb.Vec.Y > screensiz.Y-h {
		cb.Vec.Y = screensiz.Y - h
	}*/
}

func (cb *CodeBlock) SnapToGrid(sizemult float64, gsys *Gridsys) {
	if cb.Vec == nil || cb.Pos == nil {
		return
	}
	oldPos := *cb.Pos

	gridSize := int(10 * sizemult)
	// ビューポートのオフセットを考慮
	offsetX := /*int(10 * sizemult * gsys.tX)*/ 0
	offsetY := /*int(10 * sizemult * gsys.tY)*/ 0

	// カメラのズレを考慮してグリッド位置を計算
	gridX := int(math.Round(float64(cb.Vec.X-offsetX) / float64(gridSize)))
	gridY := int(math.Round(float64(cb.Vec.Y-offsetY) / float64(gridSize)))
	newPos := Vec2{X: gridX, Y: gridY}

	// 重複チェック
	if _, exists := gsys.Blocks[newPos]; exists {
		// 重複がある場合、元の位置に戻す
		cb.Vec.X = oldPos.X*gridSize - offsetX
		cb.Vec.Y = oldPos.Y*gridSize - offsetY
		return
	}

	// 重複がない場合、位置を更新
	cb.Pos = &newPos
	cb.Vec.X = gridX*gridSize - offsetX
	cb.Vec.Y = gridY*gridSize - offsetY

	// マップのキーを更新
	if gsys != nil {
		delete(gsys.Blocks, oldPos)
		gsys.Blocks[newPos] = cb
	}
}

func (s *Stroke) Update(gsys *Gridsys) {
	if !s.codeblock.dragged {
		return
	}

	if s.source.IsJustReleased() {
		s.codeblock.dragged = false
		s.codeblock.SnapToGrid(gsys.SizeMult, gsys)
		return
	}

	x, y := s.source.Position()
	x -= s.offsetX
	y -= s.offsetY

	// ビューポートのオフセットを考慮
	offsetX := int(10 * gsys.SizeMult * gsys.tX)
	offsetY := int(10 * gsys.SizeMult * gsys.tY)
	x += offsetX
	y += offsetY

	s.codeblock.MoveTo(x, y, gsys.SizeMult, Vec2{2000, 2000}, gsys)
}

func (gsys *Gridsys) Tick() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if sp := gsys.BlockAt(ebiten.CursorPosition()); sp != nil {
			s := NewStroke(strokemouse{}, sp, sp.Pos, gsys)
			gsys.strokes[s] = struct{}{}
		}
	}

	for s := range gsys.strokes {
		s.Update(gsys)
		if !s.codeblock.dragged {
			delete(gsys.strokes, s)
		}
	}

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
	grid := 10 * int(gsys.SizeMult)
	gsys.Blocks[Vec2{1, 1}] = &CodeBlock{
		Kind: StartBlock,
		Pos:  &Vec2{X: 1, Y: 1},
		Vec:  &Vec2{X: 1 * grid, Y: 1 * grid}, // ピクセル座標
	}
	gsys.Blocks[Vec2{1, 2}] = &CodeBlock{
		Kind: ForInfBlock,
		Pos:  &Vec2{X: 1, Y: 2},
		Vec:  &Vec2{X: 1 * grid, Y: 2 * grid},
	}
	gsys.Blocks[Vec2{2, 3}] = &CodeBlock{
		Kind: WalkBlock,
		Pos:  &Vec2{X: 2, Y: 3},
		Vec:  &Vec2{X: 2 * grid, Y: 3 * grid},
	}
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
	ebitenScr.DrawImage(BlockBlueImg, op)
	op.GeoM.Reset()
	op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
	gsys.DrawBlockPart(ebitenScr, color.RGBA{86, 147, 255, 255}, color.RGBA{61, 105, 181, 255}, pos.X+(gsys.SizeMult*float64(StartBlockImg.Bounds().Dx())), pos.Y, length)
	top.GeoM.Translate(pos.X+(gsys.SizeMult*10), pos.Y+gsys.SizeMult)
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

	var tempVec Vec2F
	if codeblock.dragged && codeblock.Vec != nil {
		// ドラッグ中はVec（ピクセル座標）＋ビューポートオフセット
		tempVec.X = float64(codeblock.Vec.X) + 10*gsys.SizeMult*gsys.tX
		tempVec.Y = float64(codeblock.Vec.Y) + 10*gsys.SizeMult*gsys.tY
	} else if codeblock.Pos != nil {
		// 通常時はグリッド座標＋ビューポートオフセット
		tempVec.X = 10*gsys.SizeMult*float64(codeblock.Pos.X) + 10*gsys.SizeMult*gsys.tX
		tempVec.Y = 10*gsys.SizeMult*float64(codeblock.Pos.Y) + 10*gsys.SizeMult*gsys.tY
	} else {
		tempVec = pos // デフォルトの位置
	}

	switch codeblock.Kind {
	case StartBlock:
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X, tempVec.Y)
		ebitenScr.DrawImage(StartBlockImg, op)
		op.GeoM.Reset()

		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X, tempVec.Y-(3*gsys.SizeMult))
		ebitenScr.DrawImage(StartBlockCapImg, op)

		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, tempVec.X+(gsys.SizeMult*float64(StartBlockImg.Bounds().Dx())), tempVec.Y, 100)
		top.GeoM.Translate(tempVec.X+gsys.SizeMult, tempVec.Y+gsys.SizeMult)
		text.Draw(ebitenScr, gridlocale.Start, misakiGothic2ndFace, top)
	case IfBlock:
		forinflen := gsys.SizeMult * 10 * 3
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X, tempVec.Y)
		ebitenScr.DrawImage(BracketBlockImg, op)
		gsys.DrawBlockPart(ebitenScr, color.RGBA{255, 221, 0, 255}, color.RGBA{224, 192, 0, 255}, tempVec.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), tempVec.Y, forinflen)
		op.GeoM.Reset()

		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())), tempVec.Y)
		ebitenScr.DrawImage(StartBlockImg, op)
		op.GeoM.Reset()
		top.GeoM.Translate(tempVec.X+(gsys.SizeMult*10), tempVec.Y+gsys.SizeMult)
		text.Draw(ebitenScr, gridlocale.ForInf, misakiGothic2ndFace, top)
		if nestc < 1 {
			nestc = 1
		}
		for i := 1; i < nestc+1; i++ {
			fmt.Println(i)
			op.GeoM.Reset()
			op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
			op.GeoM.Translate(tempVec.X, tempVec.Y+(10*gsys.SizeMult*float64(i)))
			ebitenScr.DrawImage(BracketBlockImg, op)
		}
		op.GeoM.Reset()
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X, tempVec.Y+(10*gsys.SizeMult*float64(nestc+1)))
		ebitenScr.DrawImage(BracketBlockEndImg, op)
		op.GeoM.Reset()
		op.GeoM.Translate(tempVec.X+float64(BracketBlockEndImg.Bounds().Dx()), tempVec.Y+(10*gsys.SizeMult*float64(nestc+1)))
		ebitenScr.DrawImage(BracketBlockEnd2Img, op)
		gsys.DrawBlockPart(
			ebitenScr,
			color.RGBA{255, 221, 0, 255},
			color.RGBA{224, 192, 0, 255},
			tempVec.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx()+1)),
			tempVec.Y+(10*gsys.SizeMult*float64(nestc+1)),
			forinflen,
		)
	case ForInfBlock:
		forinflen := gsys.SizeMult * 10 * 3
		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X, tempVec.Y)
		ebitenScr.DrawImage(BracketBlockImg, op)
		gsys.DrawBlockPart(
			ebitenScr,
			color.RGBA{255, 221, 0, 255},
			color.RGBA{224, 192, 0, 255},
			tempVec.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx())),
			tempVec.Y,
			forinflen,
		)
		op.GeoM.Reset()

		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx()+1)), tempVec.Y)
		ebitenScr.DrawImage(StartBlockImg, op)
		op.GeoM.Reset()

		top.GeoM.Translate(tempVec.X+(gsys.SizeMult*10), tempVec.Y+gsys.SizeMult)
		text.Draw(ebitenScr, gridlocale.ForInf, misakiGothic2ndFace, top)
		if nestc < 1 {
			nestc = 1
		}
		for i := 1; i < nestc+1; i++ {
			fmt.Println(i)
			op.GeoM.Reset()
			op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
			op.GeoM.Translate(tempVec.X, tempVec.Y+(10*gsys.SizeMult*float64(i)))
			ebitenScr.DrawImage(BracketBlockImg, op)
		}
		op.GeoM.Reset()

		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X, tempVec.Y+(10*gsys.SizeMult*float64(nestc+1)))
		ebitenScr.DrawImage(BracketBlockEndImg, op)
		op.GeoM.Reset()

		op.GeoM.Scale(gsys.SizeMult, gsys.SizeMult)
		op.GeoM.Translate(tempVec.X+gsys.SizeMult*float64(BracketBlockEndImg.Bounds().Dx()+1), tempVec.Y+(10*gsys.SizeMult*float64(nestc+1)))
		ebitenScr.DrawImage(BracketBlockEnd2Img, op)
		gsys.DrawBlockPart(
			ebitenScr,
			color.RGBA{255, 221, 0, 255},
			color.RGBA{224, 192, 0, 255},
			tempVec.X+(gsys.SizeMult*float64(BracketBlockEndImg.Bounds().Dx())),
			tempVec.Y+(10*gsys.SizeMult*float64(nestc+1)),
			1*gsys.SizeMult,
		)
		gsys.DrawBlockPart(
			ebitenScr,
			color.RGBA{255, 221, 0, 255},
			color.RGBA{224, 192, 0, 255},
			tempVec.X+(gsys.SizeMult*float64(BracketBlockImg.Bounds().Dx()*2+1)),
			tempVec.Y+(10*gsys.SizeMult*float64(nestc+1)),
			forinflen-(10*gsys.SizeMult),
		)
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
		gsys.DrawBlueBlock(ebitenScr, gridlocale.Walk, Vec2F{tempVec.X, tempVec.Y}, 24*gsys.SizeMult)
	case TurnRightBlock:
		gsys.DrawBlueBlock(ebitenScr, gridlocale.TurnRight, Vec2F{tempVec.X, tempVec.Y}, 36*gsys.SizeMult)
	case TurnLeftBlock:
		gsys.DrawBlueBlock(ebitenScr, gridlocale.TurnLeft, Vec2F{tempVec.X, tempVec.Y}, 36*gsys.SizeMult)
	case FlipBlock:
		gsys.DrawBlueBlock(ebitenScr, gridlocale.Flip, Vec2F{tempVec.X, tempVec.Y}, 30*gsys.SizeMult)
	}
}

func (gsys *Gridsys) DrawAllBlocks(ebitenScr *ebiten.Image, pos Vec2F) {
	for v2, some := range gsys.Blocks {
		// v2 to correct v2f
		v2f := Vec2F{
			X: 10*gsys.SizeMult*float64(v2.X) + 10*gsys.SizeMult*gsys.tX,
			Y: 10*gsys.SizeMult*float64(v2.Y) + 10*gsys.SizeMult*gsys.tY,
		}
		gsys.DrawBlock(ebitenScr, some, v2f, 1)
	}
}

func (gsys *Gridsys) Draw(ebitenScr *ebiten.Image, pos Vec2, size Vec2) {
	fmt.Println(gsys.tX, gsys.tY)
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
	// draw block
	gsys.DrawAllBlocks(ebitenScr, Vec2F{10 * gsys.SizeMult * gsys.tX, 10 * gsys.SizeMult * gsys.tY})

	// palette
	separate := 7
	edratio := 4

	vector.DrawFilledRect(ebitenScr, float32(wx/separate*edratio), 0, float32(wx/separate*separate-edratio), float32(wy), color.RGBA{150, 150, 150, 255}, false)

	//cursor
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
}
