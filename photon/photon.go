package photon

/*
 * PhotonText very alpha (codename crunchy) by plasticgaming99
 * (c)opyright plasticgaming99, 2023-2024
 */

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/plasticgaming99/gomaze/_lib/dyntypes"
	"github.com/plasticgaming99/gomaze/_lib/plastk"
	"github.com/plasticgaming99/gomaze/assets/phfonts"
	"github.com/plasticgaming99/gomaze/assets/phicons"

	"golang.design/x/clipboard"
)

/* basic */
var (
	screenWidth      = 640
	screenHeight     = 480
	mplusNormalFont  font.Face
	mplusBigFont     font.Face
	mplusSmallFont   font.Face
	HackGenFont      font.Face
	smallHackGenFont font.Face
	smallHackGenV2   *textv2.GoTextFaceSource

	PhotonText = []string{}
	undohist   = []undoPreserver{}

	photoncmd   = string("")
	cmdresult   = string("")
	clearresult = int(0)
	rellines    = int(0)

	textrepeatness = int(0)

	cursornowx    = int(1)
	cursornowy    = int(1)
	cursorxeffort = int(0)

	closewindow   = false
	clickrepeated = false
	modalmode     = false
	returncode    = string("\n")
	returntype    = string("")

	// options
	hanzenlock      = true
	hanzenlockstat  = false
	limitterenabled = true
	limitterlevel   = int(3)
	dbgmode         = false
	editmode        = int(1)

	editorforcused      = true
	commandlineforcused = false

	editingfile = string("")

	// Textures
	sideBar         *ebiten.Image
	infoBar         *ebiten.Image
	commandLine     *ebiten.Image
	cursorimg       = ebiten.NewImage(2, 15)
	topopbar        *ebiten.Image
	filesmenubutton = ebiten.NewImage(80, 20)
	linessep        *ebiten.Image
	scrollbar       *ebiten.Image
	scrollbit       *ebiten.Image

	// Texture options
	sideBarop = &ebiten.DrawImageOptions{}
)

/* Texture options */
var (
	topopBarSize    = 20
	infoBarSize     = 20
	commandlineSize = 20
	scrollbarwidth  = 18
)

var (
	synctps  = &sync.WaitGroup{}
	synctps2 = &sync.WaitGroup{}
)

type undoPreserver struct {
	textlet string
	cursorx int
	cursory int
}

func saveUndoAppend(in string) undoPreserver {
	toappend := undoPreserver{
		textlet: in,
		cursorx: cursornowx - 1,
		cursory: cursornowy - 1,
	}
	return toappend
}

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

func checkMixedKanjiLength(kantext string, length int) (int, int, int) {
	kantext = string([]rune(kantext)[0 : length-1])
	kanji := (len(kantext) - len([]rune(kantext))) / 2
	nonkanji := len([]rune(kantext)) - kanji
	tab := strings.Count(kantext, "	")
	return nonkanji, kanji - tab, tab
}

// renew image. if same size, return same image
func renewimg(image *ebiten.Image, targetwidth int, targetheight int, targetcolor color.Color) *ebiten.Image {
	widthnow := image.Bounds().Dx()
	heightnow := image.Bounds().Dy()
	if widthnow != targetwidth || heightnow != targetheight {
		newimage := ebiten.NewImage(targetwidth, targetheight)
		newimage.Fill(targetcolor)
		return newimage
	}
	return image
}

