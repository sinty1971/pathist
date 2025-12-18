package models

import (
	"errors"
	"fmt"
	"maps"
	"math/big"
	"regexp"
	"server-grpc/internal/core"
	"slices"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Timestamp は time.Time をラップし、カスタムの YAML/JSON マーシャリング/アンマーシャリングを提供します
type Timestamp struct {
	*timestamppb.Timestamp
}

// GetID 時間データからIDを生成
func (ts *Timestamp) GetID() string {
	// 秒とナノ秒からバイト配列を生成
	bytes := big.NewInt(int64(ts.GetSeconds())*1_000_000_000 + int64(ts.GetNanos())).Bytes()

	// 時刻をナノ秒精度の文字列に変換してIDを生成
	return core.ParseIdFromBytes(bytes)
}

// MarshalYAML implements yaml.Marshaler
func (ts *Timestamp) MarshalYAML() (any, error) {
	// nil チェック
	if ts == nil || ts.Timestamp.CheckValid() == nil {
		return "", errors.New("(ts *Timestamp) MarshalYAML() で値が nil もしくは不正です")
	}
	return ts.Timestamp.AsTime().Format(time.RFC3339Nano), nil
}

// UnmarshalYAML implements yaml.Unmarshaler
func (ts *Timestamp) UnmarshalYAML(unmarshal func(any) error) error {
	// 日時文字列を抽出
	var enc string
	if err := unmarshal(&enc); err != nil {
		return err
	}

	// 空文字列の場合はゼロ値を設定
	enc = strings.TrimSpace(enc)
	if enc == "" {
		ts.Timestamp = &timestamppb.Timestamp{}
		return nil
	}

	// timestamppb 内蔵の JSON アンマーシャリングを優先利用
	quoted := strconv.Quote(enc)
	if err := protojson.Unmarshal([]byte(quoted), ts.Timestamp); err != nil {
		// 互換性維持のため旧ロジックにフォールバック
		var parsed *Timestamp
		_, parseErr := ParseTimestamp(enc, parsed)
		if parseErr != nil {
			return fmt.Errorf("timestamppb.UnmarshalJSON failed: %w; fallback parse failed: %v", err, parseErr)
		}
		ts.Timestamp = &timestamppb.Timestamp{
			Seconds: parsed.Seconds,
			Nanos:   parsed.Nanos,
		}
	}
	return nil
}

// DefaultTimestampFormatsWithTZ タイムゾーン情報を含む日時フォーマットのリスト（優先順位順）
var DefaultTimestampFormatsWithTZ = []string{
	time.RFC3339Nano,
	time.RFC3339,
}

// TimestampRegexpsWithTZ タイムゾーン情報を含む日時フォーマットの正規表現マップ
var TimestampRegexpsWithTZ = make(map[string]regexp.Regexp, len(DefaultTimestampFormatsWithTZ))

// TimestampFormatsWithoutTZ タイムゾーン情報を含まない日時フォーマットのリスト（ローカルタイムゾーンを使用）
var TimestampFormatsWithoutTZ = []string{
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

// TimestampRegexpsWithoutTZ タイムゾーン情報を含まない日時フォーマットの正規表現マップ
var TimestampRegexpsWithoutTZ = make(map[string]regexp.Regexp, len(TimestampFormatsWithoutTZ))

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
	// 日時フォーマットを正規表現に変換して初期化します
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
		TimestampRegexpsWithTZ[format] = formatToRegex(format)
	}

	// Initialize patterns for formats without timezone
	for _, format := range TimestampFormatsWithoutTZ {
		TimestampRegexpsWithoutTZ[format] = formatToRegex(format)
	}
}

// findTimeString 正規表現を使用して文字列から日時部分を抽出し、
// 日時部分とそれ以外の文字列を分離して返します
func findTimeString(re *regexp.Regexp, in string, out *string) (string, error) {
	matches := re.FindStringIndex(in)
	if matches == nil {
		return in, fmt.Errorf("unable to parse date/time in the string: %s", in)
	}

	// 出力パラメータが nil の場合
	if out == nil {
		out = new(string)
	}

	// 日時部分を抽出
	*out = in[matches[0]:matches[1]]

	// 日時の前の部分を取得
	removed := strings.TrimSpace(in[:matches[0]])
	if removed != "" {
		removed = removed + " "
	}

	// 日時の後の部分を取得
	removed += strings.TrimSpace(in[matches[1]:])
	removed = strings.TrimSpace(removed)

	return removed, nil
}

