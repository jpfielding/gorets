package gorets_metadata

import (
	"time"
)

type MUpdateType struct {
	Version                 string
	Date                    time.Time
	Resource, Class, Update string
	UpdateTypes             []UpdateType
}

type UpdateType struct {
}
