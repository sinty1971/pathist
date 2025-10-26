package models

import (
	"encoding/json"
	"errors"
	"os"

	grpcv1 "grpc-backend/gen/grpc/v1"
	"grpc-backend/internal/utils"

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
	targetPath, err = utils.CleanAbsPath(targetPath)
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

// MarshalYAML は YAML 用のシリアライズを行います。
func (fi FileInfo) MarshalYAML() (any, error) {
	if fi.FileInfo == nil {
		return FileInfo{}, nil
	}
	return fi, nil
}

// UnmarshalYAML は YAML からの復元を行います。
func (fi *FileInfo) UnmarshalYAML(unmarshal func(any) error) error {
	var enc FileInfo
	if err := unmarshal(&enc); err != nil {
		return err
	}
	return fi.applyEncoding(enc)
}

// MarshalJSON は JSON Serde を従来形式で行います。
func (fi FileInfo) MarshalJSON() ([]byte, error) {
	if fi.FileInfo == nil {
		return json.Marshal(FileInfo{})
	}
	return json.Marshal(fi)
}

// UnmarshalJSON は JSON からの復元を行います。
func (fi *FileInfo) UnmarshalJSON(data []byte) error {
	var enc FileInfo
	if err := json.Unmarshal(data, &enc); err != nil {
		return err
	}
	return fi.applyEncoding(enc)
}

func (fi *FileInfo) applyEncoding(enc FileInfo) error {
	fi.FileInfo.SetPath(enc.GetPath())
	fi.FileInfo.SetIsDirectory(enc.GetIsDirectory())
	fi.FileInfo.SetSize(enc.GetSize())
	fi.FileInfo.SetModifiedTime(enc.GetModifiedTime())

	return nil
}
