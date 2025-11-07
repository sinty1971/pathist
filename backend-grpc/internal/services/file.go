package services

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	grpc "backend-grpc/gen/grpc/v1"
	grpcConnect "backend-grpc/gen/grpc/v1/grpcv1connect"
	"backend-grpc/internal/models"

	"connectrpc.com/connect"
)

// FileService の実装

// FileService exposes FileService operations via Connect handlers.
type FileService struct {
	// Embed the unimplemented handler for forward compatibility
	grpcConnect.UnimplementedFileServiceHandler

	// services は任意のgrpcサービスハンドラーへの参照
	services *Services

	// BasePath はファイルサービスの絶対パスフォルダー
	BasePath string `json:"basePath" yaml:"base_path" example:"/penguin/豊田築炉"`
}

func NewFileService(services *Services, options *ServiceOptions) *FileService {
	return &FileService{
		services: services,
		BasePath: options.FileServiceFolder,
	}
}

func (s *FileService) Cleanup() {
	// 現在はクリーンアップ処理は不要
}

// ListFileInfos は指定されたパスのファイル情報一覧を返す
func (s *FileService) ListFileInfos(
	ctx context.Context,
	req *connect.Request[grpc.GetFileInfosRequest]) (
	res *connect.Response[grpc.GetFileInfosResponse],
	err error) {
	// コンテキストを無視
	_ = ctx

	// 変数定義
	var (
		dirs []os.DirEntry
		fis  []*grpc.FileInfo = []*grpc.FileInfo{}
	)

	// リクエスト情報の取得
	relPath := req.Msg.GetPath()

	// 絶対パスを取得
	absPath, err := s.GetAbsPathFrom(relPath)
	if err != nil {
		return // naked return: res=nil, err=err
	}

	// ファイルエントリ配列を取得
	dirs, err = os.ReadDir(absPath)
	if err != nil {
		return // naked return: res=nil, err=err
	}
	// ファイルエントリが0の場合は空配列を返す

	if len(dirs) == 0 {
		res = connect.NewResponse(&grpc.GetFileInfosResponse{})
		res.Msg.SetFileInfos(fis)
		return // naked return: res=res, err=nil
	}

	// チャンネルとワーカーグループを設定
	jobsChan := make(chan int, len(dirs))
	fisChan := make(chan *grpc.FileInfo, len(dirs))
	// 並列処理用のワーカー数を決定（CPU数の2倍、最大16）
	numWorkers := min(min(runtime.NumCPU()*2, 16), len(dirs))
	// ワーカーグループを設定
	var wg sync.WaitGroup
	// ワーカーを起動
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobsChan {
				dir := dirs[index]
				fullpath := filepath.Join(absPath, dir.Name())

				fi, _ := models.NewFileInfo(fullpath)
				if fi != nil {
					fisChan <- fi.FileInfo
				}
			}
		}()
	}

	// ジョブを送信
	for i := range dirs {
		jobsChan <- i
	}
	close(jobsChan)

	// ワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(fisChan)
	}()

	// 結果を収集
	fis = make([]*grpc.FileInfo, 0, len(dirs))
	for fi := range fisChan {
		fis = append(fis, fi)
	}

	// レスポンスを更新して返す
	res = connect.NewResponse(&grpc.GetFileInfosResponse{})
	res.Msg.SetFileInfos(fis)
	return // naked return: res=res, err=nil
}

// GetAbsPathFrom BasePathに引数の相対パスを追加した絶対パスを返す
func (s *FileService) GetAbsPathFrom(relPath string) (res string, err error) {

	// 絶対パスがある場合はエラーを返す
	if strings.HasPrefix(relPath, "~/") || filepath.IsAbs(relPath) {
		return "", errors.New("絶対パスは使用できません")
	}

	res = filepath.Join(s.BasePath, relPath)

	return // naked return
}

// CopyFile はファイルまたはディレクトリをコピーする
func (s *FileService) CopyFile(relSrc, relDst string) (err error) {
	var absSrc, absDst string

	// relSrcがパスチェック及び絶対パス変換
	absSrc, err = s.GetAbsPathFrom(relSrc)
	if err != nil {
		return
	}

	// relDstのパスチェック及び絶対パス変換
	absDst, err = s.GetAbsPathFrom(relDst)
	if err != nil {
		return
	}

	// コピー元の存在確認
	srcOsFi, err := os.Stat(absSrc)
	if err != nil {
		return
	}

	// ディレクトリの場合
	if srcOsFi.IsDir() {
		err = s.absCopyDir(absSrc, absDst)
	} else {
		// ファイルの場合
		err = s.absCopyFile(absSrc, absDst)
	}

	return
}

// absCopyFile はファイルをコピーする内部関数
func (s *FileService) absCopyFile(absSrc, absDst string) (err error) {
	// コピー元ファイルを開く
	srcFile, err := os.Open(absSrc)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// コピー先のディレクトリが存在しない場合は作成
	dstDir := filepath.Dir(absDst)
	if err = os.MkdirAll(dstDir, 0755); err != nil {
		return
	}

	// コピー先ファイルを作成
	dstFile, err := os.Create(absDst)
	if err != nil {
		return
	}
	defer dstFile.Close()

	// ファイル内容をコピー
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return
	}

	// ファイル権限をコピー
	if fi, err := os.Stat(absSrc); err != nil {
		return err
	} else {
		return os.Chmod(absDst, fi.Mode())
	}
}

// absCopyDir はディレクトリを再帰的にコピーする内部関数
func (s *FileService) absCopyDir(absSrc, absDst string) error {
	// コピー元ディレクトリの情報を取得
	srcInfo, err := os.Stat(absSrc)
	if err != nil {
		return err
	}

	// コピー先ディレクトリを作成
	if err := os.MkdirAll(absDst, srcInfo.Mode()); err != nil {
		return err
	}

	// ディレクトリ内のエントリを読み取り
	entries, err := os.ReadDir(absSrc)
	if err != nil {
		return err
	}

	// 各エントリを処理
	for _, entry := range entries {
		srcPath := filepath.Join(absSrc, entry.Name())
		dstPath := filepath.Join(absDst, entry.Name())

		if entry.IsDir() {
			// サブディレクトリの場合、再帰的にコピー
			if err := s.absCopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// ファイルの場合、ファイルをコピー
			if err := s.absCopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveFile はファイルを移動する
func (s *FileService) MoveFile(relSrc, relDst string) error {
	absSrc, err := s.GetAbsPathFrom(relSrc)
	if err != nil {
		return err
	}
	absDst, err := s.GetAbsPathFrom(relDst)
	if err != nil {
		return err
	}

	// 移動先のディレクトリが存在するかチェック
	if _, err := os.Stat(absSrc); os.IsNotExist(err) {
		return errors.New("移動元のファイル/ディレクトリが存在しません: " + relSrc)
	}

	// 移動先の親ディレクトリを作成（必要に応じて）
	dstParent := filepath.Dir(absDst)
	if err := os.MkdirAll(dstParent, 0755); err != nil {
		return err
	}

	return os.Rename(absSrc, absDst)
}

// DeleteFile はファイルを削除する
func (s *FileService) DeleteFile(relPath string) error {
	absPath, err := s.GetAbsPathFrom(relPath)
	if err != nil {
		return err
	}

	return os.Remove(absPath)
}
