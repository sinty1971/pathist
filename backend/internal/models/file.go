package models

import (
	penguinv1 "penguin-backend/gen/penguin/v1"
	"penguin-backend/internal/adapters"
)

// FileInfo は serialize.FileInfoWithYAML のエイリアスです。
type FileInfo = adapters.FileInfoWithYAML

// NewFileInfo は serialize.NewFileInfo のラッパーです。
func NewFileInfo(paths ...string) (*FileInfo, error) {
	return adapters.NewFileInfo(paths...)
}

// NewFileInfoFromProto は serialize.NewFileInfoFromProto のラッパーです。
func NewFileInfoFromProto(src *penguinv1.FileInfo) *FileInfo {
	return adapters.NewFileInfoFromProto(src)
}
