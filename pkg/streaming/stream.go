//go:generate go-enum $GOFILE

package streaming

import (
	"fmt"

	"github.com/domonda/go-errs"
)

type RTSPTransport string //#enum

const (
	RTSPTransportTCP RTSPTransport = "tcp"
	RTSPTransportUDP RTSPTransport = "udp"
)

// Valid indicates if r is any of the valid values for RTSPTransport
func (r RTSPTransport) Valid() bool {
	switch r {
	case
		RTSPTransportTCP,
		RTSPTransportUDP:
		return true
	}
	return false
}

// Validate returns an error if r is none of the valid values for RTSPTransport
func (r RTSPTransport) Validate() error {
	if !r.Valid() {
		return fmt.Errorf("invalid value %#v for type streaming.RTSPTransport", r)
	}
	return nil
}

// String implements the fmt.Stringer interface for RTSPTransport
func (r RTSPTransport) String() string {
	return string(r)
}

type Stream struct {
	Name          string        `yaml:"name"`
	URL           string        `yaml:"url"`
	RTSPTransport RTSPTransport `yaml:"rtsp_transport"`
	Scale         string        `yaml:"scale"`
}

func (stream *Stream) Validate() error {
	if stream.Name == "" {
		return errs.New("missing stream name")
	}
	if stream.URL == "" {
		return errs.New("missing stream URL")
	}
	if err := stream.RTSPTransport.Validate(); err != nil {
		return err
	}
	return nil
}

func (stream *Stream) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type WithDefaultsStream Stream // we create a new temp type to avoid infinite recursion
	withDefaultsStream := WithDefaultsStream{
		RTSPTransport: RTSPTransportTCP,
	}
	if err := unmarshal(&withDefaultsStream); err != nil {
		return err
	}
	*stream = Stream(withDefaultsStream)
	if err := stream.Validate(); err != nil {
		return err
	}
	return nil
}
