import React, { useState, useEffect } from 'react';
import type { ModelsProject, ModelsTimestamp, ModelsManagedFile } from '../api/types.gen';

interface ProjectEditModalProps {
  isOpen: boolean;
  onClose: () => void;
  project: ModelsProject | null;
  onSave: (project: ModelsProject) => Promise<void>;
}

// フォーム用のデータ型（日付を文字列として扱う、ステータスは除外）
type ProjectFormData = Omit<ModelsProject, 'start_date' | 'end_date' | 'tags' | 'status'> & {
  start_date?: string;
  end_date?: string;
  tags?: string;
};

const ProjectEditModal: React.FC<ProjectEditModalProps> = ({ isOpen, onClose, project, onSave }) => {
  const [formData, setFormData] = useState<ProjectFormData>({
    id: '',
    company_name: '',
    location_name: '',
    description: '',
    tags: '',
    start_date: '',
    end_date: ''
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 日付を安全に変換する関数
  const extractDateString = (timestamp: any): string => {
    if (!timestamp) return '';
    
    // timestampが文字列の場合（通常のケース）
    if (typeof timestamp === 'string') {
      const match = timestamp.match(/^(\d{4})-(\d{2})-(\d{2})/);
      if (match) {
        return `${match[1]}-${match[2]}-${match[3]}`;
      }
    }
    
    // timestampがオブジェクトの場合（ModelsTimestamp形式）
    if (typeof timestamp === 'object') {
      const timeString = timestamp['time.Time'];
      if (timeString) {
        const match = timeString.match(/^(\d{4})-(\d{2})-(\d{2})/);
        if (match) {
          return `${match[1]}-${match[2]}-${match[3]}`;
        }
      }
    }
    
    return '';
  };

  // プロジェクトが変更されたときにフォームデータを更新
  useEffect(() => {
    if (project) {
      const startDate = extractDateString(project.start_date);
      const endDate = extractDateString(project.end_date);
      
      console.log('Project dates extracted:', { startDate, endDate });
      console.log('=== MANAGED FILES DEBUG ===');
      console.log('Project managed_files:', project.managed_files);
      console.log('Type of managed_files:', typeof project.managed_files);
      console.log('Is array:', Array.isArray(project.managed_files));
      console.log('Length:', project.managed_files?.length);
      console.log('Full project keys:', Object.keys(project));
      console.log('=== END DEBUG ===');
      
      setFormData({
        id: project.id || '',
        company_name: project.company_name || '',
        location_name: project.location_name || '',
        description: project.description || '',
        tags: Array.isArray(project.tags) ? project.tags.join(', ') : (project.tags || ''),
        start_date: startDate,
        end_date: endDate
      });
      setError(null);
    }
  }, [project]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!project) return;

    setIsLoading(true);
    setError(null);

    try {
      // フォームデータをModelsProject形式に変換
      const updatedProject: ModelsProject = {
        ...project,
        company_name: formData.company_name,
        location_name: formData.location_name,
        description: formData.description,
        tags: formData.tags ? formData.tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0) : [],
        start_date: formData.start_date ? { 'time.Time': `${formData.start_date}T00:00:00+09:00` } as ModelsTimestamp : undefined,
        end_date: formData.end_date ? { 'time.Time': `${formData.end_date}T23:59:59+09:00` } as ModelsTimestamp : undefined
      };
      
      console.log('Saving project with dates:', updatedProject);

      await onSave(updatedProject);
      onClose();
    } catch (err) {
      console.error('Error saving project:', err);
      setError(err instanceof Error ? err.message : '保存に失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>工事詳細編集</h2>
          <button type="button" className="close-button" onClick={onClose}>
            ×
          </button>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="modal-body">
          <div className="form-row">
            <div className="form-group form-group-half">
              <label htmlFor="company_name">会社名</label>
              <input
                type="text"
                id="company_name"
                name="company_name"
                value={formData.company_name}
                onChange={handleInputChange}
                className="path-input"
                required
              />
            </div>
            <div className="form-group form-group-half">
              <label htmlFor="location_name">現場名</label>
              <input
                type="text"
                id="location_name"
                name="location_name"
                value={formData.location_name}
                onChange={handleInputChange}
                className="path-input"
                required
              />
            </div>
          </div>

          <div className="form-row">
            <div className="form-group form-group-half">
              <label htmlFor="start_date">開始日</label>
              <input
                type="date"
                id="start_date"
                name="start_date"
                value={formData.start_date}
                onChange={handleInputChange}
                className="path-input"
              />
            </div>
            <div className="form-group form-group-half">
              <label htmlFor="end_date">終了日</label>
              <input
                type="date"
                id="end_date"
                name="end_date"
                value={formData.end_date}
                onChange={handleInputChange}
                className="path-input"
              />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="description">説明</label>
            <textarea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleInputChange}
              className="path-input"
              rows={2}
            />
          </div>

          <div className="form-group">
            <label htmlFor="tags">タグ</label>
            <input
              type="text"
              id="tags"
              name="tags"
              value={formData.tags}
              onChange={handleInputChange}
              className="path-input"
              placeholder="カンマ区切りで入力"
            />
          </div>

          {/* 管理ファイル一覧 */}
          {project?.managed_files && project.managed_files.length > 0 && (
            <div className="form-group">
              <label>管理ファイル ({project.managed_files.length}件)</label>
              <div className="managed-files-table">
                <div className="files-header">
                  <div className="header-item">種別</div>
                  <div className="header-item">ファイルパス</div>
                  <div className="header-item">状態</div>
                </div>
                <div className="files-body">
                  {project.managed_files.map((file: ModelsManagedFile, index: number) => (
                    <React.Fragment key={index}>
                      {file.current && (
                        <div className="file-row">
                          <div className="file-type">
                            <span className="file-icon">📎</span>
                            現在ファイル
                          </div>
                          <div className="file-path">{file.current}</div>
                          <div className="file-status">
                            <span className="status-current">使用中</span>
                          </div>
                        </div>
                      )}
                      {file.recommended && (
                        <div className="file-row">
                          <div className="file-type">
                            <span className="file-icon">💡</span>
                            推奨ファイル
                          </div>
                          <div className="file-path">{file.recommended}</div>
                          <div className="file-status">
                            <span className="status-recommended">推奨</span>
                          </div>
                        </div>
                      )}
                    </React.Fragment>
                  ))}
                </div>
              </div>
            </div>
          )}
          
          {/* 管理ファイルがない場合の表示 */}
          {project && (!project.managed_files || project.managed_files.length === 0) && (
            <div className="form-group">
              <label>管理ファイル</label>
              <div className="no-managed-files">
                <span className="no-files-icon">📄</span>
                <span>管理ファイルはありません</span>
              </div>
            </div>
          )}

          <div className="form-actions">
            <button
              type="button"
              className="cancel-button"
              onClick={onClose}
              disabled={isLoading}
            >
              キャンセル
            </button>
            <button
              type="submit"
              className="submit-button"
              disabled={isLoading}
            >
              {isLoading ? '保存中...' : '保存'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ProjectEditModal;