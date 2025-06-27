import { useState, useEffect, useRef } from 'react';
import { getProjectRecent as getProjectsRecent } from '../api/sdk.gen';
import type { ModelsProject } from '../api/types.gen';
import '../styles/gantt.css';

interface GanttItem extends ModelsProject {
  startX: number;
  width: number;
  row: number;
}

const ProjectGanttChart = () => {
  const [projects, setProjects] = useState<ModelsProject[]>([]);
  const [ganttItems, setGanttItems] = useState<GanttItem[]>([]);
  const [selectedProject, setSelectedProject] = useState<ModelsProject | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(0);
  const [viewStartDate, setViewStartDate] = useState(new Date());
  const [viewEndDate, setViewEndDate] = useState(new Date());
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const ITEMS_PER_PAGE = 10;
  const DAY_WIDTH = 30; // ピクセル/日
  const ROW_HEIGHT = 40; // ピクセル
  const MONTHS_TO_SHOW = 6; // 表示する月数

  // 工事データを読み込み
  const loadProjects = async () => {
    try {
      setLoading(true);
      setError(null);
      console.log('Loading projects...');
      
      const response = await getProjectsRecent();
      
      console.log('API response:', response);
      const projects = response.data || [];
      
      console.log('Projects:', projects);
      setProjects(projects);
    } catch (err) {
      console.error('Error loading projects:', err);
      setError(`工事データの読み込みに失敗しました: ${err instanceof Error ? err.message : 'Unknown error'}`);
    } finally {
      setLoading(false);
    }
  };

  // 初期ロード
  useEffect(() => {
    loadProjects();
  }, []);

  // 表示期間の計算
  useEffect(() => {
    const today = new Date();
    const start = new Date(today.getFullYear(), today.getMonth() - 2, 1);
    const end = new Date(today.getFullYear(), today.getMonth() + MONTHS_TO_SHOW - 2, 0);
    setViewStartDate(start);
    setViewEndDate(end);
  }, []);

  // ガントチャートアイテムの計算
  useEffect(() => {
    if (projects.length === 0) return;

    const startIndex = currentPage * ITEMS_PER_PAGE;
    const endIndex = startIndex + ITEMS_PER_PAGE;
    const pageItems = projects.slice(startIndex, endIndex);

    const items: GanttItem[] = pageItems.map((project, index) => {
      // 安全な日付処理
      let startDate: Date;
      let endDate: Date;
      
      try {
        startDate = project.start_date ? new Date(project.start_date as string) : new Date();
        // 無効な日付をチェック
        if (isNaN(startDate.getTime())) {
          startDate = new Date();
        }
      } catch {
        startDate = new Date();
      }
      
      try {
        endDate = project.end_date ? new Date(project.end_date as string) : new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
        // 無効な日付をチェック
        if (isNaN(endDate.getTime())) {
          endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
        }
      } catch {
        endDate = new Date(startDate.getTime() + 90 * 24 * 60 * 60 * 1000);
      }
      
      // 安全な計算
      const daysDiff = (startDate.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24);
      const startX = Math.max(0, daysDiff * DAY_WIDTH);
      
      const endDaysDiff = (endDate.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24);
      const endX = endDaysDiff * DAY_WIDTH;
      const width = Math.max(DAY_WIDTH, endX - startX);

      return {
        ...project,
        startX: isNaN(startX) ? 0 : startX,
        width: isNaN(width) ? DAY_WIDTH : width,
        row: index
      };
    });

    setGanttItems(items);
  }, [projects, currentPage, viewStartDate]);

  // プロジェクトクリック処理
  const handleProjectClick = (project: ModelsProject) => {
    setSelectedProject(project);
    setIsDetailModalOpen(true);
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

  // 月のヘッダーを生成
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

  // 日付フォーマット
  const formatDate = (dateString?: string | any) => {
    if (!dateString) return '';
    try {
      return new Date(dateString as string).toLocaleDateString('ja-JP');
    } catch {
      return '無効な日付';
    }
  };

  // ページネーション
  const totalPages = Math.max(1, Math.ceil(projects.length / ITEMS_PER_PAGE));
  const canGoNext = currentPage < totalPages - 1;
  const canGoPrev = currentPage > 0;

  if (loading) {
    return <div className="loading">工事データを読み込み中...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  const monthHeaders = generateMonthHeaders();
  const totalWidth = (viewEndDate.getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH;

  return (
    <div className="gantt-container">
      <h1>工程表（ガントチャート）</h1>
      
      <div className="gantt-controls">
        <div className="pagination">
          <button 
            onClick={() => setCurrentPage(currentPage - 1)} 
            disabled={!canGoPrev}
          >
            前へ
          </button>
          <span>
            {currentPage + 1} / {totalPages} ページ
          </span>
          <button 
            onClick={() => setCurrentPage(currentPage + 1)} 
            disabled={!canGoNext}
          >
            次へ
          </button>
        </div>
      </div>

      <div className="gantt-wrapper">
        <div className="gantt-sidebar">
          <div className="gantt-header-left">工事プロジェクト</div>
          {ganttItems.map((item) => (
            <div 
              key={item.id} 
              className="gantt-row-label"
              style={{ height: ROW_HEIGHT }}
            >
              <div className="project-name">{item.company_name} - {item.location_name}</div>
            </div>
          ))}
        </div>

        <div className="gantt-chart-container" ref={scrollContainerRef}>
          <div className="gantt-chart" style={{ width: totalWidth }}>
            <div className="gantt-header">
              {monthHeaders.map((header, index) => (
                <div 
                  key={index}
                  className="month-header"
                  style={{ 
                    left: header.startX, 
                    width: header.width 
                  }}
                >
                  {header.year}年{header.month}月
                </div>
              ))}
            </div>

            <div className="gantt-body">
              {/* グリッド線 */}
              <div className="gantt-grid">
                {Array.from({ length: Math.ceil(totalWidth / DAY_WIDTH) }).map((_, index) => (
                  <div 
                    key={index}
                    className="grid-line"
                    style={{ left: index * DAY_WIDTH }}
                  />
                ))}
              </div>

              {/* ガントバー */}
              {ganttItems.map((item) => (
                <div 
                  key={item.id}
                  className="gantt-bar"
                  style={{
                    left: item.startX,
                    width: item.width,
                    top: item.row * ROW_HEIGHT,
                    height: ROW_HEIGHT - 10,
                    backgroundColor: getStatusColor(item.status)
                  }}
                  onClick={() => handleProjectClick(item)}
                  title={`${item.company_name} - ${item.location_name}`}
                >
                  <span className="gantt-bar-text">
                    {item.location_name}
                  </span>
                </div>
              ))}

              {/* 今日の線 */}
              <div 
                className="today-line"
                style={{
                  left: (new Date().getTime() - viewStartDate.getTime()) / (1000 * 60 * 60 * 24) * DAY_WIDTH
                }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* 詳細モーダル */}
      {isDetailModalOpen && selectedProject && (
        <div className="modal-overlay" onClick={() => setIsDetailModalOpen(false)}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <h2>工事詳細情報</h2>
            <div className="project-details">
              <p><strong>会社名:</strong> {selectedProject.company_name}</p>
              <p><strong>現場名:</strong> {selectedProject.location_name}</p>
              <p><strong>開始日:</strong> {formatDate(selectedProject.start_date as string)}</p>
              <p><strong>終了日:</strong> {formatDate(selectedProject.end_date as string)}</p>
              <p><strong>ステータス:</strong> {selectedProject.status}</p>
              {selectedProject.description && (
                <p><strong>説明:</strong> {selectedProject.description}</p>
              )}
            </div>
            <button onClick={() => setIsDetailModalOpen(false)}>閉じる</button>
          </div>
        </div>
      )}
    </div>
  );
};

export default ProjectGanttChart;