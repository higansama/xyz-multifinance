package logger

import (
	"sync"
	"time"

	"github.com/higansama/xyz-multi-finance/config"
	ierrors "github.com/higansama/xyz-multi-finance/internal/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var singleton sync.Once

func InitLogger(cfg config.Config) {
	singleton.Do(func() {
		if cfg.App.Debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}

		zerolog.ErrorStackMarshaler = ierrors.MarshalStack
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		output := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.RFC3339
		})

		logger := log.With().Stack().Logger().Output(output)
		//if cfg.App.NoLog {
		//	logger = zerolog.Nop()
		//}

		log.Logger = logger
	})
}
