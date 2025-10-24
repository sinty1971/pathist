package models

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Timestamp は time.Time をラップし、カスタムの YAML/JSON マーシャリング/アンマーシャリングを提供します
type TimestampEx struct {
	timestamppb.Timestamp
}

type TimestampEncoding struct {
	// Time is the time value in RFC3339 format
	Time string `json:"time" yaml:"time" swaggertype:"string" format:"date-time" example:"2024-01-15T10:30:00Z"`
}

// MarshalYAML implements yaml.Marshaler
func (ts *TimestampEx) MarshalYAML() (any, error) {
	// nil チェック
	if ts == nil {
		return "", errors.New("(ts *TimestampEx) MarshalYAML() で ts 値が nil です。")
	}
	asTime := ts.AsTime()
	if asTime.IsZero() {
		return "", nil
	}
	return asTime.Format(time.RFC3339Nano), nil
}

// UnmarshalYAML implements yaml.Unmarshaler
func (ts *TimestampEx) UnmarshalYAML(unmarshal func(any) error) error {
	// 日時文字列を抽出
	var in string
	if err := unmarshal(&in); err != nil {
		return err
	}

	// 日時文字列をパース
	parsed, err := ParseTime(in, nil)
	if err != nil {
		return err
	}

	ts.Timestamp = *timestamppb.New(parsed)
	return nil
}

// DefaultTimestampFormatsWithTZ タイムゾーン情報を含む日時フォーマットのリスト（優先順位順）
var DefaultTimestampFormatsWithTZ = []string{
	time.RFC3339Nano,
	time.RFC3339,
}

// DefaultTimestampRegexpsWithTZ タイムゾーン情報を含む日時フォーマットの正規表現マップ
var DefaultTimestampRegexpsWithTZ = make(map[string]regexp.Regexp, len(DefaultTimestampFormatsWithTZ))

// DefaultTimestampFormatsWithoutTZ タイムゾーン情報を含まない日時フォーマットのリスト（ローカルタイムゾーンを使用）
var DefaultTimestampFormatsWithoutTZ = []string{
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-0102",
	"2006-01-02",
	"20060102",
	"2006/01/02",
	"2006.01.02",
	"2006/1/2",
	"2006.1.2",
}

// DefaultTimestampRegexpsWithoutTZ タイムゾーン情報を含まない日時フォーマットの正規表現マップ
var DefaultTimestampRegexpsWithoutTZ = make(map[string]regexp.Regexp, len(DefaultTimestampFormatsWithoutTZ))

// TimestampParseReplaceRule 日時フォーマットを正規表現パターンに変換するための置換ルール
var TimestampParseReplaceRule = map[string]string{
	"2006":      `\d{4}`,                 // 年
	"01":        `\d{2}`,                 // 月
	"02":        `\d{2}`,                 // 日
	"15":        `\d{2}`,                 // 時（24時間）
	"04":        `\d{2}`,                 // 分
	"05":        `\d{2}`,                 // 秒
	"999999999": `\d{9}`,                 // ナノ秒
	"Z07:00":    `(?:Z|[+-]\d{2}:\d{2})`, // タイムゾーン
	"Z0700":     `(?:Z|[+-]\d{4})`,       // タイムゾーン（コロンなし）
	"Z07":       `(?:Z|[+-]\d{2})`,       // タイムゾーン（時間のみ）
}

// init パッケージ初期化時に正規表現を初期化
func init() {
	InitializeTimestampRegexps()
}

// InitializeTimestampRegexps 日時フォーマットを正規表現に変換して初期化します
func InitializeTimestampRegexps() {
	formatToRegex := func(format string) regexp.Regexp {
		// 順序を考慮して置換を実行
		pattern := format

		// DateParseReplaceRuleからキーを取得し、文字数順でソート
		keys := slices.Collect(maps.Keys(TimestampParseReplaceRule))
		slices.SortFunc(keys, func(a, b string) int {
			return len(b) - len(a) // 文字数の大きい順
		})

		// 順序に従って置換
		for _, key := range keys {
			pattern = strings.Replace(pattern, key, TimestampParseReplaceRule[key], -1)
		}

		return *regexp.MustCompile(pattern)
	}

	// Initialize patterns for formats with timezone
	for _, format := range DefaultTimestampFormatsWithTZ {
		DefaultTimestampRegexpsWithTZ[format] = formatToRegex(format)
	}

	// Initialize patterns for formats without timezone
	for _, format := range DefaultTimestampFormatsWithoutTZ {
		DefaultTimestampRegexpsWithoutTZ[format] = formatToRegex(format)
	}
}

