package core

import (
	"encoding/json"
	"log"
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
	// GetTarget は永続化ファイルを保存するフルパスを取得します。
	// protobuf メッセージ内の target フィールドを返す実装が一般的です。
	GetTarget() string

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
	persistPath := filepath.Join(p.target.GetTarget(), p.persistFilename)

	// YAMLファイルからテキストデータを読み込む
	yamltext, err := os.ReadFile(persistPath)
	if err != nil {
		// ファイルが存在しない場合は新規作成
		return p.SavePersists()
	}

	// YAMLファイルデータをJSONマップデータに変換
	jsonmap := &map[string]any{}
	if err = yaml.Unmarshal(yamltext, jsonmap); err != nil {
		return p.SavePersists()
	}

	// バイトデータをアンマーシャル
	err = p.FromPersistJsonmap(jsonmap)
	if err != nil {
		return p.SavePersists()
	}
	return nil
}

// Save はデータを永続化ファイルに保存します。
// ファイル形式は YAML です。
func (p *Persist) SavePersists() error {
	// Persist 永続化バイトデータの取得
	jsonmap, err := p.ToPersistJsonmap()
	if err != nil {
		return err
	}

	// JSONマップをYAMLデータに変換
	yamlBytes, err := yaml.Marshal(jsonmap)
	if err != nil {
		return err
	}

	// 永続化ファイルのフルパスを取得
	persistPath := filepath.Join(p.target.GetTarget(), p.persistFilename)

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

// ToPersistJsonmap は永続化用のフィールド値をJSONマップに変換します
func (p *Persist) ToPersistJsonmap() (*map[string]any, error) {
	targetMsg := p.target.ProtoReflect()

	// camelCase キーで JSON にマーシャル
	opts := protojson.MarshalOptions{
		UseProtoNames:     true,
		EmitUnpopulated:   false,
		EmitDefaultValues: true,
	}
	jsonbytes, err := opts.Marshal(targetMsg.Interface())
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

// FromPersistJsonmap はJSONマップを永続化用のフィールドに設定します
func (p *Persist) FromPersistJsonmap(jsonmap *map[string]any) error {

	// persist_フィールドのみを抽出したマップを作成
	jsonbytes, err := json.Marshal(*jsonmap)
	if err != nil {
		return err
	}

	// 永続化用フィールドのみをアンマーシャル
	targetMsg := p.target.ProtoReflect()
	targetDsc := targetMsg.Descriptor()
	targetFields := targetDsc.Fields()

	// 一時的なメッセージを作成してJSONデータをアンマーシャル
	tempMsg := dynamicpb.NewMessage(targetDsc)
	opts := protojson.UnmarshalOptions{AllowPartial: true}
	if err := opts.Unmarshal(jsonbytes, tempMsg); err != nil {
		log.Printf("Failed to unmarshal persist jsonmap: %v", err)
		return err
	}

	// persist_フィールドのみを元のメッセージにコピー
	for i := 0; i < targetFields.Len(); i++ {
		f := targetFields.Get(i)
		log.Printf("Field Name: %s", f.Name())
		if !strings.HasPrefix(string(f.Name()), "persist") {
			continue
		}
		targetMsg.Set(f, tempMsg.Get(f))
	}
	return nil
}
