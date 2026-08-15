package ingest

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseEpochPreservesSubsecondUnits(t *testing.T) {
	tests := []struct {
		name string
		text string
		want time.Time
	}{
		{"milliseconds", "1735689600123", time.Date(2025, 1, 1, 0, 0, 0, 123000000, time.UTC)},
		{"microseconds", "1735689600123456", time.Date(2025, 1, 1, 0, 0, 0, 123456000, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTimestamp(json.Number(tt.text), nil)
			if err != nil {
				t.Fatal(err)
			}
			if !got.Equal(tt.want) {
				t.Fatalf("parseTimestamp(%s) = %s, want %s", tt.text, got.Format(time.RFC3339Nano), tt.want.Format(time.RFC3339Nano))
			}
		})
	}
}
