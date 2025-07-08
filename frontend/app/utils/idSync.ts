/**
 * バックエンドとフロントエンドのID同期システム
 * 
 * このシステムは以下の機能を提供します：
 * 1. ID生成アルゴリズムの統一（Go ↔ TypeScript）
 * 2. ID一致検証機能
 * 3. 自動・半自動同期ロジック
 * 4. 移行サポート機能
 */

import { FastID } from './fastId';

// ID生成要素の定義
export interface IdComponents {
  /** 開始日 */
  startDate: Date;
  /** 会社名 */
  companyName: string;
  /** 場所名 */
  locationName: string;
  /** その他のファイルパス情報 */
  additionalPath?: string;
}

/**
 * ID同期管理クラス
 */
export class IdSyncManager {
  private static instance: IdSyncManager;
  private validationCache = new Map<string, boolean>();
  private idMappingCache = new Map<string, string>();

  private constructor() {}

  static getInstance(): IdSyncManager {
    if (!IdSyncManager.instance) {
      IdSyncManager.instance = new IdSyncManager();
    }
    return IdSyncManager.instance;
  }

  /**
   * 工事IDの生成（バックエンドと同じアルゴリズム）
   */
  generateKojiId(components: IdComponents): string {
    const { startDate, companyName, locationName } = components;
    
    // Goと同じフォーマットで文字列を構築
    const year = startDate.getFullYear();
    const month = (startDate.getMonth() + 1).toString().padStart(2, '0');
    const day = startDate.getDate().toString().padStart(2, '0');
    const dateStr = `${year}${month}${day}`;
    
    const combined = `${dateStr}${companyName}${locationName}`;
    return FastID.fromString(combined).len5();
  }

  /**
   * フルパスからLen7 IDを生成
   */
  generatePathId(fullPath: string): string {
    return FastID.fromString(fullPath).len7();
  }

  /**
   * IDの一致検証（バックエンドとの同期確認）
   */
  async validateIdConsistency(
    frontendId: string,
    components: IdComponents,
    apiEndpoint?: string
  ): Promise<{
    isValid: boolean;
    backendId?: string;
    needsSync: boolean;
    error?: string;
  }> {
    try {
      // キャッシュから確認
      const cacheKey = `${frontendId}_${JSON.stringify(components)}`;
      if (this.validationCache.has(cacheKey)) {
        return {
          isValid: this.validationCache.get(cacheKey)!,
          needsSync: false
        };
      }

      // フロントエンドで生成されるIDを計算
      const expectedId = this.generateKojiId(components);
      const isValid = frontendId === expectedId;

      // APIエンドポイントが指定されている場合は、バックエンドとも比較
      let backendId: string | undefined;
      if (apiEndpoint) {
        try {
          const response = await fetch(apiEndpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(components)
          });
          
          if (response.ok) {
            const data = await response.json();
            backendId = data.id;
          }
        } catch (error) {
          console.warn('バックエンドIDの検証に失敗:', error);
        }
      }

      // 結果をキャッシュ
      this.validationCache.set(cacheKey, isValid);
      
      return {
        isValid,
        backendId,
        needsSync: backendId ? frontendId !== backendId : false
      };
    } catch (error) {
      return {
        isValid: false,
        needsSync: true,
        error: error instanceof Error ? error.message : '不明なエラー'
      };
    }
  }

  /**
   * 自動同期：IDの不整合を検出し、自動的に修正
   */
  async autoSync(
    currentId: string,
    components: IdComponents,
    updateCallback: (newId: string) => void
  ): Promise<{ success: boolean; newId?: string; error?: string }> {
    try {
      const correctId = this.generateKojiId(components);
      
      if (currentId !== correctId) {
        // IDの不整合を検出
        console.log(`ID不整合を検出: ${currentId} → ${correctId}`);
        
        // 自動修正
        updateCallback(correctId);
        
        return {
          success: true,
          newId: correctId
        };
      }
      
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : '同期エラー'
      };
    }
  }

  /**
   * 半自動同期：ユーザーの確認を求めながら同期
   */
  async semiAutoSync(
    currentId: string,
    components: IdComponents,
    confirmCallback: (oldId: string, newId: string) => Promise<boolean>
  ): Promise<{ success: boolean; newId?: string; cancelled?: boolean }> {
    try {
      const correctId = this.generateKojiId(components);
      
      if (currentId !== correctId) {
        // ユーザーに確認
        const shouldUpdate = await confirmCallback(currentId, correctId);
        
        if (shouldUpdate) {
          return {
            success: true,
            newId: correctId
          };
        } else {
          return {
            success: false,
            cancelled: true
          };
        }
      }
      
      return { success: true };
    } catch (error) {
      return { success: false };
    }
  }

  /**
   * フルパスID → Len7 ID 変換テーブルの生成
   */
  createIdMappingTable(fullPathIds: string[]): Map<string, string> {
    const mapping = new Map<string, string>();
    
    fullPathIds.forEach(fullPath => {
      const len7Id = this.generatePathId(fullPath);
      mapping.set(fullPath, len7Id);
      
      // 逆引き用
      this.idMappingCache.set(len7Id, fullPath);
    });
    
    return mapping;
  }

  /**
   * Len7 ID → フルパスID 逆引き
   */
  getFullPathFromLen7(len7Id: string): string | null {
    return this.idMappingCache.get(len7Id) || null;
  }

  /**
   * 一括ID変換とバリデーション
   */
  async bulkConvertAndValidate(
    items: Array<{ id: string; components: IdComponents }>
  ): Promise<Array<{
    originalId: string;
    newId: string;
    isValid: boolean;
    needsUpdate: boolean;
  }>> {
    const results = [];
    
    for (const item of items) {
      const correctId = this.generateKojiId(item.components);
      const isValid = item.id === correctId;
      
      results.push({
        originalId: item.id,
        newId: correctId,
        isValid,
        needsUpdate: !isValid
      });
    }
    
    return results;
  }

  /**
   * キャッシュクリア
   */
  clearCache(): void {
    this.validationCache.clear();
    this.idMappingCache.clear();
  }
}

