package streaming

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/domonda/go-errs"
	"gopkg.in/yaml.v3"
)

var streams []*Stream

// StartAll starts all streams from the streams file.
func StartAll(ctx context.Context) (err error) {
	log := log.With().Ctx(ctx).SubLogger()

	log.Info("Starting all").
		Log()

	fileBytes, err := config.StreamFile.ReadAllContext(ctx)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fileBytes, &streams)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		err = Start(ctx, stream)
		if err != nil {
			return err
		}
	}

	return nil
}

func Start(ctx context.Context, stream *Stream) error {
	log := log.With().Ctx(ctx).
		StructFields(stream).
		Str("httpMasterPlaylistPath", fmt.Sprintf("%s/%s/%s", HTTP_BASE_PATH, url.PathEscape(stream.Name), MASTER_PLAYLIST_NAME)).
		SubLogger()

	log.Info("Starting").
		Log()

	var args []string

	// input stream
	args = append(args, "-rtsp_transport", string(stream.RTSPTransport))
	args = append(args, "-i", stream.URL)

	// filters
	filterComplex := new(bytes.Buffer)
	{
		// watermarks
		{
			// inputs
			for _, watermark := range stream.Watermarks {
				if !watermark.File.Exists() {
					return errs.Errorf("watermark file %q does not exist", watermark.File.Path())
				}
				args = append(args, "-i", watermark.File.Path())
			}

			// apply opacity to watermarks
			for i, watermark := range stream.Watermarks {
				filterComplex.WriteString(
					fmt.Sprintf("[%[1]d]lut=a=val*%.1f[%[1]dwm];",
						i+1, // +1 because stream is first input
						watermark.Opacity,
					),
				)
			}

			// apply watermarks to stream
			for i, watermark := range stream.Watermarks {
				var prevBg string
				if i == 0 {
					prevBg = "[0]"
				} else {
					prevBg = fmt.Sprintf("[bg%d]", i)
				}

				var bg string
				if i < len(stream.Watermarks)-1 {
					bg = fmt.Sprintf("[bg%d];", i+1)
				}

				filterComplex.WriteString(
					fmt.Sprintf("%[2]s[%[1]dwm]overlay=%[3]s%[4]s",
						i+1, // +1 because stream is first input
						prevBg,
						watermark.Position.FilterComplexWatermarkPosition(),
						bg,
					),
				)
			}
		}

		// scale
		{
			if stream.Scale != "" {
				if filterComplex.Len() > 0 {
					// if there's stuff in the filter complex, watermarks are applied - append the scale
					filterComplex.WriteString(",scale=" + stream.Scale)
				} else {
					// otherwise the scale is the only filter
					filterComplex.WriteString("scale=" + stream.Scale)
				}
			}
		}
	}

	args = append(args, "-filter_complex", strings.ReplaceAll(strings.ReplaceAll(filterComplex.String(), "\n", ""), "\t", ""))

	// hls
	args = append(args,
		"-maxrate", "500k", // will create an #EXT-X-STREAM-INF entry in the `MASTER_PLAYLIST_NAME`
		"-f", "hls",
		"-hls_flags", "delete_segments",
		"-master_pl_name", MASTER_PLAYLIST_NAME,
	)

	// output
	args = append(args, "-preset", "veryfast")
	streamDir := config.Dir.Join(stream.Name)
	err := streamDir.MakeAllDirs()
	if err != nil {
		return err
	}
	if config.DeleteContentsOnStart {
		err = streamDir.RemoveDirContentsRecursiveContext(ctx)
		if err != nil {
			return err
		}
	}

	streamOut := streamDir.Join("stream.out") // TODO: change out file if one exists
	streamFile := streamDir.Join("stream.m3u8")
	args = append(args, streamFile.AbsPath())

	log.Debug("ffmpeg").
		Str("args", strings.Join(args, " ")).
		Str("out", streamOut.AbsPath()).
		Str("stream", streamFile.AbsPath()).
		Log()

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	out, err := os.Create(streamOut.AbsPath())
	if err != nil {
		return err
	}
	defer out.Close()

	cmd.Stdout = out
	cmd.Stderr = out

	if err = cmd.Start(); err != nil {
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Error("ffmpeg error").
				Err(err).
				Str("out", streamOut.AbsPath()).
				Log()
		} else {
			log.Info("ffmpeg done").
				Str("out", streamOut.AbsPath()).
				Log()
		}

		log.Info("Restarting after 3 seconds").
			Log()
		time.Sleep(time.Second * 3)

		err = Start(ctx, stream)
		if err != nil {
			log.Error("Restart error (will not restart again)").
				Err(err).
				Log()
		}
	}()

	return nil
}
