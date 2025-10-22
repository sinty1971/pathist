package models

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	grpcv1 "penguin-backend/gen/penguin/v1"
	"penguin-backend/internal/utils"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// grpcv1.FileInfo 関連のヘルパー
// このファイルは proto ファイルから自動生成された gRPC メッセージ型を
// 補完するためのヘルパー関数や型を提供します。

// NewFileInfo はファイルのフルパスから FileInfo を作成します。
func NewFileInfo(targetPath string) (*grpcv1.FileInfo, error) {
	var err error

	// 絶対パスの正規化
	targetPath, err = utils.CleanAbsPath(targetPath)
	if err != nil {
		return nil, err
	}

	// ファイル情報の取得
	osFi, err := os.Stat(targetPath)
	if err != nil {
		return nil, err
	}

	osModTime := osFi.ModTime()
	if osModTime.IsZero() {
		return nil, errors.New("file modification time is zero")
	}

	return grpcv1.FileInfo_builder{
		TargetPath:   targetPath,
		IdealPath:    "",
		IsDirectory:  osFi.IsDir(),
		Size:         osFi.Size(),
		ModifiedTime: timestamppb.New(osModTime),
	}.Build(), nil
}

// FileInfoEx は gRPC の FileInfo メッセージの拡張機能版です。
type FileInfoEx struct {
	*grpcv1.FileInfo
}

// NewFileInfoEx は*grpcv1.FileInfo から FileInfoWithEx を作成します。
func NewFileInfoEx(fileInfo *grpcv1.FileInfo) (*FileInfoEx, error) {
	return &FileInfoEx{
		FileInfo: proto.Clone(fileInfo).(*grpcv1.FileInfo),
	}, nil
}

type FileInfoEncoding struct {
	TargetPath   string `json:"targetPath" yaml:"target_path"`
	IdealPath    string `json:"idealPath" yaml:"ideal_path"`
	IsDirectory  bool   `json:"isDirectory" yaml:"-"`
	Size         int64  `json:"size" yaml:"-"`
	ModifiedTime string `json:"modifiedTime" yaml:"-"`
}

func NewFileInfoEncoding(fi *grpcv1.FileInfo) FileInfoEncoding {
	return FileInfoEncoding{
		TargetPath:  fi.GetTargetPath(),
		IdealPath:   fi.GetIdealPath(),
		IsDirectory: fi.GetIsDirectory(),
		Size:        fi.GetSize(),
	}
}

// ensureMessage は FileInfo メッセージが nil の場合に初期化を行います。
func (fi *FileInfoEx) ensureMessage() {
	if fi.FileInfo == nil {
		fi.FileInfo = &grpcv1.FileInfo{}
	}
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (fi FileInfoEx) MarshalYAML() (any, error) {
	if fi.FileInfo == nil {
		return FileInfoEncoding{}, nil
	}
	enc := NewFileInfoEncoding(fi.FileInfo)
	if ts := fi.GetModifiedTime(); ts != nil {
		enc.ModifiedTime = ts.AsTime().Format(time.RFC3339Nano)
	}
	return enc, nil
}

// UnmarshalYAML は YAML からの復元を行います。
func (fi *FileInfoEx) UnmarshalYAML(unmarshal func(any) error) error {
	var enc FileInfoEncoding
	if err := unmarshal(&enc); err != nil {
		return err
	}
	return fi.applyEncoding(enc)
}

// MarshalJSON は JSON Serde を従来形式で行います。
func (fi FileInfoEx) MarshalJSON() ([]byte, error) {
	if fi.FileInfo == nil {
		return json.Marshal(FileInfoEncoding{})
	}
	enc := NewFileInfoEncoding(fi.FileInfo)
	if ts := fi.GetModifiedTime(); ts != nil {
		enc.ModifiedTime = ts.AsTime().Format(time.RFC3339Nano)
	}
	return json.Marshal(enc)
}

// UnmarshalJSON は JSON からの復元を行います。
func (fi *FileInfoEx) UnmarshalJSON(data []byte) error {
	var enc FileInfoEncoding
	if err := json.Unmarshal(data, &enc); err != nil {
		return err
	}
	return fi.applyEncoding(enc)
}

func (fi *FileInfoEx) applyEncoding(enc FileInfoEncoding) error {
	fi.ensureMessage()
	fi.FileInfo.SetTargetPath(enc.TargetPath)
	fi.FileInfo.SetIdealPath(enc.IdealPath)
	fi.FileInfo.SetIsDirectory(enc.IsDirectory)
	fi.FileInfo.SetSize(enc.Size)

	if enc.ModifiedTime == "" {
		fi.FileInfo.SetModifiedTime(nil)
		return nil
	}

	parsed, err := ParseTime(enc.ModifiedTime, nil)
	if err != nil {
		return err
	}
	fi.FileInfo.SetModifiedTime(timestamppb.New(parsed))
	return nil
}
