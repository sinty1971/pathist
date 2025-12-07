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

// NewFileInfo FieleInfo インスタンスを作成します。
func NewFileInfo() *FileInfo {
	return &FileInfo{
		FileInfo: grpcv1.FileInfo_builder{}.Build(),
	}
}

// ParseFromPath は指定されたフルパスからファイル情報を解析して設定します
func (obj *FileInfo) ParseFromPath(fullpath string) error {
	var err error

	// 絶対パスのクリーン化
	fullpath, err = core.NormalizeAbsPath(fullpath)
	if err != nil {
		return err
	}

	// ファイル情報の取得
	osFi, err := os.Stat(fullpath)
	if err != nil {
		return err
	}

	// 最終更新時刻の取得
	osModTime := osFi.ModTime()
	if osModTime.IsZero() {
		return errors.New("ファイル最終更新日の取得に失敗しました: file modification time is zero")
	}
	modifiedTime := timestamppb.New(osModTime)

	// フィールドの設定
	obj.SetPath(fullpath)
	obj.SetIsDirectory(osFi.IsDir())
	obj.SetSize(osFi.Size())
	obj.SetModifiedTime(modifiedTime)
	return nil
}

// GetPersistData は永続化対象のオブジェクトを取得します
// Persistable インターフェースの実装
func (obj *FileInfo) GetPersistData() (map[string]any, error) {
	return map[string]any{
		"path": obj.GetPath(),
	}, nil
}
