package sl

import "log/slog"

// Err implements prettier call of err msg
// before: logger.Error("error msg", slog.String("error", err.Error()))
// after:  logger.Error("error msg", sl.Err(err))
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
