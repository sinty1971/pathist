import { useState, useEffect } from "react";
import { Link } from "react-router";
import { getProjectRecent } from "../api/sdk.gen";
import type { ModelsProject } from "../api/types.gen";
import ProjectDetailModal from "./ProjectDetailModal";
import { useProject } from "../contexts/ProjectContext";
import "../styles/business-entity-list.css";

const Projects = () => {
  const [projects, setProjects] = useState<ModelsProject[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedProject, setSelectedProject] = useState<ModelsProject | null>(
    null
  );
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [showHelp, setShowHelp] = useState(false);
  const { setProjectCount } = useProject();

  // 工事データを読み込み
  const loadProjects = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await getProjectRecent();

      if (response.data) {
        setProjects(response.data);
        setProjectCount(response.data.length);
      } else {
        setProjects([]);
        setProjectCount(0);
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
    loadProjects();
  }, []);

  // プロジェクトクリック処理
  const handleProjectClick = (project: ModelsProject) => {
    setSelectedProject(project);
    setIsEditModalOpen(true);
  };

  // プロジェクト更新処理（APIコール）
  const updateProject = async (updatedProject: ModelsProject): Promise<ModelsProject> => {
    try {
      // バックエンドに更新リクエストを送信
      const response = await fetch("http://localhost:8080/api/project/update", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(updatedProject),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "更新に失敗しました");
      }

      // レスポンスから更新されたプロジェクトデータを取得
      const savedProject = await response.json();

      // プロジェクト一覧を更新
      setProjects((prevProjects) =>
        prevProjects.map((p) => (p.id === savedProject.id ? savedProject : p))
      );

      // 更新されたプロジェクトデータを返す
      return savedProject;
    } catch (err) {
      console.error("Error updating project:", err);
      throw err; // エラーをモーダルに伝播
    }
  };

  // モーダルを閉じる
  const closeEditModal = () => {
    setIsEditModalOpen(false);
    setSelectedProject(null);
  };

  // 管理ファイルの変更が必要かチェック
  const needsFileRename = (project: ModelsProject): boolean => {
    if (!project.managed_files || project.managed_files.length === 0) {
      return false;
    }
    
    // managed_filesの中で現在の名前と推奨名が異なるものがあるかチェック
    const needsRename = project.managed_files.some(file => {
      // currentとrecommendedが両方存在し、異なる場合にtrueを返す
      return file.current && file.recommended && file.current !== file.recommended;
    });
    
    return needsRename;
  };

  // プロジェクトデータを更新
  const handleProjectUpdate = (updatedProject: ModelsProject) => {
    // 選択中のプロジェクトを更新
    setSelectedProject(updatedProject);

    // プロジェクト一覧を更新
    setProjects((prevProjects) => {
      // 既存のプロジェクトを探す（IDで照合）
      const existingIndex = prevProjects.findIndex(p => p.id === updatedProject.id);
      
      if (existingIndex !== -1) {
        // 既存のプロジェクトを更新
        const updatedProjects = [...prevProjects];
        updatedProjects[existingIndex] = updatedProject;
        return updatedProjects;
      } else {
        // フォルダー名が変更された可能性があるため、元のプロジェクトを探して削除し、新しいものを追加
        // 同じ会社名・現場名で探す
        const oldProjectIndex = prevProjects.findIndex(p => 
          p.company_name === updatedProject.company_name && 
          p.location_name === updatedProject.location_name &&
          p.id !== updatedProject.id
        );
        
        if (oldProjectIndex !== -1) {
          // 古いプロジェクトを削除して新しいものを追加
          const updatedProjects = [...prevProjects];
          updatedProjects.splice(oldProjectIndex, 1);
          updatedProjects.push(updatedProject);
          // 開始日順でソート（新しい順）
          return updatedProjects.sort((a, b) => {
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
          return [...prevProjects, updatedProject];
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
          onClick={loadProjects}
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
          to="/projects/gantt"
          className="business-entity-gantt-button"
        >
          📊 工程表を表示
        </Link>
        
        <div className="business-entity-count">
          全{projects.length}件
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
          {projects.length === 0 ? (
            <div className="business-entity-empty">
              工事データが見つかりません
            </div>
          ) : (
            <div>
              {projects.map((project, index) => (
              <div
                key={project.id || index}
                className="business-entity-item-row"
                onClick={() => handleProjectClick(project)}
                title="クリックして編集"
              >
                <div className="business-entity-item-info">
                  <div className="business-entity-item-info-date">
                    {project.start_date
                      ? new Date(
                          typeof project.start_date === 'string' 
                            ? project.start_date 
                            : (project.start_date as any)['time.Time']
                        ).toLocaleDateString("ja-JP")
                      : "未設定"}
                  </div>
                  
                  <div className="business-entity-item-info-company">
                    {project.company_name || "会社名未設定"}
                  </div>
                  
                  <div className="business-entity-item-info-location">
                    {project.location_name || "現場名未設定"}
                  </div>
                  
                  <div className="business-entity-item-info-description">
                    {project.description || ""}
                  </div>
                  
                  <div className="business-entity-item-info-date end-date" style={{ 
                    marginRight: "24px"
                  }}>
                    ～{project.end_date
                      ? new Date(
                          typeof project.end_date === 'string' 
                            ? project.end_date 
                            : (project.end_date as any)['time.Time']
                        ).toLocaleDateString("ja-JP")
                      : "未設定"}
                  </div>
                </div>
                <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                  {needsFileRename(project) && (
                    <span 
                      className="business-entity-item-rename-indicator"
                      title="管理ファイルの名前変更が必要です"
                    >
                      ⚠️
                    </span>
                  )}
                  <div
                    className={`business-entity-item-status ${
                      project.status === "進行中" ? "business-entity-item-status-ongoing" :
                      project.status === "完了" ? "business-entity-item-status-completed" :
                      project.status === "予定" ? "business-entity-item-status-planned" :
                      ""
                    }`}
                  >
                    {project.status || "未設定"}
                  </div>
                </div>
              </div>
              ))}
            </div>
          )}
        </div>
      </div>


      {/* 編集モーダル */}
      <ProjectDetailModal
        isOpen={isEditModalOpen}
        onClose={closeEditModal}
        project={selectedProject}
        onUpdate={updateProject}
        onProjectUpdate={handleProjectUpdate}
      />
    </div>
  );
};

export default Projects;
