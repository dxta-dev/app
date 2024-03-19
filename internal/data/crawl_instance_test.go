package data

import (
	"reflect"
	"testing"
	"time"
)

func TestFindGaps(t *testing.T) {

	parseTime := func(s string) time.Time {
		time, err := time.Parse("2006-01-02 10:22:12", s)
		if err != nil {
			t.Error(err)
		}
		return time
	}

	tests := []struct {
		name     string
		from	 time.Time
		to	   time.Time
		timeFrames TimeFrameSlice
		want	 TimeFrameSlice
	}{
		{
			name: "no time frames",
			from: parseTime("2019-01-01 00:00:00"),
			to: parseTime("2019-01-02 00:00:00"),
			timeFrames: nil,
			want: TimeFrameSlice{
				{
					Since: parseTime("2019-01-01 00:00:00"),
					Until: parseTime("2019-01-02 00:00:00"),
				},
			},
		},
		{
			name: "consecutive time frames",
			from: parseTime("2019-01-01 00:00:00"),
			to: parseTime("2019-01-03 00:00:00"),
			timeFrames: TimeFrameSlice{
				{
					Since: parseTime("2019-01-01 00:00:00"),
					Until: parseTime("2019-01-02 00:00:00"),
				},
				{
					Since: parseTime("2019-01-02 00:00:00"),
					Until: parseTime("2019-01-03 00:00:00"),
				},
			},
			want: TimeFrameSlice{},
		},
		{
			name: "one gap",
			from: parseTime("2019-01-01 00:00:00"),
			to: parseTime("2019-01-03 00:00:00"),
			timeFrames: TimeFrameSlice{
				{
					Since: parseTime("2019-01-01 00:00:00"),
					Until: parseTime("2019-01-02 00:00:00"),
				},
			},
			want: TimeFrameSlice{
				{
					Since: parseTime("2019-01-02 00:00:00"),
					Until: parseTime("2019-01-03 00:00:00"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindGaps(tt.from, tt.to, tt.timeFrames)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindGaps() = %v, want %v", got, tt.want)
			}
		})
	}
}
