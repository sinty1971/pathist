/**
 * ID変換のデモとベンチマーク
 */

import { FastID } from './fastId';
import { 
  analyzeIdConversion, 
  analyzePerformanceImpact, 
  analyzeImplementationImpact,
  testCollisions,
  getMigrationStrategy 
} from './idComparison';

/**
 * 実際のサンプルデータでのデモ実行
 */
export function runIdConversionDemo(): void {
  console.log('='.repeat(60));
  console.log('📊 フルパスID vs Len7 ID 変換分析');
  console.log('='.repeat(60));

  // 1. ID変換の例
  console.log('\n1️⃣ ID変換例:');
  const conversions = analyzeIdConversion();
  conversions.forEach(conv => {
    console.log(`\nフルパス: ${conv.fullPath}`);
    console.log(`  長さ: ${conv.fullPathLength}文字`);
    console.log(`  Len7: ${conv.len7Id}`);
    console.log(`  短縮: ${conv.memoryReduction}% 削減`);
    console.log(`  衝突リスク: ${conv.collisionRisk}`);
  });

  // 2. パフォーマンス影響
  console.log('\n2️⃣ パフォーマンス影響:');
  const performance = analyzePerformanceImpact(10000);
  performance.forEach(perf => {
    console.log(`\n${perf.scenario}:`);
    console.log(`  フルパスメモリ: ${(perf.fullPathMemory / 1024).toFixed(1)}KB`);
    console.log(`  Len7メモリ: ${(perf.len7Memory / 1024).toFixed(1)}KB`);
    console.log(`  メモリ削減: ${perf.memoryReduction}%`);
    console.log(`  検索速度: ${perf.searchSpeedImprovement}倍高速`);
    console.log(`  ハッシュマップ: ${perf.hashMapPerformance}`);
  });

  // 3. 実装影響
  console.log('\n3️⃣ 実装影響範囲:');
  const impacts = analyzeImplementationImpact();
  impacts.forEach(impact => {
    console.log(`\n${impact.component}:`);
    console.log(`  現在: ${impact.currentUsage}`);
    console.log(`  変更: ${impact.requiredChanges.join(', ')}`);
    console.log(`  リスク: ${impact.riskLevel}`);
    console.log(`  工数: ${impact.estimatedEffort}`);
  });

  // 4. 移行戦略
  console.log('\n4️⃣ 推奨移行戦略:');
  const strategy = getMigrationStrategy();
  strategy.forEach(phase => {
    console.log(`\nフェーズ${phase.phase}: ${phase.description}`);
    console.log(`  変更: ${phase.changes.join(', ')}`);
    console.log(`  リスク: ${phase.riskLevel}`);
    console.log(`  回復: ${phase.rollbackPlan}`);
  });
}

/**
 * 衝突テストの実行
 */
export function runCollisionTest(): void {
  console.log('\n🔍 衝突テスト実行中...');
  
  // 大量のテストデータを生成
  const testPaths: string[] = [];
  const companies = ['豊田築炉', 'ABC建設', 'XYZ工業', 'DEF製作所'];
  const locations = ['名和工場', '東海工場', '刈谷工場', '豊田工場', '岡崎工場'];
  const files = ['工事.xlsx', '図面.pdf', '写真.jpg', '仕様書.docx', '契約書.pdf'];

  // 10,000件のテストデータを生成
  for (let year = 2020; year <= 2025; year++) {
    for (let month = 1; month <= 12; month++) {
      for (let day = 1; day <= 28; day++) {
        for (const company of companies) {
          for (const location of locations) {
            const dateStr = `${year}-${month.toString().padStart(2, '0')}${day.toString().padStart(2, '0')}`;
            const basePath = `${company}/2-工事/${dateStr} ${company} ${location}`;
            testPaths.push(basePath);
            
            // ファイルも追加
            files.forEach(file => {
              testPaths.push(`${basePath}/${file}`);
            });
          }
        }
      }
    }
    
    // 10,000件を超えたら停止
    if (testPaths.length > 10000) break;
  }

  const result = testCollisions(testPaths.slice(0, 10000));
  
  console.log(`\n📈 衝突テスト結果:`);
  console.log(`  総パス数: ${result.totalPaths.toLocaleString()}`);
  console.log(`  ユニークLen7 ID数: ${result.uniqueLen7Ids.toLocaleString()}`);
  console.log(`  衝突数: ${result.collisions.length}`);
  console.log(`  衝突率: ${result.collisionRate.toFixed(4)}%`);
  
  if (result.collisions.length > 0) {
    console.log('\n⚠️ 検出された衝突:');
    result.collisions.slice(0, 5).forEach(collision => {
      console.log(`  ID: ${collision.len7Id}`);
      collision.paths.forEach(path => console.log(`    - ${path}`));
    });
  } else {
    console.log('\n✅ 衝突は検出されませんでした！');
  }
}

/**
 * パフォーマンスベンチマーク
 */
export function runPerformanceBenchmark(): void {
  console.log('\n⚡ パフォーマンスベンチマーク実行中...');
  
  const testData = [
    "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場",
    "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場/工事.xlsx",
    "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場/図面フォルダ/設計図.pdf",
  ];

  const iterations = 10000;

  // フルパス検索ベンチマーク
  console.time('フルパス検索');
  for (let i = 0; i < iterations; i++) {
    const target = testData[i % testData.length];
    testData.findIndex(path => path === target);
  }
  console.timeEnd('フルパス検索');

  // Len7 ID検索ベンチマーク
  const len7Ids = testData.map(path => FastID.fromString(path).len7());
  console.time('Len7 ID検索');
  for (let i = 0; i < iterations; i++) {
    const target = len7Ids[i % len7Ids.length];
    len7Ids.findIndex(id => id === target);
  }
  console.timeEnd('Len7 ID検索');

  // メモリ使用量比較
  const fullPathMemory = testData.join('').length * 2; // UTF-16
  const len7Memory = len7Ids.join('').length * 2;
  
  console.log(`\n💾 メモリ使用量比較 (${testData.length}件):`);
  console.log(`  フルパス: ${fullPathMemory} bytes`);
  console.log(`  Len7: ${len7Memory} bytes`);
  console.log(`  削減: ${Math.round((1 - len7Memory / fullPathMemory) * 100)}%`);
}

// ブラウザ環境での実行用
if (typeof window !== 'undefined') {
  // グローバルに公開してブラウザコンソールから実行可能に
  (window as any).runIdDemo = runIdConversionDemo;
  (window as any).runCollisionTest = runCollisionTest;
  (window as any).runPerformanceBenchmark = runPerformanceBenchmark;
}