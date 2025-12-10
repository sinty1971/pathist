package core

import (
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Persistable は永続化インターフェースを定義します。
type Persistable interface {
	// GetManagedFolder は永続化ファイルを保存するフルパスを取得します。
	GetManagedFolder() string

	// GetPersistMap は永続化対象のメッセージを取得します。
	GetPersistBytes() ([]byte, error)

	// SetPersistMessage は永続化対象のメッセージを設定します。
	SetPersistBytes(b []byte) error

	// ProtoReflect は対象メッセージのプロトリフレクションを取得します。
	ProtoReflect() protoreflect.Message
}

// Persist は永続化されるエンティティに必要な振る舞いを定義します。
type Persist struct {
	target          Persistable
	persistFilename string
}

func NewPersister(target Persistable, persistFilename string) *Persist {
	return &Persist{
		target:          target,
		persistFilename: persistFilename,
	}
}

// getPersistPath は永続化ファイルのフルパスを取得します。
func (p *Persist) getPersistPath() string {
	return filepath.Join(p.target.GetManagedFolder(), p.persistFilename)
}

// Load は永続化ファイルからデータを読み込みます。
func (p *Persist) Load() error {

	// YAMLファイルからバイトデータを読み込む
	b, err := os.ReadFile(p.getPersistPath())
	if err != nil {
		// ファイルが存在しない場合は一度 Save を呼び出してファイルを作成する
		if os.IsNotExist(err) {
			return p.Save()
		}
	}

	return p.target.SetPersistBytes(b)
}

// Save はデータを永続化ファイルに保存します。
func (p *Persist) Save() error {

	// Persist 永続化バイトデーターの取得
	b, err := p.target.GetPersistBytes()
	if err != nil {
		return err
	}

	// ファイルに書き込み
	return os.WriteFile(p.getPersistPath(), b, 0644)
}

// GetPersistValues は永続化用のフィールド値マップを取得します
func (p *Persist) GetPersistValues() map[string]protoreflect.Value {
	refMsg := p.target.ProtoReflect()
	fields := refMsg.Descriptor().Fields()

	result := make(map[string]protoreflect.Value)
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := string(field.Name())
		if strings.HasPrefix(fieldName, "persist_") {
			result[fieldName] = refMsg.Get(field)
	}
	return result
}

// GetUnpersistValues は永続化対象外のフィールド値マップを取得します
func (p *Persist) GetUnpersistValues() map[string]protoreflect.Value {
	// 例: 永続化対象外フィールドがない場合は空マップを返す
	return make(map[string]protoreflect.Value)
}
