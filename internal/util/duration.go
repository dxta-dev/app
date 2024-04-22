package util

import (
	"fmt"
	"time"
)

func FormatDurationAgo(when time.Time, now time.Time) string {
	dur := now.Sub(when)

	sixty := float64(60)
	fourteen := float64(14)
	twentyFour := float64(24)

	secondsAgo := dur.Seconds()
	minutesAgo := secondsAgo / sixty
	hoursAgo := minutesAgo / sixty
	daysAgo := hoursAgo / twentyFour

	if secondsAgo < sixty {
		return "just now"
	}

	if minutesAgo < sixty {
		return fmt.Sprintf("%vm ago", int64(minutesAgo))
	}

	if hoursAgo < twentyFour {
		return fmt.Sprintf("%vh ago", int64(hoursAgo))
	}

	if daysAgo < fourteen {
		return fmt.Sprintf("%vd ago", int64(daysAgo))
	}

	return "ages ago"
}

func FormatDurationDaysAgo(when time.Time, now time.Time) string {
	dur := now.Sub(when)

	fourteen := float64(14)
	twentyFour := float64(24)

	hoursAgo := dur.Hours()
	daysAgo := hoursAgo / twentyFour

	if hoursAgo < twentyFour {
		return "today"
	}

	if daysAgo < fourteen {
		return fmt.Sprintf("%vd ago", int64(daysAgo))
	}

	return "ages ago"
}
