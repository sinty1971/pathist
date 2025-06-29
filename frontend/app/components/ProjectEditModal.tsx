import React, { useState, useEffect } from 'react';
import type { ModelsProject, ModelsTimestamp, ModelsManagedFile } from '../api/types.gen';
import { postProjectRenameManagedFile, getProjectGetByPath } from '../api/sdk.gen';
import CallyCalendar from './CallyCalendar';

interface ProjectEditModalProps {
  isOpen: boolean;
  onClose: () => void;
  project: ModelsProject | null;
  onUpdate: (project: ModelsProject) => Promise<ModelsProject>;
  onProjectUpdate?: (project: ModelsProject) => void;
}

// フォーム用のデータ型（日付を文字列として扱う、ステータスは除外）
type ProjectFormData = Omit<ModelsProject, 'start_date' | 'end_date' | 'tags' | 'status'> & {
  start_date?: string;
  end_date?: string;
  tags?: string;
};

// DaisyUI Calendar コンポーネント
interface DatePickerProps {
  value: string;
  onChange: (dateString: string) => void;
  placeholder: string;
  disabled?: boolean;
  minDate?: string;
}

const DatePickerComponent: React.FC<DatePickerProps> = ({ value, onChange, placeholder, disabled }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<string>(value || '');

  const handleDateSelect = (dateString: string) => {
    setSelectedDate(dateString);
    onChange(dateString);
    // 少し遅延を入れて確実に閉じる
    setTimeout(() => {
      setIsOpen(false);
    }, 100);
  };


  useEffect(() => {
    setSelectedDate(value || '');
  }, [value]);

  // 外側クリックで閉じる
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Element;
      if (isOpen && !target.closest('.date-picker-container') && !target.closest('.cally-wrapper')) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen]);

  const formatDisplayDate = (dateStr: string) => {
    if (!dateStr) return placeholder;
    try {
      const date = new Date(dateStr);
      return date.toLocaleDateString('ja-JP');
    } catch {
      return dateStr;
    }
  };

  return (
    <div className="relative date-picker-container">
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        disabled={disabled}
        className="path-input w-full text-left"
        style={{ width: '100%' }}
      >
        {formatDisplayDate(selectedDate)}
      </button>
      
      {isOpen && (
        <div className="absolute top-full left-0 z-50 mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg p-3" style={{ width: '220px' }}>
          <CallyCalendar 
            onDateSelect={handleDateSelect}
            selectedDate={selectedDate}
          />
        </div>
      )}
    </div>
  );
};

