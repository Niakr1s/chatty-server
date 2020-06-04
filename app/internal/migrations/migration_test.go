package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fileNameToNum(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     int
		wantErr  bool
	}{
		{"valid 0", "0000.sql", 0, false},
		{"valid 1", "0001.sql", 1, false},
		{"valid 10", "0010.sql", 10, false},
		{"invalid", "000a.sql", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileNameToNum(tt.filename)
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
