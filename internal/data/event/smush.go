package event

import (
	"sort"

	"github.com/dxta-dev/app/internal/data"
)

func groupEventsByMergeRequest(events data.EventSlice) map[int64]data.EventSlice {
	grouped := make(map[int64]data.EventSlice)
	for _, event := range events {
		grouped[event.MergeRequestId] = append(grouped[event.MergeRequestId], event)
	}
	for _, slice := range grouped {
		sort.Sort(slice)
	}
	return grouped
}

func SmushEventSlice(events data.EventSlice) data.EventSlice {
	grouped := groupEventsByMergeRequest(events)


}





