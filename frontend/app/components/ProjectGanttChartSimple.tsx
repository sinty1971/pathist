import { useState, useEffect } from "react";
import { getProjectRecent } from "../api/sdk.gen";
import type { ModelsProject } from "../api/types.gen";
import ProjectEditModal from "./ProjectEditModal";
import { useProject } from "../contexts/ProjectContext";

const ProjectGanttChartSimple = () => {
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
      <div style={{ padding: "20px" }}>
        <div>工事データを読み込み中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: "20px" }}>
        <div
          style={{
            color: "red",
            padding: "10px",
            backgroundColor: "#ffe6e6",
            borderRadius: "4px",
          }}
        >
          {error}
        </div>
        <button
          onClick={loadProjects}
          style={{ marginTop: "10px", padding: "10px 20px" }}
        >
          再試行
        </button>
      </div>
    );
  }

  return (
    <div style={{ 
      padding: "20px", 
      paddingTop: "60px",
      flex: 1,
      overflow: "hidden",
      display: "flex",
      flexDirection: "column",
      boxSizing: "border-box"
    }}>
      {showHelp && (
        <div
          style={{
            marginBottom: "20px",
            padding: "15px",
            backgroundColor: "#f0f8ff",
            borderRadius: "4px",
            border: "1px solid #b3d9ff",
            position: "relative",
            flexShrink: 0
          }}
        >
          <button
            onClick={() => setShowHelp(false)}
            style={{
              position: "absolute",
              top: "10px",
              right: "10px",
              background: "none",
              border: "none",
              fontSize: "16px",
              cursor: "pointer",
              color: "#666"
            }}
            title="閉じる"
          >
            ×
          </button>
          
          <h3 style={{ marginTop: 0 }}>使用方法</h3>
          <p>
            📝 <strong>リストをクリック</strong>して工事情報を編集できます
          </p>
          <p>✅ 開始日・終了日・説明・タグ・会社名・現場名を編集可能</p>
          <p>💾 編集後は自動で保存されます</p>

          <h3 style={{ marginTop: "15px" }}>開発状況</h3>
          <p>✅ 工事データの取得</p>
          <p>✅ 編集モーダル機能</p>
          <p>🔄 工程表機能（次のステップ）</p>
        </div>
      )}
      
      <div
        style={{
          border: "1px solid #ddd",
          borderRadius: "8px",
          overflow: "hidden",
          position: "relative",
          height: "calc(100vh - 240px)",
          display: "flex",
          flexDirection: "column"
        }}
      >
        <div
          style={{
            backgroundColor: "#f5f5f5",
            padding: "10px 15px",
            fontWeight: "bold",
            borderBottom: "1px solid #ddd",
            display: "flex",
            alignItems: "center",
            gap: "16px",
            flexShrink: 0
          }}
        >
          <div style={{ minWidth: "90px", textAlign: "center", fontSize: "14px" }}>開始日</div>
          <div style={{ minWidth: "120px", fontSize: "14px" }}>会社名</div>
          <div style={{ minWidth: "120px", fontSize: "14px" }}>現場名</div>
          <div style={{ flex: 1 }}></div>
          <div style={{ minWidth: "90px", textAlign: "center", fontSize: "14px", marginRight: "24px" }}>終了日</div>
          <div style={{ minWidth: "80px", textAlign: "center", fontSize: "14px" }}>ステータス</div>
          <button
            onClick={() => setShowHelp(!showHelp)}
            style={{
              background: "none",
              border: "1px solid #ccc",
              borderRadius: "50%",
              width: "24px",
              height: "24px",
              cursor: "pointer",
              fontSize: "12px",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              marginLeft: "8px",
              color: "#666"
            }}
            title="使用方法を表示"
          >
            ?
          </button>
        </div>

        <div style={{ 
          flex: 1,
          overflowY: "auto",
          minHeight: 0
        }}>
          {projects.length === 0 ? (
            <div style={{ padding: "20px", textAlign: "center", color: "#666" }}>
              工事データが見つかりません
            </div>
          ) : (
            <div>
              {projects.map((project, index) => (
              <div
                key={project.id || index}
                style={{
                  padding: "15px",
                  borderBottom:
                    index < projects.length - 1 ? "1px solid #eee" : "none",
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  cursor: "pointer",
                  transition: "background-color 0.3s",
                }}
                onClick={() => handleProjectClick(project)}
                onMouseEnter={(e) =>
                  (e.currentTarget.style.backgroundColor = "#f8f9fa")
                }
                onMouseLeave={(e) =>
                  (e.currentTarget.style.backgroundColor = "transparent")
                }
                title="クリックして編集"
              >
                <div style={{ display: "flex", alignItems: "center", gap: "16px", width: "100%" }}>
                  <div style={{ 
                    fontWeight: "600", 
                    fontSize: "14px", 
                    color: "#fff",
                    backgroundColor: "#1976d2",
                    padding: "3px 8px",
                    borderRadius: "4px",
                    minWidth: "90px",
                    textAlign: "center"
                  }}>
                    {project.start_date
                      ? new Date(
                          typeof project.start_date === 'string' 
                            ? project.start_date 
                            : (project.start_date as any)['time.Time']
                        ).toLocaleDateString("ja-JP")
                      : "未設定"}
                  </div>
                  
                  <div style={{ 
                    fontWeight: "600", 
                    fontSize: "16px", 
                    minWidth: "120px"
                  }}>
                    {project.company_name || "会社名未設定"}
                  </div>
                  
                  <div style={{ 
                    fontWeight: "600", 
                    fontSize: "16px", 
                    minWidth: "120px"
                  }}>
                    {project.location_name || "現場名未設定"}
                  </div>
                  
                  <div style={{ flex: 1 }}></div>
                  
                  <div style={{ 
                    fontSize: "14px", 
                    color: "#fff",
                    backgroundColor: "#666",
                    padding: "3px 8px",
                    borderRadius: "4px",
                    minWidth: "90px",
                    textAlign: "center",
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
                      style={{ 
                        fontSize: "16px",
                        filter: "drop-shadow(0 0 3px rgba(0, 0, 0, 0.5))"
                      }}
                      title="管理ファイルの名前変更が必要です"
                    >
                      ⚠️
                    </span>
                  )}
                  <div
                    style={{
                      padding: "4px 16px",
                      borderRadius: "4px",
                      backgroundColor:
                        project.status === "進行中"
                          ? "#4CAF50"
                          : project.status === "完了"
                          ? "#9E9E9E"
                          : project.status === "予定"
                          ? "#FF9800"
                          : "#2196F3",
                      color: "white",
                      fontSize: "12px",
                      minWidth: "80px",
                      textAlign: "center"
                    }}
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
      <ProjectEditModal
        isOpen={isEditModalOpen}
        onClose={closeEditModal}
        project={selectedProject}
        onUpdate={updateProject}
        onProjectUpdate={handleProjectUpdate}
      />
    </div>
  );
};

export default ProjectGanttChartSimple;
