package routes

//go:generate gen-func-wrappers -replaceForJSON=fs.FileReader:fs.File $GOFILE

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/denelop/potok/pkg/streaming"

	"github.com/domonda/go-function"
)

// streamingGetHLSBytes wraps streaming.GetHLSBytes as function.Wrapper (generated code)
var streamingGetHLSBytes streamingGetHLSBytesT

// streamingGetHLSBytesT wraps streaming.GetHLSBytes as function.Wrapper (generated code)
type streamingGetHLSBytesT struct{}

func (streamingGetHLSBytesT) String() string {
	return "streaming.GetHLSBytes(ctx context.Context, streaming string) (hlsStream []byte, err error)"
}

func (streamingGetHLSBytesT) Name() string {
	return "GetHLSBytes"
}

func (streamingGetHLSBytesT) NumArgs() int      { return 2 }
func (streamingGetHLSBytesT) ContextArg() bool  { return true }
func (streamingGetHLSBytesT) NumResults() int   { return 2 }
func (streamingGetHLSBytesT) ErrorResult() bool { return true }

func (streamingGetHLSBytesT) ArgNames() []string {
	return []string{"ctx", "streaming"}
}

func (streamingGetHLSBytesT) ArgDescriptions() []string {
	return []string{"", ""}
}

func (streamingGetHLSBytesT) ArgTypes() []reflect.Type {
	return []reflect.Type{
		function.ReflectType[context.Context](),
		function.ReflectType[string](),
	}
}

func (streamingGetHLSBytesT) ResultTypes() []reflect.Type {
	return []reflect.Type{
		function.ReflectType[[]byte](),
		function.ReflectType[error](),
	}
}

func (streamingGetHLSBytesT) Call(ctx context.Context, args []any) (results []any, err error) {
	results = make([]any, 1)
	results[0], err = streaming.GetHLSBytes(ctx, args[0].(string)) // wrapped call
	return results, err
}

func (streamingGetHLSBytesT) CallWithStrings(ctx context.Context, strs ...string) (results []any, err error) {
	var a struct {
		streaming string
	}
	if 0 < len(strs) {
		a.streaming = strs[0]
	}
	results = make([]any, 1)
	results[0], err = streaming.GetHLSBytes(ctx, a.streaming) // wrapped call
	return results, err
}

func (streamingGetHLSBytesT) CallWithNamedStrings(ctx context.Context, strs map[string]string) (results []any, err error) {
	var a struct {
		streaming string
	}
	if str, ok := strs["streaming"]; ok {
		a.streaming = str
	}
	results = make([]any, 1)
	results[0], err = streaming.GetHLSBytes(ctx, a.streaming) // wrapped call
	return results, err
}

func (f streamingGetHLSBytesT) CallWithJSON(ctx context.Context, argsJSON []byte) (results []any, err error) {
	var a struct {
		Stream string
	}
	err = json.Unmarshal(argsJSON, &a)
	if err != nil {
		return nil, function.NewErrParseArgsJSON(err, f, argsJSON)
	}
	results = make([]any, 1)
	results[0], err = streaming.GetHLSBytes(ctx, a.Stream) // wrapped call
	return results, err
}
