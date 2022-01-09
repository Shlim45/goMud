package main

import "fmt"

var (
	CFriend    = Green
	CEnemy     = Red
	CItem      = BrightYellow
	CExit      = Yellow
	CDamageIn  = BrightRed
	CDamageOut = Green
	CHealing   = BrightGreen
	CArea      = Red
	CSigil     = Cyan
	CNormal    = White
)

var (
	Black         = Color("\033[1;30m%s\033[0m")
	Red           = Color("\033[1;31m%s\033[0m")
	Green         = Color("\033[1;32m%s\033[0m")
	Yellow        = Color("\033[1;33m%s\033[0m")
	Blue          = Color("\033[1;34m%s\033[0m")
	Magenta       = Color("\033[1;35m%s\033[0m")
	Cyan          = Color("\033[1;36m%s\033[0m")
	White         = Color("\033[1;37m%s\033[0m")
	Grey          = Color("\033[1;90m%s\033[0m")
	BrightRed     = Color("\033[1;91m%s\033[0m")
	BrightGreen   = Color("\033[1;92m%s\033[0m")
	BrightYellow  = Color("\033[1;93m%s\033[0m")
	BrightBlue    = Color("\033[1;94m%s\033[0m")
	BrightMagenta = Color("\033[1;95m%s\033[0m")
	BrightCyan    = Color("\033[1;96m%s\033[0m")
	BrightWhite   = Color("\033[1;97m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
