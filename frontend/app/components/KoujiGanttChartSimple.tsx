import { useState, useEffect } from 'react';
import { getKoujiEntries } from '../api/sdk.gen';
import type { KoujiEntryExtended } from '../types/kouji';
import KoujiEditModal from './KoujiEditModal';
import '../styles/modal.css';

const KoujiGanttChartSimple = () => {
  const [koujiEntries, setKoujiEntries] = useState<KoujiEntryExtended[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedProject, setSelectedProject] = useState<KoujiEntryExtended | null>(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  // 工事データを読み込み
  const loadKoujiEntries = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await getKoujiEntries({
        query: {}
      });
      
      const responseData = response.data as any;
      const entries = responseData.kouji_entries || [];
      setKoujiEntries(entries);
    } catch (err) {
      console.error('Error loading kouji entries:', err);
      setError(`工事データの読み込みに失敗しました: ${err instanceof Error ? err.message : String(err)}`);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadKoujiEntries();
  }, []);

  // プロジェクトクリック処理
  const handleProjectClick = (project: KoujiEntryExtended) => {
    setSelectedProject(project);
    setIsEditModalOpen(true);
  };

  // プロジェクト保存処理
  const handleProjectSave = async (updatedProject: KoujiEntryExtended) => {
    try {
      // バックエンドに保存リクエストを送信
      const response = await fetch('http://localhost:8080/api/kouji/save', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify([updatedProject]),
      });

      if (!response.ok) {
        throw new Error('保存に失敗しました');
      }

      console.log('プロジェクトが正常に保存されました');
      
      // 保存成功後、最新のデータを再取得してステータスを更新
      await loadKoujiEntries();
    } catch (err) {
      console.error('Error saving project:', err);
      alert('保存に失敗しました。もう一度お試しください。');
    }
  };

  // モーダルを閉じる
  const closeEditModal = () => {
    setIsEditModalOpen(false);
    setSelectedProject(null);
  };

  if (loading) {
    return (
      <div style={{ padding: '20px' }}>
        <h1>工程表（ガントチャート）</h1>
        <div>工事データを読み込み中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: '20px' }}>
        <h1>工程表（ガントチャート）</h1>
        <div style={{ color: 'red', padding: '10px', backgroundColor: '#ffe6e6', borderRadius: '4px' }}>
          {error}
        </div>
        <button onClick={loadKoujiEntries} style={{ marginTop: '10px', padding: '10px 20px' }}>
          再試行
        </button>
      </div>
    );
  }

  return (
    <div style={{ padding: '20px' }}>
      <h1>工程表（ガントチャート）</h1>
      
      <div style={{ marginBottom: '20px' }}>
        <p>取得した工事データ: {koujiEntries.length}件</p>
      </div>

      <div style={{ border: '1px solid #ddd', borderRadius: '8px', overflow: 'hidden' }}>
        <div style={{ 
          backgroundColor: '#f5f5f5', 
          padding: '10px', 
          fontWeight: 'bold',
          borderBottom: '1px solid #ddd'
        }}>
          工事プロジェクト一覧
        </div>
        
        {koujiEntries.length === 0 ? (
          <div style={{ padding: '20px', textAlign: 'center', color: '#666' }}>
            工事データが見つかりません
          </div>
        ) : (
          <div>
            {koujiEntries.map((entry, index) => (
              <div 
                key={entry.id || index}
                style={{ 
                  padding: '15px',
                  borderBottom: index < koujiEntries.length - 1 ? '1px solid #eee' : 'none',
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  cursor: 'pointer',
                  transition: 'background-color 0.3s'
                }}
                onClick={() => handleProjectClick(entry)}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#f8f9fa'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                title="クリックして編集"
              >
                <div>
                  <div style={{ fontWeight: 'bold' }}>
                    {entry.company_name || '会社名未設定'} - {entry.location_name || '現場名未設定'}
                  </div>
                  <div style={{ fontSize: '14px', color: '#666', marginTop: '5px' }}>
                    開始: {entry.start_date ? new Date(entry.start_date).toLocaleDateString('ja-JP') : '未設定'} | 
                    終了: {entry.end_date ? new Date(entry.end_date).toLocaleDateString('ja-JP') : '未設定'}
                  </div>
                </div>
                <div style={{
                  padding: '4px 12px',
                  borderRadius: '4px',
                  backgroundColor: entry.status === '進行中' ? '#4CAF50' : 
                                   entry.status === '完了' ? '#9E9E9E' :
                                   entry.status === '予定' ? '#FF9800' : '#2196F3',
                  color: 'white',
                  fontSize: '12px'
                }}>
                  {entry.status || '未設定'}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div style={{ marginTop: '20px', padding: '15px', backgroundColor: '#f0f8ff', borderRadius: '4px' }}>
        <h3>使用方法</h3>
        <p>📝 <strong>リストをクリック</strong>して工事情報を編集できます</p>
        <p>✅ 開始日・終了日・説明・タグ・会社名・現場名を編集可能</p>
        <p>💾 編集後は自動で保存されます</p>
        
        <h3 style={{ marginTop: '15px' }}>開発状況</h3>
        <p>✅ 工事データの取得</p>
        <p>✅ 編集モーダル機能</p>
        <p>🔄 ガントチャート機能（次のステップ）</p>
      </div>

      {/* 編集モーダル */}
      <KoujiEditModal
        isOpen={isEditModalOpen}
        onClose={closeEditModal}
        project={selectedProject}
        onSave={handleProjectSave}
      />
    </div>
  );
};

export default KoujiGanttChartSimple;