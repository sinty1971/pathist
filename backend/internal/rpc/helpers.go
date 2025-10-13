package rpc

import (
	penguinv1 "penguin-backend/gen/penguin/v1"
	"penguin-backend/internal/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoTimestamp(ts models.Timestamp) *timestamppb.Timestamp {
	if ts.Time.IsZero() {
		return nil
	}
	return timestamppb.New(ts.Time)
}

func toModelTimestamp(ts *timestamppb.Timestamp) models.Timestamp {
	if ts == nil {
		return models.Timestamp{}
	}
	return models.Timestamp{Time: ts.AsTime()}
}

func convertModelFileInfo(src *models.FileInfo) *penguinv1.FileInfo {
	if src == nil {
		return nil
	}

	return penguinv1.FileInfo_builder{
		TargetPath:   src.TargetPath,
		StandardPath: src.StandardPath,
		IsDirectory:  src.IsDirectory,
		Size:         src.Size,
		ModifiedTime: toProtoTimestamp(src.ModifiedTime),
	}.Build()
}

func convertProtoFileInfo(src *penguinv1.FileInfo) models.FileInfo {
	if src == nil {
		return models.FileInfo{}
	}

	return models.FileInfo{
		TargetPath:   src.GetTargetPath(),
		StandardPath: src.GetStandardPath(),
		IsDirectory:  src.GetIsDirectory(),
		Size:         src.GetSize(),
		ModifiedTime: toModelTimestamp(src.GetModifiedTime()),
	}
}
