# 🔄 ID同期システム 実装ガイド

バックエンドとフロントエンドの**自動・半自動ID同期**機能の実装と使用方法について説明します。

## 📋 概要

このシステムは以下の課題を解決します：

- ✅ **ID不整合問題**: 工事データ変更時のID同期ずれ
- ✅ **フルパス vs 短縮ID**: メモリ効率とパフォーマンスの最適化
- ✅ **バックエンド/フロントエンド同期**: 一貫したID生成アルゴリズム
- ✅ **自動・半自動同期**: ユーザー操作負荷の軽減

## 🏗️ システム構成

### フロントエンド実装

#### 1. 基本ID同期ユーティリティ
```typescript
// app/utils/idSync.ts
import { IdSyncManager, useIdSync } from './idSync';

const { 
  generateKojiId,     // 工事ID生成
  generatePathId,     // パスID生成（Len7）
  validateId,         // ID検証
  autoSync,           // 自動同期
  semiAutoSync        // 半自動同期（確認あり）
} = useIdSync();
```

#### 2. 自動同期React Hook
```typescript
// app/hooks/useAutoIdSync.ts
import { useAutoIdSync } from '../hooks/useAutoIdSync';

const syncResult = useAutoIdSync(
  kojiData.id,                    // 現在のID
  {                               // ID生成要素
    startDate: kojiData.startDate,
    companyName: kojiData.companyName,
    locationName: kojiData.locationName
  },
  {
    requireConfirmation: false,    // 自動同期 (true=半自動)
    onSuccess: (newId) => {        // 成功時コールバック
      console.log('ID更新:', newId);
    },
    onError: (error) => {          // エラー時コールバック
      console.error('同期エラー:', error);
    }
  }
);

// syncResult.currentId      - 現在のID
// syncResult.isSyncing      - 同期中かどうか
// syncResult.sync()         - 手動同期実行
```

#### 3. フルパス→Len7変換Hook
```typescript
// フルパスIDをLen7に一括変換
const { 
  convertToLen7,
  convertFromLen7,
  mappingTable 
} = usePathIdSync(fullPathIds, {
  autoConvert: true,
  onMappingReady: (mapping) => {
    console.log('変換テーブル準備完了:', mapping.size);
  }
});
```

### バックエンド実装

#### 1. ID同期API エンドポイント
```
POST /api/id-sync/generate-koji    # 工事ID生成
POST /api/id-sync/generate-path    # パスID生成
POST /api/id-sync/validate         # ID検証
POST /api/id-sync/bulk-convert     # 一括変換
GET  /api/id-sync/mapping          # IDマッピング取得
GET  /api/id-sync/config           # 同期設定取得
```

#### 2. ID生成アルゴリズム統一
```go
// backend/internal/models/id.go - 既存のID生成機能を使用
func generateKojiID(startDate time.Time, companyName, locationName string) (string, error) {
    dateStr := startDate.Format("20060102")
    combined := dateStr + companyName + locationName
    id := models.NewIDFromString(combined).Len5()
    return id, nil
}

// フルパス → Len7 ID変換
func generatePathID(fullPath string) string {
    return models.NewIDFromString(fullPath).Len7()
}
```

## 🔧 使用方法

### 1. 基本的な工事ID生成・検証

```typescript
// 工事データから自動ID生成
const kojiId = generateKojiId({
  startDate: new Date(2025, 5, 18),
  companyName: '豊田築炉',
  locationName: '名和工場'
});
// → 例: "3XK9P" (5文字)

// バックエンドとの一致確認
const validation = await validateId(kojiId, components);
if (!validation.isValid) {
  console.log('ID不整合:', validation.expectedId);
}
```

### 2. 自動同期の組み込み

```typescript
// 工事詳細コンポーネントでの使用例
function KojiDetailModal({ koji, onUpdate }) {
  const syncResult = useAutoIdSync(
    koji.id,
    {
      startDate: koji.startDate,
      companyName: koji.companyName,
      locationName: koji.locationName
    },
    {
      requireConfirmation: false,
      onSuccess: (newId) => {
        // 自動更新が完了した場合の処理
        onUpdate({ ...koji, id: newId });
      }
    }
  );

  // 会社名や場所名が変更されると自動的に同期チェック
  // 不整合があれば自動修正される
}
```

### 3. パスID変換（メモリ最適化）

