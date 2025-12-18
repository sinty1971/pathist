package core

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
	"gopkg.in/yaml.v3"
)

// Pathist はPathist共通フィールドを提供します。
type Pathist struct {
	// pathistableModel はPathistable インターフェイスを満たすモデルです。
	pathistableModel Pathistable

	// persistFilename は永続化ファイル名を保持します、変更不可です。
	persistFilename string

	// modelNameId は proto メッセージ名の一意な識別子です、指定及び変更不可です。
	modelNameId string
}

// Pathistable は共通フィールドを持つモデルのインターフェースを定義します。
//   - protobuf メッセージを持っていることが前提となります。
type Pathistable interface {
	// GetProtoMessage はモデルの protobuf メッセージを取得します。
	GetProtoMessage() proto.Message

	// SetId はモデルのIDを設定します。
	//  - proto メッセージ内の id フィールドを返すのが一般的であり実装不要です。
	SetId(id string)

	// GetPathistFolder は永続化ファイルを保存するフルパスを取得します。
	//  - proto メッセージ内の pathist_folder フィールドを返す実装が一般的であり実装不要です。
	GetPathistFolder() string
}

// NewPathist は Pathist インスタンスを作成します。
func NewPathist(model Pathistable, persistFilename string) *Pathist {
	// モデルのフルネームから一意なIDを生成
	fullname := model.GetProtoMessage().ProtoReflect().Descriptor().FullName()
	id := GenerateIdFromString(string(fullname))

	// persistFilename が空文字列の場合はパニックを発生させる
	filename := strings.TrimSpace(persistFilename)
	if filename == "" {
		panic("persistFilename cannot be empty")
	}

	// インスタンス作成
	return &Pathist{
		pathistableModel: model,
		persistFilename:  filename,
		modelNameId:      id,
	}
}

// G GenerateId はメッセージのIdを設定します。
//
// インスタンスの PathistFolderフィールドが事前に設定されている必要があります。
func (p *Pathist) GenerateId() (string, error) {
	//
	if p.pathistableModel.GetPathistFolder() == "" {
		return "", errors.New("pathist_folder is not set")
	}
	// ID 生成用テキストを作成してIDを生成
	text := p.modelNameId + filepath.Base(p.pathistableModel.GetPathistFolder())
	return GenerateIdFromString(text), nil
}

// LoadPersists は永続化ファイルから永続化データのみを読み込みます。
// ファイル形式は YAML です。
func (p *Pathist) LoadPersists() error {
	// YAMLファイルからテキストデータを読み込む
	yamltext, err := os.ReadFile(p.getPersistPath())
	if err != nil {
		// ファイルが存在しない場合は新規作成
		return p.SavePersists()
	}

	// YAMLファイルデータをJSONマップデータに変換
	jsonmap := &map[string]any{}
	if err = yaml.Unmarshal(yamltext, jsonmap); err != nil {
		return p.SavePersists()
	}

	// JSONマップデータから永続化データを取り込む
	err = p.SetPersistsFrom(jsonmap)
	if err != nil {
		return p.SavePersists()
	}
	return nil
}

// Save はデータを永続化ファイルに保存します。
// ファイル形式は YAML です。
func (p *Pathist) SavePersists() error {
	// JSONマップの取得
	jsonmap, err := p.GetPersistJsonMap()
	if err != nil {
		return err
	}

	// JSONマップをYAMLデータに変換
	yamlBytes, err := yaml.Marshal(jsonmap)
	if err != nil {
		return err
	}

	// ファイルに書き込み
	return os.WriteFile(p.getPersistPath(), yamlBytes, 0644)
}

// getPersistPath は永続化ファイルのフルパスを取得します。
func (p *Pathist) getPersistPath() string {
	return filepath.Join(p.pathistableModel.GetPathistFolder(), p.persistFilename)
}

// ImportPersists は別の Persist インスタンスから persist_ フィールドのみ取り込みます。
func (p *Pathist) ImportPersists(src *Pathist) error {
	// 引数チェック
	if src == nil {
		return errors.New("Pathistable src is nil")
	}

	// モデル名チェック
	if p.modelNameId != src.modelNameId {
		return errors.New("model name mismatch in src Pathistable")
	}

	// persist_ フィールドのみを更新
	destMsg := p.pathistableModel.GetProtoMessage().ProtoReflect()
	destFields := destMsg.Descriptor().Fields()
	srcRefMsg := src.pathistableModel.GetProtoMessage().ProtoReflect()
	for i := 0; i < destFields.Len(); i++ {
		f := destFields.Get(i)
		v := srcRefMsg.Get(f)
		if !strings.HasPrefix(string(f.Name()), "persist") {
			continue
		}
		destMsg.Set(f, v)
	}
	return nil
}

// GetPersistJsonMap は永続化用のフィールド値をJSONマップに変換します
func (p *Pathist) GetPersistJsonMap() (*map[string]any, error) {
	// camelCase キーで JSON にマーシャル
	jsonbytes, err := protojson.MarshalOptions{
		UseProtoNames:     true,
		EmitUnpopulated:   false,
		EmitDefaultValues: true,
	}.Marshal(p.pathistableModel.GetProtoMessage())
	if err != nil {
		return nil, err
	}

	// persist_ フィールドのみを抽出した JSON マップを作成
	jsonmap := &map[string]any{}
	json.Unmarshal(jsonbytes, jsonmap)
	for k := range *jsonmap {
		if !strings.HasPrefix(k, "persist") {
			delete(*jsonmap, k)
		}
	}

	return jsonmap, nil
}

// SetPersistsFrom はJSONマップを永続化用のフィールドに設定します
func (p *Pathist) SetPersistsFrom(jsonmap *map[string]any) error {

	// JSONマップをバイトデータに変換
	bytes, err := json.Marshal(*jsonmap)
	if err != nil {
		return err
	}

	// 代入先メッセージの取得
	destMsg := p.pathistableModel.GetProtoMessage()
	destDsc := destMsg.ProtoReflect().Descriptor()
	destFields := destDsc.Fields()

	// 一時的なメッセージを作成してJSONデータをアンマーシャル
	tempMsg := dynamicpb.NewMessage(destDsc)
	opts := protojson.UnmarshalOptions{AllowPartial: true}
	if err := opts.Unmarshal(bytes, tempMsg); err != nil {
		log.Printf("Failed to unmarshal persist jsonmap: %v", err)
		return err
	}

	// persist_フィールドのみを元のメッセージにコピー
	for i := 0; i < destFields.Len(); i++ {
		f := destFields.Get(i)
		if !strings.HasPrefix(string(f.Name()), "persist") {
			continue
		}
		destMsg.ProtoReflect().Set(f, tempMsg.Get(f))
	}
	return nil
}
