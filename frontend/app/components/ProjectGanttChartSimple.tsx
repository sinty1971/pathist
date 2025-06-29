import { useState, useEffect } from "react";
import { getProjectRecent } from "../api/sdk.gen";
import type { ModelsProject } from "../api/types.gen";
import ProjectEditModal from "./ProjectEditModal";

const ProjectGanttChartSimple = () => {
  const [projects, setProjects] = useState<ModelsProject[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedProject, setSelectedProject] = useState<ModelsProject | null>(
    null
  );
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  // 工事データを読み込み
  const loadProjects = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await getProjectRecent();

      if (response.data) {
        setProjects(response.data);
      } else {
        setProjects([]);
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
    <div style={{ padding: "20px" }}>
      <div style={{ marginBottom: "20px" }}>
        <p>取得した工事データ: {projects.length}件</p>
      </div>

      <div
        style={{
          border: "1px solid #ddd",
          borderRadius: "8px",
          overflow: "hidden",
        }}
      >
        <div
          style={{
            backgroundColor: "#f5f5f5",
            padding: "10px",
            fontWeight: "bold",
            borderBottom: "1px solid #ddd",
          }}
        >
          工程表
        </div>

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
                <div>
                  <div style={{ fontWeight: "bold" }}>
                    {project.company_name || "会社名未設定"} -{" "}
                    {project.location_name || "現場名未設定"}
                  </div>
                  <div
                    style={{
                      fontSize: "14px",
                      color: "#666",
                      marginTop: "5px",
                    }}
                  >
                    開始:{" "}
                    {project.start_date
                      ? new Date(
                          project.start_date as string
                        ).toLocaleDateString("ja-JP")
                      : "未設定"}{" "}
                    | 終了:{" "}
                    {project.end_date
                      ? new Date(project.end_date as string).toLocaleDateString(
                          "ja-JP"
                        )
                      : "未設定"}
                  </div>
                </div>
                <div
                  style={{
                    padding: "4px 12px",
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
                  }}
                >
                  {project.status || "未設定"}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div
        style={{
          marginTop: "20px",
          padding: "15px",
          backgroundColor: "#f0f8ff",
          borderRadius: "4px",
        }}
      >
        <h3>使用方法</h3>
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
