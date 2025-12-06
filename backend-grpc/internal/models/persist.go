package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v2"
)

// Persist は永続化されるエンティティに必要な振る舞いを定義します。
type Persist struct {
	// PersistFilename は永続化用のファイル名を保持します。
	PersistFilename string
}

type PersistInterface interface {
	// GetPersistPath はファイルの永続化フルパスを取得します。
	PersistPath() string

	// PersistMap は永続化対象のフィールドマップを取得します。
	PersistMap() map[string]PersistFunc
}

// PersistFunc は永続化フィールドのゲッター・セッター関数を定義します。
type PersistFunc struct {
	Getter func() any
	Setter func(any)
}

// LoadPersistData は永続化ファイルからデータを読み込みます。
func (h *Persist) LoadPersistData() error {

	// PersistInterface オブジェクトの取得
	obj, ok := any(h).(PersistInterface)
	if !ok {
		return fmt.Errorf("PersistInterface インターフェースを実装していません")
	}

	// ファイルを読み込み
	in, err := os.ReadFile(obj.PersistPath())
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return h.SavePersistData()
		}
	}

	// YAMLデコード
	var persistData map[string]any
	if err := yaml.Unmarshal(in, &persistData); err != nil {
		return fmt.Errorf("永続化データの読み込みに失敗しました: %w", err)
	}

	// 各フィールドの設定
	for k, f := range obj.PersistMap() {
		if anyVal, exist := persistData[k]; exist {
			switch v := anyVal.(type) {

			// 文字列フィールドの処理
			case string:
				f.Setter(v)

			// Json文字列として格納されている場合の処理
			default:
				if json, err := json.Marshal(v); err == nil {
					f.Setter(string(json))

				} else {
					log.Printf("Json変換ができませんでした: %v", v)
				}
			}
		} else {
			log.Printf("%s の項目の処理に失敗しました。", k)
		}
	}

	return nil
}

// SavePersistData はデータを永続化ファイルに保存します。
func (h *Persist) SavePersistData() error {

	// PersistInterface オブジェクトの取得
	obj, ok := any(h).(PersistInterface)
	if !ok {
		return fmt.Errorf("PersistInterface インターフェースを実装していません")
	}

	// データをYAMLにエンコード
	in, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("永続化データの取得に失敗しました: %w", err)
	}

	// YAMLエンコード
	out, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	// ファイルに書き込み
	return os.WriteFile(obj.PersistPath(), out, 0644)
}

// MarshalJSON シリアライズを行い、map[string]string{}を返します
func (h *Persist) MarshalJSON() ([]byte, error) {

	// Persistable オブジェクトに型変換
	obj, ok := any(h).(PersistInterface)
	if !ok {
		return nil, fmt.Errorf("Persistable インターフェースを実装していません")
	}

	persistData := make(map[string]string, len(obj.PersistMap()))
	for k, f := range obj.PersistMap() {
		anyVal := f.Getter()
		// 値を文字列に変換
		switch v := anyVal.(type) {
		case string:
			persistData[k] = v
		default:
			data, err := json.Marshal(v)
			if err != nil {
				log.Printf("Json変換ができませんでした: %v", anyVal)
				continue
			}
			persistData[k] = string(data)
		}

	}
	return json.Marshal(persistData)
}

// UnmarshalJSON は JSON からの復元を行うデフォルト実装です。
func (h *Persist) UnmarshalJSON(raw []byte) error {

	// Persistable オブジェクトに型変換
	obj, ok := any(h).(PersistInterface)
	if !ok {
		return fmt.Errorf("Persistable インターフェースを実装していません")
	}

	// 永続化データの取得
	persistData := map[string]string{}
	if err := json.Unmarshal(raw, &persistData); err != nil {
		return err
	}

	for k, v := range persistData {
		if f, exist := obj.PersistMap()[k]; exist {
			f.Setter(v)
		}
	}
	return nil
}

// MarshalYAML は YAML 用のシリアライズを行うデフォルト実装です。
func (h *Persist) MarshalYAML() (any, error) {

	// Persistable オブジェクトに型変換
	obj, ok := any(h).(PersistInterface)
	if !ok {
		return nil, fmt.Errorf("Persistable インターフェースを実装していません")
	}

	persistData := make(map[string]string, len(obj.PersistMap()))
	for k, f := range obj.PersistMap() {
		persistData[k] = f.Getter()
	}
	// Getterを使って明示的にマップ化
	return persistData, nil
}

// UnmarshalYAML は YAML からの復元を行うデフォルト実装です。
func UnmarshalYAML(obj Persist, unmarshal func(any) error) error {

	// 永続化データの取得
	var persist map[string]any
	if err := unmarshal(&persist); err != nil {
		return err
	}

	for k, v := range persist {
		f, exist := obj.PersistMap()[k]
		if !exist {
			continue
		}
		setter := f.Setter

		// 値の文字列化
		var strVal string
		switch val := v.(type) {
		case string:
			strVal = val
		default:
			data, err := json.Marshal(val)
			if err != nil {
				log.Printf("Json変換ができませんでした: %v", val)
				continue
			}
			strVal = string(data)
		}

		// セッター関数の呼び出し
		setter(strVal)
	}

	return nil
}

// PersistStringFunc は文字列型のゲッター・セッター関数を持つ PersistFunc を作成します。
func PersistStringFunc(getter func() string, setter func(string)) PersistFunc {
	return PersistFunc{
		Getter: func() any {
			return getter()
		},
		Setter: func(val any) {
			if strVal, ok := val.(string); ok {
				setter(strVal)
			}
		},
	}
}

func PersistStringSliceFunc(getter func() []string, setter func([]string)) PersistFunc {
	return PersistFunc{
		Getter: func() any {
			return getter()
		},
		Setter: func(val any) {
			if stringSlice, ok := val.([]string); ok {
				setter(stringSlice)
			}
		},
	}
}

// PersistTimestampFunc はタイムスタンプ文字列型のゲッター・セッター関数を持つ PersistFunc を作成します。
//
//	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
func PersistTimestampFunc(getter func() *timestamppb.Timestamp, setter func(*timestamppb.Timestamp)) PersistFunc {
	return PersistFunc{
		Getter: func() any {
			return getter()
		},
		Setter: func(valAny any) {
			if valTimestamp, ok := valAny.(*timestamppb.Timestamp); ok {
				setter(valTimestamp)
			}
		},
	}
}
