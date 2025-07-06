# React Router v7 ルート階層構造

> 🤖 自動生成日時: 2025/7/6 17:36:15  
> 📄 生成元: `app/routes.ts`

## Mermaid図

```mermaid
graph TD
    A[/"/" - Root] --> B[_layout.tsx]
    B --> C["/" - ホームページ（機能紹介）<br/>routes/_layout._index.tsx]
    B --> D["/files" - ファイル一覧（TreeView）<br/>routes/_layout.files.tsx]
    B --> E["/kojies" - /kojies<br/>routes/_layout.kojies.tsx]
    B --> F["/kojies/gantt" - /kojies/gantt<br/>routes/_layout.gantt.tsx]
    B --> G["/companies" - /companies<br/>routes/_layout.companies.tsx]

    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#e8f5e8
    style D fill:#e8f5e8
    style E fill:#e8f5e8
    style F fill:#e8f5e8
    style G fill:#e8f5e8

```

## ツリー構造（ASCII）

```
/
└── _layout.tsx (共通レイアウト)
    ├── / (ホーム) → routes/_layout._index.tsx
    ├── /files (ファイル一覧) → routes/_layout.files.tsx
    ├── /kojies (/kojies) → routes/_layout.kojies.tsx
    ├── /kojies/gantt (/kojies/gantt) → routes/_layout.gantt.tsx
    └── /companies (/companies) → routes/_layout.companies.tsx

```

## ルート一覧

| パス | ファイル | 説明 | タイプ |
|------|----------|------|------|
| `/` | `routes/_layout._index.tsx` | ホームページ（機能紹介） | index |
| `/files` | `routes/_layout.files.tsx` | ファイル一覧（TreeView） | route |
| `/kojies` | `routes/_layout.kojies.tsx` | /kojies | route |
| `/kojies/gantt` | `routes/_layout.gantt.tsx` | /kojies/gantt | route |
| `/companies` | `routes/_layout.companies.tsx` | /companies | route |


## 実行方法

```bash
# 図を再生成
npm run generate-routes

# または直接実行
node scripts/generate-route-diagram.js
```

## 注意事項

- このファイルは自動生成されます
- `routes.ts` を変更後は `npm run generate-routes` で更新してください
- 手動で編集しないでください（変更は失われます）
