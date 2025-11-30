package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Persistable は永続化されるエンティティに必要な振る舞いを定義します。
type Persistable interface {
	// GetPersistPath はファイルの永続化フルパスを取得します。
	GetPersistPath() string

	// GetPersistData は永続化対象のオブジェクトを取得します。
	GetPersistData() (map[string]any, error)

	// SetPersistData は永続化対象のオブジェクトを設定します。
	SetPersistData(map[string]any) error
}

// LoadPersistData は永続化ファイルからデータを読み込みます。
func LoadPersistData(obj Persistable) error {

	// 永続化ファイルのフルパスを取得
	persistPath := obj.GetPersistPath()

	// ファイルを読み込み
	in, err := os.ReadFile(persistPath)
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return SavePersistData(obj)
		}
	}

	// YAMLをデコード
	out, err := obj.GetPersistData()
	if err != nil {
		return fmt.Errorf("永続化データの取得に失敗しました: %w", err)
	}

	if err := yaml.Unmarshal(in, out); err != nil {
		return fmt.Errorf("YAMLのデコードに失敗しました: %w", err)
	}

	// エンティティにデコード結果を設定（既に out に反映されているが念のため）
	obj.SetPersistData(out)

	return nil
}

// SavePersistData はデータを永続化ファイルに保存します。
func SavePersistData(obj Persistable) error {
	// 永続化ファイルパスの取得
	persistPath := obj.GetPersistPath()

	// データをYAMLにエンコード
	in, err := obj.GetPersistData()
	if err != nil {
		return fmt.Errorf("永続化データの取得に失敗しました: %w", err)
	}

	// YAMLエンコード
	out, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	// ファイルに書き込み
	return os.WriteFile(persistPath, out, 0644)
}

// DefaultSetPersistData は永続化対象のオブジェクトを設定するデフォルト実装です。
func DefaultSetPersistData(persistData map[string]any, setterMap map[string]func(string)) error {

	// 各フィールドの設定
	for key, setter := range setterMap {
		if val, exist := persistData[key]; exist {
			switch v := val.(type) {

			// 文字列フィールドの処理
			case string:
				setter(v)

			// Json文字列として格納されている場合の処理
			default:
				if json, err := json.Marshal(v); err == nil {
					setter(string(json))

				} else {
					log.Printf("Json変換ができませんでした: %v", v)
				}
			}
		} else {
			log.Printf("%s の項目の処理に失敗しました。", key)
		}
	}

	return nil
}

// DefaultMarshalJSON は JSON 用のシリアライズを行うデフォルト実装です。
func DefaultMarshalJSON(getter func() (map[string]any, error)) ([]byte, error) {
	data, err := getter()
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

// DefaultUnmarshalJSON は JSON からの復元を行うデフォルト実装です。
func DefaultUnmarshalJSON(data []byte, setter func(map[string]any) error) error {
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	return setter(v)
}

func DefaultMarshalYAML(getter func() (map[string]any, error)) (any, error) {
	data, err := getter()
	if err != nil {
		return nil, err
	}
	// Getterを使って明示的にマップ化
	return data, nil
}

// DefaultUnmarshalYAML は YAML からの復元を行うデフォルト実装です。
func DefaultUnmarshalYAML(unmarshal func(any) error, setter func(map[string]any) error) error {
	var v map[string]any
	if err := unmarshal(&v); err != nil {
		return err
	}

	return setter(v)
}
