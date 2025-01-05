package log

const (
	// Colors
	cyan    = "\x1b[36m"
	yellow  = "\x1b[33m"
	green   = "\x1b[32m"
	magenta = "\x1b[35m"
	reset   = "\x1b[0m"

	// Styles
	bold = "\x1b[1m"
)

// Format text with color and optional bold
func Color(text string, color string, isBold bool) string {
	if isBold {
		return bold + color + text + reset
	}
	return color + text + reset
}

// Predefined formatters
func Title(text string) string {
	return bold + cyan + text + reset
}

func Section(text string) string {
	return bold + yellow + text + reset
}

func Command(text string) string {
	return cyan + text + reset
}

func Description(text string) string {
	return green + text + reset
}

func Highlight(text string) string {
	return magenta + text + reset
}
