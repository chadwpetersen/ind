package log

import (
	"encoding/json"
	"fmt"
	"log"
)

type logger struct {
	Level   Level
	Message string
	Labels  map[string]any

	deferFuncs map[string]func()
}

func newLogger(level Level, msg string, opts ...Option) *logger {
	l := &logger{
		Level:      level,
		Message:    msg,
		deferFuncs: make(map[string]func()),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *logger) Print(v ...any) {
	defer func() {
		for _, fn := range l.deferFuncs {
			fn()
		}
	}()

	if !debug && l.Level == LevelDebug {
		return
	}

	log.Print(l.Level.Colour(l.Level.String()), ": ", l.msg(v...), "\n")
}

func (l *logger) msg(v ...any) string {
	if l.Labels == nil {
		return fmt.Sprint(v...)
	}

	if !verbose && l.Level != LevelDebug {
		return fmt.Sprint(v...)
	}

	labels, err := json.MarshalIndent(l.Labels, "", "\t")
	if err != nil {
		return fmt.Sprintf("%s\n%v", fmt.Sprint(v...), l.Labels)
	}

	return fmt.Sprintf("%s\n%v", fmt.Sprint(v...), string(labels))
}
