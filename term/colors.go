package term

import "fmt"

// RGB color function
func RGB(r, g, b int) Color {
	return Color(fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b))
}

// Background RGB color function
func BgRGB(r, g, b int) Color {
	return Color(fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b))
}

// Reset is the ANSI escape code to reset all attributes.
const Reset = "\033[0m"

const (
	// Special effects
	None      = Color("-")
	Bold      = Color("\033[1m")
	Dim       = Color("\033[2m")
	Italic    = Color("\033[3m")
	Underline = Color("\033[4m")
	Blink     = Color("\033[5m")
	Reverse   = Color("\033[7m")
	Hidden    = Color("\033[8m")

	// Basic colors
	Black   = Color("\033[30m") // #000000
	Red     = Color("\033[31m") // #FF0000
	Green   = Color("\033[32m") // #00FF00
	Yellow  = Color("\033[33m") // #FFFF00
	Blue    = Color("\033[34m") // #0000FF
	Magenta = Color("\033[35m") // #FF00FF
	Cyan    = Color("\033[36m") // #00FFFF
	White   = Color("\033[37m") // #FFFFFF

	// Bright colors
	BrightBlack   = Color("\033[90m") // #808080
	BrightRed     = Color("\033[91m") // #FF6347
	BrightGreen   = Color("\033[92m") // #00FF7F
	BrightYellow  = Color("\033[93m") // #FFFF00
	BrightBlue    = Color("\033[94m") // #1E90FF
	BrightMagenta = Color("\033[95m") // #FF00FF
	BrightCyan    = Color("\033[96m") // #00FFFF
	BrightWhite   = Color("\033[97m") // #FFFFFF

	// Background colors
	BgBlack   = Color("\033[40m") // #000000
	BgRed     = Color("\033[41m") // #FF0000
	BgGreen   = Color("\033[42m") // #00FF00
	BgYellow  = Color("\033[43m") // #FFFF00
	BgBlue    = Color("\033[44m") // #0000FF
	BgMagenta = Color("\033[45m") // #FF00FF
	BgCyan    = Color("\033[46m") // #00FFFF
	BgWhite   = Color("\033[47m") // #FFFFFF

	// Bright background colors
	BgBrightBlack   = Color("\033[100m") // #808080
	BgBrightRed     = Color("\033[101m") // #FF6347
	BgBrightGreen   = Color("\033[102m") // #00FF7F
	BgBrightYellow  = Color("\033[103m") // #FFFF00
	BgBrightBlue    = Color("\033[104m") // #1E90FF
	BgBrightMagenta = Color("\033[105m") // #FF00FF
	BgBrightCyan    = Color("\033[106m") // #00FFFF
	BgBrightWhite   = Color("\033[107m") // #FFFFFF

	// Extended colors - Reds
	LightCoral = Color("\033[38;5;210m") // #F08080
	IndianRed  = Color("\033[38;5;167m") // #CD5C5C
	Crimson    = Color("\033[38;5;160m") // #DC143C
	Maroon     = Color("\033[38;5;52m")  // #800000
	DarkRed    = Color("\033[38;5;88m")  // #8B0000

	// Extended colors - Pinks
	Pink            = Color("\033[38;5;218m") // #FFC0CB
	LightPink       = Color("\033[38;5;217m") // #FFB6C1
	HotPink         = Color("\033[38;5;205m") // #FF69B4
	DeepPink        = Color("\033[38;5;198m") // #FF1493
	MediumVioletRed = Color("\033[38;5;162m") // #C71585

	// Extended colors - Oranges
	Coral      = Color("\033[38;5;209m") // #FF7F50
	Tomato     = Color("\033[38;5;203m") // #FF6347
	OrangeRed  = Color("\033[38;5;202m") // #FF4500
	DarkOrange = Color("\033[38;5;208m") // #FF8C00
	Orange     = Color("\033[38;5;214m") // #FFA500

	// Extended colors - Yellows
	Gold         = Color("\033[38;5;220m") // #FFD700
	LightYellow  = Color("\033[38;5;228m") // #FFFFE0
	LemonChiffon = Color("\033[38;5;230m") // #FFFACD
	Khaki        = Color("\033[38;5;185m") // #F0E68C
	DarkKhaki    = Color("\033[38;5;143m") // #BDB76B

	// Extended colors - Purples
	Lavender     = Color("\033[38;5;183m") // #E6E6FA
	Thistle      = Color("\033[38;5;182m") // #D8BFD8
	Plum         = Color("\033[38;5;96m")  // #DDA0DD
	Violet       = Color("\033[38;5;135m") // #EE82EE
	Orchid       = Color("\033[38;5;170m") // #DA70D6
	Fuchsia      = Color("\033[38;5;201m") // #FF00FF
	MediumOrchid = Color("\033[38;5;134m") // #BA55D3
	MediumPurple = Color("\033[38;5;141m") // #9370DB
	BlueViolet   = Color("\033[38;5;92m")  // #8A2BE2
	DarkViolet   = Color("\033[38;5;128m") // #9400D3
	DarkOrchid   = Color("\033[38;5;98m")  // #9932CC
	DarkMagenta  = Color("\033[38;5;90m")  // #8B008B
	Purple       = Color("\033[38;5;129m") // #800080
	Indigo       = Color("\033[38;5;54m")  // #4B0082

	// Extended colors - Greens
	GreenYellow       = Color("\033[38;5;154m") // #ADFF2F
	Chartreuse        = Color("\033[38;5;118m") // #7FFF00
	LawnGreen         = Color("\033[38;5;118m") // #7CFC00
	Lime              = Color("\033[38;5;46m")  // #00FF00
	LimeGreen         = Color("\033[38;5;77m")  // #32CD32
	PaleGreen         = Color("\033[38;5;120m") // #98FB98
	LightGreen        = Color("\033[38;5;119m") // #90EE90
	MediumSpringGreen = Color("\033[38;5;49m")  // #00FA9A
	SpringGreen       = Color("\033[38;5;48m")  // #00FF7F
	MediumSeaGreen    = Color("\033[38;5;77m")  // #3CB371
	SeaGreen          = Color("\033[38;5;29m")  // #2E8B57
	ForestGreen       = Color("\033[38;5;28m")  // #228B22
	DarkGreen         = Color("\033[38;5;22m")  // #006400
	YellowGreen       = Color("\033[38;5;113m") // #9ACD32
	OliveDrab         = Color("\033[38;5;64m")  // #6B8E23
	Olive             = Color("\033[38;5;58m")  // #808000
	DarkOliveGreen    = Color("\033[38;5;59m")  // #556B2F
	MediumAquamarine  = Color("\033[38;5;79m")  // #66CDAA
	DarkSeaGreen      = Color("\033[38;5;108m") // #8FBC8F

	// Extended colors - Blues/Cyans
	Aqua            = Color("\033[38;5;51m")  // #00FFFF
	LightCyan       = Color("\033[38;5;195m") // #E0FFFF
	PaleTurquoise   = Color("\033[38;5;159m") // #AFEEEE
	Aquamarine      = Color("\033[38;5;122m") // #7FFFD4
	Turquoise       = Color("\033[38;5;80m")  // #40E0D0
	MediumTurquoise = Color("\033[38;5;80m")  // #48D1CC
	DarkTurquoise   = Color("\033[38;5;44m")  // #00CED1
	LightSeaGreen   = Color("\033[38;5;37m")  // #20B2AA
	CadetBlue       = Color("\033[38;5;73m")  // #5F9EA0
	DarkCyan        = Color("\033[38;5;36m")  // #008B8B
	Teal            = Color("\033[38;5;23m")  // #008080
	LightSteelBlue  = Color("\033[38;5;152m") // #B0C4DE
	PowderBlue      = Color("\033[38;5;152m") // #B0E0E6
	LightBlue       = Color("\033[38;5;152m") // #ADD8E6
	SkyBlue         = Color("\033[38;5;117m") // #87CEEB
	LightSkyBlue    = Color("\033[38;5;117m") // #87CEFA
	DeepSkyBlue     = Color("\033[38;5;39m")  // #00BFFF
	DodgerBlue      = Color("\033[38;5;33m")  // #1E90FF
	CornflowerBlue  = Color("\033[38;5;69m")  // #6495ED
	SteelBlue       = Color("\033[38;5;67m")  // #4682B4
	RoyalBlue       = Color("\033[38;5;62m")  // #4169E1
	MediumBlue      = Color("\033[38;5;20m")  // #0000CD
	DarkBlue        = Color("\033[38;5;18m")  // #00008B
	Navy            = Color("\033[38;5;17m")  // #000080
	MidnightBlue    = Color("\033[38;5;17m")  // #191970

	// Extended colors - Browns
	Cornsilk       = Color("\033[38;5;230m") // #FFF8DC
	BlanchedAlmond = Color("\033[38;5;230m") // #FFEBCD
	Bisque         = Color("\033[38;5;224m") // #FFE4C4
	NavajoWhite    = Color("\033[38;5;223m") // #FFDAB9
	Wheat          = Color("\033[38;5;223m") // #F5DEB3
	BurlyWood      = Color("\033[38;5;180m") // #DEB887
	Tan            = Color("\033[38;5;180m") // #D2B48C
	RosyBrown      = Color("\033[38;5;138m") // #BC8F8F
	SandyBrown     = Color("\033[38;5;215m") // #F4A460
	Goldenrod      = Color("\033[38;5;178m") // #DAA520
	DarkGoldenrod  = Color("\033[38;5;136m") // #B8860B
	Peru           = Color("\033[38;5;173m") // #CD853F
	Chocolate      = Color("\033[38;5;166m") // #D2691E
	SaddleBrown    = Color("\033[38;5;94m")  // #8B4513
	Sienna         = Color("\033[38;5;94m")  // #A0522D
	Brown          = Color("\033[38;5;124m") // #A52A2A

	// Extended colors - Whites
	Snow          = Color("\033[38;5;15m")  // #FFFAFA
	HoneyDew      = Color("\033[38;5;15m")  // #F0FFF0
	MintCream     = Color("\033[38;5;15m")  // #F5FFFA
	Azure         = Color("\033[38;5;15m")  // #F0FFFF
	AliceBlue     = Color("\033[38;5;15m")  // #F0F8FF
	GhostWhite    = Color("\033[38;5;15m")  // #F8F8FF
	WhiteSmoke    = Color("\033[38;5;15m")  // #F5F5F5
	SeaShell      = Color("\033[38;5;15m")  // #FFF5EE
	Beige         = Color("\033[38;5;15m")  // #F5F5DC
	OldLace       = Color("\033[38;5;230m") // #FDF5E6
	FloralWhite   = Color("\033[38;5;15m")  // #FFFAF0
	Ivory         = Color("\033[38;5;15m")  // #FFFFF0
	AntiqueWhite  = Color("\033[38;5;230m") // #FAEBD7
	Linen         = Color("\033[38;5;255m") // #FAF0E6
	LavenderBlush = Color("\033[38;5;15m")  // #FFF0F5
	MistyRose     = Color("\033[38;5;224m") // #FFE4E1

	// Extended colors - Grays
	Gainsboro      = Color("\033[38;5;252m") // #DCDCDC
	LightGray      = Color("\033[38;5;250m") // #D3D3D3
	Silver         = Color("\033[38;5;7m")   // #C0C0C0
	DarkGray       = Color("\033[38;5;240m") // #A9A9A9
	Gray           = Color("\033[38;5;245m") // #808080
	DimGray        = Color("\033[38;5;242m") // #696969
	LightSlateGray = Color("\033[38;5;246m") // #778899
	SlateGray      = Color("\033[38;5;246m") // #708090
	DarkSlateGray  = Color("\033[38;5;238m") // #2F4F4F

	// Additional background colors
	BgLightGray      = Color("\033[48;5;250m") // #D3D3D3
	BgSilver         = Color("\033[48;5;7m")   // #C0C0C0
	BgDarkGray       = Color("\033[48;5;240m") // #A9A9A9
	BgGray           = Color("\033[48;5;245m") // #808080
	BgDimGray        = Color("\033[48;5;242m") // #696969
	BgLightSlateGray = Color("\033[48;5;246m") // #778899
	BgSlateGray      = Color("\033[48;5;246m") // #708090
	BgDarkSlateGray  = Color("\033[48;5;238m") // #2F4F4F

	// Additional background colors for common hues
	BgLightRed     = Color("\033[48;5;217m") // #FFCCCB
	BgLightGreen   = Color("\033[48;5;120m") // #90EE90
	BgLightBlue    = Color("\033[48;5;153m") // #ADD8E6
	BgLightYellow  = Color("\033[48;5;228m") // #FFFFE0
	BgLightCyan    = Color("\033[48;5;195m") // #E0FFFF
	BgLightMagenta = Color("\033[48;5;219m") // #FFB6C1
	BgPink         = Color("\033[48;5;218m") // #FFC0CB
	BgLightOrange  = Color("\033[48;5;223m") // #FFD3A5
	BgLightPurple  = Color("\033[48;5;189m") // #E6E6FA
	BgLightBrown   = Color("\033[48;5;180m") // #D2B48C

	// Miscellaneous colors
	DefaultForeground = Color("\033[39m") // Default foreground color
	DefaultBackground = Color("\033[49m") // Default background color
)
