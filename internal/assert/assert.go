package assert

import (
	"reflect"
	"runtime"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

func Never(message string, fields ...map[string]interface{}) {
	write(message, fields...)
}

func NotNil(item interface{}, message string, fields ...map[string]interface{}) {
	if item == nil || reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).IsNil() {
		typeFields := map[string]interface{}{
			"expected_type": reflect.TypeOf(item),
		}
		fields = append(fields, typeFields)
		write(message, fields...)
	}
}

func Assert(condition bool, message string, fields ...map[string]interface{}) {
	if !condition {
		write(message, fields...)
	}
}

func NoError(err error, message string, fields ...map[string]interface{}) {
	if err != nil {
		errorFields := map[string]interface{}{
			"error":         err,
			"error_message": err.Error(),
		}
		fields = append(fields, errorFields)
		write(message, fields...)
	}
}

func write(message string, fields ...map[string]interface{}) {
	event := log.Fatal()
	pc, file, line, _ := runtime.Caller(1)

	callerFields := map[string]interface{}{
		"caller_file":     file,
		"caller_line":     line,
		"caller_function": runtime.FuncForPC(pc).Name(),
	}
	fields = append(fields, callerFields)
	for _, field := range fields {
		for key, value := range field {
			event = event.Interface(key, value)
		}
	}

	stackTrace := string(debug.Stack())
	event = event.Str("stack_trace", stackTrace)

	event.Msg(message)
}
