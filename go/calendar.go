package vbml

import (
	"strconv"
	"time"
)

func MakeCalendar(month, year int, options CalendarData) Board {
	numberOfDaysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	firstDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Weekday()
	offset := int(firstDayOfMonth)

	calendarDayColor := options.DefaultDayColor
	if calendarDayColor == 0 {
		calendarDayColor = int(CharacterCodeYellow)
	}

	firstRowDays := [2]string{strconv.Itoa(1), strconv.Itoa(7 - offset)}
	secondRowDays := [2]string{strconv.Itoa(8 - offset), strconv.Itoa(14 - offset)}
	thirdRowDays := [2]string{strconv.Itoa(15 - offset), strconv.Itoa(21 - offset)}
	fourthRowDays := [2]string{strconv.Itoa(22 - offset), strconv.Itoa(28 - offset)}
	fifthStart := 29 - offset
	fifthEnd := minInt(7-offset+numberOfDaysInMonth, numberOfDaysInMonth)
	numberOfDaysInLastRow := fifthEnd - fifthStart + 1

	firstRow := buildFirstCalendarRow(firstRowDays, offset, calendarDayColor, options.HideDates)
	secondRow := buildCalendarRangeRow(secondRowDays, calendarDayColor, options.HideDates)
	thirdRow := buildCalendarRangeRow(thirdRowDays, calendarDayColor, options.HideDates)
	fourthRow := buildExactCalendarRangeRow(fourthRowDays, calendarDayColor, options.HideDates)
	fifthRow := buildFifthCalendarRow(
		fifthStart,
		fifthEnd,
		numberOfDaysInMonth,
		numberOfDaysInLastRow,
		calendarDayColor,
		options.HideDates,
	)

	calendar := Board{
		buildCalendarHeaderRow(month, year, options.HideMonthYear, options.HideSMTWTFS),
		firstRow,
		secondRow,
		thirdRow,
		fourthRow,
		fifthRow,
	}

	for day, color := range options.Days {
		if day > numberOfDaysInMonth {
			continue
		}

		todaysRow := ((day + offset - 1) / 7) + 1
		modulus := (day + offset - 1) % 7
		todaysColumn := modulus + 5
		if todaysRow > 5 {
			if modulus == 0 {
				todaysColumn = 12
			} else {
				todaysColumn = 13
			}
		}

		rowIndex := todaysRow
		if rowIndex > 5 {
			rowIndex = 5
		}
		calendar[rowIndex][todaysColumn] = color
	}

	return calendar
}

func buildCalendarHeaderRow(month, year int, hideMonthYear, hideDayOfWeek bool) []int {
	var monthYear []int
	if hideMonthYear {
		monthYear = []int{0, 0, 0, 0, 0}
	} else {
		for _, digit := range strconv.Itoa(month) {
			monthYear = append(monthYear, calendarDigitCharacterCode(digit))
		}
		monthYear = append(monthYear, int(CharacterCodeSlash))
		yearDigits := strconv.Itoa(year)
		for _, digit := range yearDigits[len(yearDigits)-2:] {
			monthYear = append(monthYear, calendarDigitCharacterCode(digit))
		}
	}

	headerSpace := 5 - len(monthYear)
	row := append([]int{}, monthYear...)
	row = append(row, make([]int, headerSpace)...)

	if hideDayOfWeek {
		row = append(row, make([]int, 7)...)
	} else {
		for _, day := range []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"} {
			row = append(row, weekdayCharacterCode(day))
		}
	}

	row = append(row, make([]int, 22-(7+5))...)
	return row
}

func buildFirstCalendarRow(days [2]string, offset, dayColor int, hideDates bool) []int {
	if days[0] == days[1] {
		row := []int{0, 0, 0, dateDigit(days[0], 0, hideDates), 0}
		row = append(row, make([]int, offset)...)
		row = append(row, repeatedInt(dayColor, 7-offset)...)
		row = append(row, make([]int, 22-12)...)
		return row
	}

	row := []int{
		0,
		dateDigit(days[0], 0, hideDates),
		calendarHyphen(hideDates),
		dateDigit(days[1], 0, hideDates),
		0,
	}
	row = append(row, make([]int, offset)...)
	row = append(row, repeatedInt(dayColor, 7-offset)...)
	row = append(row, make([]int, 22-12)...)
	return row
}

