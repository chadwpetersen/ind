package log

type Level string

const (
	LevelPass    Level = "PASS"
	LevelInfo    Level = "INFO"
	LevelWarn    Level = "WARNING"
	LevelDebug   Level = "DEBUG"
	LevelError   Level = "ERROR"
	levelUnknown Level = ""
)

var colours = map[Level]string{
	LevelPass:    "\033[1;32m",
	LevelInfo:    "\033[1;34m",
	LevelWarn:    "\033[1;33m",
	LevelDebug:   "\033[1;35m",
	LevelError:   "\033[1;31m",
	levelUnknown: "\033[0m",
}

func (l Level) String() string {
	_, ok := colours[l]
	if !ok {
		return ""
	}

	return string(l)
}

func (l Level) Colour(txt string) string {
	colour, ok := colours[l]
	if !ok {
		colour = colours[levelUnknown]
	}

	return colour + txt + colours[levelUnknown]
}
