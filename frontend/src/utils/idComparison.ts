/**
 * フルパスID vs Len7 ID 比較分析
 */

import { FastID } from './fastId';

// 現在のフルパスID例（実際のデータ）
const sampleFullPaths = [
  "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場",
  "豊田築炉/2-工事/2025-0615 豊田築炉 東海工場", 
  "豊田築炉/2-工事/2025-0620 豊田築炉 刈谷工場",
  "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場/工事.xlsx",
  "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場/図面.pdf",
];

export interface IDComparison {
  fullPath: string;
  fullPathLength: number;
  len7Id: string;
  len7Length: number;
  memoryReduction: number;
  collisionRisk: string;
}

export function analyzeIdConversion(): IDComparison[] {
  return sampleFullPaths.map(fullPath => {
    const len7Id = FastID.fromString(fullPath).len7();
    const fullPathLength = fullPath.length;
    const len7Length = len7Id.length;
    const memoryReduction = Math.round((1 - len7Length / fullPathLength) * 100);
    
    return {
      fullPath,
      fullPathLength,
      len7Id,
      len7Length,
      memoryReduction,
      collisionRisk: calculateCollisionRisk(len7Id)
    };
  });
}

function calculateCollisionRisk(len7Id: string): string {
  // 32^7 = 34,359,738,368 (約343億通り)
  const totalCombinations = Math.pow(32, 7);
  const expectedCollisionAt = Math.sqrt(totalCombinations); // 約185,000件
  
  if (expectedCollisionAt > 100000) {
    return "極めて低い（10万件以下では衝突なし）";
  } else if (expectedCollisionAt > 10000) {
    return "低い（1万件以下では衝突なし）";
  } else {
    return "要注意（少数でも衝突の可能性）";
  }
}

/**
 * メモリ使用量とパフォーマンスの比較
 */
export interface PerformanceComparison {
  scenario: string;
  fullPathMemory: number;
  len7Memory: number;
  memoryReduction: number;
  searchSpeedImprovement: number;
  hashMapPerformance: string;
}

export function analyzePerformanceImpact(fileCount: number): PerformanceComparison[] {
  const avgFullPathLength = 45; // 日本語文字含む平均的なパス長
  const len7Length = 7;
  
  const scenarios = [
    { name: "小規模プロジェクト", files: 100 },
    { name: "中規模プロジェクト", files: 1000 },
    { name: "大規模プロジェクト", files: 10000 },
    { name: "超大規模プロジェクト", files: 100000 },
  ];

  return scenarios.map(scenario => {
    const fullPathMemory = scenario.files * avgFullPathLength * 2; // UTF-16
    const len7Memory = scenario.files * len7Length * 2;
    const memoryReduction = Math.round((1 - len7Memory / fullPathMemory) * 100);
    
    // 文字列比較の計算量改善
    const searchSpeedImprovement = Math.round(avgFullPathLength / len7Length * 100) / 100;
    
    // ハッシュマップのパフォーマンス
    let hashMapPerformance = "";
    if (scenario.files < 1000) {
      hashMapPerformance = "差異なし";
    } else if (scenario.files < 10000) {
      hashMapPerformance = "軽微な改善";
    } else {
      hashMapPerformance = "大幅な改善";
    }

    return {
      scenario: scenario.name,
      fullPathMemory,
      len7Memory,
      memoryReduction,
      searchSpeedImprovement,
      hashMapPerformance
    };
  });
}

/**
 * 実装時の影響範囲分析
 */
export interface ImplementationImpact {
  component: string;
  currentUsage: string;
  requiredChanges: string[];
  riskLevel: "低" | "中" | "高";
  estimatedEffort: string;
}

export function analyzeImplementationImpact(): ImplementationImpact[] {
  return [
    {
      component: "Files.tsx (TreeView)",
      currentUsage: "fileInfo.path をIDとして直接使用",
      requiredChanges: [
        "convertToTreeNode関数でLen7 IDを生成",
        "updateNodeInTree関数の一致判定を変更",
        "findNodeById関数の検索ロジック変更",
        "React key属性の変更"
      ],
      riskLevel: "中",
      estimatedEffort: "2-3時間"
    },
    {
      component: "KojiGanttChart.tsx",
      currentUsage: "koji.id で工事を識別",
      requiredChanges: [
        "バックエンドから受け取るID形式に依存",
        "選択状態管理の変更",
        "配列操作での一致判定変更"
      ],
      riskLevel: "低",
      estimatedEffort: "1時間"
    },
    {
      component: "Kojies.tsx",
      currentUsage: "koji.id で工事を識別",
      requiredChanges: [
        "handleKojiUpdate関数の修正",
        "配列操作での一致判定変更"
      ],
      riskLevel: "低", 
      estimatedEffort: "1時間"
    },
    {
      component: "バックエンドAPI",
      currentUsage: "フルパスベースのID生成",
      requiredChanges: [
        "models/koji.go のID生成をLen7に変更",
        "models/company.go のID生成をLen7に変更",
        "既存データの移行スクリプト作成"
      ],
      riskLevel: "高",
      estimatedEffort: "4-6時間"
    }
  ];
}

/**
 * 衝突テスト用の関数
 */
export function testCollisions(samplePaths: string[]): {
  totalPaths: number;
  uniqueLen7Ids: number;
  collisions: Array<{len7Id: string; paths: string[]}>;
  collisionRate: number;
} {
  const len7Map = new Map<string, string[]>();
  
  samplePaths.forEach(path => {
    const len7Id = FastID.fromString(path).len7();
    if (!len7Map.has(len7Id)) {
      len7Map.set(len7Id, []);
    }
    len7Map.get(len7Id)!.push(path);
  });
  
  const collisions = Array.from(len7Map.entries())
    .filter(([_, paths]) => paths.length > 1)
    .map(([len7Id, paths]) => ({ len7Id, paths }));
  
  return {
    totalPaths: samplePaths.length,
    uniqueLen7Ids: len7Map.size,
    collisions,
    collisionRate: collisions.length / samplePaths.length * 100
  };
}

/**
 * 推奨される移行戦略
 */
export interface MigrationStrategy {
  phase: number;
  description: string;
  changes: string[];
  riskLevel: "低" | "中" | "高";
  rollbackPlan: string;
}

export function getMigrationStrategy(): MigrationStrategy[] {
  return [
    {
      phase: 1,
      description: "フロントエンド準備段階",
      changes: [
        "FastID クラスの実装とテスト",
        "衝突テストの実行",
        "デバッグ用の変換テーブル作成"
      ],
      riskLevel: "低",
      rollbackPlan: "新しいユーティリティを削除するだけ"
    },
    {
      phase: 2,
      description: "バックエンドID生成変更",
      changes: [
        "models/*.go のID生成をLen7に変更",
        "既存APIレスポンスのID形式変更",
        "データベース移行（該当する場合）"
      ],
      riskLevel: "高",
      rollbackPlan: "Gitリバート + データ復旧"
    },
    {
      phase: 3,
      description: "フロントエンド適用",
      changes: [
        "各コンポーネントでの ID使用箇所を変更",
        "TreeView のキー管理変更",
        "状態管理の一致判定変更"
      ],
      riskLevel: "中",
      rollbackPlan: "フロントエンドのみリバート可能"
    },
    {
      phase: 4,
      description: "検証とクリーンアップ",
      changes: [
        "全機能のテスト実行",
        "パフォーマンステスト",
        "古いID関連コードの削除"
      ],
      riskLevel: "低",
      rollbackPlan: "機能単位でのロールバック"
    }
  ];
}