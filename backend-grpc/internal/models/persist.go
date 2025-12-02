package models

import (
	grpcv1 "backend-grpc/gen/grpc/v1"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Persist は永続化されるエンティティに必要な振る舞いを定義します。
type Persist struct {
	// PersistMap は永続化対象のオブジェクトを取得します。
	//	PersistMap() map[string]PersistFunc
}

type PersisInterface interface {
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

func PersistFileInfoSliceFunc(getter func() []*grpcv1.FileInfo, setter func([]*grpcv1.FileInfo)) PersistFunc {
	return PersistFunc{
		Getter: func() any {
			return getter()
		},
		Setter: func(val any) {
			if fileInfoSlice, ok := val.([]*grpcv1.FileInfo); ok {
				setter(fileInfoSlice)
			}
		},
	}
}

// LoadPersistData は永続化ファイルからデータを読み込みます。
func (h *Persist) LoadPersistData() error {

	// Persistable に型変換
	obj, ok := any(h).(PersisInterface)
	if !ok {
		return fmt.Errorf("Persistable インターフェースを実装していません")
	}

	// 永続化ファイルのフルパスを取得
	persistPath := obj.PersistPath()

	// ファイルを読み込み
	in, err := os.ReadFile(persistPath)
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return h.SavePersistData()
		}
	}

	// YAMLをデコード
	out, err := h.MarshalYAML()
	if err != nil {
		return fmt.Errorf("永続化データの取得に失敗しました: %w", err)
	}

	if err := yaml.Unmarshal(in, out); err != nil {
		return fmt.Errorf("YAMLのデコードに失敗しました: %w", err)
	}

	// エンティティにデコード結果を設定（既に out に反映されているが念のため）
	h.SetPersistData(out)

	return nil
}

// SavePersistData はデータを永続化ファイルに保存します。
func (h *Persist) SavePersistData() error {

	// Persistable に型変換
	obj, ok := any(h).(PersisInterface)
	if !ok {
		return fmt.Errorf("Persistable インターフェースを実装していません")
	}

	// 永続化ファイルパスの取得
	persistPath := obj.PersistPath()

	// データをYAMLにエンコード
	in, err := MarshalJSON(obj)
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

// SetPersistData は永続化対象のオブジェクトを設定するデフォルト実装です。
func (h *Persist) SetPersistData(persistData map[string]any) error {

	// Persistable に型変換
	obj, ok := any(h).(PersisInterface)
	if !ok {
		return fmt.Errorf("Persistable インターフェースを実装していません")
	}

	// 各フィールドの設定
	for k, f := range obj.PersistMap() {
		if val, exist := persistData[k]; exist {
			switch v := val.(type) {

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

// MarshalJSON シリアライズを行い、map[string]string{}を返します
func MarshalJSON() (map[string]string, error) {

	persistData := map[string]string{}
	for k, f := range obj.PersistMap() {
		persistData[k] = f.Getter()
	}
	return persistData, nil
}

// UnmarshalJSON は JSON からの復元を行うデフォルト実装です。
func UnmarshalJSON(obj Persist, raw []byte) error {
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
	obj, ok := any(h).(PersisInterface)
	if !ok {
		return nil, fmt.Errorf("Persistable インターフェースを実装していません")
	}

	persistData := map[string]string{}
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
