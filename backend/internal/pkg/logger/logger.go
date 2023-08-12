package logger

import (
	"os"
	"strconv"
	"time"

	"github.com/f4mk/api/config"
	"github.com/rs/zerolog"
)

func New(cfg *config.Config) *zerolog.Logger {

	var logger zerolog.Logger

	if cfg.Environment == "production" {

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		logger = zerolog.New(os.Stdout).With().Timestamp().Str("service", cfg.Service.ServiceName).Logger()
		zerolog.SetGlobalLevel(intToLogLevel(cfg.Log.LogLevel))

	} else {

		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short

			return file + ":" + strconv.Itoa(line)
		}

		zerolog.TimeFieldFormat = time.RFC3339Nano
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Caller().Logger()
		zerolog.SetGlobalLevel(intToLogLevel(cfg.Log.LogLevel))
	}

	return &logger
}

func intToLogLevel(level int) zerolog.Level {

	switch level {
	case 0:
		return zerolog.DebugLevel
	case 1:
		return zerolog.InfoLevel
	case 2:
		return zerolog.WarnLevel
	case 3:
		return zerolog.ErrorLevel
	case 4:
		return zerolog.FatalLevel
	case 5:
		return zerolog.PanicLevel
	default:
		return zerolog.NoLevel
	}
}
