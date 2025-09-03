package reverser

import (
	"ccnu-library-mcp-go/internal/auther"
	"ccnu-library-mcp-go/pkg"
	"context"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSeat_IsFreeByTime(t *testing.T) {
	// 定义测试用例
	tests := []struct {
		name            string
		seat            Seat
		startTime       time.Time
		endTime         time.Time
		expectedFree    bool
		expectedPeriods []Period
	}{
		{
			name: "完全空闲",
			seat: Seat{
				ReserveStartTime: time.Date(2023, 10, 1, 8, 0, 0, 0, time.UTC),
				ReserveEndTime:   time.Date(2023, 10, 1, 18, 0, 0, 0, time.UTC),
				IsFree:           true,
			},
			startTime:       time.Date(2023, 10, 1, 9, 0, 0, 0, time.UTC),
			endTime:         time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			expectedFree:    true,
			expectedPeriods: []Period{{StartTime: time.Date(2023, 10, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC)}},
		},
		{
			name: "部分空闲",
			seat: Seat{
				ReserveStartTime: time.Date(2023, 10, 1, 8, 0, 0, 0, time.UTC),
				ReserveEndTime:   time.Date(2023, 10, 1, 18, 0, 0, 0, time.UTC),
				IsFree:           false,
				OccupyStates: []Period{
					{StartTime: time.Date(2023, 10, 1, 9, 0, 0, 0, time.UTC), EndTime: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC)},
				},
			},
			startTime:       time.Date(2023, 10, 1, 8, 30, 0, 0, time.UTC),
			endTime:         time.Date(2023, 10, 1, 9, 30, 0, 0, time.UTC),
			expectedFree:    false,
			expectedPeriods: []Period{{StartTime: time.Date(2023, 10, 1, 8, 30, 0, 0, time.UTC), EndTime: time.Date(2023, 10, 1, 9, 0, 0, 0, time.UTC)}},
		},
		{
			name: "完全被占用",
			seat: Seat{
				ReserveStartTime: time.Date(2023, 10, 1, 8, 0, 0, 0, time.UTC),
				ReserveEndTime:   time.Date(2023, 10, 1, 18, 0, 0, 0, time.UTC),
				IsFree:           false,
				OccupyStates: []Period{
					{StartTime: time.Date(2023, 10, 1, 8, 30, 0, 0, time.UTC), EndTime: time.Date(2023, 10, 1, 9, 30, 0, 0, time.UTC)},
				},
			},
			startTime:       time.Date(2023, 10, 1, 8, 30, 0, 0, time.UTC),
			endTime:         time.Date(2023, 10, 1, 9, 30, 0, 0, time.UTC),
			expectedFree:    false,
			expectedPeriods: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			free, periods := tt.seat.IsFreeByTime(tt.startTime, tt.endTime)
			if free != tt.expectedFree {
				t.Errorf("expected free: %v, got: %v", tt.expectedFree, free)
			}
			if !reflect.DeepEqual(periods, tt.expectedPeriods) {
				t.Errorf("expected periods: %v, got: %v", tt.expectedPeriods, periods)
			}
		})
	}
}

func LoadInfo() (string, string) {
	stuID := os.Getenv("STUID")
	pwd := os.Getenv("PASSWORD")
	return stuID, pwd
}

func TestReverser_GetSeatsByTime(t *testing.T) {
	stuID, pwd := LoadInfo()
	a := auther.NewAuther()
	r := NewReverser(a)
	_ = r.StoreStuInfo(context.Background(), stuID, pwd)

	res, err := r.GetSeatsByTime(context.Background(), "2023214414", pkg.Rooms["n1m"],
		pkg.CreateShanghaiTime(2025, 9, 3, 18, 0),
		pkg.CreateShanghaiTime(2025, 9, 3, 21, 0),
		true)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("可用座位数: %d", len(res))

}
