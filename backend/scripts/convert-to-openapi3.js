#!/usr/bin/env node

// Swagger 2.0からOpenAPI 3.0に変換するスクリプト
// 使い方: node convert-to-openapi3.js

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { parse, stringify } from 'yaml';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// OpenAPI 3.1をOpenAPI 3.0に変換（swaggo v2対応）
function convertSwagger2ToOpenAPI3(swagger2) {
    const openapi3 = {
        openapi: '3.0.3',
        info: swagger2.info,
        servers: swagger2.servers || [{
            url: 'http://localhost:8080/api'
        }],
        paths: {},
        components: {
            schemas: swagger2.components?.schemas || swagger2.definitions || {}
        }
    };

    // パスの変換（OpenAPI 3.1から3.0への変換）
    if (swagger2.paths) {
        // swaggo v2はすでにOpenAPI 3.1形式なので、そのまま使用
        openapi3.paths = swagger2.paths;
    }

    // $refの更新（必要に応じて）
    const jsonString = JSON.stringify(openapi3);
    const updatedJson = jsonString.replace(
        /#\/definitions\//g,
        '#/components/schemas/'
    );

    return JSON.parse(updatedJson);
}

// メイン処理
async function main() {
    try {
        // OpenAPI 3.1ファイルを読み込み（swaggo v2生成）
        const swagger2Path = path.join(__dirname, '../../schemas/openapi.yaml');
        const swagger2Content = fs.readFileSync(swagger2Path, 'utf8');
        const swagger2 = parse(swagger2Content);

        // OpenAPI 3.0に変換
        const openapi3 = convertSwagger2ToOpenAPI3(swagger2);

        // 出力
        const outputPath = path.join(__dirname, '../../schemas/openapi-v3.yaml');
        fs.writeFileSync(outputPath, stringify(openapi3));

        const outputJsonPath = path.join(__dirname, '../../schemas/openapi-v3.json');
        fs.writeFileSync(outputJsonPath, JSON.stringify(openapi3, null, 2));

        console.log('✅ OpenAPI 3.0への変換が完了しました:');
        console.log(`  - YAML: ${outputPath}`);
        console.log(`  - JSON: ${outputJsonPath}`);
    } catch (error) {
        console.error('❌ エラー:', error);
        process.exit(1);
    }
}

// 直接実行
main();