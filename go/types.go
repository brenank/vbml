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
	Height int `json:"height,omitempty"`
	Width  int `json:"width,omitempty"`
}

type AbsolutePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ComponentStyle struct {
	Justify          Justify           `json:"justify,omitempty"`
	Align            Align             `json:"align,omitempty"`
	Height           int               `json:"height,omitempty"`
	Width            int               `json:"width,omitempty"`
	AbsolutePosition *AbsolutePosition `json:"absolutePosition,omitempty"`
}

type CalendarData struct {
	Month           int         `json:"month,omitempty"`
	Year            int         `json:"year,omitempty"`
	DefaultDayColor int         `json:"defaultDayColor,omitempty"`
	Days            map[int]int `json:"days,omitempty"`
	HideSMTWTFS     bool        `json:"hideSMTWTFS,omitempty"`
	HideDates       bool        `json:"hideDates,omitempty"`
	HideMonthYear   bool        `json:"hideMonthYear,omitempty"`
}

type RandomColorsData struct {
	Colors []int `json:"colors,omitempty"`
}

type TemplateWrap string

const (
	TemplateWrapNormal TemplateWrap = "normal"
	TemplateWrapNever  TemplateWrap = "never"
)

type TemplatePart struct {
	Template string       `json:"template"`
	Wrap     TemplateWrap `json:"wrap,omitempty"`
}

type Component struct {
	Template      string            `json:"template,omitempty"`
	TemplateParts []TemplatePart    `json:"-"`
	RawCharacters Board             `json:"rawCharacters,omitempty"`
	Calendar      *CalendarData     `json:"calendar,omitempty"`
	RandomColors  *RandomColorsData `json:"randomColors,omitempty"`
	Style         *ComponentStyle   `json:"style,omitempty"`
}

type Input struct {
	Props      map[string]any `json:"props,omitempty"`
	Style      *BoardStyle    `json:"style,omitempty"`
	Components []Component    `json:"components"`
}

type AbsoluteComponent struct {
	Characters Board `json:"characters"`
	X          int   `json:"x"`
	Y          int   `json:"y"`
}

type CharacterCodesToStringOptions struct {
	AllowLineBreaks bool `json:"allowLineBreaks,omitempty"`
}

type ClassicOptions struct {
	ExtraHPadding        int  `json:"extraHPadding,omitempty"`
	PreserveDoubleSpaces bool `json:"preserveDoubleSpaces,omitempty"`
}