// func
func init() {
	const dpi = 144

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		var iconphoton []image.Image
		ebiten.SetVsyncEnabled(true)
		ebiten.SetScreenClearedEveryFrame(false)

		iconphotonreader := bytes.NewReader(phicons.PhotonIcon16)
		_, iconphoton16, err := ebitenutil.NewImageFromReader(iconphotonreader)
		phloginfo(err)
		iconphotonreader = bytes.NewReader(phicons.PhotonIcon32)
		_, iconphoton32, err := ebitenutil.NewImageFromReader(iconphotonreader)
		phloginfo(err)
		iconphotonreader = bytes.NewReader(phicons.PhotonIcon48)
		_, iconphoton48, err := ebitenutil.NewImageFromReader(iconphotonreader)
		phloginfo(err)
		iconphotonreader = bytes.NewReader(phicons.PhotonIcon128)
		_, iconphoton128, err := ebitenutil.NewImageFromReader(iconphotonreader)
		phloginfo(err)
		iconphoton = append(iconphoton, iconphoton16, iconphoton32, iconphoton48, iconphoton128)

		ebiten.SetWindowIcon(iconphoton)

		/*100, 250, 500, 750, 1000 or your monitor's refresh rate*/
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

		err = clipboard.Init()
		if err != nil {
			fmt.Println("**WARN** Clipboard is disabled.", err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		tt, err := opentype.Parse(phfonts.MPlus1pRegular_ttf)
		if err != nil {
			log.Fatal(err)
		}

		mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    12,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
		mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    8,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
		mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    24,
			DPI:     dpi,
			Hinting: font.HintingFull, // Use quantization to save glyph cache images.
		})
		if err != nil {
			log.Fatal(err)
		}

		// Adjust the line height.
		mplusBigFont = text.FaceWithLineHeight(mplusBigFont, 54)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		tt, err := opentype.Parse(phfonts.HackGenRegular_ttf)
		if err != nil {
			log.Fatal(err)
		}

		HackGenFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    12,
			DPI:     dpi,
			Hinting: font.HintingFull,
		})
		if err != nil {
			log.Fatal(err)
		}
		smallHackGenFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    8,
			DPI:     dpi,
			Hinting: font.HintingFull,
		})
		smallHackGenV2, err = textv2.NewGoTextFaceSource(bytes.NewReader(phfonts.HackGenRegular_ttf))
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	// load file
	wg.Add(1)
	go func() {
		if len(os.Args) >= 2 {
			phload(os.Args[1])
		}
		wg.Done()
	}()

	wg.Wait()

	// after loaded text to memory, if photontext
	// has not any strings, init photontext with
	// 1-line, 0-column text.
	wg.Add(1)
	go func() {
		if len(PhotonText) == 0 {
			PhotonText = append(PhotonText, "")
		}
		wg.Done()
	}()
	wg.Wait()

	// Execute PhotonRC when its avaliable
	{
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
		}
		phrcpath := home + "/.photonrc"

		if err != nil {
			panic(err)
		}
		_, err = os.Stat(phrcpath)
		if err == nil {
			fmt.Println("Using PhotonRC")
			photonRC, err := sliceload(phrcpath)
			if err != nil {
				fmt.Println("PhotonRC initializing failed. Using default")
			}
			for i := 0; i < len(photonRC); i++ {
				proceedcmd(photonRC[i])
			}
			photonRC = nil
		}
	}

	{
		// After executed PhotonRC, Initialize textures
		sideBar = ebiten.NewImage(60, 100)
		infoBar = ebiten.NewImage(100, infoBarSize)
		commandLine = ebiten.NewImage(100, commandlineSize)
		cursorimg = ebiten.NewImage(2, 15)
		topopbar = ebiten.NewImage(4100, topopBarSize)
		filesmenubutton = ebiten.NewImage(80, 20)
		linessep = ebiten.NewImage(2, 100)
		scrollbar = ebiten.NewImage(scrollbarwidth, 100)
		scrollbit = ebiten.NewImage(scrollbarwidth, 100)
	}

	// Fill textures
	{
		/* init sidebar image. */
		sideBar.Fill(color.RGBA{57, 57, 57, 255})
		/* init information bar image */
		infoBar.Fill(color.RGBA{87, 97, 87, 255})
		/* init commandline image */
		commandLine.Fill(color.RGBA{39, 39, 39, 255})
		/* init cursor image */
		cursorimg.Fill(color.RGBA{255, 255, 255, 5})
		/* init top-op-bar image */
		topopbar.Fill(color.RGBA{100, 100, 100, 255})
		/* init top-op-bar "files" button */
		filesmenubutton.Fill(color.RGBA{110, 110, 110, 255})
		/* init line-bar separator */
		linessep.Fill(color.RGBA{100, 100, 100, 255})
		/* init scroll-bar */
		scrollbar.Fill(color.RGBA{80, 80, 80, 255})
		/* init scroll-bit */
		scrollbit.Fill(color.RGBA{30, 30, 30, 255})

		// Init texture options
		sideBarop.GeoM.Translate(float64(0), float64(20))
	}
}

type Editor struct {
	/*counter        int
	kanjiText      string
	kanjiTextColor color.RGBA*/
	rune2Input []rune
}

func checkcurx(line int) {
	if len([]rune(PhotonText[line-1])) < cursornowx {
		if PhotonText[line-1] == "" {
			cursornowx = 1
		} else {
			cursornowx = len([]rune(PhotonText[line-1])) + 1
		}
	}
}

