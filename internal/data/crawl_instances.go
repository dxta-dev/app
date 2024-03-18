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

func FindGaps(from, to time.Time, timeFrames TimeFrameSlice) TimeFrameSlice {
	var gaps TimeFrameSlice

	sort.Sort(timeFrames)


	var until time.Time
	if len(timeFrames) == 0 {
		until = to
	} else {
		until = timeFrames[0].Since
	}


	gaps = append(gaps, TimeFrame{
		Since: from,
		Until: until,
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


