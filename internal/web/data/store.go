package data

import "context"

type Store struct {
	DbUrl string
	DriverName string
	Context context.Context
}
