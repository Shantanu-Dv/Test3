package types

import (
	"fmt"
	"time"
)

type Timestamp time.Time

const TimestampFormat = "2006-01-02 15:04:05"

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := time.Parse(TimestampFormat, string(b[1:len(b)-1]))
	if err != nil {
		return fmt.Errorf("invalid timestamp format")
	}

	*t = Timestamp(ts)
	return nil
}
