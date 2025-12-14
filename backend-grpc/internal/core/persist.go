package core

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Persistable は永続化インターフェースを定義します。
// protoreflect.Message を返す必要があるため、protobuf メッセージをラップする構造体で実装する必要があります。
type Persistable interface {
	// GetPersistDir は永続化ファイルを保存するフルパスを取得します。
	GetPersistDir() string

	// ProtoReflect は対象メッセージの protoreflect.Message を取得します。
	ProtoReflect() protoreflect.Message
}

// Persist は永続化されるエンティティに必要な振る舞いを定義します。
type Persist struct {
	target          Persistable
	persistFilename string
}

// NewPersister は Persist インスタンスを作成します。
// target を初期化する方法はこの関数のみです。
func NewPersister(target Persistable, persistFilename string) *Persist {
	// target が nil の場合はパニックを発生させる
	if target == nil {
		panic("Persistable target cannot be nil")
	}

	// persistFilename が空文字列の場合はパニックを発生させる
	if strings.TrimSpace(persistFilename) == "" {
		panic("persistFilename cannot be empty")
	}

	return &Persist{
		target:          target,
		persistFilename: persistFilename,
	}
}

// Load は永続化ファイルからデータを読み込みます。
// ファイル形式は YAML です。
func (p *Persist) LoadPersists() error {
	// 永続化ファイルのフルパスを取得
	persistPath := filepath.Join(p.target.GetPersistDir(), p.persistFilename)

	// YAMLファイルからバイトデータを読み込む
	yamlBytes, err := os.ReadFile(persistPath)
	if err != nil {
		// ファイルが存在しない場合は新規作成
		return p.SavePersists()
	}

	// YAMLデータをJSONデータに変換
	jsonBytes, err := YAMLToJSON(yamlBytes)
	if err != nil {
		return p.SavePersists()
	}

	// バイトデータをアンマーシャル
	err = p.UnmarshalPersistJSON(jsonBytes)
	if err != nil {
		return p.SavePersists()
	}
	return nil
}

// Save はデータを永続化ファイルに保存します。
// ファイル形式は YAML です。
func (p *Persist) SavePersists() error {
	// Persist 永続化バイトデータの取得
	jsonBytes, err := p.MarshalPersistJSON()
	if err != nil {
		return err
	}

	// JSONデータをYAMLデータに変換
	yamlBytes, err := JSONToYAML(jsonBytes)
	if err != nil {
		return err
	}

	// 永続化ファイルのフルパスを取得
	persistPath := filepath.Join(p.target.GetPersistDir(), p.persistFilename)

	// ファイルに書き込み
	return os.WriteFile(persistPath, yamlBytes, 0644)
}

// UpdatePersists は別の Persist インスタンスから persist_ フィールドのみを更新します
func (p *Persist) UpdatePersists(newPersist *Persist) error {
	targetMsg := p.target.ProtoReflect()
	targetFields := targetMsg.Descriptor().Fields()

	newRefMsg := newPersist.target.ProtoReflect()
	newMsgDes := newRefMsg.Descriptor()

	for i := 0; i < targetFields.Len(); i++ {
		field := targetFields.Get(i)
		fieldName := string(field.Name())
		newFieldDesc := newMsgDes.Fields().ByName(protoreflect.Name(fieldName))
		if newFieldDesc == nil {
			continue
		}
		targetMsg.Set(field, newRefMsg.Get(newFieldDesc))
	}
	return nil
}

// MarshalPersistJSON は永続化用のフィールド値をJSONデータに変換します
func (p *Persist) MarshalPersistJSON() ([]byte, error) {
	targetMsg := p.target.ProtoReflect()
	targetDsc := targetMsg.Descriptor()

	// 永続化対象フィールドのみを抽出した新しいメッセージを作成
	filtered := dynamicpb.NewMessage(targetDsc)
	fields := targetDsc.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		if !strings.HasPrefix(string(f.Name()), "persist_") {
			continue
		}
		filtered.Set(f, targetMsg.Get(f))
	}

	// JSON にシリアライズ
	opts := protojson.MarshalOptions{Multiline: true, Indent: "  ", AllowPartial: false, EmitDefaultValues: true}
	return opts.Marshal(filtered)
}

// UnmarshalPersistJSON はJSONデータを永続化用のフィールドに設定します
func (p *Persist) UnmarshalPersistJSON(b []byte) error {
	targetMsg := p.target.ProtoReflect()
	targetDsc := targetMsg.Descriptor()

	// 一時的なメッセージを作成してJSONをアンマーシャル
	jsonMsg := dynamicpb.NewMessage(targetDsc)
	opts := protojson.UnmarshalOptions{AllowPartial: false}
	if err := opts.Unmarshal(b, jsonMsg); err != nil {
		return err
	}

	// persist_フィールドのみを元のメッセージにコピー
	fields := targetDsc.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		if !strings.HasPrefix(string(f.Name()), "persist_") {
			continue
		}
		targetMsg.Set(f, jsonMsg.Get(f))
	}
	return nil
}

func (p *Persist) UnmarshalPersistYAML(b []byte) error {
	targetMsg := p.target.ProtoReflect()
	targetDsc := targetMsg.Descriptor()

	// 一時的なメッセージを作成してJSONとしてアンマーシャル（YAML形式のJSONとして扱う）
	yamlMsg := dynamicpb.NewMessage(targetDsc)
	opts := protojson.UnmarshalOptions{AllowPartial: true}
	if err := opts.Unmarshal(b, yamlMsg); err != nil {
		return err
	}

	// persist_フィールドのみを元のメッセージにコピー
	fields := targetDsc.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		if strings.HasPrefix(string(f.Name()), "persist_") {
			targetMsg.Set(f, yamlMsg.Get(f))
		}
	}
	return nil
}

// JSONToYAML はJSONバイト配列をYAMLバイト配列に変換します
func JSONToYAML(jsonData []byte) ([]byte, error) {
	// JSONをmapにアンマーシャル
	var data interface{}
	if err := yaml.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	// mapをYAMLにマーシャル
	return yaml.Marshal(data)
}

// YAMLToJSON はYAMLバイト配列をJSONバイト配列に変換します
func YAMLToJSON(yamlData []byte) ([]byte, error) {
	// YAMLをmapにアンマーシャル
	var data interface{}
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		return nil, err
	}

	// mapをJSONにマーシャル（encoding/jsonを使用）
	return yaml.Marshal(data)
}
