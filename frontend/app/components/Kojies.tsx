import { useState, useEffect } from "react";
import { Link } from "react-router";
import { getKojiRecent } from "../api/sdk.gen";
import type { ModelsKoji } from "../api/types.gen";
import KojiDetailModal from "./KojiDetailModal";
import { useKoji } from "../contexts/KojiContext";
import "../styles/business-entity-list.css";

const Kojies = () => {
  const [kojies, setKojies] = useState<ModelsKoji[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedKoji, setSelectedKoji] = useState<ModelsKoji | null>(
    null
  );
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [showHelp, setShowHelp] = useState(false);
  const { setKojiCount } = useKoji();

  // 工事データを読み込み
  const loadKojies = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await getKojiRecent();

      if (response.data) {
        setKojies(response.data);
        setKojiCount(response.data.length);
      } else {
        setKojies([]);
        setKojiCount(0);
      }
    } catch (err) {
      console.error("Error loading kouji entries:", err);
      setError(
        `工事データの読み込みに失敗しました: ${
          err instanceof Error ? err.message : String(err)
        }`
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadKojies();
  }, []);

  // 工事クリック処理
  const handleKojiClick = (koji: ModelsKoji) => {
    setSelectedKoji(koji);
    setIsEditModalOpen(true);
  };

  // 工事更新処理（APIコール）
  const updateKoji = async (updatedKoji: ModelsKoji): Promise<ModelsKoji> => {
    try {
      // バックエンドに更新リクエストを送信
      const response = await fetch("http://localhost:8080/api/kojies", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedKoji),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "更新に失敗しました");
      }

      // レスポンスから更新された工事データを取得
      const savedKoji = await response.json();

      // 工事一覧を更新
      setKojies((prevKojies) =>
        prevKojies.map((k) => (k.id === savedKoji.id ? savedKoji : k))
      );

      // 更新された工事データを返す
      return savedKoji;
    } catch (err) {
      console.error("Error updating koji:", err);
      throw err; // エラーをモーダルに伝播
    }
  };

  // モーダルを閉じる
  const closeEditModal = () => {
    setIsEditModalOpen(false);
    setSelectedKoji(null);
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

  // 工事データを更新
  const handleKojiUpdate = (updatedKoji: ModelsKoji) => {
    // 選択中の工事を更新
    setSelectedKoji(updatedKoji);

    // 工事一覧を更新
    setKojies((prevKojies) => {
      // 既存の工事を探す（IDで照合）
      const existingIndex = prevKojies.findIndex(k => k.id === updatedKoji.id);
      
      if (existingIndex !== -1) {
        // 既存の工事を更新
        const updatedKojies = [...prevKojies];
        updatedKojies[existingIndex] = updatedKoji;
        return updatedKojies;
      } else {
        // フォルダー名が変更された可能性があるため、元の工事を探して削除し、新しいものを追加
        // 同じ会社名・現場名で探す
        const oldKojiIndex = prevKojies.findIndex(k => 
          k.company_name === updatedKoji.company_name && 
          k.location_name === updatedKoji.location_name &&
          k.id !== updatedKoji.id
        );
        
        if (oldKojiIndex !== -1) {
          // 古い工事を削除して新しいものを追加
          const updatedKojies = [...prevKojies];
          updatedKojies.splice(oldKojiIndex, 1);
          updatedKojies.push(updatedKoji);
          // 開始日順でソート（新しい順）
          return updatedKojies.sort((a, b) => {
            const dateA = a.start_date ? new Date(typeof a.start_date === 'string' ? a.start_date : (a.start_date as any)['time.Time']).getTime() : 0;
            const dateB = b.start_date ? new Date(typeof b.start_date === 'string' ? b.start_date : (b.start_date as any)['time.Time']).getTime() : 0;
            
            // 開始日が設定されている方を優先
            if (dateA > 0 && dateB === 0) return -1;
            if (dateA === 0 && dateB > 0) return 1;
            
            // 両方開始日がある場合は新しい順
            if (dateA > 0 && dateB > 0) return dateB - dateA;
            
            // 両方開始日がない場合はフォルダー名で降順
            return (b.name || '').localeCompare(a.name || '');
          });
        } else {
          // 新規追加
          return [...prevKojies, updatedKoji];
        }
      }
    });
  };

  if (loading) {
    return (
      <div className="business-entity-loading">
        <div>工事データを読み込み中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="business-entity-error">
        <div className="business-entity-error-message">
          {error}
        </div>
        <button
          onClick={loadKojies}
          className="business-entity-retry-button"
        >
          再試行
        </button>
      </div>
    );
  }

  return (
    <div className="business-entity-container">
      <div className="business-entity-controls">
        <Link
          to="/kojies/gantt"
          className="business-entity-gantt-button"
        >
          📊 工程表を表示
        </Link>
        
        <div className="business-entity-count">
          全{kojies.length}件
        </div>
      </div>

      {showHelp && (
        <div className="business-entity-help-box">
          <button
            onClick={() => setShowHelp(false)}
            className="business-entity-help-close"
            title="閉じる"
          >
            ×
          </button>
          
          <h3 style={{ marginTop: 0 }}>使用方法</h3>
          <p>
            📝 <strong>リストをクリック</strong>して工事情報を編集できます
          </p>
          <p>✅ 開始日・終了日・説明・タグ・会社名・現場名を編雈可能</p>
          <p>💾 編集後は自動で保存されます</p>

          <h3 style={{ marginTop: "15px" }}>開発状況</h3>
          <p>✅ 工事データの取得</p>
          <p>✅ 編集モーダル機能</p>
          <p>🔄 工程表機能（次のステップ）</p>
        </div>
      )}
      
      <div className="business-entity-list-container">
        <div className="business-entity-list-header">
          <div className="business-entity-header-date">開始日</div>
          <div className="business-entity-header-company">会社名</div>
          <div className="business-entity-header-location">現場名</div>
          <div className="business-entity-header-spacer"></div>
          <div className="business-entity-header-date" style={{ marginRight: "24px" }}>終了日</div>
          <div className="business-entity-header-status">ステータス</div>
          <button
            onClick={() => setShowHelp(!showHelp)}
            className="business-entity-help-button"
            title="使用方法を表示"
          >
            ?
          </button>
        </div>

        <div className="business-entity-scroll-area">
          {kojies.length === 0 ? (
            <div className="business-entity-empty">
              工事データが見つかりません
            </div>
          ) : (
            <div>
              {kojies.map((koji, index) => (
              <div
                key={koji.id || index}
                className="business-entity-item-row"
                onClick={() => handleKojiClick(koji)}
                title="クリックして編集"
              >
                <div className="business-entity-item-info">
                  <div className="business-entity-item-info-date">
                    {koji.start_date
                      ? new Date(
                          typeof koji.start_date === 'string' 
                            ? koji.start_date 
                            : (koji.start_date as any)['time.Time']
                        ).toLocaleDateString("ja-JP")
                      : "未設定"}
                  </div>
                  
                  <div className="business-entity-item-info-company">
                    {koji.company_name || "会社名未設定"}
                  </div>
                  
                  <div className="business-entity-item-info-location">
                    {koji.location_name || "現場名未設定"}
                  </div>
                  
                  <div className="business-entity-item-info-description">
                    {koji.description || ""}
                  </div>
                  
                  <div className="business-entity-item-info-date end-date" style={{ 
                    marginRight: "24px"
                  }}>
                    ～{koji.end_date
                      ? new Date(
                          typeof koji.end_date === 'string' 
                            ? koji.end_date 
                            : (koji.end_date as any)['time.Time']
                        ).toLocaleDateString("ja-JP")
                      : "未設定"}
                  </div>
                </div>
                <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                  {needsFileRename(koji) && (
                    <span 
                      className="business-entity-item-rename-indicator"
                      title="管理ファイルの名前変更が必要です"
                    >
                      ⚠️
                    </span>
                  )}
                  <div
                    className={`business-entity-item-status ${
                      koji.status === "進行中" ? "business-entity-item-status-ongoing" :
                      koji.status === "完了" ? "business-entity-item-status-completed" :
                      koji.status === "予定" ? "business-entity-item-status-planned" :
                      ""
                    }`}
                  >
                    {koji.status || "未設定"}
                  </div>
                </div>
              </div>
              ))}
            </div>
          )}
        </div>
      </div>


      {/* 編集モーダル */}
      <KojiDetailModal
        isOpen={isEditModalOpen}
        onClose={closeEditModal}
        koji={selectedKoji}
        onUpdate={updateKoji}
        onKojiUpdate={handleKojiUpdate}
      />
    </div>
  );
};

export default Kojies;