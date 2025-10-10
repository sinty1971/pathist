/**
 * 自動ID同期 React Hook
 * 
 * このHookは工事データのIDが変更された際に、自動的にコンポーネントの状態を
 * 同期させるための機能を提供します。
 */

import { useEffect, useCallback, useState } from 'react';
import { useIdSync, IdComponents } from '@/utils/idSync';

export interface AutoIdSyncOptions {
  /** 自動同期を有効にするかどうか（デフォルト: true） */
  enabled?: boolean;
  /** 同期前の確認を求めるかどうか（デフォルト: false） */
  requireConfirmation?: boolean;
  /** 同期失敗時のカスタムエラーハンドラー */
  onError?: (error: string) => void;
  /** 同期成功時のカスタムハンドラー */
  onSuccess?: (newId: string) => void;
}

export interface AutoIdSyncResult {
  /** 現在のID */
  currentId: string;
  /** ID同期中かどうか */
  isSyncing: boolean;
  /** 同期エラー */
  syncError: string | null;
  /** 手動同期を実行する関数 */
  sync: () => Promise<void>;
  /** IDを更新する関数 */
  updateId: (newId: string) => void;
  /** 最後の同期時刻 */
  lastSyncTime: Date | null;
}

/**
 * 自動ID同期Hook
 * 
 * @param initialId 初期ID
 * @param components ID生成に必要な要素
 * @param options 同期オプション
 * @returns 同期結果と制御関数
 */
export function useAutoIdSync(
  initialId: string,
  components: IdComponents,
  options: AutoIdSyncOptions = {}
): AutoIdSyncResult {
  const {
    enabled = true,
    requireConfirmation = false,
    onError,
    onSuccess
  } = options;

  const [currentId, setCurrentId] = useState(initialId);
  const [isSyncing, setIsSyncing] = useState(false);
  const [syncError, setSyncError] = useState<string | null>(null);
  const [lastSyncTime, setLastSyncTime] = useState<Date | null>(null);

  const { autoSync, semiAutoSync, validateId } = useIdSync();

  // ID同期を実行
  const sync = useCallback(async () => {
    if (!enabled || isSyncing) return;

    setIsSyncing(true);
    setSyncError(null);

    try {
      // バリデーション
      const validation = await validateId(currentId, components);
      
      if (!validation.isValid && validation.needsSync) {
        // 同期が必要
        if (requireConfirmation) {
          // 半自動同期（確認あり）
          const result = await semiAutoSync(currentId, components, async (oldId, newId) => {
            return confirm(`工事IDを更新しますか？\n${oldId} → ${newId}`);
          });

          if (result.success && result.newId) {
            setCurrentId(result.newId);
            setLastSyncTime(new Date());
            onSuccess?.(result.newId);
          } else if (result.cancelled) {
            // ユーザーがキャンセルした場合は何もしない
          }
        } else {
          // 自動同期（確認なし）
          const result = await autoSync(currentId, components, (newId) => {
            setCurrentId(newId);
            setLastSyncTime(new Date());
            onSuccess?.(newId);
          });

          if (!result.success && result.error) {
            setSyncError(result.error);
            onError?.(result.error);
          }
        }
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'ID同期でエラーが発生しました';
      setSyncError(errorMessage);
      onError?.(errorMessage);
    } finally {
      setIsSyncing(false);
    }
  }, [currentId, components, enabled, isSyncing, requireConfirmation, autoSync, semiAutoSync, validateId, onError, onSuccess]);

  // ID更新関数
  const updateId = useCallback((newId: string) => {
    setCurrentId(newId);
    setSyncError(null);
  }, []);

  // components変更時の自動同期
  useEffect(() => {
    if (enabled && !isSyncing) {
      const timer = setTimeout(sync, 100); // 100ms遅延で同期実行
      return () => clearTimeout(timer);
    }
  }, [components, enabled, isSyncing, sync]);

  // 初期IDが変更された場合の同期
  useEffect(() => {
    if (initialId !== currentId && !isSyncing) {
      setCurrentId(initialId);
      // 新しいIDで同期チェック
      const timer = setTimeout(sync, 100);
      return () => clearTimeout(timer);
    }
  }, [initialId, currentId, isSyncing, sync]);

  return {
    currentId,
    isSyncing,
    syncError,
    sync,
    updateId,
    lastSyncTime
  };
}

/**
 * 複数の工事データに対する一括ID同期Hook
 */
export function useBulkIdSync<T extends { id: string }>(
  items: T[],
  getComponents: (item: T) => IdComponents,
  options: AutoIdSyncOptions = {}
) {
  const [syncResults, setSyncResults] = useState<Map<string, AutoIdSyncResult>>(new Map());
  const [globalSyncing, setGlobalSyncing] = useState(false);

  const { bulkConvertAndValidate } = useIdSync();

  const syncAll = useCallback(async () => {
    if (globalSyncing || items.length === 0) return;

    setGlobalSyncing(true);

    try {
      const validationData = items.map(item => ({
        id: item.id,
        components: getComponents(item)
      }));

      const results = await bulkConvertAndValidate(validationData);
      
      // 更新が必要なアイテムを特定
      const needsUpdate = results.filter(r => r.needsUpdate);
      
      if (needsUpdate.length > 0) {
        const updateCount = needsUpdate.length;
        const shouldUpdate = !options.requireConfirmation || 
          confirm(`${updateCount}件の工事IDを更新しますか？`);

        if (shouldUpdate) {
          // 更新を実行
          needsUpdate.forEach(result => {
            options.onSuccess?.(result.newId);
          });
        }
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : '一括同期でエラーが発生しました';
      options.onError?.(errorMessage);
    } finally {
      setGlobalSyncing(false);
    }
  }, [items, getComponents, globalSyncing, bulkConvertAndValidate, options]);

  return {
    syncAll,
    globalSyncing,
    syncResults: Array.from(syncResults.values())
  };
}

/**
 * ファイルパスID同期Hook（フルパス → Len7変換用）
 */
export function usePathIdSync(
  fullPathIds: string[],
  options: { 
    autoConvert?: boolean;
    onMappingReady?: (mapping: Map<string, string>) => void;
  } = {}
) {
  const [mappingTable, setMappingTable] = useState<Map<string, string>>(new Map());
  const [isConverting, setIsConverting] = useState(false);

  const { createMappingTable, getFullPathFromLen7 } = useIdSync();

  // フルパスIDリストの変換
  useEffect(() => {
    if (fullPathIds.length === 0) return;

    setIsConverting(true);
    
    try {
      const mapping = createMappingTable(fullPathIds);
      setMappingTable(mapping);
      options.onMappingReady?.(mapping);
    } catch (error) {
      console.error('パスID変換エラー:', error);
    } finally {
      setIsConverting(false);
    }
  }, [fullPathIds, createMappingTable, options]);

  // ヘルパー関数
  const convertToLen7 = useCallback((fullPath: string): string => {
    return mappingTable.get(fullPath) || fullPath;
  }, [mappingTable]);

  const convertFromLen7 = useCallback((len7Id: string): string => {
    return getFullPathFromLen7(len7Id) || len7Id;
  }, [getFullPathFromLen7]);

  return {
    mappingTable,
    isConverting,
    convertToLen7,
    convertFromLen7
  };
}