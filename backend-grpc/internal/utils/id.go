package utils

import (
	"math/big"

	"golang.org/x/crypto/blake2b"
)

// RadixTable は53進数変換用の文字テーブル
// 0, I, O, Q, i, j, o, u, vを除く53文字（=10+26+26-9読み間違いを防ぐため）
// 53^6 = 22,164,361,129 通り（約220億通り）
const RadixTable = "123456789ABCDEFGHJKLMNPRSTUVWXYZabcdefghklmnpqrstwxyz"

// GenerateIdFromBytes はバイト配列からハッシュ文字列IDを生成
func GenerateIdFromBytes(data []byte) string {

	// バイト配列からBLAKE2b-256ハッシュを計算し下位128ビットを取得
	// BLAKE2b-256を使用（GoにはBLAKE3の標準実装がないため）
	hash := blake2b.Sum256(data)

	// 下位128ビットを使用
	hashLower128 := new(big.Int).SetBytes(hash[16:])

	// 生成文字列を計算、長さは６文字固定
	length := 6
	bytes := make([]byte, length)
	value := new(big.Int).Set(hashLower128)
	base := big.NewInt(53)
	mod := new(big.Int)

	for i := length - 1; i >= 0; i-- {
		value.DivMod(value, base, mod)
		bytes[i] = RadixTable[mod.Int64()]
	}

	return string(bytes)
}

// GenerateIdFromString は文字列からハッシュ文字列IDを生成
func GenerateIdFromString(str string) string {
	return GenerateIdFromBytes([]byte(str))
}
