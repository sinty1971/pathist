package models

import (
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/blake2b"
)

// RadixTable は32進数変換用の文字テーブル
// 0, I, O, Qを除く32文字（読み間違いを防ぐため）
const RadixTable = "123456789ABCDEFGHJKLMNPRSTUVWXYZ"

// ID はBLAKE2bハッシュを基にした一意識別子を生成する構造体
type ID struct {
	hash128 *big.Int
}

// NewIDFromBytes はバイト配列からIDを生成
func NewIDFromBytes(data []byte) *ID {
	// BLAKE2b-256を使用（GoにはBLAKE3の標準実装がないため）
	hash := blake2b.Sum256(data)
	// 下位128ビットを使用
	lower128 := hash[16:]

	return &ID{
		hash128: new(big.Int).SetBytes(lower128),
	}
}

// NewIDFromString は文字列からIDを生成
func NewIDFromString(str string) *ID {
	return NewIDFromBytes([]byte(str))
}

// NewIDFromTime はtime.TimeからIDを生成
func NewIDFromTime(t time.Time) *ID {
	// 時刻をナノ秒精度の文字列に変換してIDを生成
	timeStr := t.Format(time.RFC3339Nano)
	return NewIDFromString(timeStr)
}

// NewIDFromStateInfo フォルダー名からIDを生成
func NewIDFromKoujiProject(koujiProject KoujiEntry) *ID {
	// フォルダー名で一意性を確保
	data := fmt.Sprintf("%d%s%s", koujiProject.FileEntry.ID, koujiProject.CompanyName, koujiProject.LocationName)
	return NewIDFromString(data)
}

// Full25 は25文字のIDを返す（32^25 = 約10^37通り）
func (id *ID) Full25() string {
	return id.toBase32(25)
}

// Len7 は7文字のIDを返す（32^7 = 約34億通り）
func (id *ID) Len7() string {
	return id.toBase32(7)
}

// Len5 は5文字のIDを返す（32^5 = 約3,355万通り）
func (id *ID) Len5() string {
	return id.toBase32(5)
}

// Len3 は3文字のIDを返す（32^3 = 32,768通り）
func (id *ID) Len3() string {
	return id.toBase32(3)
}

// toBase32 は指定された長さの32進数文字列を生成
func (id *ID) toBase32(length int) string {
	bytes := make([]byte, length)
	value := new(big.Int).Set(id.hash128)
	base := big.NewInt(32)
	mod := new(big.Int)

	for i := length - 1; i >= 0; i-- {
		value.DivMod(value, base, mod)
		bytes[i] = RadixTable[mod.Int64()]
	}

	return string(bytes)
}

// String はデフォルトで5文字のIDを返す
func (id *ID) String() string {
	return id.Len5()
}
