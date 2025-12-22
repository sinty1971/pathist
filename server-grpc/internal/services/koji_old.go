package services

// import (
// 	"server-grpc/internal/models"
// 	"server-grpc/internal/utils"
// 	"errors"
// 	"fmt"
// 	"os"
// 	"path"
// 	"path/filepath"
// 	"regexp"
// 	"sort"
// 	"strings"
// 	"sync"
// )

// // GetKoji は指定されたパスから工事を取得する
// // folderName は工事フォルダー名 (Ex. '2025-0618 豊田築炉 名和工場')
// // 工事を返す
// func (ks *KojiServiceOld) GetKoji(folderName string) (*models.Koji, error) {
// 	// フォルダー名からフルパスを作成
// 	folderPath := filepath.Join(ks.TargetFolder, folderName)

// 	// 工事データモデルを作成
// 	koji, err := models.NewKoji(folderPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("工事データモデルの作成に失敗しました: %v", err)
// 	}

// 	// データベースファイルから読み込む
// 	err = ks.LoadDatabaseFile(koji)
// 	if err != nil {
// 		return nil, fmt.Errorf("データベースファイルからの読み込みに失敗しました: %v", err)
// 	}

// 	// 必須ファイルの設定
// 	err = ks.UpdateRequiredFiles(koji)
// 	if err != nil {
// 		return nil, fmt.Errorf("必須ファイルの設定に失敗しました: %v", err)
// 	}

// 	return koji, nil
// }

// type getKojiesMode int

// const (
// 	GetKoujiesModeRecent getKojiesMode = 1 << iota
// 	GetKoujiesModeOld
// 	GetKojiesModeSyncDatabase
// )

// func getKojiesModeFunc(modes ...getKojiesMode) getKojiesMode {
// 	opt := getKojiesMode(GetKoujiesModeRecent)
// 	for _, mode := range modes {
// 		opt |= mode
// 	}
// 	return opt
// }

// // LoadDatabaseFile は工事データベースファイルからデータを取り込みます
// func (ks *KojiServiceOld) LoadDatabaseFile(target *models.Koji) error {
// 	// 工事データベースファイルからデータを読み込む
// 	database, err := ks.DatabaseService.Load(target)
// 	if err != nil {
// 		// ファイルが存在しない場合はデフォルト値を設定
// 		target.UpdateStatus()

// 		// 新規ファイルの場合は非同期で保存
// 		go func() {
// 			ks.DatabaseService.Save(target)
// 		}()
// 		return nil
// 	}

// 	// データベースファイルにしか保持されない内容をtargetに反映
// 	target.Description = database.Description
// 	target.EndDate = database.EndDate
// 	target.Tags = database.Tags

// 	// 計算が必要な項目の更新
// 	target.UpdateStatus()

// 	return nil
// }

// // UpdateRequiredFiles は工事の必須ファイルを更新します
// // データベースファイルの情報を優先するので、値が空の場合のみ更新する
// func (ks *KojiServiceOld) UpdateRequiredFiles(koji *models.Koji) error {
// 	// 必須ファイルがすでに設定されている場合は更新しない
// 	if len(koji.RequiredFiles) > 0 {
// 		return errors.New("必須ファイルがすでに設定されています")
// 	}

// 	// 工事フォルダー内のファイルを取得
// 	entries, err := os.ReadDir(koji.TargetFolder)
// 	if err != nil {
// 		// エラーの場合はデフォルトの必須ファイルを設定
// 		requiredFile, err := models.NewFileInfo(path.Join(koji.TargetFolder, "工事.xlsx"))
// 		if err != nil {
// 			return err
// 		}
// 		koji.RequiredFiles = []models.FileInfo{*requiredFile}
// 		return nil
// 	}

// 	// 必須ファイルの設定（現在は1個のみ）
// 	var requiredFile *models.FileInfo
// 	var found bool

// 	// 日付パターンのコンパイル（1回のみ）
// 	datePattern := regexp.MustCompile(`^20\d{2}-\d{4}`)

// 	// 工事ファイルの検索
// 	for _, file := range entries {
// 		filename := file.Name()
// 		// エクセルファイル以外はスキップ
// 		if !utils.FilenameIsExcel(filename) {
// 			continue
// 		}

// 		// ファイル名の解析
// 		extension := filepath.Ext(filename)
// 		filenameWithoutExt := strings.TrimSuffix(filename, extension)

// 		// 条件に合うファイルかチェック
// 		if !datePattern.MatchString(filenameWithoutExt) && filenameWithoutExt != "工事" {
// 			continue
// 		}

// 		actualPath := path.Join(koji.TargetFolder, file.Name())
// 		standardName := koji.GetFolderName() + extension
// 		standardPath := path.Join(koji.TargetFolder, standardName)
// 		requiredFile, err = models.NewFileInfo(actualPath, standardPath)
// 		if err != nil {
// 			continue
// 		}
// 		found = true
// 		break
// 	}

// 	// 工事ファイルが見つからない場合はひな形を設定
// 	if !found {
// 		standardPath := path.Join(koji.TargetFolder, koji.GetFolderName()+".xlsx")
// 		requiredFile, err = models.NewFileInfo("工事.xlsx", standardPath)

// 		// テンプレートファイルのコピーは省略（高速化のため）
// 		// 必要に応じて別途実行する
// 	}

// 	// 管理ファイルの設定
// 	koji.RequiredFiles = []models.FileInfo{*requiredFile}

// 	return nil
// }

// // Update は工事情報 FileInfoの情報を更新
// // kojiの詳細情報で @profile.yaml を更新
// // 移動元フォルダー名：FileInfo.Name
// // 移動先フォルダー名：StartDate, CompanyName, LocationNameから生成されたフォルダー名
// func (ks *KojiServiceOld) Update(target *models.Koji) error {
// 	// 更新前の工事情報を保存
// 	temp := *target

// 	// フォルダー名の変更
// 	if temp.UpdateFolderPath() {
// 		// 工事IDの更新（フォルダー名の変更に伴う）
// 		temp.UpdateID()

// 		// 必須ファイル情報の更新（フォルダー名の変更に伴う）
// 		temp.RequiredFiles = []models.FileInfo{}
// 		if err := ks.UpdateRequiredFiles(&temp); err != nil {
// 			return err
// 		}

// 		// ファイルの移動
// 		// TODO: ファイル移動の実装
// 		// if err := ks.RootService.MoveFile(temp.GetFolderName(), target.GetFolderName()); err != nil {
// 		// 	return err
// 		// }
// 	}

// 	// フォルダー名の変更完了後、情報を更新
// 	*target = temp

// 	// 計算が必要な項目の更新
// 	target.UpdateStatus()

// 	// 更新後の工事情報を属性ファイルに反映
// 	return ks.DatabaseService.Save(target)
// }
