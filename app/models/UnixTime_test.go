package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnixTime_ToSQLTimeStamp(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
	}{
		{"now", time.Now()},
		{"now+6", time.Now().Add(time.Hour * 6)},
		{"now+12", time.Now().Add(time.Hour * 12)},
		{"now+18", time.Now().Add(time.Hour * 18)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				tt.t.Year(), tt.t.Month(), tt.t.Day(),
				tt.t.Hour(), tt.t.Minute(), tt.t.Second())
			converted := UnixTime(tt.t).ToSQLTimeStamp()
			assert.Equal(t, want, converted)
		})
	}
}
