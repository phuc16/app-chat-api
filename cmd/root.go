package cmd

import (
	"app/build"
	"app/config"
	"app/errors"
	"app/http"
	"app/pkg/apperror"
	"app/pkg/logger"
	"app/pkg/mongodb"
	"app/pkg/trace"
	"app/repository"
	"app/service"
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			errors.WrapError(&err)
			log.Println(err.(*apperror.Error).StackTrace())
		}()

		ctx := context.Background()

		err = config.Load()
		if err != nil {
			return
		}
		logger.InitLogger(config.Cfg.Logger.ToLoggerConfig())

		logger.For(ctx).Infof("Application info %s", build.Info().String())

		otelShutdown, err := trace.SetupOTelSDK(ctx, config.Cfg.OTel.ToTraceConfig())
		if err != nil {
			return
		}
		defer func() {
			if err := otelShutdown(context.Background()); err != nil {
				logger.For(ctx).Errorf("Shutting down tracer provider %v", err)
			}
		}()

		conn, err := mongodb.NewMongoDBConn(ctx, config.Cfg.DB.URI)
		if err != nil {
			return
		}
		repo := repository.NewRepo(conn)
		err = repo.InitIndex(ctx)
		if err != nil {
			return
		}

		otpSvc := service.NewOtpService(repo, repo)
		userSvc := service.NewUserService(otpSvc, repo, repo)

		httpSrv := http.NewServer(userSvc, otpSvc)

		quit := make(chan error)
		go func() {
			err := httpSrv.Start()
			if err != nil {
				quit <- err
			}
		}()

		err = <-quit
		logger.For(ctx).Errorf("Shutting down %v", err)

		return
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