func (g *Editor) Update() error {
	synctps.Add(1)
	if plastk.MenuBarDetectClickedByID("fileswritebutton") {
		proceedcmd("write")
	}
	if plastk.MenuBarDetectClickedByID("filesexitbutton") {
		proceedcmd("q")
	}
	if plastk.MenuBarDetectClickedByID("editredobutton") {
		proceedcmd("q")
	}
	// Update Text-info
	// photonlines = len(photontext)
	/*\
	 * detect cursor key actions
	\*/
	// Insert text
	if editorforcused && !(ebiten.IsKeyPressed(ebiten.KeyControl)) {
		g.rune2Input = ebiten.AppendInputChars(g.rune2Input[:0])
		// Detect left side
		if cursornowx == 1 {
			PhotonText[cursornowy-1] = string(g.rune2Input) + PhotonText[cursornowy-1]
		} else
		// Detect right side
		if cursornowx-1 == len([]rune(PhotonText[cursornowy-1])) {
			PhotonText[cursornowy-1] = PhotonText[cursornowy-1] + string(g.rune2Input)
		} else
		// Other, Insert
		{
			PhotonText[cursornowy-1] = string([]rune(PhotonText[cursornowy-1])[:cursornowx-1]) + string(g.rune2Input) + string([]rune(PhotonText[cursornowy-1])[cursornowx-1:])
		}
		// Move cursornowx. with cjk support yay!
		cursornowx += len(g.rune2Input)
	}

	if editorforcused {
		// Check commandline is called
		if (ebiten.IsKeyPressed(ebiten.KeyControl)) && (ebiten.IsKeyPressed(ebiten.KeyShift)) && (ebiten.IsKeyPressed(ebiten.KeyC)) {
			editorforcused = false
			commandlineforcused = true
		} else
		// Check upper text.
		if (repeatingKeyPressed(ebiten.KeyUp)) && (cursornowy > 1) {
			checkcurx(cursornowy - 1)
			cursornowy--
		} else
		// Check lower text.
		if (repeatingKeyPressed(ebiten.KeyDown)) && (cursornowy < len(PhotonText)) {
			checkcurx(cursornowy + 1)
			cursornowy++
		} else if (repeatingKeyPressed(ebiten.KeyLeft)) && (cursornowx > 1) {
			cursornowx--
		} else if (repeatingKeyPressed(ebiten.KeyRight)) && (cursornowx <= len([]rune(PhotonText[cursornowy-1]))) {
			cursornowx++
		} else if (ebiten.IsKeyPressed(ebiten.KeyControl)) && (repeatingKeyPressed(ebiten.KeyC)) {
			fmt.Println("c pressed")
		} else if (repeatingKeyPressed(ebiten.KeyBackquote)) && (hanzenlock) {
			if !hanzenlockstat {
				hanzenlockstat = true
			} else {
				hanzenlockstat = false
			}
		} else if repeatingKeyPressed(ebiten.KeyHome) {
			cursornowx = 1
		} else if repeatingKeyPressed(ebiten.KeyEnd) {
			cursornowx = len([]rune(PhotonText[cursornowy-1])) + 1
		} else if ebiten.IsKeyPressed(ebiten.KeyControl) && repeatingKeyPressed(ebiten.KeyV) {
			testslice := strings.Split(string(clipboard.Read(clipboard.FmtText)), "\n")
			firsttext := string([]rune(PhotonText[cursornowy-1])[:cursornowx-1])
			lasttext := string([]rune(PhotonText[cursornowy-1])[cursornowx-1:])
			//{photontext[cursornowy-1] = string(g.rune2Input)}

			fmt.Println(testslice)
			if len(testslice) == 1 {
				PhotonText[cursornowy-1] = firsttext + testslice[0] + lasttext
				cursornowx = len([]rune(firsttext + testslice[0]))
				fmt.Println("one")
			} else {
				for i := 0; i < len(testslice); i++ {
					if i == 0 {
						PhotonText[cursornowy-1] = firsttext + testslice[i]
					} else {
						{
							PhotonText = append(PhotonText[:cursornowy], append([]string{testslice[i]}, PhotonText[cursornowy:]...)...)
							/*PhotonText[cursornowy] = string([]rune(PhotonText[cursornowy-1])[cursornowx-1:])
							PhotonText[cursornowy-1] = string([]rune(PhotonText[cursornowy-1])[:cursornowx-1])*/
							cursornowy++
						}
						if i == len(testslice)-1 {
							cursornowx = len([]rune(testslice[i])) + 1
						}
					}
					fmt.Println(i)
				}
			}
		} else if repeatingKeyPressed(ebiten.KeyTab) {
			/*PhotonText[cursornowy-1] = PhotonText[cursornowy-1] + string(g.rune2Input) (legacy impl) */
			// Detect text input
			// Detect left side
			if cursornowx == 1 {
				PhotonText[cursornowy-1] = string("	") + PhotonText[cursornowy-1]
			} else
			// Detect right side
			if cursornowx-1 == len([]rune(PhotonText[cursornowy-1])) {
				PhotonText[cursornowy-1] = PhotonText[cursornowy-1] + string("	")
			} else
			// Other, Insert
			{
				PhotonText[cursornowy-1] = string([]rune(PhotonText[cursornowy-1])[:cursornowx-1]) + string("	") + string([]rune(PhotonText[cursornowy-1])[cursornowx-1:])
			}
			// Move cursornowx. with cjk support yay!
			cursornowx += len("	")
		} else
		// New line
		if (repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter)) && !hanzenlockstat {
			{
				PhotonText = append(PhotonText[:cursornowy], append([]string{""}, PhotonText[cursornowy:]...)...)
				PhotonText[cursornowy] = string([]rune(PhotonText[cursornowy-1])[cursornowx-1:])
				PhotonText[cursornowy-1] = string([]rune(PhotonText[cursornowy-1])[:cursornowx-1])
				cursornowy++
				cursornowx = 1
			}
			cursornowx = 1
		} else
		// Line deletion.
		if repeatingKeyPressed(ebiten.KeyBackspace) && !((len(PhotonText[0]) == 0) && (cursornowy == 1)) && !hanzenlockstat {
			if (PhotonText[cursornowy-1] == "") && (len(PhotonText) != 1) {
				cursornowx = len([]rune(PhotonText[cursornowy-2])) + 1
				if cursornowy-1 < len(PhotonText)-1 {
					copy(PhotonText[cursornowy-1:], PhotonText[cursornowy:])
				}
				PhotonText[len(PhotonText)-1] = ""
				PhotonText = PhotonText[:len(PhotonText)-1]
				cursornowx = len([]rune(PhotonText[cursornowy-2])) + 1
				cursornowy--
			} else {
				if !((cursornowx == 1) && (cursornowy == 1)) || (cursornowx-1 == len([]rune(PhotonText[cursornowy-1]))) {
					if cursornowx == 1 {
						cursornowx = len([]rune(PhotonText[cursornowy-2])) + 1
						PhotonText[cursornowy-2] = PhotonText[cursornowy-2] + PhotonText[cursornowy-1]
						if cursornowy-1 < len(PhotonText)-1 {
							copy(PhotonText[cursornowy-1:], PhotonText[cursornowy:])
						}
						PhotonText[len(PhotonText)-1] = ""
						PhotonText = PhotonText[:len(PhotonText)-1]
						cursornowy--
					} else
					// normal deletion
					if cursornowx-1 == len([]rune(PhotonText[cursornowy-1])) {
						// convert 2 rune
						runes := []rune(PhotonText[cursornowy-1])
						//save last char
						////undohist = append(undohist, saveUndoAppend(string(runes[len(runes)-1:])))
						////fmt.Println(undohist)
						// delete last char
						runes = runes[:len(runes)-1]
						// convert rune 2 string and insert
						PhotonText[cursornowy-1] = string(runes)
						// Move to left
						cursornowx--
					} else {
						// Convert to rune
						runes := []rune(PhotonText[cursornowy-1])[:cursornowx-1]
						//save last char
						////undohist = append(undohist, saveUndoAppend(string(runes[len(runes)-1:])))
						// Delete last
						runes = runes[:len(runes)-1]
						// Convert to string and insert
						PhotonText[cursornowy-1] = string(runes) + string([]rune(PhotonText[cursornowy-1])[cursornowx-1:])
						// Move to left
						cursornowx--
					}
				}
			}
		}
	} else
	// If command-line is forcused
	if commandlineforcused {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			cmdresult = proceedcmd(photoncmd)
			clearresult += 10
			photoncmd = ""
			editorforcused = true
			commandlineforcused = false
		}
		if (len([]rune(photoncmd)) >= 1) && (repeatingKeyPressed(ebiten.KeyBackspace)) {
			cmdrune := []rune(photoncmd)[:len([]rune(photoncmd))-1]
			photoncmd = string(cmdrune)
		} else {
			// detect text input
			g.rune2Input = ebiten.AppendInputChars(g.rune2Input[:0])

			// insert text
			if string(g.rune2Input) != "" {
				photoncmd += string(g.rune2Input)
			}
		}
	}

	/*\
	 * detect mouse wheel actions.
	\*/
	_, dy := ebiten.Wheel()
	if (dy > 0) && (rellines > 0) {
		rellines -= 3
	} else if (dy < 0) && (rellines+3 < len(PhotonText)) {
		rellines += 3
	}

	/*\
	 * Detect touch on buttons.
	\*/
	/*mousex, mousey := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		clickrepeated = true
	}
	if clickrepeated && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if mousex < 80 && mousey < 20 {
			fmt.Println("test!!!")
		}
		clickrepeated = false
	}*/

	/*\
	 * Detect cursor's position, and changes cursor shape
	\*/
	curx, cury := ebiten.CursorPosition()
	windowx, windowy := ebiten.WindowSize()
	if ((60 < curx) && (curx < windowx)) && ((20 < cury) && (cury < windowy-20)) {
		ebiten.SetCursorShape(ebiten.CursorShapeText)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	if repeatingKeyPressed(ebiten.KeyA) {
		fmt.Println("a pressed")
	}

	if closewindow {
		return fmt.Errorf("close")
	}
	synctps.Done()

	return nil
}

