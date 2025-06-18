// Package utils 汎用的なユーティリティ関数を提供します
package utils

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"time"
)

// DefaultTimestampRegexpsWithTZ タイムゾーン情報を含む日時フォーマットの正規表現マップ
var DefaultTimestampRegexpsWithTZ = make(map[string]regexp.Regexp, len(DefaultTimestampFormatsWithTZ))

// DefaultTimestampFormatsWithTZ タイムゾーン情報を含む日時フォーマットのリスト（優先順位順）
var DefaultTimestampFormatsWithTZ = []string{
	time.RFC3339Nano,
	time.RFC3339,
}

// DefaultTimestampRegexpsWithoutTZ タイムゾーン情報を含まない日時フォーマットの正規表現マップ
var DefaultTimestampRegexpsWithoutTZ = make(map[string]regexp.Regexp, len(DefaultTimestampFormatsWithoutTZ))

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

// ParseTime 様々な日時フォーマットの文字列をパースしてtime.Timeを返します
// タイムゾーンが指定されていない場合はサーバーのローカルタイムゾーンを使用します
func ParseTime(s string) (time.Time, error) {
	ts, _, err := ParseTimeAndRest(s)
	return ts, err
}

// findTimeStringAndRest 正規表現を使用して文字列から日時部分を抽出し、
// 日時部分とそれ以外の文字列を分離して返します
func findTimeStringAndRest(re *regexp.Regexp, s string) (*string, *string) {
	matches := re.FindStringIndex(s)
	if matches == nil {
		return nil, nil
	}
	// 日時部分を抽出
	dateStr := s[matches[0]:matches[1]]
	// 日時の前の部分を取得
	prefix := strings.TrimSpace(s[:matches[0]])
	// 日時の後の部分を取得
	suffix := strings.TrimSpace(s[matches[1]:])
	// 日時の前の部分と後の部分を結合（間にスペースを追加）
	var restStr string
	if prefix != "" && suffix != "" {
		restStr = prefix + " " + suffix
	} else {
		restStr = prefix + suffix
	}

	return &dateStr, &restStr
}

// ParseTimeAndRest 文字列から日時をパースし、time.Time型と残りの文字列を返します
// 例: "2025-0618 豊田築炉 名和工場" → time.Time(2025-06-18), "豊田築炉 名和工場"
func ParseTimeAndRest(s string) (time.Time, string, error) {
	// タイムゾーン付きのフォーマットを試行（配列順序で）
	for _, format := range DefaultTimestampFormatsWithTZ {
		re := DefaultTimestampRegexpsWithTZ[format]
		dateStr, restStr := findTimeStringAndRest(&re, s)
		if dateStr == nil {
			continue
		}
		if t, err := time.Parse(format, *dateStr); err == nil {
			return t, *restStr, nil
		}
	}

	// タイムゾーンなしのフォーマットを試行（配列順序で）
	for _, format := range DefaultTimestampFormatsWithoutTZ {
		re := DefaultTimestampRegexpsWithoutTZ[format]
		dateStr, restStr := findTimeStringAndRest(&re, s)
		if dateStr == nil {
			continue
		}

		if t, err := time.ParseInLocation(format, *dateStr, time.Local); err == nil {
			return t, *restStr, nil
		}
	}

	return time.Time{}, s, fmt.Errorf("unable to parse date/time in the string: %s", s)
}

// ParseRFC3339Nano RFC3339Nano形式の日時文字列をパースしてtime.Timeを返します
// 文字列内の任意の位置にある日時を抽出可能です
func ParseRFC3339Nano(s string) (time.Time, error) {
	re := DefaultTimestampRegexpsWithTZ[time.RFC3339Nano]
	dateStr, _ := findTimeStringAndRest(&re, s)
	if dateStr == nil {
		return time.Time{}, fmt.Errorf("unable to parse date/time in the string: %s", s)
	}
	t, err := time.Parse(time.RFC3339Nano, *dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse date/time in the string: %s", s)
	}
	return t, nil
}
