package ext

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Persistable は永続化されるエンティティに必要な振る舞いを定義します。
type Persistable interface {
	// GetPersistPath はファイルの永続化フルパスを取得します。
	GetPersistPath() string

	// GetPersistData は永続化対象のオブジェクトを取得します。
	GetPersistData() map[string]any

	// SetPersistData は永続化対象のオブジェクトを設定します。
	SetPersistData(map[string]any) error
}

// PersistService はファイルベースのデータ永続化のサービスを提供します。
type PersistService[T Persistable] struct {
	persistData T
}

// PersistService のコンストラクタ
func CreatePersistService[T Persistable](persistData T) *PersistService[T] {
	return &PersistService[T]{
		persistData: persistData,
	}
}

// Load は永続化ファイルからデータを読み込みます。
func (srv *PersistService[T]) Load() error {

	// 永続化ファイルのフルパスを取得
	persistPath := srv.persistData.GetPersistPath()

	// ファイルを読み込み
	in, err := os.ReadFile(persistPath)
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return srv.Save()
		}
	}

	// YAMLをデコード
	out := srv.persistData.GetPersistData()

	if err := yaml.Unmarshal(in, out); err != nil {
		return fmt.Errorf("YAMLのデコードに失敗しました: %w", err)
	}

	// エンティティにデコード結果を設定（既に out に反映されているが念のため）
	srv.persistData.SetPersistData(out)

	return nil
}

// Save はデータを永続化ファイルに保存します。
func (srv *PersistService[T]) Save() error {
	// 永続化ファイルパスの取得
	persistPath := srv.persistData.GetPersistPath()

	// データをYAMLにエンコード
	in := srv.persistData.GetPersistData()
	out, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	// ファイルに書き込み
	return os.WriteFile(persistPath, out, 0644)
}
