//go:generate go-enum $GOFILE

package streaming

import (
	"fmt"

	"github.com/domonda/go-errs"
	"github.com/ungerik/go-fs"
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
	Name          string             `yaml:"name"`
	URL           string             `yaml:"url"`
	RTSPTransport RTSPTransport      `yaml:"rtsp_transport"`
	Scale         string             `yaml:"scale"`
	Watermarks    []*StreamWatermark `yaml:"watermarks"`
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
	for _, overlay := range stream.Watermarks {
		if err := overlay.Validate(); err != nil {
			return err
		}
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

type StreamWatermarkPosition string //#enum

const (
	StreamWatermarkPositionTopLeft     StreamWatermarkPosition = "top-left"
	StreamWatermarkPositionCenter      StreamWatermarkPosition = "center"
	StreamWatermarkPositionBottomRight StreamWatermarkPosition = "bottom-right"
)

// Valid indicates if s is any of the valid values for StreamWatermarkPosition
func (s StreamWatermarkPosition) Valid() bool {
	switch s {
	case
		StreamWatermarkPositionTopLeft,
		StreamWatermarkPositionCenter,
		StreamWatermarkPositionBottomRight:
		return true
	}
	return false
}

// Validate returns an error if s is none of the valid values for StreamWatermarkPosition
func (s StreamWatermarkPosition) Validate() error {
	if !s.Valid() {
		return fmt.Errorf("invalid value %#v for type streaming.StreamWatermarkPosition", s)
	}
	return nil
}

// String implements the fmt.Stringer interface for StreamWatermarkPosition
func (s StreamWatermarkPosition) String() string {
	return string(s)
}

func (s StreamWatermarkPosition) FilterComplexWatermarkPosition() string {
	switch s {
	case StreamWatermarkPositionTopLeft:
		return "0:0"
	case StreamWatermarkPositionCenter:
		return "(main_w-overlay_w)/2:(main_h-overlay_h)/2"
	case StreamWatermarkPositionBottomRight:
		return "main_w-overlay_w:main_h-overlay_h"
	}
	panic(fmt.Sprintf("invalid stream overlay position %s", s))
}

type StreamWatermark struct {
	File     fs.File                 `yaml:"file"`
	Opacity  float32                 `yaml:"opacity"`
	Position StreamWatermarkPosition `yaml:"position"`
}

func (overlay *StreamWatermark) Validate() error {
	if overlay.File == "" {
		return errs.New("missing stream overlay file")
	}
	if overlay.Opacity == 0 {
		return errs.New("stream overlay opacity cannot be 0 ")
	}
	if overlay.Opacity < 0 {
		return errs.New("stream overlay opacity cannot be negative")
	}
	if overlay.Opacity > 1 {
		return errs.New("stream overlay opacity cannot be greater than 1")
	}
	if err := overlay.Position.Validate(); err != nil {
		return err
	}
	return nil
}
