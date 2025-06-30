package activities

import (
	"sync"
)

type DBActivities struct {
	connections sync.Map
}

func InitDBActivities() *DBActivities {
	return &DBActivities{}
}
