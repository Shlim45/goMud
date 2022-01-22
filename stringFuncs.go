package main

import (
	"strings"
	"time"
)

// before - given a string (value) and a substring(a), return substring before a string.
func before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

// after - given a string (value) and a substring(a), return substring after a string.
func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

// TimeString - given a time, return the MySQL standard string representation
func TimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999999")
}

func MudASCIILogo() string {
	var logo strings.Builder
	logo.WriteString("\r\n\r\n")
	logo.WriteString("          _____                    _____                    _____\r\n")
	logo.WriteString("         /\\    \\                  /\\    \\                  /\\    \\ \r\n")
	logo.WriteString("        /::\\    \\                /::\\    \\                /::\\    \\ \r\n")
	logo.WriteString("       /::::\\    \\              /::::\\    \\              /::::\\    \\ \r\n")
	logo.WriteString("      /::::::\\    \\            /::::::\\    \\            /::::::\\    \\ \r\n")
	logo.WriteString("     /:::/\\:::\\    \\          /:::/\\:::\\    \\          /:::/\\:::\\    \\ \r\n")
	logo.WriteString("    /:::/  \\:::\\    \\        /:::/__\\:::\\    \\        /:::/  \\:::\\    \\ \r\n")
	logo.WriteString("   /:::/    \\:::\\    \\      /::::\\   \\:::\\    \\      /:::/    \\:::\\    \\ \r\n")
	logo.WriteString("  /:::/    / \\:::\\    \\    /::::::\\   \\:::\\    \\    /:::/    / \\:::\\    \\ \r\n")
	logo.WriteString(" /:::/    /   \\:::\\ ___\\  /:::/\\:::\\   \\:::\\    \\  /:::/    /   \\:::\\    \\ \r\n")
	logo.WriteString("/:::/____/     \\:::|    |/:::/  \\:::\\   \\:::\\____\\/:::/____/     \\:::\\____\\ \r\n")
	logo.WriteString("\\:::\\    \\     /:::|____|\\::/    \\:::\\   \\::/    /\\:::\\    \\      \\::/    /\r\n")
	logo.WriteString(" \\:::\\    \\   /:::/    /  \\/____/ \\:::\\   \\/____/  \\:::\\    \\      \\/____/\r\n")
	logo.WriteString("  \\:::\\    \\ /:::/    /            \\:::\\    \\       \\:::\\    \\ \r\n")
	logo.WriteString("   \\:::\\    /:::/    /              \\:::\\____\\       \\:::\\    \\ \r\n")
	logo.WriteString("    \\:::\\  /:::/    /                \\::/    /        \\:::\\    \\ \r\n")
	logo.WriteString("     \\:::\\/:::/    /                  \\/____/          \\:::\\    \\ \r\n")
	logo.WriteString("      \\::::::/    /                                     \\:::\\    \\ \r\n")
	logo.WriteString("       \\::::/    /                                       \\:::\\____\\ \r\n")
	logo.WriteString("        \\::/____/                                         \\::/    /\r\n")
	logo.WriteString("         ¯¯                                                \\/____/\r\n")
	logo.WriteString("\r\n\r\n")
	return Green(logo.String())
}