var (
	prevcurx = int(0)
	prevcury = int(0)
	prevrell = int(0)
	phburst  = int(0)
)

var (
	textx        int
	cursorxstart int
)

var hackgen4info *textv2.GoTextFace

func init() {
	hackgen4info = &textv2.GoTextFace{
		Source:   smallHackGenV2,
		Size:     16,
		Language: language.Japanese,
	}
}

func (g *Editor) Draw(screen *ebiten.Image) {
	//synctps.Wait()

	// frame limitter, bad
	if (prevcurx == cursornowx && prevcury == cursornowy && prevrell == rellines) && limitterenabled {
		time.Sleep(time.Duration(phburst) * time.Millisecond)
		if phburst < limitterlevel {
			phburst += 1
		}
	} else {
		if 0 <= phburst {
			phburst = 0
		}
	}

	prevcurx, prevcury, prevrell = cursornowx, cursornowy, rellines

	screenWidth, screenHeight := ebiten.WindowSize()

	screenHeight -= topopBarSize

	if commandlineforcused || cmdresult != "" {
		screenHeight -= 20
	}

	// Init screen
	screen.Fill(color.RGBA{61, 61, 61, 255})

	sideBar = renewimg(sideBar, 60, screenWidth, color.RGBA{57, 57, 57, 255})
	screen.DrawImage(sideBar, sideBarop)

	// Draw left information text
	leftinfotxt := ""
	leftinfotxt = "PhotonText alpha "
	if hanzenlockstat {
		leftinfotxt += "Hanzenlock "
	}

	// Draw right information text
	rightinfotext := " " + returntype + " " + strconv.Itoa(cursornowy) + ":" + strconv.Itoa(cursornowx)

	// draw editor text
	Maxtext := int(math.Ceil(((float64(screenHeight) - 20) / 18)) - 1)
	if int(Maxtext) >= len(PhotonText) {
		textrepeatness = len(PhotonText) - 1
	} else {
		textrepeatness = int(Maxtext) - 1
	}

	// start line loop
	for printext := 0; printext < len(PhotonText[rellines:]); {
		if printext > int(Maxtext) || (len(PhotonText)-rellines) == 0 {
			break
		}
		//slicedtext := []rune(PhotonText[printext+rellines])
		textx = 55
		text.Draw(screen, strconv.Itoa(printext+rellines+1), smallHackGenFont, textx+10, ((printext + 2) * 18), color.White)
		textx = textx + (9 * len(strconv.Itoa(int(Maxtext)+rellines)))
		textx += 20
		cursorxstart = textx + 0
		spaceSplit := strings.Split(PhotonText[printext+rellines], " ")
		var samplue []textv2.Glyph

		//textv2.Draw(screen, PhotonText[printext+rellines], hackgen4info, op)

		for i := 0; i < len(spaceSplit); i++ {
			samplue = textv2.AppendGlyphs(nil, tab2space(PhotonText[printext+rellines]), hackgen4info, nil)
		}
		for i, gl := range samplue {
			op := &ebiten.DrawImageOptions{}
			//op.GeoM.Translate(90, ((float64(printext)+1)*18)+5)
			syntstat, endsat := ezSyntaxHighlight(PhotonText[printext+rellines])
			if i >= endsat {
				op.ColorScale.Reset()
			} else if syntstat == "good" {
				op.ColorScale.Scale(0, 0.5, 0, 0)
			}
			if gl.Image == nil {
				continue
			}
			op.GeoM.Translate(float64(textx-1), ((float64(printext)+1)*18)+5)
			op.GeoM.Translate(gl.X, gl.Y)

			screen.DrawImage(gl.Image, op)
		}
		// start column loop
		/*for textrepeat := 0; textrepeat < len(slicedtext); {
			if string("	") == string(slicedtext[textrepeat]) {
				textx += 30
			} else if len(string(slicedtext[textrepeat])) != 1 {
				// If multi-byte text, print bigger
				text.Draw(screen, string(slicedtext[textrepeat]), smallHackGenFont, textx-1, ((printext + 2) * 18), color.White)
				textx += 15
			} else {
				// If not, print normally
				text.Draw(screen, string(slicedtext[textrepeat]), smallHackGenFont, textx, ((printext + 2) * 18), color.White)
				textx += 9
			}
			textrepeat++
		}*/
		printext++
	}

	// draw cursor
	//nonkanj, kanj, tabs := checkMixedKanjiLength(PhotonText[cursornowy-1], cursornowx)
	//cursorproceedx := (nonkanj*9 + kanj*15 + tabs*36) + cursorxstart

	cursorop := &ebiten.DrawImageOptions{}
	//cursorop.GeoM.Translate(float64(cursorproceedx), float64((cursornowy-(rellines))*18)+5)
	txtx, _ := textv2.Measure(PhotonText[cursornowy-1][:cursornowx-1], hackgen4info, 0)
	cursorop.GeoM.Translate(float64(cursorxstart)+txtx, float64((cursornowy-(rellines))*18)+5)
	screen.DrawImage(cursorimg, cursorop)

	// Draw scroll bar base
	scrollbar = renewimg(scrollbar, scrollbarwidth, screenHeight, color.RGBA{80, 80, 80, 255})
	scrollbarop := &ebiten.DrawImageOptions{}
	scrollbarop.GeoM.Translate(float64(screenWidth)-float64(scrollbarwidth), float64(topopBarSize))
	screen.DrawImage(scrollbar, scrollbarop)

	// Draw scroll bit
	/* init scroll-bit */
	var textsize int
	{
		textsize = len(PhotonText) + Maxtext
	}
	scrollbartext := float64(screenHeight-20) / float64((float64(textsize) / float64(Maxtext)))
	if scrollbartext < 1 {
		scrollbartext = 1
	}
	scrollbit = renewimg(scrollbit, scrollbarwidth, int(scrollbartext), color.RGBA{30, 30, 30, 255}) //ebiten.NewImage(25, int(scrollbartext))

	scrollbitop := &ebiten.DrawImageOptions{}
	scrollbitop.GeoM.Translate(float64(screenWidth)-float64(scrollbarwidth), float64((float64(screenHeight-20)/float64(textsize))*float64(rellines)+20))
	screen.DrawImage(scrollbit, scrollbitop)

	// Draw lines separator
	linessepop := &ebiten.DrawImageOptions{}
	linessepop.GeoM.Translate(float64(cursorxstart-5), 0)
	linessep = renewimg(linessep, 2, screen.Bounds().Dy(), color.RGBA{100, 100, 100, 255})
	screen.DrawImage(linessep, linessepop)

	// Draw info-bar
	infoBarop := &ebiten.DrawImageOptions{}
	infoBarop.GeoM.Translate(0, float64(screenHeight))
	infoBar = renewimg(infoBar, screen.Bounds().Dx(), 100, color.RGBA{87, 97, 87, 255})
	screen.DrawImage(infoBar, infoBarop)

	{
		zurasu, padding := textv2.Measure(rightinfotext, hackgen4info, 0)
		infoBarTextY := float64(screenHeight + ((infoBarSize - int(padding)) / 2))
		//text.Draw(screen, leftinfotxt, smallHackGenFont, 5, screenHeight+infoBarSize-4, color.White)
		op := &textv2.DrawOptions{}
		op.GeoM.Translate(2, infoBarTextY)
		textv2.Draw(screen, leftinfotxt, hackgen4info, op)
		op = &textv2.DrawOptions{}
		op.GeoM.Translate(float64(screenWidth-5)-zurasu, infoBarTextY)
		//text.Draw(screen, rightinfotext, smallHackGenFont, screenWidth-((len(rightinfotext))*10), screenHeight+infoBarSize-4, color.White)
		textv2.Draw(screen, rightinfotext, hackgen4info, op)
	}
	{
		//[]string{"Files", "Save"}, []string{"Edit", "Undo"}, []string{"View"}
		var mbfiles []plastk.MenuBarColumn
		{
			filestmp := plastk.MenuBarColumn{
				ColumnType: "dropdown",
				ColumnName: "Files",
			}
			savetmp := plastk.MenuBarColumn{
				ColumnType: "button",
				ColumnName: "Save",
				ColumnID:   "filessavebutton",
			}
			exittmp := plastk.MenuBarColumn{
				ColumnType: "button",
				ColumnName: "Exit",
				ColumnID:   "filesexitbutton",
			}

			mbfiles = append(mbfiles, filestmp, savetmp, exittmp)
		}
		var mbedit []plastk.MenuBarColumn
		{
			edittmp := plastk.MenuBarColumn{
				ColumnType: "dropdown",
				ColumnName: "Edit",
			}
			undotmp := plastk.MenuBarColumn{
				ColumnType: "button",
				ColumnName: "Undo",
				ColumnID:   "editundobutton",
			}
			redotmp := plastk.MenuBarColumn{
				ColumnType: "button",
				ColumnName: "Redo",
				ColumnID:   "editredobutton",
			}
			mbedit = append(mbedit, edittmp, undotmp, redotmp)
		}
		var mbabout []plastk.MenuBarColumn
		mbabout = append(mbabout, plastk.MenuBarColumn{
			ColumnType: "button",
			ColumnName: "About",
		},
		)
		plastk.DrawMenuBar(screen, color.RGBA{100, 100, 100, 255}, hackgen4info, 20, mbfiles, mbedit, mbabout)
	}

	// draw command-line
	if commandlineforcused || cmdresult != "" {
		commandlineop := &ebiten.DrawImageOptions{}
		commandlineop.GeoM.Translate(float64(0), float64(screenHeight+commandlineSize))
		commandLine = plastk.ReNewImg(commandLine, screenWidth, commandlineSize, color.RGBA{39, 39, 39, 255})
		screen.DrawImage(commandLine, commandlineop)
		if commandlineforcused {
			text.Draw(screen, photoncmd, smallHackGenFont, 5, screenHeight+15+commandlineSize, color.White)
		} else {
			text.Draw(screen, cmdresult, smallHackGenFont, 5, screenHeight+15+commandlineSize, color.White)
		}
	}

	// Draw info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))

	/* Benchmark
	if len(os.Args) >= 2 {
		if os.Args[1] == "bench" {
			os.Exit(0)
		}
	}*/
}

