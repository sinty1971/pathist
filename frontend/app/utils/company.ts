/**
 * 会社関連のユーティリティ関数
 */

import { getBusinessCategories } from "../api/sdk.gen";

/**
 * 業種カテゴリーのマッピング（フォールバック用）
 */
const defaultCategoryMap: Record<number, string> = {
  0: "特別",
  1: "下請会社",
  2: "築炉会社",
  3: "一人親方",
  4: "元請け",
  5: "リース会社",
  6: "販売会社",
  7: "販売会社",
  8: "求人会社",
  9: "その他"
};

/**
 * APIから取得したカテゴリーマッピング（キャッシュ）
 */
let cachedCategoryMap: Record<number, string> | null = null;

/**
 * APIからカテゴリーマッピングを取得する
 */
async function loadCategoryMap(): Promise<Record<number, string>> {
  if (cachedCategoryMap) {
    return cachedCategoryMap;
  }

  try {
    const response = await getBusinessCategories();
    if (response.data) {
      // APIレスポンスは Array<{[key: string]: string}> 形式
      const categoryMap: Record<number, string> = {};
      response.data.forEach(item => {
        for (const [key, value] of Object.entries(item)) {
          const numKey = parseInt(key, 10);
          if (!isNaN(numKey)) {
            categoryMap[numKey] = value;
          }
        }
      });
      cachedCategoryMap = categoryMap;
      return categoryMap;
    }
  } catch (error) {
    console.warn("Failed to load categories from API, using default mapping:", error);
  }

  // APIが失敗した場合はデフォルトマッピングを使用
  cachedCategoryMap = defaultCategoryMap;
  return defaultCategoryMap;
}

/**
 * カテゴリー番号から業種名を取得する
 * @param category カテゴリー番号
 * @returns 業種名
 */
export function getCategoryName(category: number | undefined | null): string {
  // categoryが0の場合も正しく処理されるようにする
  if (category === undefined || category === null) {
    return "業種未設定";
  }
  
  // キャッシュされたマッピングを使用
  const categoryMap = cachedCategoryMap || defaultCategoryMap;
  
  // 0は有効な値なので、hasOwnPropertyで確認
  if (categoryMap.hasOwnProperty(category)) {
    return categoryMap[category];
  }
  
  return "業種未設定";
}

/**
 * 全ての業種カテゴリーを取得する
 * @returns 業種カテゴリーの配列
 */
export function getAllCategories(): Array<{ value: number; label: string }> {
  const categoryMap = cachedCategoryMap || defaultCategoryMap;
  return Object.entries(categoryMap).map(([value, label]) => ({
    value: parseInt(value, 10),
    label
  }));
}

/**
 * 初期化時にカテゴリーマッピングを読み込む
 */
export async function initializeCategories(): Promise<void> {
  await loadCategoryMap();
}