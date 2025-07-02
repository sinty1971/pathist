import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  TextField,
  Button,
  Box,
  Typography,
  Alert,
  Paper,
  IconButton,
  CircularProgress
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import WarningIcon from '@mui/icons-material/Warning';
import AttachFileIcon from '@mui/icons-material/AttachFile';
import LightbulbIcon from '@mui/icons-material/Lightbulb';
import type { ModelsProject, ModelsTimestamp, ModelsManagedFile } from '../api/types.gen';
import { postProjectRenameManagedFile, getProjectGetByPath } from '../api/sdk.gen';
import { CalendarPicker } from './CalendarPicker';

interface ProjectDetailModalProps {
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

// CalendarPickerコンポーネントをインポート済み

const ProjectDetailModal: React.FC<ProjectDetailModalProps> = ({ isOpen, onClose, project, onUpdate, onProjectUpdate }) => {
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
    
    // フォームデータを更新（UIの表示のみ）
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    
    // リアルタイムでのチェックや更新は一切行わない
    // すべての更新処理はEnterキーまたはフォーカスアウト時のみ実行
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
  const handleNonFilenameBlur = (fieldName: string) => () => {
    // ファイル名変更が必要な場合は即時更新しない
    if (hasFilenameChanges) {
      return;
    }
    // description と tags のみ対象
    if (fieldName === 'description' || fieldName === 'tags') {
      handleFieldUpdateWithData(formData);
    }
  };

  // ファイル名関連フィールド用のblurハンドラー
  const handleFilenameBlur = (fieldName: string) => () => {
    // ファイル名関連フィールドの変更をチェック
    if (fieldName === 'company_name' || fieldName === 'location_name') {
      checkFilenameChanges(formData);
    }
  };

  // Enterキー押下時の処理
  const handleKeyDown = (e: React.KeyboardEvent, fieldName?: string) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleFieldUpdate(fieldName);
    }
  };

  // ファイル名関連フィールド用のキーダウンハンドラー
  const handleFilenameKeyDown = (e: React.KeyboardEvent, fieldName: string) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      // Enterキーでファイル名変更チェックを実行
      if (fieldName === 'company_name' || fieldName === 'location_name') {
        checkFilenameChanges(formData);
      }
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


  return (
    <Dialog
      open={isOpen}
      onClose={onClose}
      maxWidth="md"
      fullWidth
      slotProps={{
        paper: {
          sx: {
            borderRadius: 2,
            minHeight: '500px'
          }
        }
      }}
    >
      <DialogTitle
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          pb: 1
        }}
      >
        <Typography variant="h6" component="h2">
          工事詳細編集
        </Typography>
        <IconButton
          onClick={onClose}
          size="small"
          sx={{
            color: 'grey.500'
          }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <DialogContent dividers>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Box sx={{ mt: 1 }}>
          {/* 会社名・現場名セクション */}
          <Box sx={{ display: 'flex', gap: 2, mb: 3, flexDirection: { xs: 'column', sm: 'row' } }}>
            <Box sx={{ flex: 1 }}>
              <Paper
                elevation={0}
                sx={{
                  p: 2,
                  backgroundColor: 'primary.50',
                  border: '2px solid',
                  borderColor: 'primary.200',
                  borderRadius: 2
                }}
              >
                <Typography
                  variant="subtitle2"
                  sx={{
                    color: 'primary.main',
                    fontWeight: 600,
                    mb: 1
                  }}
                >
                  会社名
                </Typography>
                <TextField
                  fullWidth
                  size="small"
                  name="company_name"
                  value={formData.company_name}
                  onChange={handleInputChange}
                  onKeyDown={(e) => handleFilenameKeyDown(e, 'company_name')}
                  onBlur={handleFilenameBlur('company_name')}
                  required
                  disabled={isLoading}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      backgroundColor: 'white'
                    }
                  }}
                />
              </Paper>
            </Box>
            <Box sx={{ flex: 1 }}>
              <Paper
                elevation={0}
                sx={{
                  p: 2,
                  backgroundColor: 'primary.50',
                  border: '2px solid',
                  borderColor: 'primary.200',
                  borderRadius: 2
                }}
              >
                <Typography
                  variant="subtitle2"
                  sx={{
                    color: 'primary.main',
                    fontWeight: 600,
                    mb: 1
                  }}
                >
                  現場名
                </Typography>
                <TextField
                  fullWidth
                  size="small"
                  name="location_name"
                  value={formData.location_name}
                  onChange={handleInputChange}
                  onKeyDown={(e) => handleFilenameKeyDown(e, 'location_name')}
                  onBlur={handleFilenameBlur('location_name')}
                  required
                  disabled={isLoading}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      backgroundColor: 'white'
                    }
                  }}
                />
              </Paper>
            </Box>
          </Box>

          {/* 日付選択セクション */}
          <Box sx={{ display: 'flex', gap: 2, mb: 3, flexDirection: { xs: 'column', sm: 'row' } }}>
            <Box sx={{ flex: 1 }}>
              <Paper
                elevation={0}
                sx={{
                  p: 2,
                  backgroundColor: 'primary.50',
                  border: '2px solid',
                  borderColor: 'primary.200',
                  borderRadius: 2
                }}
              >
                <Typography
                  variant="subtitle2"
                  sx={{
                    color: 'primary.main',
                    fontWeight: 600,
                    mb: 1
                  }}
                >
                  開始日
                </Typography>
                <CalendarPicker
                  value={formData.start_date || ''}
                  onChange={(dateString) => handleDaisyDateChange(dateString, 'start_date')}
                  placeholder="開始日を選択"
                  disabled={isLoading}
                />
              </Paper>
            </Box>
            <Box sx={{ flex: 1 }}>
              <Paper
                elevation={0}
                sx={{
                  p: 2,
                  backgroundColor: 'grey.50',
                  border: '1px solid',
                  borderColor: 'grey.200',
                  borderRadius: 2
                }}
              >
                <Typography
                  variant="subtitle2"
                  sx={{
                    color: 'text.primary',
                    fontWeight: 500,
                    mb: 1
                  }}
                >
                  終了日
                </Typography>
                <CalendarPicker
                  value={formData.end_date || ''}
                  onChange={(dateString) => handleDaisyDateChange(dateString, 'end_date')}
                  placeholder="終了日を選択"
                  disabled={isLoading}
                  minDate={formData.start_date}
                />
              </Paper>
            </Box>
          </Box>

          {/* 説明・タグセクション */}
          <Box sx={{ mb: 3 }}>
            <Box sx={{ mb: 2 }}>
              <Typography
                variant="subtitle2"
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  mb: 1
                }}
              >
                説明
              </Typography>
              <TextField
                fullWidth
                multiline
                rows={3}
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                onKeyDown={(e) => handleKeyDown(e, 'description')}
                onBlur={handleNonFilenameBlur('description')}
                disabled={isLoading}
                placeholder="工事の詳細説明を入力してください"
              />
            </Box>
            <Box>
              <Typography
                variant="subtitle2"
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  mb: 1
                }}
              >
                タグ
              </Typography>
              <TextField
                fullWidth
                name="tags"
                value={formData.tags}
                onChange={handleInputChange}
                onKeyDown={(e) => handleKeyDown(e, 'tags')}
                onBlur={handleNonFilenameBlur('tags')}
                placeholder="カンマ区切りで入力"
                disabled={isLoading}
              />
            </Box>
          </Box>

          {/* ファイル名変更警告とボタン */}
          {hasFilenameChanges && (
            <Alert
              severity="warning"
              sx={{ mb: 3 }}
              icon={<WarningIcon />}
              action={
                <Button
                  color="warning"
                  variant="contained"
                  size="small"
                  onClick={handleFilenameUpdate}
                  disabled={isLoading}
                  startIcon={isLoading ? <CircularProgress size={16} /> : null}
                >
                  {isLoading ? '更新中...' : 'ファイル名を更新'}
                </Button>
              }
            >
              <Typography variant="body2" fontWeight={500}>
                ファイル名の変更が必要です。このボタンを押下して全ての変更をまとめて更新してください。
              </Typography>
            </Alert>
          )}

          {/* 管理ファイル一覧 */}
          {currentProject?.managed_files && currentProject.managed_files.length > 0 && (
            <Box sx={{ mb: 3 }}>
              <Typography
                variant="subtitle2"
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  mb: 2
                }}
              >
                管理ファイル
              </Typography>
              <Paper
                elevation={0}
                sx={{
                  border: '1px solid',
                  borderColor: 'grey.200',
                  borderRadius: 2,
                  backgroundColor: 'grey.50',
                  p: 2
                }}
              >
                {currentProject.managed_files.map((file: ModelsManagedFile, index: number) => {
                  const needsRename = file.current && file.recommended && file.current !== file.recommended;
                  return (
                    <Paper
                      key={index}
                      elevation={0}
                      sx={{
                        mb: index < currentProject.managed_files!.length - 1 ? 2 : 0,
                        p: 2,
                        backgroundColor: 'white',
                        border: '1px solid',
                        borderColor: 'grey.200',
                        borderRadius: 1
                      }}
                    >
                      {/* 現在のファイル */}
                      <Box sx={{ display: 'flex', alignItems: 'center', mb: file.recommended ? 1 : 0 }}>
                        <AttachFileIcon sx={{ mr: 1, color: 'grey.600', fontSize: '1rem' }} />
                        <Typography variant="body2" fontWeight={500} sx={{ mr: 1 }}>
                          現在:
                        </Typography>
                        <Typography variant="body2" sx={{ fontFamily: 'monospace', color: 'grey.700' }}>
                          {file.current}
                        </Typography>
                      </Box>
                      
                      {/* 推奨ファイル */}
                      {file.recommended && (
                        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                          <Box sx={{ display: 'flex', alignItems: 'center' }}>
                            <LightbulbIcon sx={{ mr: 1, color: 'warning.main', fontSize: '1rem' }} />
                            <Typography variant="body2" fontWeight={500} sx={{ mr: 1 }}>
                              推奨:
                            </Typography>
                            <Typography variant="body2" sx={{ fontFamily: 'monospace', color: 'primary.main' }}>
                              {file.recommended}
                            </Typography>
                          </Box>
                          
                          {/* 変更ボタン */}
                          {needsRename && (
                            <Button
                              variant="contained"
                              size="small"
                              onClick={handleRenameFiles}
                              disabled={isRenaming || isLoading || hasFilenameChanges}
                              startIcon={isRenaming ? <CircularProgress size={16} /> : null}
                              title={hasFilenameChanges ? 'ファイル名の変更が必要です。上のボタンで更新してください。' : ''}
                            >
                              {isRenaming ? '変更中...' : '変更'}
                            </Button>
                          )}
                        </Box>
                      )}
                    </Paper>
                  );
                })}
              </Paper>
            </Box>
          )}
          
          {/* 管理ファイルがない場合の表示 */}
          {currentProject && (!currentProject.managed_files || currentProject.managed_files.length === 0) && (
            <Box sx={{ mb: 3 }}>
              <Typography
                variant="subtitle2"
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  mb: 2
                }}
              >
                管理ファイル
              </Typography>
              <Paper
                elevation={0}
                sx={{
                  border: '1px solid',
                  borderColor: 'grey.200',
                  borderRadius: 2,
                  backgroundColor: 'grey.50',
                  p: 3,
                  textAlign: 'center'
                }}
              >
                <AttachFileIcon sx={{ fontSize: '2rem', color: 'grey.400', mb: 1 }} />
                <Typography variant="body2" color="text.secondary">
                  管理ファイルはありません
                </Typography>
              </Paper>
            </Box>
          )}
        </Box>
      </DialogContent>
    </Dialog>
  );
};

export default ProjectDetailModal;