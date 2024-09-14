package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/je4/certloader/v2/pkg/loader"
	"github.com/je4/miniresolver/v2/configs"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/miniresolver/v2/pkg/service"
	"github.com/je4/trustutil/v2/pkg/grpchelper"
	"github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/zLogger"
	ublogger "gitlab.switch.ch/ub-unibas/go-ublogger"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var cfg = flag.String("config", "", "location of toml configuration file")

func main() {
	flag.Parse()
	var cfgFS fs.FS
	var cfgFile string
	if *cfg != "" {
		cfgFS = os.DirFS(filepath.Dir(*cfg))
		cfgFile = filepath.Base(*cfg)
	} else {
		cfgFS = configs.ConfigFS
		cfgFile = "miniresolver.toml"
	}
	conf := &MiniResolverConfig{
		LocalAddr:          "localhost:7777",
		LogLevel:           "DEBUG",
		BufferSize:         1024,
		ServiceExpiration:  config.Duration(5 * time.Minute),
		NotFoundExpiration: config.Duration(4 * time.Second),
	}
	if err := LoadMiniResolverConfig(cfgFS, cfgFile, conf); err != nil {
		log.Fatalf("cannot load toml from [%v] %s: %v", cfgFS, cfgFile, err)
	}

	// create logger instance
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("cannot get hostname: %v", err)
	}

	var loggerTLSConfig *tls.Config
	var loggerLoader io.Closer
	if conf.Log.Stash.TLS != nil {
		loggerTLSConfig, loggerLoader, err = loader.CreateClientLoader(conf.Log.Stash.TLS, nil)
		if err != nil {
			log.Fatalf("cannot create client loader: %v", err)
		}
		defer loggerLoader.Close()
	}

	_logger, _logstash, _logfile := ublogger.CreateUbMultiLoggerTLS(conf.Log.Level, conf.Log.File,
		ublogger.SetDataset(conf.Log.Stash.Dataset),
		ublogger.SetLogStash(conf.Log.Stash.LogstashHost, conf.Log.Stash.LogstashPort, conf.Log.Stash.Namespace, conf.Log.Stash.LogstashTraceLevel),
		ublogger.SetTLS(conf.Log.Stash.TLS != nil),
		ublogger.SetTLSConfig(loggerTLSConfig),
	)
	if _logstash != nil {
		defer _logstash.Close()
	}
	if _logfile != nil {
		defer _logfile.Close()
	}

	l2 := _logger.With().Timestamp().Str("host", hostname).Str("addr", conf.LocalAddr).Logger() //.Output(output)
	var logger zLogger.ZLogger = &l2

	srv := service.NewMiniResolver(conf.BufferSize, time.Duration(conf.ServiceExpiration), conf.ProxyAddr, logger)
	defer srv.Close()

	tlsConfig, l, err := loader.CreateServerLoader(true, &conf.TLS, nil, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server loader")
	}
	defer l.Close()

	grpcServer, err := grpchelper.NewServer(conf.LocalAddr, tlsConfig, nil, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server")
	}
	addr := grpcServer.GetAddr()
	l2 = _logger.With().Str("addr", addr).Logger() //.Output(output)
	logger = &l2

	pb.RegisterMiniResolverServer(grpcServer, srv)

	grpcServer.Startup()
	if conf.ProxyAddr != "" {
		if err := srv.StartProxy(); err != nil {
			logger.Error().Err(err).Msg("cannot start proxy")
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	fmt.Println("press ctrl+c to stop server")
	s := <-done
	fmt.Println("got signal:", s)

	defer grpcServer.Shutdown()
	if conf.ProxyAddr != "" {
		defer func() {
			if err := srv.StopProxy(); err != nil {
				logger.Error().Err(err).Msg("cannot stop proxy")
			}
		}()
	}
}
