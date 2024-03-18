package data

import (
	"sort"
	"time"
)

type TimeFrame struct {
	Since time.Time
	Until time.Time
}

type TimeFrameSlice []TimeFrame

func (tfs TimeFrameSlice) Len() int {
	return len(tfs)
}

func (tfs TimeFrameSlice) Less(i, j int) bool {
	return tfs[i].Since.Before(tfs[j].Since)
}

func (tfs TimeFrameSlice) Swap(i, j int) {
	tfs[i], tfs[j] = tfs[j], tfs[i]
}

func GetCrawlInstances(from, to time.Time) TimeFrameSlice {
	return nil
}

func FindGaps(timeFrames TimeFrameSlice) TimeFrameSlice {
	var gaps TimeFrameSlice

	sort.Sort(timeFrames)

	gaps = append(gaps, TimeFrame{
		Since: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
		Until: timeFrames[0].Since,
	})

	for i := 1; i < len(timeFrames)-1; i++ {
		if timeFrames[i-1].Until.Before(timeFrames[i].Since) {
			gaps = append(gaps, TimeFrame{
				Since: timeFrames[i-1].Until,
				Until: timeFrames[i].Since,
			})
		}
	}

	return gaps

}


