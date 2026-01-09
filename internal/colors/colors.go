package colors

import "fmt"

type TextDecoration string
type ForegroundColor string
type BackgroundColor string

const (
	Reset         TextDecoration = "\033[0m"
	Bold          TextDecoration = "\033[01m"
	Disable       TextDecoration = "\033[02m"
	Underline     TextDecoration = "\033[04m"
	Reverse       TextDecoration = "\033[07m"
	Strikethrough TextDecoration = "\033[09m"
	Invisible     TextDecoration = "\033[08m"
)

const (
	FgBlack      ForegroundColor = "\033[30m"
	FgRed        ForegroundColor = "\033[31m"
	FgGreen      ForegroundColor = "\033[32m"
	FgOrange     ForegroundColor = "\033[33m"
	FgBlue       ForegroundColor = "\033[34m"
	FgPurple     ForegroundColor = "\033[35m"
	FgCyan       ForegroundColor = "\033[36m"
	FgLightGrey  ForegroundColor = "\033[37m"
	FgDarkGrey   ForegroundColor = "\033[90m"
	FgLightRed   ForegroundColor = "\033[91m"
	FgLightGreen ForegroundColor = "\033[92m"
	FgYellow     ForegroundColor = "\033[93m"
	FgLightBlue  ForegroundColor = "\033[94m"
	FgPink       ForegroundColor = "\033[95m"
	FgLightCyan  ForegroundColor = "\033[96m"
)

const (
	BgBlack     BackgroundColor = "\033[40m"
	BgRed       BackgroundColor = "\033[41m"
	BgGreen     BackgroundColor = "\033[42m"
	BgOrange    BackgroundColor = "\033[43m"
	BgBlue      BackgroundColor = "\033[44m"
	BgPurple    BackgroundColor = "\033[45m"
	BgCyan      BackgroundColor = "\033[46m"
	BgLightGrey BackgroundColor = "\033[47m"
)

func Yellow(prompt string) string {
	return fmt.Sprintf("%s%s%s", FgYellow, prompt, Reset)
}
