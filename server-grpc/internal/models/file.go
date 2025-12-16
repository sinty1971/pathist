package models

import (
	"errors"
	"os"

	grpcv1 "server-grpc/gen/grpc/v1"
	"server-grpc/internal/core"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// grpcv1.FileInfo 関連のヘルパー
// このファイルは proto ファイルから自動生成された gRPC メッセージ型を
// 補完するためのヘルパー関数や型を提供します。

// FileInfoEx は gRPC の File メッセージの拡張機能版です。
type File struct {
	*grpcv1.File
}

// NewFile File インスタンスを作成します。
func NewFile() *File {
	return &File{
		File: grpcv1.File_builder{}.Build(),
	}
}

// ParseFromPath は指定されたフルパスからファイル情報を解析して設定します
func (obj *File) ParseFromPath(target string) error {
	var err error

	// 絶対パスの正規化
	normalized, err := core.NormalizeAbsPath(target)
	if err != nil {
		return err
	}

	// ファイル情報の取得
	osFi, err := os.Stat(normalized)
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
	obj.SetTarget(normalized)
	obj.SetSize(osFi.Size())
	obj.SetModifiedTime(modifiedTime)
	return nil
}

// GetPersistData は永続化対象のオブジェクトを取得します
// Persistable インターフェースの実装
func (obj *File) GetPersistData() (map[string]any, error) {
	return map[string]any{
		"target": obj.GetTarget(),
	}, nil
}
