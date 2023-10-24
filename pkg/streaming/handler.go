package streaming

import (
	"net/http"

	"github.com/domonda/golog"
	"github.com/gorilla/mux"
	"github.com/ungerik/go-httpx/httperr"
)

func Handler(res http.ResponseWriter, req *http.Request) {
	requestID := golog.GetOrCreateRequestUUID(req)
	res.Header().Set("X-Request-ID", golog.FormatUUID(requestID))
	req = golog.RequestWithAttribs(req, golog.UUID{Key: "requestID", Val: requestID})

	ctx := req.Context()

	vars := mux.Vars(req)
	streamName := vars["streamName"]
	file := vars["file"]
	if streamName == "" || file == "" {
		httperr.NotFound.ServeHTTP(res, req)
		return
	}

	log.InfoCtx(ctx, "Handling").
		Request(req).
		Str("streamName", streamName).
		Str("file", file).
		Log()

	for _, stream := range streams {
		if stream.Name == streamName {
			file := config.Dir.Join(streamName, file)
			if !file.Exists() {
				httperr.NotFound.ServeHTTP(res, req)
				return
			}

			switch file.Ext() {
			case ".m3u8":
				res.Header().Add("content-type", "application/vnd.apple.mpegurl")
			case ".ts":
				res.Header().Add("content-type", "video/mp2t")
			case ".m4s":
				res.Header().Add("content-type", "video/mp4")
			case ".m4a":
				res.Header().Add("content-type", "audio/mp4")
			case ".m4v":
				res.Header().Add("content-type", "video/mp4")
			case ".mp4":
				res.Header().Add("content-type", "video/mp4")
			default:
				httperr.NotFound.ServeHTTP(res, req)
				return
			}

			file.WriteTo(res)
			return
		}
	}

	httperr.NotFound.ServeHTTP(res, req)
}