func (g *Editor) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenWidth, screenHeight := ebiten.WindowSize()
	return screenWidth, screenHeight
}

func ezSyntaxHighlight(txt string) (string, int) {
	splitted := strings.Split(txt, " ")
	ret := len([]rune(splitted[0]))
	if splitted[0] == "package" {
		return "good", ret
	}
	if splitted[0] == "import" {
		return "good", ret
	}
	if splitted[0] == "func" {
		return "good", ret
	}
	if splitted[0] == "fn" {
		return "good", ret
	}
	if splitted[0] == "var" {
		return "good", ret
	}
	return "", 0
}

/*func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("PhotonText(kari)")
	ebiten.SetTPS(140)

	go func() {
		now := time.Now()
		fmt.Println("PhotonText will booted with no error(s)")

	loginloop:
		err := client.Login("1199337296307163146")
		if err != nil {
			time.Sleep(20 * time.Second)
			goto loginloop
		}

		success := bool(false)
	activityloop:
		state := strconv.Itoa(cursornowy) + string(":") + strconv.Itoa(cursornowx)
		err = client.SetActivity(client.Activity{
			Details:    "Coding with PhotonText",
			State:      state,
			LargeImage: "photon2",
			LargeText:  "PhotonText Logo",
			Timestamps: &client.Timestamps{
				Start: &now,
			},
		})
		if err != nil {
			fmt.Println(err)
			goto loginloop
		}
		if !success {
			fmt.Println("rich presence active")
			success = true
		}
		time.Sleep(1 * time.Second)
		goto activityloop
	}()

	go func() {
		for {
			looping := false
			if clearresult == 0 {
				if !looping {
					cmdresult = ""
				}
				time.Sleep(1 * time.Second)
				looping = true
			} else if clearresult <= 10 {
				looping = false
				time.Sleep(1 * time.Second)
				clearresult--
			} else if clearresult > 10 {
				looping = false
				clearresult = 10
			}
		}
	}()

	go func() {
		for {
			if editingfile == "" {
				ebiten.SetWindowTitle("PhotonText(kari)")
			} else {
				ebiten.SetWindowTitle(fmt.Sprint(editingfile, " - PhotonText(kari)"))
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	if err := ebiten.RunGame(&Editor{}); err != nil {
		log.Fatal(err)
	}
}*/

