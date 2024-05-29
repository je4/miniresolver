package main

import (
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	loaderConfig "github.com/je4/trustutil/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/zLogger"
	"io/fs"
	"os"
)

type MiniResolverConfig struct {
	LocalAddr          string                 `toml:"localaddr"`
	TLS                loaderConfig.TLSConfig `toml:"tls"`
	LogFile            string                 `toml:"logfile"`
	LogLevel           string                 `toml:"loglevel"`
	ServiceExpiration  config.Duration        `toml:"serviceExpiration"`
	NotFoundExpiration config.Duration        `toml:"notFoundExpiration"`
	BufferSize         int                    `toml:"bufferSize"`
	Log                zLogger.Config         `toml:"log"`
}

func LoadMiniResolverConfig(fSys fs.FS, fp string, conf *MiniResolverConfig) error {
	if _, err := fs.Stat(fSys, fp); err != nil {
		path, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "cannot get current working directory")
		}
		fSys = os.DirFS(path)
		fp = "miniresolver.toml"
	}
	data, err := fs.ReadFile(fSys, fp)
	if err != nil {
		return errors.Wrapf(err, "cannot read file [%v] %s", fSys, fp)
	}
	_, err = toml.Decode(string(data), conf)
	if err != nil {
		return errors.Wrapf(err, "error loading config file %v", fp)
	}
	return nil
}
