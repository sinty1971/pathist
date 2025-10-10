#!/usr/bin/env node

import { fileURLToPath } from 'url';
import path from 'path';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const legacyRoutesPath = path.join(__dirname, '../app/routes.ts');
const legacyOutputPath = path.join(__dirname, '../route-structure-generated.md');

console.log('⚠️  Next.js へ移行したため、このスクリプトは利用できません。');
console.log('   Next.js のルート構造は app ディレクトリ内のフォルダ構成に対応します。');
console.log('   必要であれば新しいルート図生成スクリプトを追加してください。');
console.log(`   旧スクリプトは React Router 用 ( ${legacyRoutesPath} → ${legacyOutputPath} ) でした。`);
