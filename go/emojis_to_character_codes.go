package vbml

import "strings"

var emojiToCharacterCodeReplacer = strings.NewReplacer(
	"🟥", "{63}",
	"🟧", "{64}",
	"🟨", "{65}",
	"🟩", "{66}",
	"🟦", "{67}",
	"🟪", "{68}",
	"⬜", "{69}",
	"⬛", "{70}",
	"ß", "SS",
	"❤️", "{62}",
)

func emojisToCharacterCodes(template string) string {
	return emojiToCharacterCodeReplacer.Replace(template)
}
