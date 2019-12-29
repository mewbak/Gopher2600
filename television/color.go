package television

// VideoBlack is the PixelSignal value that indicates no VCS pixel is to be shown
const VideoBlack ColorSignal = -1

// a color is made of a number of color components
type color struct {
	red   byte
	green byte
	blue  byte
}

// colors is the entire palette
type colors []color

// the entire palette is made up of many colors
var colorsNTSC = colors{}
var colorsPAL = colors{}
var colorsAlt = colors{}

// the VideoBlack signal results in the following color
var videoBlack = color{red: 0, green: 0, blue: 0}

// getColor translates a signals to the color type
func getColor(spec *Specification, col ColorSignal) color {
	// we're usng the ColorSignal to index an array so we need to be extra
	// careful to make sure the value is valid. if it's not a valid index then
	// assume the intention was video black
	if col == VideoBlack {
		return videoBlack
	}
	return spec.Colors[col]
}

// List of valid AltColorSignals
const (
	AltColBackground AltColorSignal = iota
	AltColBall
	AltColPlayfield
	AltColPlayer0
	AltColPlayer1
	AltColMissile0
	AltColMissile1
	altColCount
)

// getAltColor translates an alternative color signals to the color type
func getAltColor(altCol AltColorSignal) color {
	// we're usng the AltColorSignal to index an array so we need to be extra
	// careful to make sure the value is valid. if it's not a valid index then
	// assume the intention is video black
	if altCol > altColCount {
		return videoBlack
	}
	return colorsAlt[altCol]
}
