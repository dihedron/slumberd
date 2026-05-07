package timex

import (
	"slices"
	"testing"
	"time"
)

// TestDuration tests the Duration type.
func TestDuration(t *testing.T) {
	var d Duration
	tests := []struct {
		input string
		want  Duration
	}{
		{"1h", Duration(time.Hour)},
		{"1m", Duration(time.Minute)},
		{"1s", Duration(time.Second)},
	}
	for _, tt := range tests {
		d.UnmarshalFlag(tt.input)
		if d != tt.want {
			t.Errorf("expected %v, got %v", tt.want, d)
		}
	}
}

func TestDurationString(t *testing.T) {
	tests := []struct {
		input Duration
		want  string
	}{
		{Duration(time.Hour), "1h0m0s"},
		{Duration(time.Minute), "1m0s"},
		{Duration(time.Second), "1s"},
	}
	for _, tt := range tests {
		if tt.input.String() != tt.want {
			t.Errorf("expected %v, got %v", tt.want, tt.input.String())
		}
	}
}

func TestDurationMarshalJSON(t *testing.T) {
	tests := []struct {
		input Duration
		want  []byte
	}{
		{Duration(time.Hour), []byte("\"1h0m0s\"")},
		{Duration(time.Minute), []byte("\"1m0s\"")},
		{Duration(time.Second), []byte("\"1s\"")},
	}
	for _, tt := range tests {
		b, err := tt.input.MarshalJSON()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if !slices.Equal(b, tt.want) {
			t.Errorf("expected %v, got %v", tt.want, b)
		}
	}
}

func TestDurationUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input []byte
		want  Duration
	}{
		{[]byte("\"1h\""), Duration(time.Hour)},
		{[]byte("\"1m\""), Duration(time.Minute)},
		{[]byte("\"1s\""), Duration(time.Second)},
	}
	for _, tt := range tests {
		var d Duration
		if err := d.UnmarshalJSON(tt.input); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if d != tt.want {
			t.Errorf("expected %v, got %v", tt.want, d)
		}
	}
}

func TestDurationMarshalText(t *testing.T) {
	tests := []struct {
		input Duration
		want  []byte
	}{
		{Duration(time.Hour), []byte("1h0m0s")},
		{Duration(time.Minute), []byte("1m0s")},
		{Duration(time.Second), []byte("1s")},
	}
	for _, tt := range tests {
		b, err := tt.input.MarshalText()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if !slices.Equal(b, tt.want) {
			t.Errorf("expected %v, got %v", tt.want, b)
		}
	}
}

func TestDurationUnmarshalText(t *testing.T) {
	tests := []struct {
		input []byte
		want  Duration
	}{
		{[]byte("1h"), Duration(time.Hour)},
		{[]byte("1m"), Duration(time.Minute)},
		{[]byte("1s"), Duration(time.Second)},
	}
	for _, tt := range tests {
		var d Duration
		if err := d.UnmarshalText(tt.input); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if d != tt.want {
			t.Errorf("expected %v, got %v", tt.want, d)
		}
	}
}

func TestDurationMarshalYAML(t *testing.T) {

	tests := []struct {
		input Duration
		want  string
	}{
		{Duration(time.Hour), "1h0m0s"},
		{Duration(time.Minute), "1m0s"},
		{Duration(time.Second), "1s"},
	}
	for _, tt := range tests {
		b, err := tt.input.MarshalYAML()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if b.(string) != tt.want {
			t.Errorf("expected %v, got %v", tt.want, b)
		}
	}
}