func proceedcmd(command string) string {
	return ProceedCmd(command)
}

// Proceed command
func ProceedCmd(command string) (returnstr string) {
	command2slice := strings.Split(command, " ")
	if len(command2slice) >= 1 {
		cmd := command2slice[0]
		// Save override
		if cmd == "w" || cmd == "wr" || cmd == "wri" || cmd == "writ" || cmd == "write" {
			if len(command2slice) >= 2 {
				return "Too many arguments for command: " + cmd
			} else {
				phsave(editingfile)
				return fmt.Sprint("Saved to ", editingfile)
			}
		} else
		//
		if cmd == "q" || cmd == "qu" || cmd == "qui" || cmd == "quit" {
			ebiten.SetWindowClosingHandled(true)
			closewindow = true
		} else if cmd == "wq" {
			proceedcmd("w")
			proceedcmd("q")
		} else if cmd == "version" {
			return "PhotonText rolling " + runtime.Version()
		} else
		// Save with other name.
		if cmd == "sav" || cmd == "save" || cmd == "savea" || cmd == "saveas" {
			if len(command2slice) == 1 {
				return "Too few arguments for command: " + cmd
			} else if len(command2slice) >= 3 {
				return "Too many arguments for command: " + cmd
			} else /* when 2 args */ {
				if strings.HasPrefix(command2slice[1], "~") {
					home, err := os.UserHomeDir()
					if err != nil {
						fmt.Println(err)
					}
					savepath := home + command2slice[1][1:]
					phsave(savepath)
					return fmt.Sprint("Saved to ", savepath)
				} else {
					phsave(command2slice[1])
					return fmt.Sprint("Saved to ", command2slice[1])
				}
			}
		} else
		// Toggle VSync
		if command2slice[0] == "togglevsync" {
			ebiten.SetVsyncEnabled(!ebiten.IsVsyncEnabled())
			return "Toggled VSync"
		} else if command2slice[0] == "set" {
			if len(command2slice) == 1 {
				return "Too few arguments for command: " + cmd
			} else if len(command2slice) >= 3 {
				return "Too many arguments for command: " + cmd
			} else {
				var2set := strings.Split(command2slice[1], "=")[1]
				if 1 < len(var2set) {
					switch strings.Split(command2slice[1], "=")[0] {
					case "vsync":
						if dyntypes.IsDynTypeMatch(var2set, "bool") {
							ebiten.SetVsyncEnabled(dyntypes.DynBool(var2set))
						}
						return strconv.FormatBool(ebiten.IsVsyncEnabled())
					case "rellines":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							rellines = dyntypes.DynInt(var2set)
						}
					case "topopbarsize":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							topopBarSize = dyntypes.DynInt(var2set)
						}
					case "infobarsize":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							infoBarSize = dyntypes.DynInt(var2set)
						}
					case "commandlinesize":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							commandlineSize = dyntypes.DynInt(var2set)
						}
					case "limitter":
						if dyntypes.IsDynTypeMatch(var2set, "bool") {
							limitterenabled = dyntypes.DynBool(var2set)
						}
					case "limitterlevel":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							limitterlevel = dyntypes.DynInt(var2set)
						}
					case "modalmode":
						if dyntypes.IsDynTypeMatch(var2set, "bool") {
							modalmode = dyntypes.DynBool(var2set)
						}
					case "ebiten-tps":
						if dyntypes.IsDynTypeMatch(var2set, "int") {
							ebiten.SetTPS(dyntypes.DynInt(var2set))
						}
					default:
						return "No internal variables named " + (strings.Split(command2slice[1], "="))[0]
					}
				}
			}
		} else
		// If not command is avaliable
		{
			return fmt.Sprintf("Not an editor command: %s", command2slice[0])
		}
	} else {
		return "No command was input."
	}
	return
}

