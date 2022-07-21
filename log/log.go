package log

// Pass sends a pass message to logs.
func Pass(message string, opts ...Option) {
	newLogger(LevelPass, message, opts...).Print(message)
}

// Info sends an info message to logs.
func Info(message string, opts ...Option) {
	newLogger(LevelInfo, message, opts...).Print(message)
}

// Warn sends a warning message to logs.
func Warn(message string, opts ...Option) {
	newLogger(LevelWarn, message, opts...).Print(message)
}

// Debug sends a debug message to logs.
func Debug(message string, opts ...Option) {
	newLogger(LevelDebug, message, opts...).Print(message)
}

// Error sends an error message to logs.
func Error(message string, err error, opts ...Option) {
	// Ensure we will always alert and add
	// the error as a label.
	opts = append(opts, withError(err))
	newLogger(LevelError, message, opts...).Print(message)
}
