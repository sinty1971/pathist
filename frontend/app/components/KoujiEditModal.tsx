import { useState, useEffect } from 'react';
import type { KoujiEntryExtended } from '../types/kouji';

interface KoujiEditModalProps {
  isOpen: boolean;
  onClose: () => void;
  project: KoujiEntryExtended | null;
  onSave: (updatedProject: KoujiEntryExtended) => void;
}

const KoujiEditModal = ({ isOpen, onClose, project, onSave }: KoujiEditModalProps) => {
  const [formData, setFormData] = useState<Partial<KoujiEntryExtended>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tagInput, setTagInput] = useState('');

  // プロジェクトデータをフォームにセット
  useEffect(() => {
    if (project && isOpen) {
      setFormData({
        ...project,
        start_date: project.start_date ? formatDateForInput(project.start_date) : '',
        end_date: project.end_date ? formatDateForInput(project.end_date) : '',
      });
      setTagInput(project.tags ? project.tags.join(', ') : '');
      setError(null);
    }
  }, [project, isOpen]);

  // 日付を入力フィールド用にフォーマット (YYYY-MM-DD)
  const formatDateForInput = (dateString: string): string => {
    try {
      const date = new Date(dateString);
      return date.toISOString().split('T')[0];
    } catch {
      return '';
    }
  };

  // フォーム値の更新
  const handleInputChange = (field: keyof KoujiEntryExtended, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
  };

  // タグの処理
  const handleTagChange = (value: string) => {
    setTagInput(value);
    const tags = value.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0);
    setFormData(prev => ({
      ...prev,
      tags
    }));
  };

  // 保存処理
  const handleSave = async () => {
    if (!project) return;

    setIsLoading(true);
    setError(null);

    try {
      // 日付の検証
      if (formData.start_date && formData.end_date) {
        const startDate = new Date(formData.start_date);
        const endDate = new Date(formData.end_date);
        
        if (startDate > endDate) {
          setError('開始日は終了日より前である必要があります');
          setIsLoading(false);
          return;
        }
      }

      // ISO形式に変換
      const updatedProject: KoujiEntryExtended = {
        ...project,
        ...formData,
        start_date: formData.start_date ? `${formData.start_date}T00:00:00` : project.start_date,
        end_date: formData.end_date ? `${formData.end_date}T23:59:59` : project.end_date,
      };

      onSave(updatedProject);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : '保存に失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  // 必須フィールドの検証
  const isFormValid = () => {
    return formData.company_name && formData.location_name;
  };

  if (!isOpen || !project) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content large-modal" onClick={e => e.stopPropagation()}>
        <div className="modal-header">
          <h2>工事プロジェクト編集</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>

        <div className="modal-body">
          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <div className="form-grid">
            <div className="form-group">
              <label htmlFor="company_name">会社名 *</label>
              <input
                id="company_name"
                type="text"
                value={formData.company_name || ''}
                onChange={e => handleInputChange('company_name', e.target.value)}
                placeholder="例: 豊田築炉"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="location_name">現場名 *</label>
              <input
                id="location_name"
                type="text"
                value={formData.location_name || ''}
                onChange={e => handleInputChange('location_name', e.target.value)}
                placeholder="例: 名和工場"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="start_date">開始日</label>
              <input
                id="start_date"
                type="date"
                value={formData.start_date || ''}
                onChange={e => handleInputChange('start_date', e.target.value)}
              />
            </div>

            <div className="form-group">
              <label htmlFor="end_date">終了日</label>
              <input
                id="end_date"
                type="date"
                value={formData.end_date || ''}
                onChange={e => handleInputChange('end_date', e.target.value)}
              />
            </div>

            <div className="form-group full-width">
              <label htmlFor="description">説明</label>
              <textarea
                id="description"
                value={formData.description || ''}
                onChange={e => handleInputChange('description', e.target.value)}
                placeholder="工事の詳細説明を入力してください"
                rows={4}
              />
            </div>

            <div className="form-group full-width">
              <label htmlFor="tags">タグ</label>
              <input
                id="tags"
                type="text"
                value={tagInput}
                onChange={e => handleTagChange(e.target.value)}
                placeholder="タグをカンマ区切りで入力 (例: 工事, 豊田築炉, 名和工場)"
              />
              {formData.tags && formData.tags.length > 0 && (
                <div className="tags-preview">
                  {formData.tags.map((tag, index) => (
                    <span key={index} className="tag-item">
                      {tag}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="modal-footer">
          <button 
            className="cancel-button" 
            onClick={onClose}
            disabled={isLoading}
          >
            キャンセル
          </button>
          <button 
            className="save-button" 
            onClick={handleSave}
            disabled={isLoading || !isFormValid()}
          >
            {isLoading ? '保存中...' : '保存'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default KoujiEditModal;