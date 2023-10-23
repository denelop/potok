package stream

import (
	"context"
)

func GetHLSBytes(ctx context.Context, stream string) (hls []byte, err error) {
	log.InfoCtx(ctx, "Getting HLS stream bytes").
		Str("stream", stream).
		Log()

	return []byte{}, nil
}
