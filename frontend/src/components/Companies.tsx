'use client';

import { useState, useEffect, useRef } from "react";
import type { Company } from "@/api/types.gen";
import { getCompanies } from "@/api/sdk.gen";
import CompanyDetailModal from "@/components/CompanyDetailModal";
import {
  getCategoryName,
  initializeCategories,
  loadCategories,
} from "@/utils/company";
import "@/styles/business-entity-list.css";

const Companies = () => {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showHelp, setShowHelp] = useState(false);
  const [selectedCompany, setSelectedCompany] = useState<Company | null>(
    null
  );
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [selectedCategories, setSelectedCategories] = useState<Set<string>>(
    new Set()
  );
  const [showFilterDropdown, setShowFilterDropdown] = useState(false);
  const filterDropdownRef = useRef<HTMLDivElement>(null);
  const [allCategories, setAllCategories] = useState<
    Array<{ value: string; label: string }>
  >([]);

  // 会社データを読み込み
  const loadCompanies = async () => {
    try {
      // ローディング中のフラグをセット
      setLoading(true);
      setError(null);

      // 会社データを取得
      const response = await getCompanies();

      // 取得したデータをセット
      if (response.data) {
        const normalizedCompanies = response.data.map((company) => ({
          ...company,
          category:
            company.category !== undefined &&
            company.category !== null &&
            String(company.category).trim() !== ""
              ? String(company.category)
              : "",
        }));
        setCompanies(normalizedCompanies);
      } else {
        // データがない場合は空配列をセット
        setCompanies([]);
      }
    } catch (err) {
      setError(
        `会社データの読み込みに失敗しました: ${
          err instanceof Error ? err.message : String(err)
        }`
      );
    } finally {
      // ローディング中のフラグを解除
      setLoading(false);
    }
  };

  // コンポーネントがマウントされた時に実行
  useEffect(() => {
    const initialize = async () => {
      // カテゴリーマッピングを初期化
      await initializeCategories();
      // カテゴリーデータを取得
      const categories = await loadCategories();
      setAllCategories(
        categories
          .filter(
            (cat) =>
              cat.code !== undefined &&
              cat.code !== null &&
              String(cat.code).trim() !== ""
          )
          .map((cat) => ({
            value: String(cat.code),
            label: cat.label ?? "業種未設定",
          }))
      );
      // 会社データを読み込み
      await loadCompanies();
    };
    initialize();
  }, []);

  // ドロップダウンの外側をクリックしたときに閉じる
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        filterDropdownRef.current &&
        !filterDropdownRef.current.contains(event.target as Node)
      ) {
        setShowFilterDropdown(false);
      }
    };

    if (showFilterDropdown) {
      document.addEventListener("mousedown", handleClickOutside);
      return () => {
        document.removeEventListener("mousedown", handleClickOutside);
      };
    }
  }, [showFilterDropdown]);

  // フィルタリングされた会社リスト
  const filteredCompanies = companies.filter((company) => {
    if (selectedCategories.size === 0) return true; // 何も選択されていない場合は全て表示
    const categoryKey =
      company.category !== undefined &&
      company.category !== null &&
      String(company.category).trim() !== ""
        ? String(company.category)
        : null;
    return (
      categoryKey !== null && selectedCategories.has(categoryKey)
    );
  });

  // カテゴリーの選択/解除
  const toggleCategory = (category: string) => {
    const newSelected = new Set(selectedCategories);
    if (newSelected.has(category)) {
      newSelected.delete(category);
    } else {
      newSelected.add(category);
    }
    setSelectedCategories(newSelected);
  };

  // 全選択/全解除
  const toggleAllCategories = () => {
    if (
      allCategories.length > 0 &&
      selectedCategories.size === allCategories.length
    ) {
      setSelectedCategories(new Set()); // 全解除
    } else {
      setSelectedCategories(new Set(allCategories.map((c) => c.value))); // 全選択
    }
  };

  // 会社クリック処理
  const handleCompanyClick = (company: Company) => {
    setSelectedCompany(company);
    setIsDetailModalOpen(true);
  };

  // モーダルを閉じる
  const closeDetailModal = async () => {
    setIsDetailModalOpen(false);
    setSelectedCompany(null);
    
    // フォルダー名が変更された可能性があるため、会社一覧を再読み込み
    await loadCompanies();
  };

  // 会社データを更新
  const handleCompanyUpdate = (updatedCompany: Company) => {
    const normalized = {
      ...updatedCompany,
      category:
        updatedCompany.category !== undefined &&
        updatedCompany.category !== null &&
        String(updatedCompany.category).trim() !== ""
          ? String(updatedCompany.category)
          : "",
    };
    setSelectedCompany(normalized);
    setCompanies((prevCompanies) =>
      prevCompanies.map((c) =>
        c.id === normalized.id ? normalized : c
      )
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
        <div className="business-entity-error-message">{error}</div>
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
          {selectedCategories.size > 0 &&
            ` / 絞り込み: ${filteredCompanies.length}件`}
        </div>

        <div
          ref={filterDropdownRef}
          style={{
            display: "flex",
            gap: "16px",
            alignItems: "center",
            position: "relative",
          }}
        >
          <label style={{ fontSize: "14px", fontWeight: 500 }}>
            業種で絞り込み:
          </label>
          <button
            onClick={() => setShowFilterDropdown(!showFilterDropdown)}
            style={{
              padding: "6px 12px",
              borderRadius: "4px",
              border: "1px solid #ccc",
              fontSize: "14px",
              cursor: "pointer",
              backgroundColor: "white",
              display: "flex",
              alignItems: "center",
              gap: "4px",
              minWidth: "200px",
            }}
          >
            {selectedCategories.size === 0
              ? "すべての業種"
              : selectedCategories.size === allCategories.length
              ? "すべて選択中"
              : `${selectedCategories.size}件選択中`}
            <span style={{ marginLeft: "auto" }}>▼</span>
          </button>

          {showFilterDropdown && (
            <div
              style={{
                position: "absolute",
                top: "100%",
                left: "120px",
                marginTop: "4px",
                backgroundColor: "white",
                border: "1px solid #ccc",
                borderRadius: "4px",
                boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
                zIndex: 1000,
                minWidth: "200px",
                maxHeight: "300px",
                overflowY: "auto",
              }}
            >
              <div style={{ padding: "8px", borderBottom: "1px solid #eee" }}>
                <label
                  style={{
                    display: "flex",
                    alignItems: "center",
                    cursor: "pointer",
                  }}
                >
                  <input
                    type="checkbox"
                    checked={
                      selectedCategories.size === allCategories.length
                    }
                    onChange={toggleAllCategories}
                    style={{ marginRight: "8px" }}
                  />
                  <strong>すべて選択</strong>
                </label>
              </div>
              {allCategories.map((category) => (
                <div key={category.value} style={{ padding: "8px" }}>
                  <label
                    style={{
                      display: "flex",
                      alignItems: "center",
                      cursor: "pointer",
                    }}
                  >
                    <input
                      type="checkbox"
                      checked={selectedCategories.has(category.value)}
                      onChange={() => toggleCategory(category.value)}
                      style={{ marginRight: "8px" }}
                    />
                    {category.label}
                  </label>
                </div>
              ))}
            </div>
          )}

          {selectedCategories.size > 0 && (
            <button
              onClick={() => setSelectedCategories(new Set())}
              style={{
                padding: "4px 8px",
                fontSize: "12px",
                borderRadius: "4px",
                border: "1px solid #ccc",
                backgroundColor: "#f5f5f5",
                cursor: "pointer",
              }}
            >
              クリア
            </button>
          )}
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
            📝 <strong>会社情報表示</strong>
          </p>
          <p>
            🏢 <strong>業種別分類</strong>
          </p>
          <p>
            📞 <strong>連絡先情報</strong>
          </p>
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
          ) : filteredCompanies.length === 0 ? (
            <div className="business-entity-empty">
              選択した業種の会社が見つかりません
            </div>
          ) : (
            <div>
              {filteredCompanies.map((company, index) => (
                <div
                  key={company.id || index}
                  className="business-entity-item-row"
                  onClick={() => handleCompanyClick(company)}
                  title="クリックして詳細表示"
                  style={{ cursor: "pointer" }}
                >
                  <div className="business-entity-item-info">
                    <div className="business-entity-item-info-company">
                      {company.shortName || "会社名未設定"}
                    </div>

                    <div className="business-entity-item-info-company">
                      {(() => {
                        const matched = allCategories.find(
                          (cat) => cat.value === String(company.category)
                        );
                        if (matched) {
                          return matched.label;
                        }
                        const numericCategory = Number(company.category);
                        return getCategoryName(
                          Number.isNaN(numericCategory)
                            ? undefined
                            : numericCategory
                        );
                      })()}
                    </div>

                    <div className="business-entity-item-info-spacer"></div>
                  </div>

                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "8px",
                    }}
                  >
                    {company.tags && company.tags.length > 0 && (
                      <div
                        style={{
                          display: "flex",
                          gap: "4px",
                          flexWrap: "wrap",
                        }}
                      >
                        {company.tags.slice(0, 3).map((tag, tagIndex) => (
                          <span
                            key={tagIndex}
                            style={{
                              fontSize: "0.7em",
                              background: "#e3f2fd",
                              color: "#1976d2",
                              padding: "2px 6px",
                              borderRadius: "10px",
                              whiteSpace: "nowrap",
                            }}
                          >
                            {tag}
                          </span>
                        ))}
                        {company.tags.length > 3 && (
                          <span style={{ fontSize: "0.7em", color: "#666" }}>
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
        onCompanyUpdate={handleCompanyUpdate}
      />
    </div>
  );
};

export default Companies;
