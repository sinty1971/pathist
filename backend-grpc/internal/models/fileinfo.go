package models

import (
	"errors"
	"os"

	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/ext"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// grpcv1.FileInfo 関連のヘルパー
// このファイルは proto ファイルから自動生成された gRPC メッセージ型を
// 補完するためのヘルパー関数や型を提供します。

// FileInfoEx は gRPC の FileInfo メッセージの拡張機能版です。
type FileInfo struct {
	*grpcv1.FileInfo
}

// NewFileInfo はファイルのフルパスから FileInfo を作成します。
func NewFileInfo(targetPath string) (*FileInfo, error) {
	var err error

	// 絶対パスのクリーン化
	targetPath, err = ext.NormalizeAbsPath(targetPath)
	if err != nil {
		return nil, err
	}

	// ファイル情報の取得
	osFi, err := os.Stat(targetPath)
	if err != nil {
		return nil, err
	}

	// 最終更新時刻の取得
	osModTime := osFi.ModTime()
	if osModTime.IsZero() {
		return nil, errors.New("ファイル最終更新日の取得に失敗しました: file modification time is zero")
	}

	return &FileInfo{
		FileInfo: grpcv1.FileInfo_builder{
			Path:         targetPath,
			IsDirectory:  osFi.IsDir(),
			Size:         osFi.Size(),
			ModifiedTime: timestamppb.New(osModTime),
		}.Build()}, nil
}

// GetPersistData は永続化対象のオブジェクトを取得します
// Persistable インターフェースの実装
func (obj *FileInfo) GetPersistData() (map[string]any, error) {
	return map[string]any{
		"path": obj.GetPath(),
	}, nil
}

// SetPersistData は永続化対象のオブジェクトを設定します
// Persistable インターフェースの実装
func (obj *FileInfo) SetPersistData(persistData map[string]any) error {
	// マップとセッター関数の対応表
	setterMap := map[string]func(string){
		"path": obj.SetPath,
	}

	// デフォルトの文字列フィールド設定処理を呼び出し
	return ext.DefaultSetPersistData(persistData, setterMap)
}

// MarshalJSON は JSON Serde を従来形式で行います。
func (obj FileInfo) MarshalJSON() ([]byte, error) {
	return ext.DefaultMarshalJSON(obj.GetPersistData)
}

// UnmarshalJSON は JSON からの復元を行います。
func (obj *FileInfo) UnmarshalJSON(data []byte) error {
	return ext.DefaultUnmarshalJSON(data, obj.SetPersistData)
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (obj FileInfo) MarshalYAML() (any, error) {
	return ext.DefaultMarshalYAML(obj.GetPersistData)
}

// UnmarshalYAML は YAML からの復元を行います。
func (obj *FileInfo) UnmarshalYAML(unmarshal func(any) error) error {
	return ext.DefaultUnmarshalYAML(unmarshal, obj.SetPersistData)
}

// applyEncoding は別の FileInfo の内容を適用します。
func (obj *FileInfo) applyEncoding(enc FileInfo) error {
	obj.FileInfo.SetPath(enc.GetPath())
	obj.FileInfo.SetIsDirectory(enc.GetIsDirectory())
	obj.FileInfo.SetSize(enc.GetSize())
	obj.FileInfo.SetModifiedTime(enc.GetModifiedTime())

	return nil
}