// findTimeString 正規表現を使用して文字列から日時部分を抽出し、
// 日時部分とそれ以外の文字列を分離して返します
func findTimeString(re *regexp.Regexp, in string, out *string) (string, error) {
	matches := re.FindStringIndex(in)
	if matches == nil {
		return "", fmt.Errorf("unable to parse date/time in the string: %s", in)
	}
	// 日時部分を抽出
	timeStr := in[matches[0]:matches[1]]

	// 出力パラメータが指定されている場合
	if out != nil {

		// 日時の前の部分を取得
		*out = strings.TrimSpace(in[:matches[0]])
		if *out != "" {
			*out = *out + " "
		}

		// 日時の後の部分を取得
		*out += strings.TrimSpace(in[matches[1]:])
		*out = strings.TrimSpace(*out)
	}

	return timeStr, nil
}

// ParseTime 文字列から日時をパースし、time.Time型と残りの文字列を返します
// 例: "2025-0618 豊田築炉 名和工場" → time.Time(2025-06-18), "豊田築炉 名和工場"
func ParseTime(in string, out *string) (time.Time, error) {
	// タイムゾーン付きのフォーマットを試行（配列順序で）
	for _, format := range DefaultTimestampFormatsWithTZ {
		//
		re := DefaultTimestampRegexpsWithTZ[format]
		dateStr, err := findTimeString(&re, in, out)
		if err != nil {
			continue
		}
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// タイムゾーンなしのフォーマットを試行（配列順序で）
	for _, format := range DefaultTimestampFormatsWithoutTZ {
		re := DefaultTimestampRegexpsWithoutTZ[format]
		dateStr, err := findTimeString(&re, in, out)
		if err != nil {
			continue
		}

		if t, err := time.ParseInLocation(format, dateStr, time.Local); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date/time in the string: %s", in)
}

// FormatTime 日時を指定のフォーマットに変換します
func FormatTime(layout string, t time.Time) (string, error) {
	if _, exists := DefaultTimestampRegexpsWithTZ[layout]; exists {
		return t.Format(layout), nil
	}
	if _, exists := DefaultTimestampRegexpsWithoutTZ[layout]; exists {
		return t.Format(layout), nil
	}
	return "", fmt.Errorf("invalid layout: %s", layout)
}

// ParseRFC3339Nano RFC3339Nano形式の日時文字列をパースしてtime.Timeを返します
// 文字列内の任意の位置にある日時を抽出可能です
func ParseRFC3339Nano(in string) (time.Time, error) {
	// 日時部分を抽出
	var out string
	re := DefaultTimestampRegexpsWithTZ[time.RFC3339Nano]
	timeStr, err := findTimeString(&re, in, &out)
	if err != nil {
		return time.Time{}, err
	}

	// 日時部分をパース
	t, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// ParseToTimestamp parses various date/time string formats and returns a TimestampEx
// When no timezone is specified, it uses the server's local timezone
func ParseToTimestamp(in string, out *string) (TimestampEx, error) {
	t, err := ParseTime(in, out)
	if err != nil {
		return TimestampEx{}, err
	}
	return TimestampEx{Timestamp: *timestamppb.New(t)}, nil
}

// MarshalJSON implements json.Marshaler
// タイムスタンプをJSON形式の文字列に変換します
func (ts *TimestampEx) MarshalJSON() ([]byte, error) {
	if ts.Timestamp.AsTime().IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + ts.Timestamp.AsTime().Format(time.RFC3339Nano) + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
// バイト文字列を受け取り、タイムスタンプをパースします
func (ts *TimestampEx) UnmarshalJSON(data []byte) error {
	in := string(data)
	if len(in) >= 2 && in[0] == '"' && in[len(in)-1] == '"' {
		in = in[1 : len(in)-1]
	}

	if in == "" {
		ts.Timestamp = *timestamppb.New(time.Time{})
		return nil
	}

	// タイムスタンプ文字列のパース
	parsed, err := ParseTime(in, nil)
	if err != nil {
		return err
	}

	ts.Timestamp = *timestamppb.New(parsed)
	return nil
}

func (ts *TimestampEx) Format(layout string) (string, error) {
	return FormatTime(layout, ts.Timestamp.AsTime())
}
