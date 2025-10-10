/**
 * 会社関連のユーティリティ関数
 */

import { getBusinessCompaniesCategories } from "@/api/sdk.gen";
import type { ModelsCompanyCategoryInfo } from "@/api/types.gen";

/**
 * 業種カテゴリーのデフォルトデータ（フォールバック用）
 */
const defaultCategories: ModelsCompanyCategoryInfo[] = [
  { code: 0, name: "特別" },
  { code: 1, name: "下請会社" },
  { code: 2, name: "築炉会社" },
  { code: 3, name: "一人親方" },
  { code: 4, name: "元請け" },
  { code: 5, name: "リース会社" },
  { code: 6, name: "販売会社" },
  { code: 7, name: "販売会社" },
  { code: 8, name: "求人会社" },
  { code: 9, name: "その他" }
];

/**
 * カテゴリーキャッシュ
 */
let cachedCategories: ModelsCompanyCategoryInfo[] | null = null;

/**
 * APIからカテゴリー一覧を取得する
 */
export async function loadCategories(forceRefresh: boolean = false): Promise<ModelsCompanyCategoryInfo[]> {
  // キャッシュが存在し、強制更新でない場合はキャッシュを返す
  if (!forceRefresh && cachedCategories) {
    return cachedCategories;
  }

  try {
    const response = await getBusinessCompaniesCategories();
    if (response.data) {
      // APIレスポンスは ModelsCompanyCategoryInfo[] 形式
      cachedCategories = response.data;
      return response.data;
    }
  } catch (error) {
    console.warn("APIから会社カテゴリーの取得に失敗しました。デフォルトカテゴリーを使用します:", error);
  }

  // APIが失敗した場合はデフォルトカテゴリーを使用
  cachedCategories = defaultCategories;
  return defaultCategories;
}

/**
 * カテゴリー番号から業種名を取得する（キャッシュ優先）
 * @param category カテゴリー番号
 * @returns 業種名
 */
export function getCategoryName(category: number | undefined | null): string {
  // categoryが0の場合も正しく処理されるようにする
  if (category === undefined || category === null) {
    return "業種未設定";
  }
  
  // キャッシュされたカテゴリーを優先して使用
  const categories = cachedCategories || defaultCategories;
  const found = categories.find(cat => cat.code === category);
  return found?.name || "業種未設定";
}

/**
 * カテゴリー番号から業種名を取得する（APIから最新データを取得）
 * @param category カテゴリー番号
 * @returns 業種名
 */
export async function getCategoryNameFromAPI(category: number | undefined | null): Promise<string> {
  // categoryが0の場合も正しく処理されるようにする
  if (category === undefined || category === null) {
    return "業種未設定";
  }
  
  try {
    // APIから最新のカテゴリーを取得
    const categories = await loadCategories();
    const found = categories.find(cat => cat.code === category);
    return found?.name || "業種未設定";
  } catch (error) {
    // APIが失敗した場合はデフォルトカテゴリーを使用
    const found = defaultCategories.find(cat => cat.code === category);
    return found?.name || "業種未設定";
  }
}


/**
 * カテゴリーを初期化する
 */
export async function initializeCategories(): Promise<void> {
  await loadCategories();
}

/**
 * カテゴリーキャッシュを更新する（設定モーダルで使用）
 */
export async function refreshCategoriesCache(): Promise<void> {
  await loadCategories(true);
}

/**
 * カテゴリーキャッシュをクリアする
 */
export function clearCategoriesCache(): void {
  cachedCategories = null;
}