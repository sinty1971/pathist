package core

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"gopkg.in/yaml.v3"
)

// PathistInterface は共通フィールドを持つモデルのインターフェースを定義します。
//
// 注意事項：protobuf メッセージを持っていることが前提となります。
// 以下のメソッドを実装する必要があります。
//
//   - GetMessage() protoreflect.Message
//   - GetModelName() string
//   - SetId(id string)
//   - GetTarget() string
type PathistInterface interface {
	// GetMessage はモデルの protobuf メッセージを取得します。
	GetMessage() protoreflect.Message

	// GetModelName はモデル名を取得します。
	GetModelName() string

	// SetId はモデルのIDを設定します。
	// protobuf メッセージ内の id フィールドを返すのが一般的です。
	SetId(id string)

	// GetTarget は永続化ファイルを保存するフルパスを取得します。
	// protobuf メッセージ内の target フィールドを返す実装が一般的です。
	GetTarget() string
}

// PathistModel はPathist共通モデルフィールドを提供します。
type PathistModel struct {
	// core はモデルインターフェースを保持します。
	PathistInterface

	// persistFilename は永続化ファイル名を保持します。
	persistFilename string
}

// NewPathistModel は Model インスタンスを作成します。
func NewPathistModel(persistFilename string) *PathistModel {
	model := &PathistModel{}

	// persistFilename が空文字列の場合はパニックを発生させる
	if filename := strings.TrimSpace(persistFilename); filename == "" {
		panic("persistFilename cannot be empty")
	} else {
		model.persistFilename = filename
	}

	return model
}

// SetMessageId はメッセージのIdを設定します。
//
// インスタンスの Targetフィールドが事前に設定されている必要があります。
func (m *PathistModel) SetMessageId() error {
	//
	if m.GetTarget() == "" {
		return errors.New("target is not set")
	}

	// ID 生成用テキストを作成してIDを生成
	text := m.GetModelName() + filepath.Base(m.GetTarget())
	id := GenerateIdFromString(text)
	m.SetId(id)
	return nil
}

// LoadPersists は永続化ファイルから永続化データのみを読み込みます。
// ファイル形式は YAML です。
func (m *PathistModel) LoadPersists() error {
	// 永続化ファイルのフルパスを取得
	persistPath := filepath.Join(m.GetTarget(), m.persistFilename)

	// YAMLファイルからテキストデータを読み込む
	yamltext, err := os.ReadFile(persistPath)
	if err != nil {
		// ファイルが存在しない場合は新規作成
		return m.SavePersists()
	}

	// YAMLファイルデータをJSONマップデータに変換
	jsonmap := &map[string]any{}
	if err = yaml.Unmarshal(yamltext, jsonmap); err != nil {
		return m.SavePersists()
	}

	// バイトデータをアンマーシャル
	err = m.FromJsonMap(jsonmap)
	if err != nil {
		return m.SavePersists()
	}
	return nil
}

// Save はデータを永続化ファイルに保存します。
// ファイル形式は YAML です。
func (m *PathistModel) SavePersists() error {
	// Persist 永続化バイトデータの取得
	jsonmap, err := m.ToJsonMap()
	if err != nil {
		return err
	}

	// JSONマップをYAMLデータに変換
	yamlBytes, err := yaml.Marshal(jsonmap)
	if err != nil {
		return err
	}

	// 永続化ファイルのフルパスを取得
	persistPath := filepath.Join(m.GetTarget(), m.persistFilename)

	// ファイルに書き込み
	return os.WriteFile(persistPath, yamlBytes, 0644)
}

// UpdatePersists は別の Persist インスタンスから persist_ フィールドのみを更新します
func (m *PathistModel) UpdatePersists(newData PathistInterface) error {
	// 引数チェック
	if newData == nil {
		return errors.New("newPersistable is nil")
	}

	// モデル名チェック
	if m.GetModelName() != newData.GetModelName() {
		return errors.New("model name mismatch in UpdatePersists")
	}

	// persist_ フィールドのみを更新
	curMsg := m.GetMessage()
	curFields := curMsg.Descriptor().Fields()
	newRefMsg := newData.GetMessage()
	for i := 0; i < curFields.Len(); i++ {
		curField := curFields.Get(i)
		newValue := newRefMsg.Get(curField)
		if !strings.HasPrefix(string(curField.Name()), "persist") {
			continue
		}
		curMsg.Set(curField, newValue)
	}
	return nil
}

// ToJsonMap は永続化用のフィールド値をJSONマップに変換します
func (m *PathistModel) ToJsonMap() (*map[string]any, error) {
	targetMsg := m.GetMessage()

	// camelCase キーで JSON にマーシャル
	jsonbytes, err := protojson.MarshalOptions{
		UseProtoNames:     true,
		EmitUnpopulated:   false,
		EmitDefaultValues: true,
	}.Marshal(targetMsg.Interface())
	if err != nil {
		return nil, err
	}

	jsonmap := &map[string]any{}
	json.Unmarshal(jsonbytes, jsonmap)
	for k := range *jsonmap {
		if !strings.HasPrefix(k, "persist") {
			delete(*jsonmap, k)
		}
	}

	return jsonmap, nil
}

// FromJsonMap はJSONマップを永続化用のフィールドに設定します
func (m *PathistModel) FromJsonMap(src *map[string]any) error {

	// JSONマップをバイトデータに変換
	srcJsonBytes, err := json.Marshal(*src)
	if err != nil {
		return err
	}

	// 代入先メッセージの取得
	destMsg := m.GetMessage()
	destDsc := destMsg.Descriptor()
	destFields := destDsc.Fields()

	// 一時的なメッセージを作成してJSONデータをアンマーシャル
	tempMsg := dynamicpb.NewMessage(destDsc)
	opts := protojson.UnmarshalOptions{AllowPartial: true}
	if err := opts.Unmarshal(srcJsonBytes, tempMsg); err != nil {
		log.Printf("Failed to unmarshal persist jsonmap: %v", err)
		return err
	}

	// persist_フィールドのみを元のメッセージにコピー
	for i := 0; i < destFields.Len(); i++ {
		f := destFields.Get(i)
		if !strings.HasPrefix(string(f.Name()), "persist") {
			continue
		}
		destMsg.Set(f, tempMsg.Get(f))
	}
	return nil
}
