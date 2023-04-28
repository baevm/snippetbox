package main

import (
	"snippetbox/internal/assert"
	"testing"
	"time"
)

func Test_humanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2023, 4, 28, 18, 23, 15, 0, time.UTC),
			want: "28 Apr 2023 at 18:23",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2023, 4, 28, 18, 23, 15, 0, time.FixedZone("CET", 1*60*60)),
			want: "28 Apr 2023 at 17:23",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
