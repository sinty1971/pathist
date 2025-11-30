package models

import (
	"errors"
	"os"

	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/core"

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
	targetPath, err = core.NormalizeAbsPath(targetPath)
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
	// フィールド名とセッター関数の対応表
	setterMap := map[string]func(string){
		"path": obj.SetPath,
	}

	// デフォルトのフィールド設定処理を呼び出し
	return core.DefaultSetPersistData(persistData, setterMap)
}

// MarshalJSON は JSON Serde を従来形式で行います。
func (obj FileInfo) MarshalJSON() ([]byte, error) {
	return core.DefaultMarshalJSON(obj.GetPersistData)
}

// UnmarshalJSON は JSON からの復元を行います。
func (obj *FileInfo) UnmarshalJSON(data []byte) error {
	return core.DefaultUnmarshalJSON(data, obj.SetPersistData)
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (obj FileInfo) MarshalYAML() (any, error) {
	return core.DefaultMarshalYAML(obj.GetPersistData)
}

// UnmarshalYAML は YAML からの復元を行います。
func (obj *FileInfo) UnmarshalYAML(unmarshal func(any) error) error {
	return core.DefaultUnmarshalYAML(unmarshal, obj.SetPersistData)
}
