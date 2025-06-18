import { useState, useEffect } from 'react';
import KoujiEntriesGrid from './KoujiProjectGrid';
import { getKoujiEntries } from '../api';
import type { KoujiEntryExtended } from '../types/kouji';

const KoujiProjectPage = () => {
  const [showHelp, setShowHelp] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [saveMessage, setSaveMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [koujiEntries, setKoujiEntries] = useState<KoujiEntryExtended[]>([]);
  const [hasChanges, setHasChanges] = useState(false);

  // 工事データをロード
  const loadKoujiEntries = async () => {
    try {
      const response = await getKoujiEntries({});
      const responseData = response.data as any;
      setKoujiEntries(responseData.kouji_entries || []);
    } catch (err) {
      console.error('Error loading kouji entries:', err);
      setError('工事データの読み込みに失敗しました');
    }
  };

  // 初期ロード
  useEffect(() => {
    loadKoujiEntries();
  }, []);

  // 工事データを更新
  const updateKoujiEntry = (updatedEntry: KoujiEntryExtended) => {
    setKoujiEntries(prev => 
      prev.map(entry => 
        entry.id === updatedEntry.id ? updatedEntry : entry
      )
    );
    setHasChanges(true);
  };

  const handleSaveKoujiEntries = async () => {
    setIsSaving(true);
    setSaveMessage(null);
    setError(null);
    
    try {
      // 編集された工事データをサーバーに送信
      const response = await fetch('http://localhost:8080/api/kouji-entries/save', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(koujiEntries),
      });
      
      if (!response.ok) {
        throw new Error('保存に失敗しました');
      }
      
      const result = await response.json();
      setSaveMessage(result.message);
      setHasChanges(false);
    } catch (err) {
      console.error('Error saving kouji entries:', err);
      setError(err instanceof Error ? err.message : '保存に失敗しました');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <h1>工事管理</h1>
        <div className="header-actions">
          <button 
            type="button" 
            onClick={handleSaveKoujiEntries}
            disabled={isSaving}
            className="save-button"
          >
            {isSaving ? '保存中...' : hasChanges ? '工事データ保存 *' : '工事データ保存'}
          </button>
          <button 
            className="help-button"
            onClick={() => setShowHelp(!showHelp)}
          >
            ヘルプ
          </button>
        </div>
      </div>

      {showHelp && (
        <div className="help-section">
          <h3>工事プロジェクトについて</h3>
          <ul>
            <li>フォルダー名は「YYYY-MMDD 会社名 現場名」の形式で命名してください</li>
            <li>例: 「2025-0618 豊田築炉 名和工場」</li>
            <li>ステータスは日付に基づいて自動的に判定されます：
              <ul>
                <li><span style={{ color: '#FF9800' }}>予定</span>: 開始日が未来の場合</li>
                <li><span style={{ color: '#4CAF50' }}>進行中</span>: 現在進行中の場合</li>
                <li><span style={{ color: '#9E9E9E' }}>完了</span>: 終了日を過ぎた場合</li>
              </ul>
            </li>
            <li>プロジェクトの期間は開始日から3ヶ月間と仮定されます</li>
          </ul>
        </div>
      )}

      {error && <div className="error">{error}</div>}
      {saveMessage && <div className="success">{saveMessage}</div>}

      <KoujiEntriesGrid 
        koujiEntries={koujiEntries}
        onUpdateEntry={updateKoujiEntry}
        onReload={loadKoujiEntries}
      />
    </div>
  );
};

export default KoujiProjectPage;