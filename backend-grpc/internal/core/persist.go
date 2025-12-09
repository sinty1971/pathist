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
	Path   string
	Object Persist
}

func NewPersister(persistPath string, persistObj Persist) *Persister {
	return &Persister{
		Path:   persistPath,
		Object: persistObj,
	}
}

// Load は永続化ファイルからデータを読み込みます。
func (h *Persister) Load() error {

	// YAMLファイルからバイトデータを読み込む
	b, err := os.ReadFile(h.Path)
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return h.Save()
		}
	}

	h.Object.SetPersistBytes(b)
	return nil
}

// Save はデータを永続化ファイルに保存します。
func (h *Persister) Save() error {

	// Persist 永続化バイトデーターの取得
	b := h.Object.GetPersistBytes()

	// ファイルに書き込み
	return os.WriteFile(h.Path, b, 0644)
}