func buildCalendarRangeRow(days [2]string, dayColor int, hideDates bool) []int {
	row := make([]int, 0, 22)
	startDigits := splitCalendarDigits(days[0])
	endDigits := splitCalendarDigits(days[1])

	if len(startDigits) > 1 {
		row = append(row, digitOrBlank(startDigits, 0, hideDates), digitOrBlank(startDigits, 1, hideDates))
	} else {
		row = append(row, 0, digitOrBlank(startDigits, 0, hideDates))
	}

	row = append(row, calendarHyphen(hideDates))
	if len(endDigits) > 1 {
		row = append(row, digitOrBlank(endDigits, 0, hideDates), digitOrBlank(endDigits, 1, hideDates))
	} else {
		row = append(row, digitOrBlank(endDigits, 0, hideDates), 0)
	}

	row = append(row, repeatedInt(dayColor, 7)...)
	row = append(row, make([]int, 22-12)...)
	return row
}

func buildExactCalendarRangeRow(days [2]string, dayColor int, hideDates bool) []int {
	startDigits := splitCalendarDigits(days[0])
	endDigits := splitCalendarDigits(days[1])

	row := []int{
		digitOrBlank(startDigits, 0, hideDates),
		digitOrBlank(startDigits, 1, hideDates),
		calendarHyphen(hideDates),
		digitOrBlank(endDigits, 0, hideDates),
		digitOrBlank(endDigits, 1, hideDates),
	}
	row = append(row, repeatedInt(dayColor, 7)...)
	row = append(row, make([]int, 22-12)...)
	return row
}

func buildFifthCalendarRow(
	start,
	end,
	numberOfDaysInMonth,
	numberOfDaysInLastRow,
	dayColor int,
	hideDates bool,
) []int {
	if start > numberOfDaysInMonth {
		return make([]int, 22)
	}

	startDigits := splitCalendarDigits(strconv.Itoa(start))
	endDigits := splitCalendarDigits(strconv.Itoa(end))
	hideRange := hideDates || start == end

	row := []int{
		digitOrBlank(startDigits, 0, hideDates),
		digitOrBlank(startDigits, 1, hideDates),
		calendarHyphen(hideRange),
		digitOrBlank(endDigits, 0, hideRange),
		digitOrBlank(endDigits, 1, hideRange),
	}
	row = append(row, repeatedInt(dayColor, numberOfDaysInLastRow)...)
	row = append(row, make([]int, 22-(5+numberOfDaysInLastRow))...)
	return row
}

func splitCalendarDigits(value string) []rune {
	return []rune(value)
}

func digitOrBlank(digits []rune, index int, hide bool) int {
	if hide || index >= len(digits) {
		return 0
	}
	return calendarDigitCharacterCode(digits[index])
}

func dateDigit(value string, index int, hide bool) int {
	return digitOrBlank(splitCalendarDigits(value), index, hide)
}

func calendarHyphen(hide bool) int {
	if hide {
		return 0
	}
	return int(CharacterCodeHyphen)
}

func repeatedInt(value, count int) []int {
	repeated := make([]int, count)
	for index := range repeated {
		repeated[index] = value
	}
	return repeated
}

func calendarDigitCharacterCode(digit rune) int {
	if digit == '0' {
		return int(CharacterCodeZero)
	}
	return int(digit-'0') + 26
}

func weekdayCharacterCode(day string) int {
	switch day {
	case "Sun":
		return int(CharacterCodeLetterS)
	case "Mon":
		return int(CharacterCodeLetterM)
	case "Tue":
		return int(CharacterCodeLetterT)
	case "Wed":
		return int(CharacterCodeLetterW)
	case "Thu":
		return int(CharacterCodeLetterT)
	case "Fri":
		return int(CharacterCodeLetterF)
	case "Sat":
		return int(CharacterCodeLetterS)
	default:
		return 0
	}
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}
