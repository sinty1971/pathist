#!/usr/bin/env node

// Swagger 2.0からOpenAPI 3.0に変換するスクリプト
// 使い方: node convert-to-openapi3.js

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import { parse, stringify } from 'yaml';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Swagger 2.0をOpenAPI 3.0に変換
function convertSwagger2ToOpenAPI3(swagger2) {
    const openapi3 = {
        openapi: '3.0.3',
        info: swagger2.info,
        servers: [{
            url: `http://${swagger2.host}${swagger2.basePath || ''}`
        }],
        paths: {},
        components: {
            schemas: swagger2.definitions || {}
        }
    };

    // パスの変換
    if (swagger2.paths) {
        Object.keys(swagger2.paths).forEach(path => {
            openapi3.paths[path] = {};
            
            Object.keys(swagger2.paths[path]).forEach(method => {
                const operation = swagger2.paths[path][method];
                const convertedOp = {
                    ...operation,
                    responses: {}
                };

                // パラメータの変換
                if (operation.parameters) {
                    convertedOp.parameters = operation.parameters.filter(p => p.in !== 'body');
                    
                    // bodyパラメータをrequestBodyに変換
                    const bodyParam = operation.parameters.find(p => p.in === 'body');
                    if (bodyParam) {
                        convertedOp.requestBody = {
                            description: bodyParam.description,
                            required: bodyParam.required,
                            content: {
                                'application/json': {
                                    schema: bodyParam.schema
                                }
                            }
                        };
                    }
                }

                // レスポンスの変換
                if (operation.responses) {
                    Object.keys(operation.responses).forEach(status => {
                        const response = operation.responses[status];
                        convertedOp.responses[status] = {
                            description: response.description
                        };

                        if (response.schema) {
                            convertedOp.responses[status].content = {
                                'application/json': {
                                    schema: response.schema
                                }
                            };
                        }
                    });
                }

                // consumesとproducesの削除（OpenAPI 3.0では使用しない）
                delete convertedOp.consumes;
                delete convertedOp.produces;

                openapi3.paths[path][method] = convertedOp;
            });
        });
    }

    // $refの更新
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
        // Swagger 2.0ファイルを読み込み
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