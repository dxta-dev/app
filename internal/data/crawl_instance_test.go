package data

import (
	"reflect"
	"testing"
	"time"
)

func TestFindGaps(t *testing.T) {

	parseTime := func(s string) time.Time {
		time, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			t.Error(err)
		}
		return time
	}

	tests := []struct {
		name         string
		from         time.Time
		to           time.Time
		timeFrames   map[int64]TimeFrameSlice
		wantByRepoId map[int64]TimeFrameSlice
	}{
		{
			name: "no time frames",
			from: parseTime("2019-01-01 00:00:00"),
			to:   parseTime("2019-01-02 00:00:00"),
			timeFrames: map[int64]TimeFrameSlice{
				1: nil,
			},
			wantByRepoId: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-02 00:00:00"),
					},
				},
			},
		},
		{
			name: "consecutive time frames",
			from: parseTime("2019-01-01 00:00:00"),
			to:   parseTime("2019-01-03 00:00:00"),
			timeFrames: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-02 00:00:00"),
					},
					{
						Since: parseTime("2019-01-02 00:00:00"),
						Until: parseTime("2019-01-03 00:00:00"),
					},
				},
			},
			wantByRepoId: map[int64]TimeFrameSlice{
				1: nil,
			},
		},
		{
			name: "one gap",
			from: parseTime("2019-01-01 00:00:00"),
			to:   parseTime("2019-01-03 00:00:00"),
			timeFrames: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-02 00:00:00"),
					},
				},
			},
			wantByRepoId: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-02 00:00:00"),
						Until: parseTime("2019-01-03 00:00:00"),
					},
				},
			},
		},
		{
			name: "one time frame fully contained within another",
			from: parseTime("2019-01-01 00:00:00"),
			to:   parseTime("2019-01-04 00:00:00"),
			timeFrames: map[int64]TimeFrameSlice{
				2: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-04 00:00:00"),
					},
					{
						Since: parseTime("2019-01-02 00:00:00"),
						Until: parseTime("2019-01-03 00:00:00"),
					},
				},
			},
			wantByRepoId: map[int64]TimeFrameSlice{
				2: nil,
			},
		},
		{
			name: "gap for repository 1 and 2",
			from: parseTime("2019-01-01 00:00:00"),
			to:   parseTime("2019-01-03 00:00:00"),
			timeFrames: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-02 00:00:00"),
					},
				},
				2: {
					{
						Since: parseTime("2019-01-02 00:00:00"),
						Until: parseTime("2019-01-03 00:00:00"),
					},
				},
			},
			wantByRepoId: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-02 00:00:00"),
						Until: parseTime("2019-01-03 00:00:00"),
					},
				},
				2: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-02 00:00:00"),
					},
				},
			},
		},
		{
			name: "gap for repository 1 and NOT for repository 2",
			from: parseTime("2019-01-01 00:00:00"),
			to:   parseTime("2019-01-03 00:00:00"),
			timeFrames: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-02 00:00:00"),
					},
				},
				2: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-02 00:00:00"),
					},
					{
						Since: parseTime("2019-01-02 00:00:00"),
						Until: parseTime("2019-01-03 00:00:00"),
					},
				},
			},
			wantByRepoId: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-02 00:00:00"),
						Until: parseTime("2019-01-03 00:00:00"),
					},
				},
				2: nil,
			},
		},
		{
			name: "completely non-overlapping time frames",
			from: parseTime("2019-01-01 00:00:00"),
			to:   parseTime("2019-01-04 00:00:00"),
			timeFrames: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-01 00:00:00"),
						Until: parseTime("2019-01-02 00:00:00"),
					},
					{
						Since: parseTime("2019-01-03 00:00:00"),
						Until: parseTime("2019-01-04 00:00:00"),
					},
				},
			},
			wantByRepoId: map[int64]TimeFrameSlice{
				1: {
					{
						Since: parseTime("2019-01-02 00:00:00"),
						Until: parseTime("2019-01-03 00:00:00"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for repoId, timeFrames := range tt.timeFrames {
				got := FindGaps(tt.from, tt.to, timeFrames)
				want, ok := tt.wantByRepoId[repoId]
				if !ok {
					t.Errorf("Missing expected result for repository ID %d", repoId)
					continue
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("FindGaps() = %v, want %v", got, want)
				}
			}
		})
	}
}
