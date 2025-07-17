// it provides misaki font, a free 8x8 bitmap outline font
package fonts

import (
	_ "embed"
)

var (
	//go:embed misaki_gothic.ttf
	MisakiGothicFont []byte
	//go:embed misaki_gothic_2nd.ttf
	MisakiGothic2ndFont []byte
	//go:embed misaki_mincho.ttf
	MisakiMincho []byte
)
