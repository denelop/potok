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
	if config.WatermarkFile != "" {
		args = append(args, "-i", config.WatermarkFile.AbsPath())
		filterComplex.WriteString(strings.ReplaceAll(`
[1]lut=a=val*0.3[opacity];
[0][opacity]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2
`, "\n", ""))
	}
	if stream.Scale != "" {
		if filterComplex.Len() > 0 {
			filterComplex.WriteString(",")
		}
		filterComplex.WriteString("scale=" + stream.Scale)
	}
	args = append(args, "-filter_complex", filterComplex.String())

	// hls
	args = append(args,
		"-maxrate", "500k", // will create an #EXT-X-STREAM-INF entry in the `MASTER_PLAYLIST_NAME`
		"-f", "hls",
		"-hls_flags", "delete_segments+append_list",
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
