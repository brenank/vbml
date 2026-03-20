package vbml

type Board [][]int

type Justify string

const (
	JustifyCenter    Justify = "center"
	JustifyLeft      Justify = "left"
	JustifyRight     Justify = "right"
	JustifyJustified Justify = "justified"
)

type Align string

const (
	AlignCenter    Align = "center"
	AlignTop       Align = "top"
	AlignBottom    Align = "bottom"
	AlignJustified Align = "justified"
	AlignAbsolute  Align = "absolute"
)

type CharacterCode int

const (
	CharacterCodeBlank CharacterCode = 0

	CharacterCodeLetterA CharacterCode = 1
	CharacterCodeLetterB CharacterCode = 2
	CharacterCodeLetterC CharacterCode = 3
	CharacterCodeLetterD CharacterCode = 4
	CharacterCodeLetterE CharacterCode = 5
	CharacterCodeLetterF CharacterCode = 6
	CharacterCodeLetterG CharacterCode = 7
	CharacterCodeLetterH CharacterCode = 8
	CharacterCodeLetterI CharacterCode = 9
	CharacterCodeLetterJ CharacterCode = 10
	CharacterCodeLetterK CharacterCode = 11
	CharacterCodeLetterL CharacterCode = 12
	CharacterCodeLetterM CharacterCode = 13
	CharacterCodeLetterN CharacterCode = 14
	CharacterCodeLetterO CharacterCode = 15
	CharacterCodeLetterP CharacterCode = 16
	CharacterCodeLetterQ CharacterCode = 17
	CharacterCodeLetterR CharacterCode = 18
	CharacterCodeLetterS CharacterCode = 19
	CharacterCodeLetterT CharacterCode = 20
	CharacterCodeLetterU CharacterCode = 21
	CharacterCodeLetterV CharacterCode = 22
	CharacterCodeLetterW CharacterCode = 23
	CharacterCodeLetterX CharacterCode = 24
	CharacterCodeLetterY CharacterCode = 25
	CharacterCodeLetterZ CharacterCode = 26

	CharacterCodeOne   CharacterCode = 27
	CharacterCodeTwo   CharacterCode = 28
	CharacterCodeThree CharacterCode = 29
	CharacterCodeFour  CharacterCode = 30
	CharacterCodeFive  CharacterCode = 31
	CharacterCodeSix   CharacterCode = 32
	CharacterCodeSeven CharacterCode = 33
	CharacterCodeEight CharacterCode = 34
	CharacterCodeNine  CharacterCode = 35
	CharacterCodeZero  CharacterCode = 36

	CharacterCodeExclamationMark  CharacterCode = 37
	CharacterCodeAtSign           CharacterCode = 38
	CharacterCodePoundSign        CharacterCode = 39
	CharacterCodeDollarSign       CharacterCode = 40
	CharacterCodeLeftParenthesis  CharacterCode = 41
	CharacterCodeRightParenthesis CharacterCode = 42
	CharacterCodeHyphen           CharacterCode = 44
	CharacterCodePlusSign         CharacterCode = 46
	CharacterCodeAmpersand        CharacterCode = 47
	CharacterCodeEqualsSign       CharacterCode = 48
	CharacterCodeSemicolon        CharacterCode = 49
	CharacterCodeColon            CharacterCode = 50
	CharacterCodeSingleQuote      CharacterCode = 52
	CharacterCodeDoubleQuote      CharacterCode = 53
	CharacterCodePercentSign      CharacterCode = 54
	CharacterCodeComma            CharacterCode = 55
	CharacterCodePeriod           CharacterCode = 56
	CharacterCodeSlash            CharacterCode = 59
	CharacterCodeQuestionMark     CharacterCode = 60
	CharacterCodeDegreeSign       CharacterCode = 62

	CharacterCodeRed    CharacterCode = 63
	CharacterCodeOrange CharacterCode = 64
	CharacterCodeYellow CharacterCode = 65
	CharacterCodeGreen  CharacterCode = 66
	CharacterCodeBlue   CharacterCode = 67
	CharacterCodeViolet CharacterCode = 68
	CharacterCodeWhite  CharacterCode = 69
	CharacterCodeBlack  CharacterCode = 70
	CharacterCodeFilled CharacterCode = 71
)

const (
	flagshipBoardHeight = 6
	flagshipBoardWidth  = 22
)

type BoardStyle struct {
	Height int
	Width  int
}

type AbsolutePosition struct {
	X int
	Y int
}

type ComponentStyle struct {
	Justify          Justify
	Align            Align
	Height           int
	Width            int
	AbsolutePosition *AbsolutePosition
}

type CalendarData struct {
	Month           int
	Year            int
	DefaultDayColor int
	Days            map[int]int
	HideSMTWTFS     bool
	HideDates       bool
	HideMonthYear   bool
}

type RandomColorsData struct {
	Colors []int
}

type Component struct {
	Template      string
	RawCharacters Board
	Calendar      *CalendarData
	RandomColors  *RandomColorsData
	Style         *ComponentStyle
}

type Input struct {
	Props      map[string]any
	Style      *BoardStyle
	Components []Component
}

type AbsoluteComponent struct {
	Characters Board
	X          int
	Y          int
}

type CharacterCodesToStringOptions struct {
	AllowLineBreaks bool
}

type ClassicOptions struct {
	ExtraHPadding        int
	PreserveDoubleSpaces bool
}
