package core

import (
	"os"
)

// Persist は永続化インターフェースを定義します。
type Persist interface {
	// GetPersistPath はファイルの永続化フルパスを取得します。
	GetPersistPath() string

	// GetPersistMap は永続化対象のメッセージを取得します。
	GetPersistBytes() []byte

	// SetPersistMessage は永続化対象のメッセージを設定します。
	SetPersistBytes(b []byte)
}

// Persister は永続化されるエンティティに必要な振る舞いを定義します。
type Persister struct {
}

// Load は永続化ファイルからデータを読み込みます。
func (h *Persister) Load(p Persist) error {

	// YAMLファイルからバイトデータを読み込む
	b, err := os.ReadFile(p.GetPersistPath())
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return h.Save(p)
		}
	}

	p.SetPersistBytes(b)
	return nil
}

// Save はデータを永続化ファイルに保存します。
func (h *Persister) Save(p Persist) error {

	// Persist 永続化バイトデーターの取得
	b := p.GetPersistBytes()

	// ファイルに書き込み
	return os.WriteFile(p.GetPersistPath(), b, 0644)
}
