package models

import (
	penguinv1 "penguin-backend/gen/penguin/v1"
	"penguin-backend/internal/serialize"
)

// FileInfo は serialize.FileInfoWithYAML のエイリアスです。
type FileInfo = serialize.FileInfoWithYAML

// NewFileInfo は serialize.NewFileInfo のラッパーです。
func NewFileInfo(paths ...string) (*FileInfo, error) {
	return serialize.NewFileInfo(paths...)
}

// NewFileInfoFromProto は serialize.NewFileInfoFromProto のラッパーです。
func NewFileInfoFromProto(src *penguinv1.FileInfo) *FileInfo {
	return serialize.NewFileInfoFromProto(src)
}
