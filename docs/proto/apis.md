# Penguin API ドキュメント

Protocol Buffers Editions 2023を使用したgRPC/Connect API仕様書

## 概要

このドキュメントは、豊田築炉（Toyoda Chikuro）プロジェクトのAPI仕様を定義します。

- **Package**: `grpc.v1`
- **Edition**: Protocol Buffers Editions 2023
- **Go Package**: `server-grpc/gen/grpc/v1;grpcv1`

## データ型

### FileInfo

ファイルまたはディレクトリの情報を表すメッセージ

| フィールド | 型 | 説明 |
|---|---|---|
| `id` | `string` | ファイルID |
| `path` | `string` | ファイルパス |
| `is_directory` | `bool` | ディレクトリかどうか |
| `size` | `int64` | ファイルサイズ（バイト） |
| `modified_time` | `google.protobuf.Timestamp` | 最終更新時刻 |

### Company

会社エンティティと内部情報を表すメッセージ

| フィールド | 型 | 説明 |
|---|---|---|
| `id` | `string` | 会社ID |
| `managed_folder` | `string` | 管理フォルダーパス |
| `short_name` | `string` | 会社略称 |
| `category_index` | `int32` | カテゴリーインデックス |
| `inside_ideal_path` | `string` | 理想的な内部パス |
| `inside_legal_name` | `string` | 正式な会社名 |
| `inside_postal_code` | `string` | 郵便番号 |
| `inside_address` | `string` | 住所 |
| `inside_phone` | `string` | 電話番号 |
| `inside_email` | `string` | メールアドレス |
| `inside_website` | `string` | ウェブサイトURL |
| `inside_tags` | `repeated string` | タグ一覧 |
| `inside_required_files` | `repeated FileInfo` | 必要ファイル一覧 |

### CompanyCategory

会社カテゴリーを表すメッセージ

| フィールド | 型 | 説明 |
|---|---|---|
| `index` | `int32` | カテゴリーインデックス |
| `label` | `string` | カテゴリーラベル |

### Koji

工事プロジェクトと内部情報を表すメッセージ

| フィールド | 型 | 説明 |
|---|---|---|
| `id` | `string` | 工事ID |
| `status` | `string` | 工事ステータス |
| `managed_folder` | `string` | 管理フォルダーパス |
| `start` | `google.protobuf.Timestamp` | 開始日時 |
| `company_name` | `string` | 会社名 |
| `location_name` | `string` | 場所名 |
| `inside_ideal_path` | `string` | 理想的な内部パス |
| `inside_end` | `google.protobuf.Timestamp` | 終了日時 |
| `inside_description` | `string` | 説明 |
| `inside_tags` | `repeated string` | タグ一覧 |
| `inside_required_files` | `repeated FileInfo` | 必要ファイル一覧 |

## サービス

### FileService

ファイル管理操作を提供するサービス

#### GetFileInfos

指定されたパスのファイル情報一覧を取得

**リクエスト**: `GetFileInfosRequest`
- `path` (string): 取得するパス

**レスポンス**: `GetFileInfosResponse`
- `file_infos` (repeated FileInfo): ファイル情報一覧

#### GetFileBasePath

ファイルベースパスを取得

**リクエスト**: `GetFileBasePathRequest` (空)

**レスポンス**: `GetFileBasePathResponse`
- `base_path` (string): ベースパス

### CompanyService

会社管理操作を提供するサービス

#### GetCompanyMapById

IDをキーとした会社マップを取得

**リクエスト**: `GetCompanyMapByIdRequest` (空)

**レスポンス**: `GetCompanyMapByIdResponse`
- `company_map_by_id` (map<string, Company>): 会社IDマップ

#### GetCompanyById

指定されたIDの会社情報を取得

**リクエスト**: `GetCompanyByIdRequest`
- `id` (string): 会社ID

**レスポンス**: `GetCompanyByIdResponse`
- `company` (Company): 会社情報

#### UpdateCompany

会社情報を更新

**リクエスト**: `UpdateCompanyRequest`
- `current_company_id` (string): 現在の会社ID
- `updated_company` (Company): 更新後の会社情報

**レスポンス**: `UpdateCompanyResponse`
- `company_map_by_id` (map<string, Company>): 更新後の会社IDマップ

#### GetCompanyCategories

会社カテゴリー一覧を取得

**リクエスト**: `GetCompanyCategoriesRequest` (空)

**レスポンス**: `GetCompanyCategoriesResponse`
- `categories` (repeated CompanyCategory): カテゴリー一覧

### KojiService

工事プロジェクト管理操作を提供するサービス

#### GetKojiMapById

IDをキーとした工事マップを取得

**リクエスト**: `GetKojiMapByIdRequest` (空)

**レスポンス**: `GetKojiMapByIdResponse`
- `koji_map_by_id` (map<string, Koji>): 工事IDマップ

#### GetKojiById

指定されたIDの工事情報を取得

**リクエスト**: `GetKojiByIdRequest`
- `id` (string): 工事ID

**レスポンス**: `GetKojiByIdResponse`
- `koji` (Koji): 工事情報

#### UpdateKoji

工事情報を更新

**リクエスト**: `UpdateKojiRequest`
- `current_koji_id` (string): 現在の工事ID
- `updated_koji` (Koji): 更新後の工事情報

**レスポンス**: `UpdateKojiResponse`
- `koji_map_by_id` (map<string, Koji>): 更新後の工事IDマップ

## Connect RPC エンドポイント

各サービスは以下のHTTPエンドポイントでアクセス可能：

### FileService
- `POST /grpc.v1.FileService/GetFileInfos`
- `POST /grpc.v1.FileService/GetFileBasePath`

### CompanyService
- `POST /grpc.v1.CompanyService/GetCompanyMapById`
- `POST /grpc.v1.CompanyService/GetCompanyById`
- `POST /grpc.v1.CompanyService/UpdateCompany`
- `POST /grpc.v1.CompanyService/GetCompanyCategories`

### KojiService
- `POST /grpc.v1.KojiService/GetKojiMapById`
- `POST /grpc.v1.KojiService/GetKojiById`
- `POST /grpc.v1.KojiService/UpdateKoji`

## 使用例

### cURL例

```bash
# ファイル情報取得
curl -X POST http://localhost:9090/grpc.v1.FileService/GetFileInfos \
  -H "Content-Type: application/json" \
  -d '{"path": "/example/path"}'

# 会社情報取得
curl -X POST http://localhost:9090/grpc.v1.CompanyService/GetCompanyById \
  -H "Content-Type: application/json" \
  -d '{"id": "company123"}'
```

### Connect-ES例 (TypeScript)

```typescript
import { createClient } from "@connectrpc/connect";
import { FileService } from "./gen/grpc/v1/toyotachikuro_connect";

const client = createClient(FileService, transport);

const response = await client.getFileInfos({
  path: "/example/path"
});
```

---

*このドキュメントは Protocol Buffers Editions 2023 に基づいて手動で作成されました。*
*最終更新: 2025年11月9日*