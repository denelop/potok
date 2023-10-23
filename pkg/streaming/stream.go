package streaming

import "github.com/domonda/go-errs"

type Stream struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func (stream *Stream) Validate() error {
	if stream.Name == "" {
		return errs.New("missing stream name")
	}
	if stream.URL == "" {
		return errs.New("missing stream URL")
	}
	return nil
}
