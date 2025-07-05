import { useState, useEffect } from "react";
import type { ModelsCompany } from "../api/types.gen";
import { getCompanyList } from "../api/sdk.gen";
import CompanyDetailModal from "./CompanyDetailModal";
import "../styles/business-entity-list.css";

const Companies = () => {
  const [companies, setCompanies] = useState<ModelsCompany[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showHelp, setShowHelp] = useState(false);
  const [selectedCompany, setSelectedCompany] = useState<ModelsCompany | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);

  // 会社データを読み込み
  const loadCompanies = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await getCompanyList();

      if (response.data) {
        setCompanies(response.data);
      } else {
        setCompanies([]);
      }
    } catch (err) {
      console.error("Error loading company entries:", err);
      setError(
        `会社データの読み込みに失敗しました: ${
          err instanceof Error ? err.message : String(err)
        }`
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadCompanies();
  }, []);

  // 会社クリック処理
  const handleCompanyClick = (company: ModelsCompany) => {
    setSelectedCompany(company);
    setIsDetailModalOpen(true);
  };

  // モーダルを閉じる
  const closeDetailModal = () => {
    setIsDetailModalOpen(false);
    setSelectedCompany(null);
  };

  // 会社情報更新処理（将来の実装用）
  const updateCompany = async (updatedCompany: ModelsCompany): Promise<ModelsCompany> => {
    // 未実装
    console.log('会社情報更新:', updatedCompany);
    return updatedCompany;
  };

  // 会社データを更新
  const handleCompanyUpdate = (updatedCompany: ModelsCompany) => {
    setSelectedCompany(updatedCompany);
    setCompanies(prevCompanies => 
      prevCompanies.map(c => c.id === updatedCompany.id ? updatedCompany : c)
    );
  };

  if (loading) {
    return (
      <div className="business-entity-loading">
        <div>会社データを読み込み中...</div>
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
          onClick={loadCompanies}
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
        <div className="business-entity-count">
          全{companies.length}件
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
          <p>📝 <strong>会社情報表示</strong></p>
          <p>🏢 <strong>業種別分類</strong></p>
          <p>📞 <strong>連絡先情報</strong></p>
          <p>✅ 会社データの取得</p>
          <p>✅ フォルダー名からの自動解析</p>
        </div>
      )}
      
      <div className="business-entity-list-container">
        <div className="business-entity-list-header">
          <div className="business-entity-header-company">会社名</div>
          <div className="business-entity-header-company">業種</div>
          <div className="business-entity-header-spacer"></div>
          <button
            onClick={() => setShowHelp(!showHelp)}
            className="business-entity-help-button"
            title="使用方法を表示"
          >
            ?
          </button>
        </div>

        <div className="business-entity-scroll-area">
          {companies.length === 0 ? (
            <div className="business-entity-empty">
              会社データが見つかりません
            </div>
          ) : (
            <div>
              {companies.map((company, index) => (
              <div
                key={company.id || index}
                className="business-entity-item-row"
                onClick={() => handleCompanyClick(company)}
                title="クリックして詳細表示"
                style={{ cursor: 'pointer' }}
              >
                <div className="business-entity-item-info">
                  <div className="business-entity-item-info-company">
                    {company.short_name || company.name || "会社名未設定"}
                  </div>
                  
                  <div className="business-entity-item-info-company">
                    {company.business_type || "業種未設定"}
                  </div>
                  
                  <div className="business-entity-item-info-spacer"></div>
                </div>
                
                <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                  {company.tags && company.tags.length > 0 && (
                    <div style={{ display: 'flex', gap: '4px', flexWrap: 'wrap' }}>
                      {company.tags.slice(0, 3).map((tag, tagIndex) => (
                        <span 
                          key={tagIndex} 
                          style={{
                            fontSize: '0.7em',
                            background: '#e3f2fd',
                            color: '#1976d2',
                            padding: '2px 6px',
                            borderRadius: '10px',
                            whiteSpace: 'nowrap'
                          }}
                        >
                          {tag}
                        </span>
                      ))}
                      {company.tags.length > 3 && (
                        <span style={{ fontSize: '0.7em', color: '#666' }}>
                          +{company.tags.length - 3}
                        </span>
                      )}
                    </div>
                  )}
                  
                  <div className="business-entity-item-status business-entity-item-status-completed">
                    {company.website ? "🌐" : "📁"}
                  </div>
                </div>
              </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* 詳細モーダル */}
      <CompanyDetailModal
        isOpen={isDetailModalOpen}
        onClose={closeDetailModal}
        company={selectedCompany}
        onUpdate={updateCompany}
        onCompanyUpdate={handleCompanyUpdate}
      />
    </div>
  );
};

export { Companies };