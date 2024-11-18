package assert

import (
	"reflect"

	"github.com/rs/zerolog/log"
)

func Never(message string, fields ...map[string]interface{}) {
	write(message, fields...)
}

func NotNil(item interface{}, message string, fields ...map[string]interface{}) {
	if item == nil || reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).IsNil() {
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
			err.Error(): err,
		}
		fields = append(fields, errorFields)
		write(message, fields...)
	}
}

func write(message string, fields ...map[string]interface{}) {
	event := log.Fatal()
	event.Msg(message)
	for _, field := range fields {
		for key, value := range field {
			event = event.Interface(key, value)
		}
	}
	event.Send()

	panic(message)
}
