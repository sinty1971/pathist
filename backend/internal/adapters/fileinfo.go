package adapters

import (
	"encoding/json"
	"os"
	"time"

	penguinv1 "penguin-backend/gen/penguin/v1"
	"penguin-backend/internal/utils"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewFileInfo はファイルのフルパスから FileInfo を作成します。
func NewFileInfo(targetPath string) (*penguinv1.FileInfo, error) {
	var err error

	targetPath, err = utils.CleanAbsPath(targetPath)
	if err != nil {
		return nil, err
	}

	osFi, err := os.Lstat(targetPath)
	if err != nil {
		return nil, err
	}

	builder := penguinv1.FileInfo_builder{
		TargetPath:  targetPath,
		IdealPath:   "",
		IsDirectory: osFi.IsDir(),
		Size:        osFi.Size(),
	}
	if !osFi.ModTime().IsZero() {
		builder.ModifiedTime = timestamppb.New(osFi.ModTime())
	}

	return builder.Build(), nil
}

// FileInfoWithYAML は gRPC の FileInfo メッセージをラップし、
// YAML/JSON でのシリアライズに対応するための型です。
type FileInfoWithYAML struct {
	*penguinv1.FileInfo
}

// NewFileInfoYAML はファイルのフルパスから FileInfoWithYAML を作成します。
// 引数が 1 つの場合は対象パス、2 つの場合は対象パスと標準パスを指定します。
func NewFileInfoYAML(targetPath string) (*FileInfoWithYAML, error) {

	if fi, err := NewFileInfo(targetPath); err != nil {
		return nil, err
	} else {
		return &FileInfoWithYAML{
			FileInfo: proto.Clone(fi).(*penguinv1.FileInfo),
		}, nil
	}
}

type fileInfoEncoding struct {
	TargetPath   string `json:"targetPath" yaml:"target_path"`
	IdealPath    string `json:"idealPath" yaml:"ideal_path"`
	IsDirectory  bool   `json:"isDirectory" yaml:"-"`
	Size         int64  `json:"size" yaml:"-"`
	ModifiedTime string `json:"modifiedTime" yaml:"-"`
}

func (fi *FileInfoWithYAML) ensureMessage() {
	if fi.FileInfo == nil {
		fi.FileInfo = &penguinv1.FileInfo{}
	}
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (fi FileInfoWithYAML) MarshalYAML() (any, error) {
	if fi.FileInfo == nil {
		return fileInfoEncoding{}, nil
	}
	enc := fileInfoEncoding{
		TargetPath:  fi.GetTargetPath(),
		IdealPath:   fi.GetIdealPath(),
		IsDirectory: fi.GetIsDirectory(),
		Size:        fi.GetSize(),
	}
	if ts := fi.GetModifiedTime(); ts != nil {
		enc.ModifiedTime = ts.AsTime().Format(time.RFC3339Nano)
	}
	return enc, nil
}

// UnmarshalYAML は YAML からの復元を行います。
func (fi *FileInfoWithYAML) UnmarshalYAML(unmarshal func(any) error) error {
	var enc fileInfoEncoding
	if err := unmarshal(&enc); err != nil {
		return err
	}
	return fi.applyEncoding(enc)
}

// MarshalJSON は JSON Serde を従来形式で行います。
func (fi FileInfoWithYAML) MarshalJSON() ([]byte, error) {
	if fi.FileInfo == nil {
		return json.Marshal(fileInfoEncoding{})
	}
	enc := fileInfoEncoding{
		TargetPath:  fi.GetTargetPath(),
		IdealPath:   fi.GetIdealPath(),
		IsDirectory: fi.GetIsDirectory(),
		Size:        fi.GetSize(),
	}
	if ts := fi.GetModifiedTime(); ts != nil {
		enc.ModifiedTime = ts.AsTime().Format(time.RFC3339Nano)
	}
	return json.Marshal(enc)
}

// UnmarshalJSON は JSON からの復元を行います。
func (fi *FileInfoWithYAML) UnmarshalJSON(data []byte) error {
	var enc fileInfoEncoding
	if err := json.Unmarshal(data, &enc); err != nil {
		return err
	}
	return fi.applyEncoding(enc)
}

func (fi *FileInfoWithYAML) applyEncoding(enc fileInfoEncoding) error {
	fi.ensureMessage()
	fi.FileInfo.SetTargetPath(enc.TargetPath)
	fi.FileInfo.SetIdealPath(enc.IdealPath)
	fi.FileInfo.SetIsDirectory(enc.IsDirectory)
	fi.FileInfo.SetSize(enc.Size)

	if enc.ModifiedTime == "" {
		fi.FileInfo.SetModifiedTime(nil)
		return nil
	}

	parsed, err := utils.ParseTime(enc.ModifiedTime, nil)
	if err != nil {
		return err
	}
	fi.FileInfo.SetModifiedTime(timestamppb.New(parsed))
	return nil
}
