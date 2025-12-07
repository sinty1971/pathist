package models

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

// Persist は永続化されるエンティティに必要な振る舞いを定義します。
type Persist struct {
	// PersistFilename は永続化用のファイル名を保持します。
	PersistFilename string
}

type Persister interface {
	// GetPersistPath はファイルの永続化フルパスを取得します。
	PersistPath() string

	// GetPersistMap は永続化対象のメッセージを取得します。
	GetPersistMessage() proto.Message

	// SetPersistMessage は永続化対象のメッセージを設定します。
	SetPersistMessage(msg proto.Message)
}

// LoadPersistData は永続化ファイルからデータを読み込みます。
func (h *Persist) LoadPersistData() error {

	// PersistInterface 実装オブジェクトの取得
	obj, ok := any(h).(Persister)
	if !ok {
		return fmt.Errorf("PersistInterface インターフェースを実装していません")
	}

	// YAMLファイルからバイトデータを読み込む
	b, err := os.ReadFile(obj.PersistPath())
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return h.SavePersistData()
		}
	}

	m := obj.GetPersistMessage()
	proto.Unmarshal(b, m)
	obj.SetPersistMessage(m)

	return nil
}

// SavePersistData はデータを永続化ファイルに保存します。
func (h *Persist) SavePersistData() error {

	// PersistInterface オブジェクトの取得
	obj, ok := any(h).(Persister)
	if !ok {
		return fmt.Errorf("PersistInterface インターフェースを実装していません")
	}

	// データをYAMLにエンコード
	b, err := proto.Marshal(obj.GetPersistMessage())
	if err != nil {
		return fmt.Errorf("永続化データの取得に失敗しました: %w", err)
	}

	// ファイルに書き込み
	return os.WriteFile(obj.PersistPath(), b, 0644)
}