func tab2space(inp string) string {
	return strings.ReplaceAll(inp, "	", "   ")
}

func phloginfo(pherror error) {
	if pherror != nil {
		fmt.Println("error:", pherror)
	}
}

// file load/save
func phload(inputpath string) {
	file, err := os.ReadFile(inputpath)
	if err != nil {
		panic(err)
	}
	editingfile, err = filepath.Abs(inputpath)
	if err != nil {
		panic(err)
	}
	ftext := string(file)

	// Check CRLF(dos) First, if not, Use LF(*nix).
	if strings.Contains(ftext, "\r\n") {
		PhotonText = strings.Split(ftext, "\r\n")
		returncode = "\r\n"
		returntype = "CRLF"
	} else {
		PhotonText = strings.Split(ftext, "\n")
		returncode = "\n"
		returntype = "LF"
	}
}

func sliceload(inputpath string) ([]string, error) {
	var slice2load []string
	var sliceerr error
	file, err := os.ReadFile(inputpath)
	if err != nil {
		sliceerr = err
	}
	ftext := string(file)

	// Check CRLF(dos) First, if not, Use LF(*nix).
	if strings.Contains(ftext, "\r\n") {
		slice2load = strings.Split(ftext, "\r\n")
		returncode = "\r\n"
		returntype = "CRLF"
	} else {
		slice2load = strings.Split(ftext, "\n")
		returncode = "\n"
		returntype = "LF"
	}
	return slice2load, sliceerr
}

func phsave(dir string) {
	output := strings.Join(PhotonText, returncode)
	runeout := []rune(output)
	err := os.WriteFile(dir, []byte(string(runeout)), 0644)
	if err != nil {
		fmt.Println(err, "Save failed")
	}
}
