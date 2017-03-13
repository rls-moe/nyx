package resources

import (
	"go.rls.moe/nyx/resources/snowflakes"
	"time"
)

var fountain = snowflakes.Generator{
	StartTime: time.Date(
		2017, 03, 11,
		11, 12, 29,
		0, time.UTC).Unix(),
}

func getID() (int, error) {
	id, err := fountain.NewID()
	return int(id), err
}
