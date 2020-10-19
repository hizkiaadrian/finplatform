package dateutils

import (
	"strings"
	"time"
)

type DateOnly struct {
	time.Time
}

func (t *DateOnly) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006-01-02", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}