const ProjectEditModal: React.FC<ProjectEditModalProps> = ({ isOpen, onClose, project, onUpdate, onProjectUpdate }) => {
  const [formData, setFormData] = useState<ProjectFormData>({
    id: '',
    company_name: '',
    location_name: '',
    description: '',
    tags: '',
    start_date: '',
    end_date: ''
  });
  const [currentProject, setCurrentProject] = useState<ModelsProject | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isRenaming, setIsRenaming] = useState(false);
  const [hasFilenameChanges, setHasFilenameChanges] = useState(false);
  const [initialFilenameData, setInitialFilenameData] = useState<{
    start_date: string;
    company_name: string;
    location_name: string;
  }>({ start_date: '', company_name: '', location_name: '' });

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
      
      const companyName = project.company_name || '';
      const locationName = project.location_name || '';
      
      setCurrentProject(project);
      setFormData({
        id: project.id || '',
        company_name: companyName,
        location_name: locationName,
        description: project.description || '',
        tags: Array.isArray(project.tags) ? project.tags.join(', ') : (project.tags || ''),
        start_date: startDate,
        end_date: endDate
      });
      
      // 初期ファイル名関連データを保存
      setInitialFilenameData({
        start_date: startDate,
        company_name: companyName,
        location_name: locationName
      });
      
      setHasFilenameChanges(false);
      setError(null);
    }
  }, [project]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    
    // フォームデータを更新
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    
    // ファイル名関連フィールドの変更をチェック
    if (name === 'company_name' || name === 'location_name') {
      checkFilenameChanges({ ...formData, [name]: value });
    }
  };
  
  // 推奨ファイル名を生成する関数
  const generateRecommendedFileName = (originalFileName: string, formData: ProjectFormData): string => {
    if (!formData.start_date || !formData.company_name || !formData.location_name) {
      return originalFileName;
    }

    // 日付をYYYY-MMDD形式に変換
    const datePart = formData.start_date.replace(/-/g, '').substring(0, 8); // YYYYMMDD
    const formattedDate = `${datePart.substring(0, 4)}-${datePart.substring(4, 8)}`; // YYYY-MMDD
    
    // 新しいプレフィックスを作成
    const newPrefix = `${formattedDate} ${formData.company_name} ${formData.location_name}`;
    
    // 既存のファイル名から拡張子を取得
    const fileExtension = originalFileName.includes('.') 
      ? '.' + originalFileName.split('.').pop() 
      : '';
    
    // 既存のファイル名が日付形式で始まっているかチェック
    const datePattern = /^\d{4}-\d{4}\s+/;
    
    if (datePattern.test(originalFileName)) {
      // 既存の日付形式を新しいプレフィックスに置換
      const afterPrefix = originalFileName.replace(/^\d{4}-\d{4}\s+[^\s]+\s+[^\s]+\s*/, '');
      return afterPrefix 
        ? `${newPrefix} ${afterPrefix}`
        : `${newPrefix}${fileExtension}`;
    } else {
      // 日付形式でない場合は、プレフィックスを追加
      const nameWithoutExt = fileExtension 
        ? originalFileName.substring(0, originalFileName.lastIndexOf('.'))
        : originalFileName;
      return `${newPrefix} ${nameWithoutExt}${fileExtension}`;
    }
  };

  // 管理ファイルの推奨名を更新する関数
  const updateManagedFileRecommendations = (formData: ProjectFormData) => {
    if (!currentProject?.managed_files) return;

    const updatedManagedFiles = currentProject.managed_files.map(file => {
      if (file.current) {
        const recommendedName = generateRecommendedFileName(file.current, formData);
        return {
          ...file,
          recommended: recommendedName
        };
      }
      return file;
    });

    setCurrentProject(prev => prev ? {
      ...prev,
      managed_files: updatedManagedFiles
    } : null);
  };

  // ファイル名関連の変更をチェックする関数
  const checkFilenameChanges = (currentFormData: ProjectFormData) => {
    const hasChanges = 
      currentFormData.start_date !== initialFilenameData.start_date ||
      currentFormData.company_name !== initialFilenameData.company_name ||
      currentFormData.location_name !== initialFilenameData.location_name;
    
    setHasFilenameChanges(hasChanges);
    
    // ファイル名関連の変更がある場合、管理ファイルの推奨名を更新
    if (hasChanges) {
      updateManagedFileRecommendations(currentFormData);
    }
  };

  // フィールド変更完了時の処理（指定されたフォームデータを使用）
  const handleFieldUpdateWithData = async (useFormData: ProjectFormData) => {
    if (!project) return;

    setIsLoading(true);
    setError(null);

    try {
      // 指定されたフォームデータをModelsProject形式に変換
      const updatedProject: ModelsProject = {
        ...project,
        company_name: useFormData.company_name,
        location_name: useFormData.location_name,
        description: useFormData.description,
        tags: useFormData.tags ? useFormData.tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0) : [],
        start_date: useFormData.start_date ? { 'time.Time': `${useFormData.start_date}T00:00:00+09:00` } as ModelsTimestamp : undefined,
        end_date: useFormData.end_date ? { 'time.Time': `${useFormData.end_date}T23:59:59+09:00` } as ModelsTimestamp : undefined
      };
      

      // 更新前のフォルダー名を保存
      const originalFolderName = project.name;
      
      // 更新して更新されたプロジェクトデータを取得
      const savedProject = await onUpdate(updatedProject);
      
      // フォルダー名が変更されたかチェック
      const folderNameChanged = originalFolderName && savedProject.name && originalFolderName !== savedProject.name;
      
      // 更新後、最新のプロジェクトデータを再取得（管理ファイル情報を含む）
      if (savedProject.name) {
        const latestProjectResponse = await getProjectGetByPath({
          path: {
            path: savedProject.name
          }
        });
        
        if (latestProjectResponse.data) {
          const latestProject = latestProjectResponse.data;
          
          // フォルダー名が変更された場合はモーダルを閉じる
          if (folderNameChanged) {
            // 親コンポーネントに最新データを渡してモーダルを閉じる
            if (onProjectUpdate) {
              onProjectUpdate(latestProject);
            }
            
            // モーダルを閉じる
            setTimeout(() => {
              onClose();
            }, 100);
            
            return;
          }
          
          // フォルダー名が変更されていない場合は通常の更新処理
          // ローカルステートを更新（管理ファイル表示用）
          setCurrentProject(latestProject);
          
          // 成功した場合、最新のプロジェクトデータを更新
          if (onProjectUpdate) {
            onProjectUpdate(latestProject);
          }
          
          // フォームデータも最新データで同期
          const startDate = extractDateString(latestProject.start_date);
          const endDate = extractDateString(latestProject.end_date);
          
          setFormData({
            id: latestProject.id || '',
            company_name: latestProject.company_name || '',
            location_name: latestProject.location_name || '',
            description: latestProject.description || '',
            tags: Array.isArray(latestProject.tags) ? latestProject.tags.join(', ') : (latestProject.tags || ''),
            start_date: startDate,
            end_date: endDate
          });
        }
      } else {
        // name がない場合は従来の処理
        if (onProjectUpdate) {
          onProjectUpdate(savedProject);
        }
        
        const startDate = extractDateString(savedProject.start_date);
        const endDate = extractDateString(savedProject.end_date);
        
        setFormData({
          id: savedProject.id || '',
          company_name: savedProject.company_name || '',
          location_name: savedProject.location_name || '',
          description: savedProject.description || '',
          tags: Array.isArray(savedProject.tags) ? savedProject.tags.join(', ') : (savedProject.tags || ''),
          start_date: startDate,
          end_date: endDate
        });
      }
    } catch (err) {
      console.error('Error updating field:', err);
      setError(err instanceof Error ? err.message : 'フィールドの更新に失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  // フィールド変更完了時の処理（Enter押下またはフォーカスアウト）
  const handleFieldUpdate = async (fieldName?: string) => {
    // ファイル名変更が必要な場合は即時更新しない
    if (hasFilenameChanges) {
      return;
    }
    // ファイル名関連フィールドの場合は更新しない
    if (fieldName === 'company_name' || fieldName === 'location_name') {
      return;
    }
    handleFieldUpdateWithData(formData);
  };

  // 非ファイル名関連フィールド用のblurハンドラー
  const handleNonFilenameBlur = () => {
    // ファイル名変更が必要な場合は即時更新しない
    if (hasFilenameChanges) {
      return;
    }
    handleFieldUpdateWithData(formData);
  };

  // Enterキー押下時の処理
  const handleKeyDown = (e: React.KeyboardEvent, fieldName?: string) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleFieldUpdate(fieldName);
    }
  };

  // ファイル名関連フィールド用のキーダウンハンドラー
  const handleFilenameKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      // ファイル名関連フィールドではEnterでも更新しない
    }
  };

  // DaisyUI日付ピッカー用のハンドラー
  const handleDaisyDateChange = (dateString: string, fieldName: 'start_date' | 'end_date') => {
    // フォームデータを更新
    const newFormData = {
      ...formData,
      [fieldName]: dateString
    };
    setFormData(newFormData);
    
    // 開始日の場合はファイル名関連なので即座に更新しない
    if (fieldName === 'start_date') {
      checkFilenameChanges(newFormData);
      return;
    }
    
    // 終了日の場合でも、ファイル名変更が必要な場合は即座に更新しない
    if (hasFilenameChanges) {
      return;
    }
    
    // 終了日の場合は即座に更新
    if (dateString) {
      handleFieldUpdateWithData(newFormData);
    }
  };
  
  // ファイル名変更ボタンのハンドラー
  const handleFilenameUpdate = async () => {
    await handleFieldUpdateWithData(formData);
    setHasFilenameChanges(false);
    // 新しい値を初期値として保存
    setInitialFilenameData({
      start_date: formData.start_date || '',
      company_name: formData.company_name || '',
      location_name: formData.location_name || ''
    });
  };

  const handleRenameFiles = async () => {
    if (!currentProject || !currentProject.managed_files) return;

    setIsRenaming(true);
    setError(null);

    try {
      // 現在のファイル名のリストを取得
      const currentFiles = currentProject.managed_files
        .filter(file => file.current && file.recommended && file.current !== file.recommended)
        .map(file => file.current as string);

      if (currentFiles.length === 0) {
        setError('変更対象のファイルがありません');
        return;
      }

      // API呼び出し
      const response = await postProjectRenameManagedFile({
        body: {
          project: currentProject,
          currents: currentFiles
        }
      });

      if (response.data) {
        // プロジェクトデータを再取得
        if (currentProject.name) {
          const updatedProjectResponse = await getProjectGetByPath({
            path: {
              path: currentProject.name
            }
          });
          
          if (updatedProjectResponse.data && onProjectUpdate) {
            // ローカルステートも更新
            setCurrentProject(updatedProjectResponse.data);
            // プロジェクトデータを更新
            onProjectUpdate(updatedProjectResponse.data);
            // フォームデータも更新
            const startDate = extractDateString(updatedProjectResponse.data.start_date);
            const endDate = extractDateString(updatedProjectResponse.data.end_date);
            
            setFormData({
              id: updatedProjectResponse.data.id || '',
              company_name: updatedProjectResponse.data.company_name || '',
              location_name: updatedProjectResponse.data.location_name || '',
              description: updatedProjectResponse.data.description || '',
              tags: Array.isArray(updatedProjectResponse.data.tags) ? updatedProjectResponse.data.tags.join(', ') : (updatedProjectResponse.data.tags || ''),
              start_date: startDate,
              end_date: endDate
            });
          }
        }
      } else if (response.error) {
        throw new Error('ファイル名の変更に失敗しました');
      }
    } catch (err) {
      console.error('Error renaming files:', err);
      setError(err instanceof Error ? err.message : 'ファイル名の変更に失敗しました');
    } finally {
      setIsRenaming(false);
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

        <div className="modal-body">
          <div className="form-row">
            <div className="form-group form-group-half">
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', backgroundColor: '#f8f9fa', padding: '8px', borderRadius: '4px', border: '2px solid #e3f2fd' }}>
                <label style={{ minWidth: '60px', margin: 0, fontWeight: '600', color: '#1976d2' }}>開始日</label>
                <DatePickerComponent
                  value={formData.start_date || ''}
                  onChange={(dateString) => handleDaisyDateChange(dateString, 'start_date')}
                  placeholder="開始日を選択"
                  disabled={isLoading}
                />
              </div>
            </div>
            <div className="form-group form-group-half">
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <label style={{ minWidth: '60px', margin: 0 }}>終了日</label>
                <DatePickerComponent
                  value={formData.end_date || ''}
                  onChange={(dateString) => handleDaisyDateChange(dateString, 'end_date')}
                  placeholder="終了日を選択"
                  disabled={isLoading}
                  minDate={formData.start_date}
                />
              </div>
            </div>
          </div>

          <div className="form-row">
            <div className="form-group form-group-half">
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', backgroundColor: '#f8f9fa', padding: '8px', borderRadius: '4px', border: '2px solid #e3f2fd' }}>
                <label htmlFor="company_name" style={{ minWidth: '60px', margin: 0, fontWeight: '600', color: '#1976d2' }}>会社名</label>
                <input
                  type="text"
                  id="company_name"
                  name="company_name"
                  value={formData.company_name}
                  onChange={handleInputChange}
                  onKeyDown={handleFilenameKeyDown}
                  className="path-input"
                  required
                  disabled={isLoading}
                  style={{ fontWeight: '500' }}
                />
              </div>
            </div>
            <div className="form-group form-group-half">
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', backgroundColor: '#f8f9fa', padding: '8px', borderRadius: '4px', border: '2px solid #e3f2fd' }}>
                <label htmlFor="location_name" style={{ minWidth: '60px', margin: 0, fontWeight: '600', color: '#1976d2' }}>現場名</label>
                <input
                  type="text"
                  id="location_name"
                  name="location_name"
                  value={formData.location_name}
                  onChange={handleInputChange}
                  onKeyDown={handleFilenameKeyDown}
                  className="path-input"
                  required
                  disabled={isLoading}
                  style={{ fontWeight: '500' }}
                />
              </div>
            </div>
          </div>

          <div className="form-group">
            <div style={{ display: 'flex', alignItems: 'flex-start', gap: '8px' }}>
              <label htmlFor="description" style={{ minWidth: '60px', margin: 0, marginTop: '8px' }}>説明</label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                onKeyDown={(e) => handleKeyDown(e, 'description')}
                onBlur={handleNonFilenameBlur}
                className="path-input"
                rows={3}
                disabled={isLoading}
                style={{ resize: 'vertical', width: '100%' }}
              />
            </div>
          </div>

          <div className="form-group">
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <label htmlFor="tags" style={{ minWidth: '60px', margin: 0 }}>タグ</label>
              <input
                type="text"
                id="tags"
                name="tags"
                value={formData.tags}
                onChange={handleInputChange}
                onKeyDown={(e) => handleKeyDown(e, 'tags')}
                onBlur={handleNonFilenameBlur}
                className="path-input"
                placeholder="カンマ区切りで入力"
                disabled={isLoading}
                style={{ width: '100%' }}
              />
            </div>
          </div>

          {/* ファイル名変更警告とボタン */}
          {hasFilenameChanges && (
            <div style={{ 
              backgroundColor: '#fff3cd', 
              border: '1px solid #ffeaa7', 
              borderRadius: '4px', 
              padding: '12px', 
              margin: '20px 0',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between'
            }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <span style={{ fontSize: '18px' }}>⚠️</span>
                <span style={{ color: '#856404', fontWeight: '500' }}>
                  ファイル名の変更が必要です。このボタンを押下して全ての変更をまとめて更新してください。
                </span>
              </div>
              <button
                type="button"
                onClick={handleFilenameUpdate}
                disabled={isLoading}
                style={{
                  backgroundColor: '#ff9800',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  padding: '8px 16px',
                  fontWeight: '600',
                  cursor: 'pointer',
                  marginLeft: '12px',
                  whiteSpace: 'nowrap'
                }}
              >
                {isLoading ? '更新中...' : 'ファイル名を更新'}
              </button>
            </div>
          )}

          {/* 管理ファイル一覧 */}
          {currentProject?.managed_files && currentProject.managed_files.length > 0 && (
            <div className="form-group">
              <label>管理ファイル</label>
              <div style={{ 
                border: '1px solid #ddd', 
                borderRadius: '6px', 
                backgroundColor: '#fafafa',
                padding: '16px'
              }}>
                {currentProject.managed_files.map((file: ModelsManagedFile, index: number) => {
                  const needsRename = file.current && file.recommended && file.current !== file.recommended;
                  return (
                    <div key={index} style={{
                      marginBottom: index < currentProject.managed_files!.length - 1 ? '16px' : '0',
                      padding: '12px',
                      backgroundColor: 'white',
                      borderRadius: '4px',
                      border: '1px solid #e0e0e0'
                    }}>
                      {/* 現在のファイル */}
                      <div style={{ 
                        display: 'flex', 
                        alignItems: 'center', 
                        marginBottom: file.recommended ? '8px' : '0',
                        fontSize: '14px'
                      }}>
                        <span style={{ marginRight: '8px' }}>📎</span>
                        <span style={{ fontWeight: '500', marginRight: '8px' }}>現在:</span>
                        <span style={{ fontFamily: 'monospace', color: '#666' }}>{file.current}</span>
                      </div>
                      
                      {/* 推奨ファイル */}
                      {file.recommended && (
                        <div style={{ 
                          display: 'flex', 
                          alignItems: 'center', 
                          justifyContent: 'space-between',
                          fontSize: '14px'
                        }}>
                          <div style={{ display: 'flex', alignItems: 'center' }}>
                            <span style={{ marginRight: '8px' }}>💡</span>
                            <span style={{ fontWeight: '500', marginRight: '8px' }}>推奨:</span>
                            <span style={{ fontFamily: 'monospace', color: '#0066cc' }}>{file.recommended}</span>
                          </div>
                          
                          {/* 変更ボタン */}
                          {needsRename && (
                            <button
                              type="button"
                              onClick={handleRenameFiles}
                              disabled={isRenaming || isLoading}
                              className="btn btn-primary btn-sm"
                            >
                              {isRenaming ? '変更中...' : '変更'}
                            </button>
                          )}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          )}
          
          {/* 管理ファイルがない場合の表示 */}
          {currentProject && (!currentProject.managed_files || currentProject.managed_files.length === 0) && (
            <div className="form-group">
              <label>管理ファイル</label>
              <div className="no-managed-files">
                <span className="no-files-icon">📄</span>
                <span>管理ファイルはありません</span>
              </div>
            </div>
          )}

          <div className="modal-action">
            <button
              type="button"
              className="btn"
              onClick={onClose}
              disabled={isLoading}
            >
              閉じる
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProjectEditModal;