package log

import (
	"fmt"

	"github.com/chadwpetersen/ind/alert"
)

// Option is a log option.
type Option func(l *logger)

// WithAlert instructs the logger to alert.
func WithAlert() Option {
	return func(l *logger) {
		// Skip if already alerted.
		if l.deferFuncs["alert"] != nil {
			return
		}

		l.deferFuncs["alert"] = func() {
			alert.Say(l.Message)
		}
	}
}

// WithLabels adds label values to the logger output.
func WithLabels(labels map[string]any) Option {
	return func(l *logger) {
		l.Labels = labels
	}
}

// withError adds an error as a label to the logger output
// and creates a detailed alert.
//
// This is an internal function and should not be used
// outside the log package.
func withError(err error) Option {
	return func(l *logger) {
		if l.Labels != nil {
			l.Labels["err"] = fmt.Sprintf("%v", err)
			return
		}

		l.Labels = map[string]any{
			"err": fmt.Sprintf("%v", err),
		}

		l.deferFuncs["alert"] = func() {
			alert.Say(fmt.Sprintf("%s %v", l.Message, err))
		}
	}
}
