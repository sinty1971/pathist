#!/usr/bin/env node

/**
 * React Router v7 ルート構造図生成スクリプト
 * routes.ts から自動的にMermaid図とASCII図を生成
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const routesPath = path.join(__dirname, '../app/routes.ts');
const outputPath = path.join(__dirname, '../route-structure-generated.md');

// ルートファイルを読み込み
const routesContent = fs.readFileSync(routesPath, 'utf-8');

// ルート情報を抽出（簡易パーサー）
function parseRoutes(content) {
  const routes = [];
  const lines = content.split('\n');
  
  let inLayoutBlock = false;
  let layoutFile = '';
  
  for (const line of lines) {
    const trimmed = line.trim();
    
    // layout行を検出
    if (trimmed.startsWith('layout(')) {
      inLayoutBlock = true;
      const match = trimmed.match(/layout\("([^"]+)"/);
      if (match) {
        layoutFile = match[1];
      }
      continue;
    }
    
    // layout block内のルート
    if (inLayoutBlock) {
      // index route
      if (trimmed.startsWith('index(')) {
        const match = trimmed.match(/index\("([^"]+)"/);
        if (match) {
          routes.push({
            path: '/',
            file: match[1],
            type: 'index',
            layout: layoutFile
          });
        }
      }
      
      // 通常のroute
      if (trimmed.startsWith('route(')) {
        const match = trimmed.match(/route\("([^"]+)",\s*"([^"]+)"/);
        if (match) {
          routes.push({
            path: `/${match[1]}`,
            file: match[2],
            type: 'route',
            layout: layoutFile
          });
        }
      }
      
      // layout blockの終了
      if (trimmed.includes('])')) {
        inLayoutBlock = false;
      }
    }
  }
  
  return routes;
}

// ルート情報に日本語説明を追加
function addRouteDescriptions(routes) {
  const descriptions = {
    '/': 'ホームページ（機能紹介）',
    '/files': 'ファイル一覧（TreeView）',
    '/projects': '工程表（プロジェクト管理）',
    '/gantt': 'ガントチャート（タイムライン）'
  };
  
  return routes.map(route => ({
    ...route,
    description: descriptions[route.path] || route.path
  }));
}

// Mermaid図を生成
function generateMermaidDiagram(routes) {
  let mermaid = 'graph TD\n';
  mermaid += '    A[/"/" - Root] --> B[_layout.tsx]\n';
  
  routes.forEach((route, index) => {
    const nodeId = String.fromCharCode(67 + index); // C, D, E, F...
    const pathDisplay = route.path === '/' ? '/' : route.path;
    const label = `"${pathDisplay}" - ${route.description}`;
    mermaid += `    B --> ${nodeId}[${label}<br/>${route.file}]\n`;
  });
  
  mermaid += '\n    style A fill:#e1f5fe\n';
  mermaid += '    style B fill:#f3e5f5\n';
  
  routes.forEach((_, index) => {
    const nodeId = String.fromCharCode(67 + index);
    mermaid += `    style ${nodeId} fill:#e8f5e8\n`;
  });
  
  return mermaid;
}

// ASCII図を生成
function generateASCIIDiagram(routes) {
  let ascii = '/\n';
  ascii += '└── _layout.tsx (共通レイアウト)\n';
  
  routes.forEach((route, index) => {
    const isLast = index === routes.length - 1;
    const prefix = isLast ? '    └──' : '    ├──';
    const pathDisplay = route.path === '/' ? '/ (ホーム)' : `${route.path} (${route.description.split('（')[0]})`;
    ascii += `${prefix} ${pathDisplay} → ${route.file}\n`;
  });
  
  return ascii;
}

// テーブルを生成
function generateTable(routes) {
  let table = '| パス | ファイル | 説明 | タイプ |\n';
  table += '|------|----------|------|------|\n';
  
  routes.forEach(route => {
    table += `| \`${route.path}\` | \`${route.file}\` | ${route.description} | ${route.type} |\n`;
  });
  
  return table;
}

// メイン処理
function generateRouteDiagram() {
  try {
    console.log('📊 ルート構造図を生成中...');
    
    const routes = parseRoutes(routesContent);
    const routesWithDesc = addRouteDescriptions(routes);
    
    const mermaidDiagram = generateMermaidDiagram(routesWithDesc);
    const asciiDiagram = generateASCIIDiagram(routesWithDesc);
    const table = generateTable(routesWithDesc);
    
    const timestamp = new Date().toLocaleString('ja-JP');
    
    const output = `# React Router v7 ルート階層構造

> 🤖 自動生成日時: ${timestamp}  
> 📄 生成元: \`app/routes.ts\`

## Mermaid図

\`\`\`mermaid
${mermaidDiagram}
\`\`\`

## ツリー構造（ASCII）

\`\`\`
${asciiDiagram}
\`\`\`

## ルート一覧

${table}

## 実行方法

\`\`\`bash
# 図を再生成
npm run generate-routes

# または直接実行
node scripts/generate-route-diagram.js
\`\`\`

## 注意事項

- このファイルは自動生成されます
- \`routes.ts\` を変更後は \`npm run generate-routes\` で更新してください
- 手動で編集しないでください（変更は失われます）
`;
    
    fs.writeFileSync(outputPath, output, 'utf-8');
    
    console.log('✅ ルート構造図を生成しました:');
    console.log(`   📄 ${outputPath}`);
    console.log(`   📊 ${routesWithDesc.length} routes processed`);
    
  } catch (error) {
    console.error('❌ エラーが発生しました:', error.message);
    process.exit(1);
  }
}

// スクリプト実行
if (import.meta.url === `file://${process.argv[1]}`) {
  generateRouteDiagram();
}

export { generateRouteDiagram };