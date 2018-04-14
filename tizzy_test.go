package tizzy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadLocation(t *testing.T) {
	data := []struct {
		Location   string
		ShouldFail bool
		Timezone   int
	}{
		{"Africa/Asmara", false, 10800},
		{"Africa/Tokyo", true, 0},
		{"America/New_York", false, -14400},
		{"Europe/Berlin", false, 7200},
		{"UTC", false, 0},
		{"", false, 0}, //UTC
		{" ", true, 0},
		{"America", true, 0},
		{"../Backslashes/are/not/allowed", true, 0},
		{"/No/Leading/Slashes", true, 0},
	}

	for _, e := range data {
		loc, err := LoadLocation(e.Location)
		if e.ShouldFail {
			require.Error(t, err)
			continue
		}

		require.NoError(t, err)

		_, tz := time.Now().In(loc).Zone()
		require.Equal(t, e.Timezone, tz)
	}

}
