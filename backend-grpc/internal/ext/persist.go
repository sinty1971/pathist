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

	// GetPersistInfo は永続化対象のオブジェクトを取得します。
	GetPersistInfo() any

	// SetPersistInfo は永続化対象のオブジェクトを設定します。
	SetPersistInfo(any)
}

// PersistService はファイルベースのデータ永続化のサービスを提供します。
type PersistService[T Persistable] struct {
	persistableObject T
}

// PersistService のコンストラクタ
func CreatePersistService[T Persistable](persistable T) *PersistService[T] {
	return &PersistService[T]{
		persistableObject: persistable,
	}
}

// LoadPersistInfo は永続化ファイルからデータを読み込みます。
func (s *PersistService[T]) LoadPersistInfo() error {

	// 永続化ファイルのフルパスを取得
	persistPath := s.persistableObject.GetPersistPath()

	// ファイルを読み込み
	in, err := os.ReadFile(persistPath)
	if err != nil {
		// ファイルが存在しない場合は一度 SavePersittInfo を呼び出して初期ファイルを作成する
		if os.IsNotExist(err) {
			if err := s.SavePersistInfo(); err != nil {
				return fmt.Errorf("初期永続化ファイルの作成に失敗しました: %w", err)
			}
			// 再度読み込みを試みる
			_, err = os.ReadFile(persistPath)
			if err != nil {
				return fmt.Errorf("永続化ファイルの読み込みに失敗しました: %w", err)
			}
			return err
		}
	}

	// YAMLをデコード
	out := s.persistableObject.GetPersistInfo()

	if err := yaml.Unmarshal(in, out); err != nil {
		return fmt.Errorf("YAMLのデコードに失敗しました: %w", err)
	}

	// エンティティにデコード結果を設定（既に out に反映されているが念のため）
	s.persistableObject.SetPersistInfo(out)

	return nil
}

// SavePersistInfo はデータを永続化ファイルに保存します。
func (s *PersistService[T]) SavePersistInfo() error {
	// 永続化ファイルパスの取得
	persistPath := s.persistableObject.GetPersistPath()

	// データをYAMLにエンコード
	in := s.persistableObject.GetPersistInfo()
	out, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	// ファイルに書き込み
	return os.WriteFile(persistPath, out, 0644)
}
