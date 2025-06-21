package models

import (
	"fmt"
	"penguin-backend/internal/utils"
	"time"
)

// Timestamp wraps time.Time with custom YAML/JSON marshaling/unmarshaling
// @Description Timestamp in RFC3339 format
// swagger:model
type Timestamp struct {
	time.Time `swaggertype:"string" format:"date-time" example:"2024-01-15T10:30:00Z"`
}

// NewTimestamp creates a new Timestamp from time.Time
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{Time: t}
}

// ParseTimestamp parses various date/time string formats and returns a Timestamp
// When no timezone is specified, it uses the server's local timezone
func ParseTimestamp(s string) (Timestamp, error) {
	t, _, err := utils.ParseTimeAndRest(s)
	return NewTimestamp(t), err
}

// MarshalYAML implements yaml.Marshaler
func (ts Timestamp) MarshalYAML() (any, error) {
	if ts.Time.IsZero() {
		return "", nil
	}
	return ts.Time.Format(time.RFC3339Nano), nil
}

// UnmarshalYAML implements yaml.Unmarshaler
func (ts *Timestamp) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	// Try parsing with RFC3339Nano first
	parsed, err := utils.ParseTime(str)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	ts.Time = parsed
	return nil
}

// MarshalJSON implements json.Marshaler
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	if ts.Time.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + ts.Time.Format(time.RFC3339Nano) + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	if str == "" {
		ts.Time = time.Time{}
		return nil
	}

	// タイムスタンプ文字列のパース
	parsed, err := utils.ParseTime(str)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	ts.Time = parsed
	return nil
}

// ParseAndRest parses a timestamp string and returns the timestamp and the rest string
func ParseTimestampAndRest(s string, ts *Timestamp) (string, error) {
	t, rest, err := utils.ParseTimeAndRest(s)
	if err != nil {
		return "", err
	}
	ts.Time = t
	return rest, nil
}
