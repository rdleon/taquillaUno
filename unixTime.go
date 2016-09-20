package main

import (
	"fmt"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	var stamp string

	stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format(time.UnixDate))
	bs := []byte(stamp)

	return bs, nil
}

func (t JSONTime) UnmarshalJSON(data []byte) (err error) {
	var u time.Time

	u, err = time.Parse(time.UnixDate, string(data))
	t = JSONTime(u)

	return
}
