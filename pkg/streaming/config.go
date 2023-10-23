package streaming

import (
	"github.com/denelop/potok/pkg/env"
	rootlog "github.com/domonda/golog/log"
	"github.com/ungerik/go-fs"
)

var (
	log = rootlog.NewPackageLogger("stream")

	config struct {
		StreamFile            fs.File `env:"STREAMING_STREAMS_FILE,required"`
		Dir                   fs.File `env:"STREAMING_DIR,required"`
		WatermarkFile         fs.File `env:"STREAMING_WATERMARK_FILE"`
		DeleteContentsOnStart bool    `env:"STREAMING_DELETE_CONTENTS_ON_START"`
	}
)

func init() {
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	log.Debug("Configured").
		StructFields(config).
		Log()
}
