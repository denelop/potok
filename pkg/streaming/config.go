package streaming

import (
	"github.com/denelop/potok/pkg/env"
	rootlog "github.com/domonda/golog/log"
	"github.com/ungerik/go-fs"
)

const HTTP_PATH_STREAM_NAME_PARAM = "streamName"
const HTTP_PATH_FILE_PARAM = "file"
const HTTP_BASE_PATH = "/streaming"
const HTTP_PATH = HTTP_BASE_PATH + "/{" + HTTP_PATH_STREAM_NAME_PARAM + "}/{" + HTTP_PATH_FILE_PARAM + "}"

const MASTER_PLAYLIST_NAME = "playlist.m3u8"

var (
	log = rootlog.NewPackageLogger("stream")

	config struct {
		StreamFile            fs.File `env:"STREAMING_STREAMS_FILE,required"`
		Dir                   fs.File `env:"STREAMING_DIR,required"`
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
