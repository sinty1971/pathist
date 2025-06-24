package models

import (
	"penguin-backend/internal/utils"
	"time"
)

// Timestamp wraps time.Time with custom YAML/JSON marshaling/unmarshaling
// @Description Timestamp in RFC3339 format
// swagger:model
type Timestamp struct {
	time.Time `swaggertype:"string" format:"date-time" example:"2024-01-15T10:30:00Z"`
}

// ParseTimestamp parses various date/time string formats and returns a Timestamp
// When no timezone is specified, it uses the server's local timezone
func ParseTimestamp(in string, out *string) (Timestamp, error) {
	t, err := utils.ParseTime(in, out)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp{Time: t}, nil
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
	// 日時文字列を抽出
	var in string
	if err := unmarshal(&in); err != nil {
		return err
	}

	// 日時文字列をパース
	parsed, err := utils.ParseTime(in, nil)
	if err != nil {
		return err
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
	in := string(data)
	if len(in) >= 2 && in[0] == '"' && in[len(in)-1] == '"' {
		in = in[1 : len(in)-1]
	}

	if in == "" {
		ts.Time = time.Time{}
		return nil
	}

	// タイムスタンプ文字列のパース
	parsed, err := utils.ParseTime(in, nil)
	if err != nil {
		return err
	}

	ts.Time = parsed
	return nil
}

func (ts *Timestamp) Format(layout string) (string, error) {
	return utils.FormatTime(layout, ts.Time)
}
