package date

import (
	"strings"
	"time"
)

const HoursInADay = 24

func GetDateByOffsetDay(offset time.Duration, t time.Time) string {
	timeStr := t.Add(offset * HoursInADay * time.Hour).String()
	return strings.Split(timeStr, ".")[0]
}
