package timeutil

import (
	"time"

	"github.com/labstack/echo/v4"
)

const TimezoneKey = "timezone"

func SetTZ(c echo.Context, tz string) {
	c.Set(TimezoneKey, tz)
}

func GetTZ(c echo.Context) *time.Location {
	tz := c.Get(TimezoneKey)
	switch tzStr := tz.(type) {
	case string:
		loc, err := time.LoadLocation(tzStr)
		if err != nil {
			loc = time.UTC
		}
		return loc
	}
	return time.UTC
}

func NormalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func ConvertToUTC(t time.Time, loc *time.Location) time.Time {
	local := t.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(),
		local.Hour(), local.Minute(), local.Second(), local.Nanosecond(), time.UTC)
}

func FromLocalToUTC(date time.Time, hour, minute int, loc *time.Location) time.Time {
	local := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, loc)
	return local.UTC()
}
