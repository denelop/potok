package routes

//go:generate gen-func-wrappers -replaceForJSON=fs.FileReader:fs.File $GOFILE

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/denelop/potok/pkg/stream"

	"github.com/domonda/go-function"
)

// streamGetHLSBytes wraps stream.GetHLSBytes as function.Wrapper (generated code)
var streamGetHLSBytes streamGetHLSBytesT

// streamGetHLSBytesT wraps stream.GetHLSBytes as function.Wrapper (generated code)
type streamGetHLSBytesT struct{}

func (streamGetHLSBytesT) String() string {
	return "stream.GetHLSBytes(ctx context.Context, stream string) (hlsStream []byte, err error)"
}

func (streamGetHLSBytesT) Name() string {
	return "GetHLSBytes"
}

func (streamGetHLSBytesT) NumArgs() int      { return 2 }
func (streamGetHLSBytesT) ContextArg() bool  { return true }
func (streamGetHLSBytesT) NumResults() int   { return 2 }
func (streamGetHLSBytesT) ErrorResult() bool { return true }

func (streamGetHLSBytesT) ArgNames() []string {
	return []string{"ctx", "stream"}
}

func (streamGetHLSBytesT) ArgDescriptions() []string {
	return []string{"", ""}
}

func (streamGetHLSBytesT) ArgTypes() []reflect.Type {
	return []reflect.Type{
		function.ReflectType[context.Context](),
		function.ReflectType[string](),
	}
}

func (streamGetHLSBytesT) ResultTypes() []reflect.Type {
	return []reflect.Type{
		function.ReflectType[[]byte](),
		function.ReflectType[error](),
	}
}

func (streamGetHLSBytesT) Call(ctx context.Context, args []any) (results []any, err error) {
	results = make([]any, 1)
	results[0], err = stream.GetHLSBytes(ctx, args[0].(string)) // wrapped call
	return results, err
}

func (streamGetHLSBytesT) CallWithStrings(ctx context.Context, strs ...string) (results []any, err error) {
	var a struct {
		stream string
	}
	if 0 < len(strs) {
		a.stream = strs[0]
	}
	results = make([]any, 1)
	results[0], err = stream.GetHLSBytes(ctx, a.stream) // wrapped call
	return results, err
}

func (streamGetHLSBytesT) CallWithNamedStrings(ctx context.Context, strs map[string]string) (results []any, err error) {
	var a struct {
		stream string
	}
	if str, ok := strs["stream"]; ok {
		a.stream = str
	}
	results = make([]any, 1)
	results[0], err = stream.GetHLSBytes(ctx, a.stream) // wrapped call
	return results, err
}

func (f streamGetHLSBytesT) CallWithJSON(ctx context.Context, argsJSON []byte) (results []any, err error) {
	var a struct {
		Stream string
	}
	err = json.Unmarshal(argsJSON, &a)
	if err != nil {
		return nil, function.NewErrParseArgsJSON(err, f, argsJSON)
	}
	results = make([]any, 1)
	results[0], err = stream.GetHLSBytes(ctx, a.Stream) // wrapped call
	return results, err
}
