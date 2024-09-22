package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	customerentrypoint "github.com/higansama/xyz-multi-finance/customer/entrypoint"

	"github.com/higansama/xyz-multi-finance/config"
	"github.com/higansama/xyz-multi-finance/internal/infrastructure"
	"github.com/higansama/xyz-multi-finance/internal/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	err := config.InitConfig("")
	if err != nil {
		log.Panic().Err(err).Send()
	}
	logger.InitLogger(config.Cfg)

	infra := infrastructure.NewInfrastructure(config.Cfg)
	infra, err, infraCleanup := infra.InitInfrastructure(ctx)
	if err != nil {
		log.Panic().Err(err).Send()
	}
	defer infraCleanup()

	if !config.Cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	// engine.HTMLRender = template.HTMLTemplate{}

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowCredentials = true
	corsCfg.AllowHeaders = []string{"*"}
	engine.Use(cors.New(corsCfg))

	// Profiling tools
	if config.Cfg.App.Debug {
		//pprof.Register(engine)
	}

	err = customerentrypoint.RegisterModuleCustomer(infra, engine)
	if err != nil {
		panic(err)
	}

	go func() {
		// Run server on separate go routine for Go < 1.18 to make sure another
		// deffered func in main working.
		err = engine.Run(config.Cfg.App.Host + ":" + config.Cfg.App.Port)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	for {
		select {
		case err = <-infra.ErrorCh:
			log.Panic().Err(err).Send()
			return
		case <-ctx.Done():
			log.Info().Msg("Server exiting")
			return
		}
	}

}
