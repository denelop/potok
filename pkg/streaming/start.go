package streaming

import (
	"context"
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
		err = stream.Validate()
		if err != nil {
			return err
		}
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
		SubLogger()

	log.Info("Starting").
		Log()

	var args []string

	// input stream
	if stream.RTSPTransport != "" {
		args = append(args, "-rtsp_transport", stream.RTSPTransport)
	}
	args = append(args, "-i", stream.URL)

	// watermark and scaling
	if config.WatermarkFile != "" {
		args = append(args, "-i", config.WatermarkFile.AbsPath())
		args = append(args, "-filter_complex", strings.ReplaceAll(`
[1]lut=a=val*0.3[opacity];
[0][opacity]overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2,scale=-1:720
`, "\n", ""))
	} else {
		args = append(args, "-filter_complex", "[0]scale=-1:720")
	}

	// hls
	args = append(args,
		"-f", "hls",
		"-hls_time", "1",
		"-hls_list_size", "3",
		"-hls_flags", "delete_segments+append_list",
	)
	args = append(args, "-preset", "veryfast")

	// output
	streamDir := config.Dir.Join(stream.Name)
	err := streamDir.MakeAllDirs()
	if err != nil {
		return err
	}
	err = streamDir.RemoveDirContentsRecursiveContext(ctx)
	if err != nil {
		return err
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
