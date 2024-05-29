package main

import (
	"flag"
	"fmt"
	"github.com/je4/miniresolver/v2/configs"
	pb "github.com/je4/miniresolver/v2/pkg/miniresolverproto"
	"github.com/je4/miniresolver/v2/pkg/service"
	"github.com/je4/trustutil/v2/pkg/grpchelper"
	"github.com/je4/trustutil/v2/pkg/loader"
	"github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/rs/zerolog"
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
		ServiceExpiration:  config.Duration(10 * time.Minute),
		NotFoundExpiration: config.Duration(10 * time.Second),
	}
	if err := LoadMiniResolverConfig(cfgFS, cfgFile, conf); err != nil {
		log.Fatalf("cannot load toml from [%v] %s: %v", cfgFS, cfgFile, err)
	}
	// create logger instance
	var out io.Writer = os.Stdout
	if conf.LogFile != "" {
		fp, err := os.OpenFile(conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("cannot open logfile %s: %v", conf.LogFile, err)
		}
		defer fp.Close()
		out = fp
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("cannot get hostname: %v", err)
	}

	output := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
	_logger := zerolog.New(output).With().Timestamp().Str("service", "miniresolver"). /*.Array("addrs", zLogger.StringArray(addrStr))*/ Str("host", hostname).Str("addr", conf.LocalAddr).Logger()
	_logger.Level(zLogger.LogLevel(conf.LogLevel))
	var logger zLogger.ZLogger = &_logger

	srv := service.NewMiniResolver(conf.BufferSize, time.Duration(conf.ServiceExpiration), logger)

	tlsConfig, l, err := loader.CreateServerLoader(true, &conf.TLS, nil, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server loader")
	}
	defer l.Close()

	grpcServer, err := grpchelper.NewServer(conf.LocalAddr, tlsConfig, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server")
	}
	pb.RegisterMiniResolverServer(grpcServer, srv)

	grpcServer.Startup()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	fmt.Println("press ctrl+c to stop server")
	s := <-done
	fmt.Println("got signal:", s)

	defer grpcServer.Shutdown()

}
