'use client';

import { useState, useEffect, useRef } from 'react';
import type { ModelsKoji } from '@/api/types.gen';
import KojiDetailModal from './KojiDetailModal';
import { kojiConnectClient } from '@/services/kojiConnect';
import '../styles/koji-gantt.css';
import '../styles/utilities.css';

interface GanttItem extends ModelsKoji {
  startX: number;
  width: number;
  row: number;
}

const KojiGanttChart = () => {
  const [kojies, setKojies] = useState<ModelsKoji[]>([]);
  const [ganttItems, setGanttItems] = useState<GanttItem[]>([]);
  const [selectedKoji, setSelectedKoji] = useState<ModelsKoji | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [viewStartDate, setViewStartDate] = useState(new Date());
  const [viewEndDate, setViewEndDate] = useState(new Date());
  // const [isDetailModalOpen, setIsDetailModalOpen] = useState(false); // 将来使用予定
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [visibleKojies, setVisibleKojies] = useState<ModelsKoji[]>([]);
  const [shouldReloadOnClose, setShouldReloadOnClose] = useState(false);
  
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [itemsPerPage, setItemsPerPage] = useState(10);
  const MIN_ITEMS = 5;
const DAY_WIDTH = 10; // ピクセル/日
const ROW_HEIGHT = 40; // ピクセル

const getKojiKey = (koji: ModelsKoji | null | undefined): string => {
  if (!koji) return "";
  const candidate = (koji as Record<string, unknown>).path;
  if (typeof candidate === "string" && candidate.length > 0) {
    return candidate;
  }
  if (koji.targetFolder && koji.targetFolder.length > 0) {
    return koji.targetFolder;
  }
  return koji.id ?? "";
};

  // 工事データを読み込み
  const loadKojies = async () => {
    try {
      setLoading(true);
      setError(null);
      const list = await kojiConnectClient.list();
      setKojies(list);
    } catch (err) {
      console.error('Error loading kojies:', err);
      setError(`工事データの読み込みに失敗しました: ${err instanceof Error ? err.message : 'Unknown error'}`);
    } finally {
      setLoading(false);
    }
  };

  // 初期ロード
  useEffect(() => {
    loadKojies();
    // ページが表示されるたびに初期スクロール状態をリセット
    setHasInitialScrolled(false);
  }, []);

  // 初回ロード時のみ今日の位置にスクロール
  const [hasInitialScrolled, setHasInitialScrolled] = useState(false);
  
  useEffect(() => {
    if (scrollContainerRef.current && viewStartDate && viewEndDate && ganttItems.length > 0 && !hasInitialScrolled) {
      // 少し遅延させてレンダリング完了後にスクロール
      setTimeout(() => {
        if (scrollContainerRef.current) {
          // 今日の位置を正しく計算（正常順序：viewStartDateから今日までの日数）
          const todayX = (new Date().getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH;
          // 今日の位置を画面中央に表示するためのスクロール位置を計算
          const containerWidth = scrollContainerRef.current.clientWidth;
          const scrollPosition = Math.max(0, todayX - containerWidth / 2);
          scrollContainerRef.current.scrollLeft = scrollPosition;
          setHasInitialScrolled(true);
        }
      }, 300);
    }
  }, [viewStartDate, viewEndDate, ganttItems, hasInitialScrolled]);

  // 表示期間の計算（全工事の最小～最大期間）- 正常順序
  useEffect(() => {
    if (kojies.length === 0) return;

    let minDate = new Date();
    let maxDate = new Date();
    let hasValidDate = false;

    kojies.forEach(koji => {
      try {
        const startDate = koji.startDate ? new Date(koji.startDate as string) : null;
        const endDate = koji.endDate ? new Date(koji.endDate as string) : null;

        if (startDate && !isNaN(startDate.getTime())) {
          if (!hasValidDate || startDate < minDate) {
            minDate = startDate;
          }
          hasValidDate = true;
        }

        if (endDate && !isNaN(endDate.getTime())) {
          if (!hasValidDate || endDate > maxDate) {
            maxDate = endDate;
          }
          hasValidDate = true;
        }
      } catch (error) {
        // 無効な日付は無視
      }
    });

    if (hasValidDate) {
      // 前後に1ヶ月の余裕を追加（正常順序：minDateを開始、maxDateを終了）
      const start = new Date(minDate.getFullYear(), minDate.getMonth() - 1, 1);
      const end = new Date(maxDate.getFullYear(), maxDate.getMonth() + 1, 0);
      setViewStartDate(start);
      setViewEndDate(end);
    } else {
      // 有効な日付がない場合はデフォルト期間（正常順序）
      const today = new Date();
      const start = new Date(today.getFullYear(), today.getMonth() - 6, 1);
      const end = new Date(today.getFullYear(), today.getMonth() + 6, 0);
      setViewStartDate(start);
      setViewEndDate(end);
    }
  }, [kojies]);

  // 表示範囲内で一番古い工事を基準にして、その開始日以降の工事を開始日昇順で10個表示
  const updateVisibleKojies = (scrollLeft: number = 0) => {
    if (kojies.length === 0 || !scrollContainerRef.current) return;

    // 現在の画面表示範囲を計算
    const containerWidth = scrollContainerRef.current.clientWidth;
    const visibleStartDays = scrollLeft / DAY_WIDTH;
    const visibleEndDays = (scrollLeft + containerWidth) / DAY_WIDTH;
    
    // 表示範囲の開始日と終了日を計算
    const visibleStartDate = new Date(viewStartDate.getTime() + visibleStartDays * 24 * 60 * 60 * 1000);
    const visibleEndDate = new Date(viewStartDate.getTime() + visibleEndDays * 24 * 60 * 60 * 1000);
    
    // 現在の画面表示範囲に含まれる工事をフィルタリング
    const relevantKojies = kojies.filter(koji => {
      try {
        const kojiStart = koji.startDate ? new Date(koji.startDate as string) : new Date();
        const kojiEnd = koji.endDate ? new Date(koji.endDate as string) : new Date(kojiStart.getTime() + 90 * 24 * 60 * 60 * 1000);
        
        // 工事が画面表示範囲と重複するかチェック
        return (kojiStart <= visibleEndDate && kojiEnd >= visibleStartDate);
      } catch {
        return false;
      }
    });

    // 表示範囲内の工事を開始日順でソートして、一番古い工事を取得
    const sortedRelevantKojies = relevantKojies.sort((a, b) => {
      const dateA = a.startDate ? new Date(a.startDate as string).getTime() : 0;
      const dateB = b.startDate ? new Date(b.startDate as string).getTime() : 0;
      return dateA - dateB; // 古い順
    });

    // 表示範囲内で一番古い工事の開始日を基準にする
    let baselineDate: number;
    if (sortedRelevantKojies.length > 0) {
      baselineDate = sortedRelevantKojies[0].startDate 
        ? new Date(sortedRelevantKojies[0].startDate as string).getTime()
        : 0;
    } else {
      // 表示範囲内に工事がない場合は、表示範囲の開始日以前で最も近い工事を基準にする
      const visibleStartTime = visibleStartDate.getTime();
      
      // 表示範囲の開始日以前の工事を取得
      const kojiesBeforeVisible = kojies.filter(koji => {
        const kojiStartDate = koji.startDate 
          ? new Date(koji.startDate as string).getTime()
          : 0;
        return kojiStartDate <= visibleStartTime;
      });
      
      if (kojiesBeforeVisible.length > 0) {
        // 表示範囲の開始日に最も近い工事を選択（開始日が最も新しい工事）
        const closestKoji = kojiesBeforeVisible.sort((a, b) => {
          const dateA = a.startDate ? new Date(a.startDate as string).getTime() : 0;
          const dateB = b.startDate ? new Date(b.startDate as string).getTime() : 0;
          return dateB - dateA; // 新しい順（降順）
        })[0];
        
        baselineDate = closestKoji.startDate 
          ? new Date(closestKoji.startDate as string).getTime()
          : 0;
      } else {
        // 表示範囲の開始日以前に工事がない場合は、全工事の最初の工事を基準にする
        const allKojiesSorted = [...kojies].sort((a, b) => {
          const dateA = a.startDate ? new Date(a.startDate as string).getTime() : 0;
          const dateB = b.startDate ? new Date(b.startDate as string).getTime() : 0;
          return dateA - dateB;
        });
        
        if (allKojiesSorted.length === 0) return;
        
        baselineDate = allKojiesSorted[0].startDate 
          ? new Date(allKojiesSorted[0].startDate as string).getTime()
          : 0;
      }
    }

    // 基準日以降の全工事を開始日昇順で取得
    const allKojiesSorted = [...kojies].sort((a, b) => {
      const dateA = a.startDate ? new Date(a.startDate as string).getTime() : 0;
      const dateB = b.startDate ? new Date(b.startDate as string).getTime() : 0;
      return dateA - dateB; // 古い順（昇順）
    });

    const kojiesFromBaselineDate = allKojiesSorted.filter(koji => {
      const kojiStartDate = koji.startDate 
        ? new Date(koji.startDate as string).getTime()
        : 0;
      return kojiStartDate >= baselineDate;
    });

    // 開始日昇順で表示件数分を取得
    let finalKojies = kojiesFromBaselineDate.slice(0, itemsPerPage);
    
    // もし表示件数が不足の場合は、開始日が最も新しいものから表示件数分を抽出
    if (finalKojies.length < itemsPerPage) {
      // 全工事を開始日の新しい順（降順）でソート
      const allKojiesDescending = [...kojies].sort((a, b) => {
        const dateA = a.startDate ? new Date(a.startDate as string).getTime() : 0;
        const dateB = b.startDate ? new Date(b.startDate as string).getTime() : 0;
        return dateB - dateA; // 新しい順（降順）
      });
      
      // 最新の表示件数分を取得して、開始日の古い順（昇順）に並び替え
      finalKojies = allKojiesDescending.slice(0, itemsPerPage).sort((a, b) => {
        const dateA = a.startDate ? new Date(a.startDate as string).getTime() : 0;
        const dateB = b.startDate ? new Date(b.startDate as string).getTime() : 0;
        return dateA - dateB; // 古い順（昇順）
      });
    }

    setVisibleKojies(finalKojies);
  };

  // スクロールイベントハンドラ
  const handleScroll = () => {
    if (scrollContainerRef.current) {
      updateVisibleKojies(scrollContainerRef.current.scrollLeft);
    }
  };

  // 工事データ変更時に初期表示を更新
  useEffect(() => {
    if (kojies.length > 0 && scrollContainerRef.current) {
      updateVisibleKojies(scrollContainerRef.current.scrollLeft);
    }
  }, [kojies, viewStartDate, itemsPerPage]);

  // 画面高さに基づいて表示件数を計算
  const calculateItemsPerPage = () => {
    if (!scrollContainerRef.current) return MIN_ITEMS;
    
    // ガントチャートエリアの高さを取得
    const containerHeight = scrollContainerRef.current.clientHeight;
    // ヘッダー分を除いた有効な高さ
    const availableHeight = containerHeight - 55; // ヘッダー高さ（月ヘッダー30px + 日付ヘッダー25px）
    // 行数を計算（最低5個、最大は画面に収まる範囲）
    const maxItems = Math.floor(availableHeight / ROW_HEIGHT);
    
    return Math.max(MIN_ITEMS, maxItems);
  };

  // ウィンドウリサイズ時に表示工事と表示件数を再計算
  useEffect(() => {
    const handleResize = () => {
      if (scrollContainerRef.current) {
        const newItemsPerPage = calculateItemsPerPage();
        setItemsPerPage(newItemsPerPage);
        updateVisibleKojies(scrollContainerRef.current.scrollLeft);
      }
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [kojies, viewStartDate]);

  // 初回レンダリング後に表示件数を計算
  useEffect(() => {
    if (scrollContainerRef.current && ganttItems.length > 0) {
      const newItemsPerPage = calculateItemsPerPage();
      setItemsPerPage(newItemsPerPage);
    }
  }, [ganttItems.length]);

  // ガントチャートアイテムの計算
  useEffect(() => {
    if (visibleKojies.length === 0) return;

    const items: GanttItem[] = visibleKojies.map((koji, index) => {
      // 安全な日付処理
      let startDate: Date;
      let endDate: Date;
      
      try {
        startDate = koji.startDate ? new Date(koji.startDate as string) : new Date();
        // 無効な日付をチェック
        if (isNaN(startDate.getTime())) {
          startDate = new Date();
        }
      } catch {
        startDate = new Date();
      }
      
      try {
        endDate = koji.endDate ? new Date(koji.endDate as string) : new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
        // 無効な日付をチェック
        if (isNaN(endDate.getTime())) {
          endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
        }
      } catch {
        endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
      }
      
      // 安全な計算（正常順序：左が古い、右が新しい）
      const daysDiff = (startDate.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24);
      const startX = Math.max(0, daysDiff * DAY_WIDTH);
      
      const endDaysDiff = (endDate.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24);
      const endX = endDaysDiff * DAY_WIDTH;
      const width = Math.max(DAY_WIDTH, endX - startX);

      return {
        ...koji,
        startX: isNaN(startX) ? 0 : startX,
        width: isNaN(width) ? DAY_WIDTH : width,
        row: index
      };
    });

    setGanttItems(items);
  }, [visibleKojies, viewStartDate]);

  // 工事クリック処理（詳細表示）- 将来使用予定
  // const handleKojiClick = (koji: ModelsKoji) => {
  //   setSelectedKoji(koji);
  //   setIsDetailModalOpen(true);
  // };

  // 工事編集処理
  const handleKojiEdit = (koji: ModelsKoji) => {
    setSelectedKoji(koji);
    setIsEditModalOpen(true);
  };

  // 工事名エリアクリック処理（中央に移動）- 正常順序
  const handleKojiNameClick = (koji: ModelsKoji) => {
    if (!scrollContainerRef.current) return;
    
    try {
      const kojiStart = koji.startDate ? new Date(koji.startDate as string) : new Date();
      const kojiEnd = koji.endDate ? new Date(koji.endDate as string) : new Date(kojiStart.getTime() + 90 * 24 * 60 * 60 * 1000);
      
      // 工事の中央位置を計算（正常順序）
      const kojiMiddle = new Date((kojiStart.getTime() + kojiEnd.getTime()) / 2);
      const middleX = (kojiMiddle.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH;
      
      // 画面中央に表示するためのスクロール位置を計算
      const containerWidth = scrollContainerRef.current.clientWidth;
      const scrollPosition = Math.max(0, middleX - containerWidth / 2);
      scrollContainerRef.current.scrollLeft = scrollPosition;
    } catch (error) {
      console.error('Error calculating koji center:', error);
    }
  };

  // 今日へ移動（正常順序）
  const scrollToToday = () => {
    if (!scrollContainerRef.current) return;
    
    const todayX = (new Date().getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH;
    const containerWidth = scrollContainerRef.current.clientWidth;
    const scrollPosition = Math.max(0, todayX - containerWidth / 2);
    scrollContainerRef.current.scrollLeft = scrollPosition;
  };

  // 工事更新処理（APIコール）
  const updateKoji = async (updatedKoji: ModelsKoji): Promise<ModelsKoji> => {
    try {
      const saved = await kojiConnectClient.update(updatedKoji);
      setKojies((prevKojies) =>
        prevKojies.map((k) => (getKojiKey(k) === getKojiKey(saved) ? saved : k))
      );
      return saved;
    } catch (err) {
      console.error("Error updating koji:", err);
      throw err;
    }
  };

  // 工事データを更新
  const handleKojiUpdate = (updatedKoji: ModelsKoji) => {
    // 選択中の工事を更新
    setSelectedKoji(updatedKoji);
    
    // ファイル名が変更された可能性がある場合、モーダルを閉じた後に再読み込みをする
    if (getKojiKey(selectedKoji) !== getKojiKey(updatedKoji)) {
      setShouldReloadOnClose(true);
    }

    // 工事一覧を更新
    setKojies((prevKojies) => {
      // 既存の工事を探す（pathで照合）
      const existingIndex = prevKojies.findIndex(k => getKojiKey(k) === getKojiKey(updatedKoji));
      
      if (existingIndex !== -1) {
        // 既存の工事を更新
        const updatedKojies = [...prevKojies];
        updatedKojies[existingIndex] = updatedKoji;
        return updatedKojies;
      } else {
        // pathが変わった場合（フォルダー名変更時など）
        // 選択中の工事のpathで元の工事を探す
        const oldKojiIndex = prevKojies.findIndex(k => 
          getKojiKey(k) === getKojiKey(selectedKoji)
        );
        
        if (oldKojiIndex !== -1) {
          // 古い工事を削除して新しいものを追加
          const updatedKojies = [...prevKojies];
          updatedKojies.splice(oldKojiIndex, 1);
          updatedKojies.push(updatedKoji);
          
          // 開始日順でソート（古い順）
          return updatedKojies.sort((a, b) => {
            const dateA = a.startDate ? new Date(typeof a.startDate === 'string' ? a.startDate : (a.startDate as any)['time.Time']).getTime() : 0;
            const dateB = b.startDate ? new Date(typeof b.startDate === 'string' ? b.startDate : (b.startDate as any)['time.Time']).getTime() : 0;
            
            // 開始日が設定されている方を優先
            if (dateA > 0 && dateB === 0) return -1;
            if (dateA === 0 && dateB > 0) return 1;
            
            // 両方開始日がある場合は古い順
            if (dateA > 0 && dateB > 0) return dateA - dateB;
            
            // 両方開始日がない場合はフォルダー名で昇順
            return getKojiKey(a).localeCompare(getKojiKey(b));
          });
        } else {
          // 新規追加
          return [...prevKojies, updatedKoji];
        }
      }
    });
  };

  // ステータスによる色の取得
  const getStatusColor = (status?: string) => {
    switch (status) {
      case '進行中':
        return '#4CAF50';
      case '完了':
        return '#9E9E9E';
      case '予定':
        return '#FF9800';
      default:
        return '#2196F3';
    }
  };

  // 管理ファイルの変更が必要かチェック
  const needsFileRename = (koji: ModelsKoji): boolean => {
    if (!koji.managed_files || koji.managed_files.length === 0) {
      return false;
    }
    
    // managed_filesの中で現在の名前と推奨名が異なるものがあるかチェック
    const needsRename = koji.managed_files.some(file => {
      // currentとrecommendedが両方存在し、異なる場合にtrueを返す
      return file.current && file.recommended && file.current !== file.recommended;
    });
    
    return needsRename;
  };

  // 月のヘッダーを生成（正常順序・1日基準）
  const generateMonthHeaders = () => {
    const headers = [];
    const current = new Date(viewStartDate);
    
    while (current <= viewEndDate) {
      const monthStart = new Date(current.getFullYear(), current.getMonth(), 1);
      const monthEnd = new Date(current.getFullYear(), current.getMonth() + 1, 0);
      const daysInMonth = monthEnd.getDate();
      
      headers.push({
        year: monthStart.getFullYear(),
        month: monthStart.getMonth() + 1,
        width: daysInMonth * DAY_WIDTH,
        startX: (monthStart.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH
      });
      
      current.setMonth(current.getMonth() + 1);
    }
    
    return headers;
  };

  // 3日毎の日付ヘッダーを生成（正常順序）
  const generateDayHeaders = () => {
    const headers = [];
    const current = new Date(viewStartDate);
    
    while (current <= viewEndDate) {
      const day = current.getDate();
      if ((day === 1 || day % 3 === 1) && day !== 31) { // 1日、4日、7日、10日...（31日は除外）
        headers.push({
          date: day,
          month: current.getMonth() + 1,
          year: current.getFullYear(),
          startX: (current.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH,
          width: DAY_WIDTH * 3 // 3日分の幅
        });
      }
      current.setDate(current.getDate() + 1);
    }
    
    return headers;
  };

  // 月境界線を生成（月の1日の位置）
  const generateMonthBoundaries = () => {
    const boundaries = [];
    const current = new Date(viewStartDate.getFullYear(), viewStartDate.getMonth() + 1, 1); // 次月の1日から開始
    
    while (current <= viewEndDate) {
      const startX = (current.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH;
      boundaries.push({
        startX,
        year: current.getFullYear(),
        month: current.getMonth() + 1
      });
      current.setMonth(current.getMonth() + 1);
    }
    
    return boundaries;
  };

  // 日付フォーマット - 将来使用予定
  // const formatDate = (dateString?: string | any) => {
  //   if (!dateString) return '';
  //   try {
  //     return new Date(dateString as string).toLocaleDateString('ja-JP');
  //   } catch {
  //     return '無効な日付';
  //   }
  // };


  if (loading) {
    return <div className="loading">工事データを読み込み中...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  const monthHeaders = generateMonthHeaders();
  const dayHeaders = generateDayHeaders();
  const monthBoundaries = generateMonthBoundaries();
  const totalWidth = (viewEndDate.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH;

  return (
    <div className="gantt-container">
      <div style={{ 
        marginBottom: "20px",
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center"
      }}>
        <button 
          onClick={scrollToToday}
          className="gantt-today-button"
        >
          📅 今日へ移動
        </button>
        
        <div style={{ 
          fontSize: "16px",
          color: "#666",
          fontWeight: "500"
        }}>
          表示中: {ganttItems.length}件 / 全{kojies.length}件
        </div>
      </div>

      <div className="gantt-wrapper">
        <div className="gantt-sidebar">
          <div className="gantt-header-left">会社名</div>
          {ganttItems.map((item, index) => {
            // 今日の日付が工事期間に含まれるかチェック
            const today = new Date();
            let startDate: Date;
            let endDate: Date;
            
            try {
              startDate = item.startDate ? new Date(item.startDate as string) : new Date();
              if (isNaN(startDate.getTime())) startDate = new Date();
            } catch {
              startDate = new Date();
            }
            
            try {
              endDate = item.endDate ? new Date(item.endDate as string) : new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
              if (isNaN(endDate.getTime())) endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
            } catch {
              endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
            }
            
            const isActiveKoji = today >= startDate && today <= endDate;
            
            return (
              <div 
                key={`${item.id}-${index}`} 
                className={`gantt-row-label ${isActiveKoji ? 'gantt-row-label-active' : ''}`}
                style={{ 
                  height: ROW_HEIGHT
                }}
              >
                <div 
                  className={`koji-name koji-name-clickable ${isActiveKoji ? 'koji-name-active' : ''}`}
                  onClick={() => handleKojiNameClick(item)}
                  title="クリックして工事期間の中央に移動"
                >
                  {item.companyName}
                </div>
              </div>
            );
          })}
        </div>

        <div className="gantt-chart-container" ref={scrollContainerRef} onScroll={handleScroll} style={{ backgroundColor: '#f5f5f5' }}>
          <div className="gantt-chart" style={{ width: totalWidth, backgroundColor: '#f5f5f5' }}>
            {/* 月ヘッダー */}
            <div className="gantt-header month-header-row" style={{ height: "30px" }}>
              {monthHeaders.map((header, index) => (
                <div 
                  key={index}
                  className="month-header month-header-content"
                  style={{ 
                    left: header.startX, 
                    width: header.width
                  }}
                >
                  {header.year}年{header.month}月
                </div>
              ))}
            </div>
            
            {/* 日付ヘッダー */}
            <div className="gantt-header day-header-row" style={{ height: "25px", borderBottom: "1px solid #333" }}>
              {dayHeaders.map((header, index) => (
                <div 
                  key={index}
                  className="day-header day-header-content"
                  style={{ 
                    left: header.startX, 
                    width: header.width
                  }}
                >
                  {header.date}
                </div>
              ))}
            </div>

            {/* 月境界線（太線） - 月ヘッダーから開始 */}
            <div className="gantt-month-boundaries" style={{ top: 0, height: '100%' }}>
              {monthBoundaries.map((boundary, index) => (
                <div 
                  key={`month-boundary-${index}`}
                  className="month-boundary-line"
                  style={{ 
                    left: boundary.startX,
                    top: 0,
                    height: Math.max(400, itemsPerPage * ROW_HEIGHT + 55) // ヘッダー分も含む
                  }}
                  title={`${boundary.year}年${boundary.month}月開始`}
                />
              ))}
            </div>

            <div className="gantt-body">
              {/* グリッド線（3日毎） */}
              <div className="gantt-grid">
                {dayHeaders.map((header, index) => (
                  <div 
                    key={index}
                    className="grid-line"
                    style={{ 
                      left: header.startX,
                      height: Math.max(400, itemsPerPage * ROW_HEIGHT)
                    }}
                  />
                ))}
              </div>

              {/* 水平線（工事行ごと） */}
              <div className="gantt-horizontal-grid">
                {ganttItems.map((_, index) => (
                  <div
                    key={`horizontal-${index}`}
                    className="horizontal-grid-line"
                    style={{ 
                      top: (index + 1) * ROW_HEIGHT,
                      width: '100%'
                    }}
                  />
                ))}
              </div>

              {/* 工事期間バー */}
              {ganttItems.map((item, index) => (
                <div 
                  key={`${item.id}-${index}`}
                  className="gantt-bar"
                  style={{
                    left: item.startX,
                    width: item.width,
                    top: index * ROW_HEIGHT + 10,
                    height: ROW_HEIGHT - 15,
                    backgroundColor: getStatusColor(item.status)
                  }}
                  onClick={() => handleKojiEdit(item)}
                  title={`${item.companyName} - ${item.locationName} (クリックして編集)`}
                >
                  <span className="gantt-bar-text">
                    {item.locationName}
                  </span>
                  {needsFileRename(item) && (
                    <span 
                      className="gantt-bar-rename-indicator"
                      title="管理ファイルの名前変更が必要です"
                    >
                      ⚠️
                    </span>
                  )}
                </div>
              ))}

              {/* 今日の範囲 */}
              <div 
                className="today-area"
                style={{
                  left: Math.floor((new Date().setHours(0, 0, 0, 0) - viewStartDate.getTime()) / (1000 * 60 * 60 * 24)) * DAY_WIDTH,
                  width: DAY_WIDTH,
                  height: Math.max(400, itemsPerPage * ROW_HEIGHT), // 動的な高さ
                  backgroundColor: 'rgba(255, 192, 203, 0.3)', // 薄いピンク
                  position: 'absolute',
                  top: 0,
                  pointerEvents: 'none',
                  zIndex: 1
                }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* 編集モーダル */}
      <KojiDetailModal
        isOpen={isEditModalOpen}
        onClose={() => {
          setIsEditModalOpen(false);
          // フォルダ名が変更された場合のみ工事一覧を再読み込み
          if (shouldReloadOnClose) {
            loadKojies();
            setShouldReloadOnClose(false);
          }
        }}
        koji={selectedKoji}
        onUpdate={updateKoji}
        onKojiUpdate={handleKojiUpdate}
      />
    </div>
  );
};

export default KojiGanttChart;
