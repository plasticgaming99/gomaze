// not random numbers, random functions
package randoms

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func RepeatingKeyPressed(key ebiten.Key) bool {
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

func RepeatingKeyPressedHiFreq(key ebiten.Key) bool {
	var (
		delay    = ebiten.TPS() / 1
		interval = ebiten.TPS() / 4
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