```typescript
// TreeViewでフルパスをLen7に変換
const pathMapping = usePathIdSync(allFilePaths);

// TreeNodeでIDをLen7に変換
const node = {
  id: pathMapping.convertToLen7(fileInfo.path),  // 70-85%メモリ削減
  name: fileInfo.name,
  path: fileInfo.path  // 表示用は元のパスを保持
};

// 逆変換（必要時）
const originalPath = pathMapping.convertFromLen7(node.id);
```

### 4. 一括変換・検証

```typescript
// 複数の工事データを一括処理
const { syncAll, globalSyncing } = useBulkIdSync(
  kojies,
  (koji) => ({
    startDate: koji.startDate,
    companyName: koji.companyName,
    locationName: koji.locationName
  }),
  {
    requireConfirmation: true,  // 一括更新前に確認
    onSuccess: (newId) => {
      // 各項目の更新処理
    }
  }
);

// 一括同期実行
await syncAll();
```

## 📊 パフォーマンス効果

### メモリ使用量削減
```
フルパス ID: "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場" (28文字)
Len7 ID:     "3XK9P7M"                                (7文字)
削減率:      75% ↓
```

### 検索速度向上
- 文字列比較: **4-8倍高速**
- ハッシュマップ操作: **大幅改善**（大規模データ時）

### 衝突リスク
- Len7 (32^7): **約343億通り** → 10万件以下では衝突リスクほぼゼロ
- Len5 (32^5): **約3,355万通り** → 工事IDに最適

## 🛠️ 実装手順

### フェーズ1: 準備段階
1. ✅ フロントエンド同期ユーティリティ実装
2. ✅ バックエンドID同期API実装
3. ✅ デモコンポーネント作成

### フェーズ2: 段階的適用
1. **工事詳細モーダル**で自動同期を有効化
2. **工事一覧**でID不整合検出機能を追加
3. **ガントチャート**で同期機能を組み込み

### フェーズ3: パスID最適化（オプション）
1. **Files.tsx**でLen7 ID変換を有効化
2. **TreeView**のkey管理を最適化
3. **メモリ使用量**の測定・検証

### フェーズ4: 本格運用
1. **全コンポーネント**で同期機能を有効化
2. **パフォーマンステスト**実行
3. **エラーハンドリング**強化

## 🔍 デバッグとテスト

### デモコンポーネント使用方法
```typescript
// app/components/IdSyncDemo.tsx をページに追加
import { IdSyncDemo } from '../components/IdSyncDemo';

function DebugPage() {
  return <IdSyncDemo />;
}
```

### 手動テストコマンド
```bash
# フロントエンド開発サーバー
just frontend

# バックエンドAPI
just backend

# gRPC サービス仕様
less proto/penguin/v1/penguin.proto
```

### APIテスト例
```bash
# 工事一覧の取得 (fileclient CLI)
go run backend/cmd/fileclient/main.go -base-url http://localhost:9090

# grpcurl を使った呼び出し例 (KojiService.ListKojies)
grpcurl -plaintext \
  -import-path proto \
  -proto penguin/v1/penguin.proto \
  localhost:9090 penguin.v1.KojiService/ListKojies
```

## ⚠️ 注意事項

### 移行時の考慮点
1. **既存データとの互換性**: 段階的な移行を推奨
2. **バックアップ**: 重要データの事前バックアップ
3. **ロールバック計画**: 問題発生時の復旧手順

### パフォーマンス監視
1. **メモリ使用量**: 変換前後の比較測定
2. **API応答時間**: ID生成・検証の速度測定  
3. **衝突検出**: 実運用でのID重複監視

### セキュリティ
1. **ID推測困難性**: BLAKE2bハッシュによる強固な生成
2. **入力検証**: 不正なパラメータのチェック
3. **レート制限**: API呼び出しの制限（将来検討）

## 📚 関連ファイル

### フロントエンド
- `app/utils/idSync.ts` - メインの同期ロジック
- `app/hooks/useAutoIdSync.ts` - React Hook実装
- `app/components/IdSyncDemo.tsx` - デモ・テスト用コンポーネント

### バックエンド  
- `internal/handlers/id_sync.go` - ID同期API
- `internal/endpoints/id_sync.go` - ルート定義
- `internal/models/id.go` - 既存のID生成機能

### 設計資料
- `frontend/app/utils/idComparison.ts` - パフォーマンス分析
- `frontend/app/utils/idDemo.ts` - ベンチマーク・衝突テスト

---

このシステムにより、バックエンドとフロントエンドのID管理が統一され、自動・半自動での同期によってユーザーの手間を大幅に削減できます。
