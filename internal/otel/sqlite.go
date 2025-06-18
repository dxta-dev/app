package otel

import (
	"github.com/XSAM/otelsql"
	_ "github.com/libsql/libsql-client-go/libsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func GetDriverName() string {
	driverName, err := otelsql.Register("libsql", otelsql.WithAttributes(
		semconv.ServiceName("turso"),
	))

	if err != nil {
		driverName = "libsql"
	}

	return driverName
}