/**
 * React Hook for ID同期
 */
export function useIdSync() {
  const syncManager = IdSyncManager.getInstance();
  
  return {
    generateKojiId: (components: IdComponents) => 
      syncManager.generateKojiId(components),
    
    generatePathId: (fullPath: string) => 
      syncManager.generatePathId(fullPath),
    
    validateId: (id: string, components: IdComponents, apiEndpoint?: string) => 
      syncManager.validateIdConsistency(id, components, apiEndpoint),
    
    autoSync: (id: string, components: IdComponents, updateCallback: (newId: string) => void) => 
      syncManager.autoSync(id, components, updateCallback),
    
    semiAutoSync: (id: string, components: IdComponents, confirmCallback: (oldId: string, newId: string) => Promise<boolean>) => 
      syncManager.semiAutoSync(id, components, confirmCallback),
    
    createMappingTable: (fullPathIds: string[]) => 
      syncManager.createIdMappingTable(fullPathIds),
    
    getFullPathFromLen7: (len7Id: string) => 
      syncManager.getFullPathFromLen7(len7Id),
    
    bulkConvertAndValidate: (items: Array<{ id: string; components: IdComponents }>) => 
      syncManager.bulkConvertAndValidate(items),
    
    clearCache: () => syncManager.clearCache()
  };
}

/**
 * 使用例とデモンストレーション
 */
export const idSyncExamples = {
  // 基本的な使用例
  basic: () => {
    const sync = IdSyncManager.getInstance();
    
    const components: IdComponents = {
      startDate: new Date(2025, 5, 18), // 2025-06-18
      companyName: '豊田築炉',
      locationName: '名和工場'
    };
    
    const id = sync.generateKojiId(components);
    console.log('Generated ID:', id);
    
    return id;
  },

  // 自動同期のデモ
  autoSyncDemo: async () => {
    const sync = IdSyncManager.getInstance();
    
    const components: IdComponents = {
      startDate: new Date(2025, 5, 18),
      companyName: '豊田築炉',
      locationName: '名和工場'
    };
    
    const wrongId = 'WRONG';
    
    const result = await sync.autoSync(wrongId, components, (newId) => {
      console.log(`ID自動修正: ${wrongId} → ${newId}`);
    });
    
    return result;
  },

  // 半自動同期のデモ
  semiAutoSyncDemo: async () => {
    const sync = IdSyncManager.getInstance();
    
    const components: IdComponents = {
      startDate: new Date(2025, 5, 18),
      companyName: '豊田築炉',
      locationName: '名和工場'
    };
    
    const wrongId = 'WRONG';
    
    const result = await sync.semiAutoSync(wrongId, components, async (oldId, newId) => {
      console.log(`ID変更確認: ${oldId} → ${newId}`);
      return confirm(`IDを ${oldId} から ${newId} に変更しますか？`);
    });
    
    return result;
  }
};