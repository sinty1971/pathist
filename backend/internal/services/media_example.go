package services

// このファイルは将来のMediaDataService実装時の参考例です

/*
import (
	"penguin-backend/internal/models"
	"path/filepath"
)

// MediaDataService はメディアファイル管理サービスを提供する
type MediaDataService struct {
	FileService  *FileService
	CacheService *CacheService  // サムネイルキャッシュ等
	MediaPath    string
}

// NewMediaDataService はMediaDataServiceを初期化する
func NewMediaDataService(mediaPath string) (*MediaDataService, error) {
	fileService, err := NewFileService(mediaPath)
	if err != nil {
		return nil, err
	}

	return &MediaDataService{
		FileService:  fileService,
		CacheService: NewCacheService(),
		MediaPath:    mediaPath,
	}, nil
}

// ListMediaFiles はメディアファイル一覧を取得する
func (s *MediaDataService) ListMediaFiles(path string) ([]models.MediaFile, error) {
	// 実装例
	// 1. ファイル一覧を取得
	// 2. メディアファイルのみフィルタリング（画像、動画、音声）
	// 3. メタデータを付与
	return []models.MediaFile{}, nil
}

// UploadMediaFile はメディアファイルをアップロードする
func (s *MediaDataService) UploadMediaFile(file []byte, filename string) (*models.MediaFile, error) {
	// 実装例
	// 1. ファイルタイプの検証
	// 2. ファイルサイズの検証
	// 3. ファイルを保存
	// 4. メタデータを抽出
	// 5. データベースに記録
	return &models.MediaFile{}, nil
}

// GetMediaMetadata はメディアファイルのメタデータを取得する
func (s *MediaDataService) GetMediaMetadata(id string) (*models.MediaMetadata, error) {
	// 実装例
	// 1. ファイルを特定
	// 2. EXIF情報等を抽出
	// 3. ファイルサイズ、解像度等を取得
	return &models.MediaMetadata{}, nil
}

// GenerateThumbnail はサムネイルを生成する
func (s *MediaDataService) GenerateThumbnail(id string, size int) (*models.Thumbnail, error) {
	// 実装例
	// 1. キャッシュを確認
	// 2. なければ生成
	// 3. キャッシュに保存
	// 4. サムネイルパスを返す
	return &models.Thumbnail{}, nil
}

// GetSupportedFormats はサポートされているメディア形式を返す
func (s *MediaDataService) GetSupportedFormats() []string {
	return []string{
		// 画像
		".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp",
		// 動画
		".mp4", ".avi", ".mov", ".wmv", ".flv",
		// 音声
		".mp3", ".wav", ".flac", ".aac", ".ogg",
	}
}
*/