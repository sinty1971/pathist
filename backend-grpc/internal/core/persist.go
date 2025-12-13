package core

import (
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Persistable は永続化インターフェースを定義します。
type Persistable interface {
	// GetPersistFolder は永続化ファイルを保存するフルパスを取得します。
	GetPersistFolder() string

	// ProtoReflect は対象メッセージのプロトリフレクションを取得します。
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

// getPersistPath は永続化ファイルのフルパスを取得します。
func (p *Persist) getPersistPath() string {
	return filepath.Join(p.target.GetPersistFolder(), p.persistFilename)
}

// Load は永続化ファイルからデータを読み込みます。
func (p *Persist) LoadPersists() error {

	// YAMLファイルからバイトデータを読み込む
	b, err := os.ReadFile(p.getPersistPath())
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return p.SavePersists()
		}
	}

	return p.target.SetPersists(b)
}

// Save はデータを永続化ファイルに保存します。
func (p *Persist) SavePersists() {

	// Persist 永続化バイトデータの取得
	b, err := p.target.GetPersists()
	if err != nil {
		return err
	}

	// ファイルに書き込み
	return os.WriteFile(p.getPersistPath(), b, 0644)
}

func (p *Persist) UpdatePersists(newPersist *Persist) error {
	refMsg := p.target.ProtoReflect()
	fields := refMsg.Descriptor().Fields()

	newRefMsg := newPersist.target.ProtoReflect()
	newMsgDes := newRefMsg.Descriptor()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := string(field.Name())
		if strings.HasPrefix(fieldName, "persist_") {
			newFieldDesc := newMsgDes.Fields().ByName(protoreflect.Name(fieldName))
			if newFieldDesc == nil {
				continue
			}
			refMsg.Set(field, newRefMsg.Get(newFieldDesc))
		}
	}
	return nil
}

// GetPersistValues は永続化用のフィールド値マップを取得します
func (p *Persist) MarshalPersistJSON() ([]byte, error) {
	ref := p.target.ProtoReflect()
	desc := ref.Descriptor()

	filtered := dynamicpb.NewMessage(desc)
	fields := desc.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		if strings.HasPrefix(string(f.Name()), "persist_") {
			filtered.Set(f, ref.Get(f))
		}
	}

	opts := protojson.MarshalOptions{Multiline: true, Indent: "  "}
	return opts.Marshal(filtered)
}

func (p *Persist) SetPersistValues(values map[string]protoreflect.Value) error {
	refMsg := p.target.ProtoReflect()
	msgDesc := refMsg.Descriptor()

	for fieldName, value := range values {
		if strings.HasPrefix(fieldName, "persist_") {
			fieldDesc := msgDesc.Fields().ByName(protoreflect.Name(fieldName))
			if fieldDesc == nil {
				continue
			}

			// フィールドの型に応じて値を設定
			refMsg.Set(fieldDesc, value)
		}
	}

	return nil
}

// GetUnpersistValues は永続化対象外のフィールド値マップを取得します
func (p *Persist) GetUnpersistValues() map[string]protoreflect.Value {
	// 例: 永続化対象外フィールドがない場合は空マップを返す
	return make(map[string]protoreflect.Value)
}