// ParseTimestamp 文字列から日時をパースし、Timestamp型と残りの文字列を返します
//
// 例: "2025-0618 豊田築炉 名和工場" → Timestamp(2025-06-18), "豊田築炉 名和工場"
func ParseTimestamp(in string, out *Timestamp) (string, error) {
	// 出力パラメータが nil の場合は新規作成
	if out == nil {
		out = &Timestamp{Timestamp: &timestamppb.Timestamp{Seconds: 0, Nanos: 0}}
	}

	// 空文字列の場合はゼロ値を設定
	in = strings.TrimSpace(in)
	if in == "" {
		return "", nil
	}

	// timestamppb 内蔵の JSON アンマーシャリングを優先利用
	quoted := strconv.Quote(in)
	pbts := timestamppb.Timestamp{}
	if err := protojson.Unmarshal([]byte(quoted), &pbts); err == nil {
		*out = Timestamp{Timestamp: &timestamppb.Timestamp{
			Seconds: pbts.GetSeconds(),
			Nanos:   pbts.GetNanos(),
		}}
		return "", nil
	}

	// タイムゾーン付きのフォーマットを試行（配列順序で）
	for _, format := range DefaultTimestampFormatsWithTZ {
		//
		re := TimestampRegexpsWithTZ[format]
		dateStr := new(string)
		removed, err := findTimeString(&re, in, dateStr)
		if err != nil {
			continue
		}
		if t, err := time.Parse(format, *dateStr); err == nil {
			*out = Timestamp{Timestamp: timestamppb.New(t)}
			return removed, nil
		}
	}

	// タイムゾーンなしのフォーマットを試行（配列順序で）
	for _, format := range TimestampFormatsWithoutTZ {
		re := TimestampRegexpsWithoutTZ[format]
		dateStr := new(string)
		removed, err := findTimeString(&re, in, dateStr)
		if err != nil {
			continue
		}

		if t, err := time.ParseInLocation(format, *dateStr, time.Local); err == nil {
			*out = Timestamp{Timestamp: timestamppb.New(t)}
			return removed, nil
		}
	}

	return in, fmt.Errorf("unable to parse date/time in the string: %s", in)
}

// ParseRFC3339Nano RFC3339Nano形式の日時文字列をパースしてtime.Timeを返します
// 文字列内の任意の位置にある日時を抽出可能です
func ParseRFC3339Nano(in string, out *time.Time) (string, error) {
	if out == nil {
		out = new(time.Time)
	}

	// 作業編集の宣言
	var (
		err     error
		timeStr = new(string)
		removed string
	)

	// 日時部分を抽出
	re := TimestampRegexpsWithTZ[time.RFC3339Nano]
	if removed, err = findTimeString(&re, in, timeStr); err != nil {
		return in, err
	}

	// 日時部分をパース
	if *out, err = time.Parse(time.RFC3339Nano, *timeStr); err != nil {
		return in, err
	}
	return removed, nil
}

// MarshalJSON implements json.Marshaler
// タイムスタンプをJSON形式の文字列に変換します
func (ts *Timestamp) MarshalJSON() ([]byte, error) {
	if ts.Timestamp.AsTime().IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + ts.Timestamp.AsTime().Format(time.RFC3339Nano) + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
// バイト文字列を受け取り、タイムスタンプをパースします
func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	in := string(data)
	if len(in) >= 2 && in[0] == '"' && in[len(in)-1] == '"' {
		in = in[1 : len(in)-1]
	}

	if in == "" {
		ts.Timestamp = timestamppb.New(time.Time{})
		return nil
	}

	// タイムスタンプ文字列のパース
	var parsed *Timestamp
	_, err := ParseTimestamp(in, parsed)
	if err != nil {
		return err
	}

	ts.Timestamp = parsed.Timestamp
	return nil
}

// FormatTime 日時を指定のフォーマットに変換します
func (ts *Timestamp) FormatTime(layout string) (string, error) {
	if _, exists := TimestampRegexpsWithTZ[layout]; exists {
		return ts.AsTime().Format(layout), nil
	}
	if _, exists := TimestampRegexpsWithoutTZ[layout]; exists {
		return ts.AsTime().Format(layout), nil
	}
	return "", fmt.Errorf("invalid layout: %s", layout)
}
