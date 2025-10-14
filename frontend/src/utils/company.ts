/**
 * 会社関連のユーティリティ関数
 */

import { companyConnectClient } from "@/services/companyConnect";
import type { CompanyCategoryInfo } from "@/types/models";

/**
 * 業種カテゴリーのデフォルトデータ（フォールバック用）
 */
const defaultCategories: CompanyCategoryInfo[] = [
  { code: "0", label: "特別" },
  { code: "1", label: "下請会社" },
  { code: "2", label: "築炉会社" },
  { code: "3", label: "一人親方" },
  { code: "4", label: "元請け" },
  { code: "5", label: "リース会社" },
  { code: "6", label: "販売会社" },
  { code: "7", label: "販売会社" },
  { code: "8", label: "求人会社" },
  { code: "9", label: "その他" }
];

/**
 * カテゴリーキャッシュ
 */
let cachedCategories: CompanyCategoryInfo[] | null = null;

/**
 * APIからカテゴリー一覧を取得する
 */
export async function loadCategories(forceRefresh: boolean = false): Promise<CompanyCategoryInfo[]> {
  // キャッシュが存在し、強制更新でない場合はキャッシュを返す
  if (!forceRefresh && cachedCategories) {
    return cachedCategories;
  }

  try {
    const categories = await companyConnectClient.listCategories();
    cachedCategories = categories.map((cat) => ({
      code: String(cat.code),
      label: cat.label ?? "業種未設定",
    }));
    return cachedCategories;
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
export function getCategoryName(category: string | number | undefined | null): string {
  if (category === undefined || category === null) {
    return "業種未設定";
  }
  const categoryKey = String(category);
  // キャッシュされたカテゴリーを優先して使用
  const categories = cachedCategories || defaultCategories;
  const found = categories.find(cat => String(cat.code) === categoryKey);
  return found?.label || "業種未設定";
}

/**
 * カテゴリー番号から業種名を取得する（APIから最新データを取得）
 * @param category カテゴリー番号
 * @returns 業種名
 */
export async function getCategoryNameFromAPI(category: string | number | undefined | null): Promise<string> {
  // categoryが0の場合も正しく処理されるようにする
  if (category === undefined || category === null) {
    return "業種未設定";
  }
  
  try {
    // APIから最新のカテゴリーを取得
    const categories = await loadCategories();
    const categoryKey = String(category);
    const found = categories.find(cat => String(cat.code) === categoryKey);
    return found?.label || "業種未設定";
  } catch (error) {
    // APIが失敗した場合はデフォルトカテゴリーを使用
    const categoryKey = String(category);
    const found = defaultCategories.find(cat => String(cat.code) === categoryKey);
    return found?.label || "業種未設定";
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
