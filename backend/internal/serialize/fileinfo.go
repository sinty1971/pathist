package serialize

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	penguinv1 "penguin-backend/gen/penguin/v1"
	"penguin-backend/internal/utils"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FileInfoWithYAML は gRPC の FileInfo メッセージをラップし、
// YAML/JSON でのシリアライズに対応するための型です。
type FileInfoWithYAML struct {
	*penguinv1.FileInfo
}

// NewFileInfo はファイルのフルパスから FileInfoWithYAML を作成します。
// 引数が 1 つの場合は対象パス、2 つの場合は対象パスと標準パスを指定します。
func NewFileInfo(paths ...string) (*FileInfoWithYAML, error) {
	if len(paths) == 0 || len(paths) > 2 {
		return nil, errors.New("引数が１つまたは２つ必要です")
	}

	targetPath, err := utils.CleanAbsPath(paths[0])
	if err != nil {
		return nil, err
	}

	osInfo, err := os.Stat(targetPath)
	if err != nil {
		return nil, err
	}

	var standardPath string
	if len(paths) == 2 {
		standardPath, err = utils.CleanAbsPath(paths[1])
		if err != nil {
			return nil, err
		}
	}

	builder := penguinv1.FileInfo_builder{
		TargetPath:  targetPath,
		IdealPath:   standardPath,
		IsDirectory: osInfo.IsDir(),
		Size:        osInfo.Size(),
	}
	if !osInfo.ModTime().IsZero() {
		builder.ModifiedTime = timestamppb.New(osInfo.ModTime())
	}

	return &FileInfoWithYAML{FileInfo: builder.Build()}, nil
}

// NewFileInfoFromProto は proto メッセージからラッパーを生成します（クローンを作成）。
func NewFileInfoFromProto(src *penguinv1.FileInfo) *FileInfoWithYAML {
	if src == nil {
		return nil
	}
	return &FileInfoWithYAML{
		FileInfo: proto.Clone(src).(*penguinv1.FileInfo),
	}
}

// CloneProto は内部の FileInfo をディープコピーして返します。
func (fi *FileInfoWithYAML) CloneProto() *penguinv1.FileInfo {
	if fi == nil || fi.FileInfo == nil {
		return nil
	}
	return proto.Clone(fi.FileInfo).(*penguinv1.FileInfo)
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
