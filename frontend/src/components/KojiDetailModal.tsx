'use client';

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
import type { ModelsKoji, ModelsTimestamp, ModelsManagedFile } from '@/api/types.gen';
import { putBusinessKojiesAssistFiles, getKojiesByPath } from '@/api/sdk.gen';
import { CalendarPicker } from './CalendarPicker';

interface KojiDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  koji: ModelsKoji | null;
  onUpdate: (koji: ModelsKoji) => Promise<ModelsKoji>;
  onKojiUpdate?: (koji: ModelsKoji) => void;
}

// フォーム用のデータ型（日付を文字列として扱う、ステータスは除外）
type KojiFormData = Omit<ModelsKoji, 'startDate' | 'endDate' | 'tags' | 'status'> & {
  startDate?: string;
  endDate?: string;
  tags?: string;
};

// CalendarPickerコンポーネントをインポート済み

const KojiDetailModal: React.FC<KojiDetailModalProps> = ({ isOpen, onClose, koji, onUpdate, onKojiUpdate }) => {
  const [formData, setFormData] = useState<KojiFormData>({
    id: '',
    companyName: '',
    locationName: '',
    description: '',
    tags: '',
    startDate: '',
    endDate: ''
  });
  const [currentKoji, setCurrentKoji] = useState<ModelsKoji | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isRenaming, setIsRenaming] = useState(false);
  const [hasFilenameChanges, setHasFilenameChanges] = useState(false);
  const [initialFilenameData, setInitialFilenameData] = useState<{
    startDate: string;
    companyName: string;
    locationName: string;
  }>({ startDate: '', companyName: '', locationName: '' });

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

  // 工事が変更されたときにフォームデータを更新
  useEffect(() => {
    if (koji) {
      const startDate = extractDateString(koji.startDate);
      const endDate = extractDateString(koji.endDate);
      
      const companyName = koji.companyName || '';
      const locationName = koji.locationName || '';
      
      setCurrentKoji(koji);
      setFormData({
        id: koji.id || '',
        companyName: companyName,
        locationName: locationName,
        description: koji.description || '',
        tags: Array.isArray(koji.tags) ? koji.tags.join(', ') : (koji.tags || ''),
        startDate: startDate,
        endDate: endDate
      });
      
      // 初期ファイル名関連データを保存
      setInitialFilenameData({
        startDate: startDate,
        companyName: companyName,
        locationName: locationName
      });
      
      setHasFilenameChanges(false);
      setError(null);
    }
  }, [koji]);

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
  const generateRecommendedFileName = (originalFileName: string, formData: KojiFormData): string => {
    if (!formData.startDate || !formData.companyName || !formData.locationName) {
      return originalFileName;
    }

    // 日付をYYYY-MMDD形式に変換
    const datePart = formData.startDate.replace(/-/g, '').substring(0, 8); // YYYYMMDD
    const formattedDate = `${datePart.substring(0, 4)}-${datePart.substring(4, 8)}`; // YYYY-MMDD
    
    // 新しいプレフィックスを作成
    const newPrefix = `${formattedDate} ${formData.companyName} ${formData.locationName}`;
    
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
  const updateManagedFileRecommendations = (formData: KojiFormData) => {
    if (!currentKoji?.managed_files) return;

    const updatedManagedFiles = currentKoji.managed_files.map(file => {
      if (file.current) {
        const recommendedName = generateRecommendedFileName(file.current, formData);
        return {
          ...file,
          recommended: recommendedName
        };
      }
      return file;
    });

    setCurrentKoji(prev => prev ? {
      ...prev,
      managed_files: updatedManagedFiles
    } : null);
  };

  // ファイル名関連の変更をチェックする関数
  const checkFilenameChanges = (currentFormData: KojiFormData) => {
    const hasChanges = 
      currentFormData.startDate !== initialFilenameData.startDate ||
      currentFormData.companyName !== initialFilenameData.companyName ||
      currentFormData.locationName !== initialFilenameData.locationName;
    
    setHasFilenameChanges(hasChanges);
    
    // ファイル名関連の変更がある場合、管理ファイルの推奨名を更新
    if (hasChanges) {
      updateManagedFileRecommendations(currentFormData);
    }
  };

  // フィールド変更完了時の処理（指定されたフォームデータを使用）
  const handleFieldUpdateWithData = async (useFormData: KojiFormData) => {
    if (!koji) return;

    setIsLoading(true);
    setError(null);

    try {
      // 指定されたフォームデータをModelsKoji形式に変換
      const updatedKoji: ModelsKoji = {
        ...koji,
        companyName: useFormData.companyName,
        locationName: useFormData.locationName,
        description: useFormData.description,
        tags: useFormData.tags ? useFormData.tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0) : [],
        startDate: useFormData.startDate ? { 'time.Time': `${useFormData.startDate}T00:00:00+09:00` } as ModelsTimestamp : undefined,
        endDate: useFormData.endDate ? { 'time.Time': `${useFormData.endDate}T23:59:59+09:00` } as ModelsTimestamp : undefined
      };
      

      // 更新前のフォルダー名を保存
      const originalFolderName = koji.name;
      
      // 更新して更新された工事データを取得
      const savedKoji = await onUpdate(updatedKoji);
      
      // フォルダー名が変更されたかチェック
      const folderNameChanged = originalFolderName && savedKoji.name && originalFolderName !== savedKoji.name;
      
      // ローカルステートを更新（管理ファイル表示用）
      setCurrentKoji(savedKoji);
      
      // 親コンポーネントに最新データを渡す
      if (onKojiUpdate) {
        onKojiUpdate(savedKoji);
      }
      
      // フォームデータも最新データで同期
      const startDate = extractDateString(savedKoji.startDate);
      const endDate = extractDateString(savedKoji.endDate);
      
      setFormData({
        id: savedKoji.id || '',
        companyName: savedKoji.companyName || '',
        locationName: savedKoji.locationName || '',
        description: savedKoji.description || '',
        tags: Array.isArray(savedKoji.tags) ? savedKoji.tags.join(', ') : (savedKoji.tags || ''),
        startDate: startDate,
        endDate: endDate
      });
      
      // フォルダー名が変更された場合はモーダルを閉じる
      if (folderNameChanged) {
        setTimeout(() => {
          onClose();
        }, 100);
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
    if (fieldName === 'companyName' || fieldName === 'locationName') {
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
    if (fieldName === 'companyName' || fieldName === 'locationName') {
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
      if (fieldName === 'companyName' || fieldName === 'locationName') {
        checkFilenameChanges(formData);
      }
    }
  };

  // DaisyUI日付ピッカー用のハンドラー
  const handleDaisyDateChange = (dateString: string, fieldName: 'startDate' | 'endDate') => {
    // フォームデータを更新
    const newFormData = {
      ...formData,
      [fieldName]: dateString
    };
    setFormData(newFormData);
    
    // 開始日の場合はファイル名関連なので即座に更新しない
    if (fieldName === 'startDate') {
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
      startDate: formData.startDate || '',
      companyName: formData.companyName || '',
      locationName: formData.locationName || ''
    });
  };

  const handleRenameFiles = async () => {
    if (!currentKoji || !currentKoji.managed_files) return;

    setIsRenaming(true);
    setError(null);

    try {
      // 現在のファイル名のリストを取得
      const currentFiles = currentKoji.managed_files
        .filter(file => file.current && file.recommended && file.current !== file.recommended)
        .map(file => file.current as string);

      if (currentFiles.length === 0) {
        setError('変更対象のファイルがありません');
        return;
      }

      // API呼び出し
      const response = await putBusinessKojiesAssistFiles({
        body: {
          koji: currentKoji,
          currents: currentFiles
        }
      });

      if (response.data) {
        // APIレスポンスが更新された工事データ（または文字列配列）を含む
        if (typeof response.data === 'object' && response.data.id) {
          // 工事データが返された場合
          const updatedKoji = response.data;
          
          // ローカルステートも更新
          setCurrentKoji(updatedKoji);
          
          // 工事データを更新
          if (onKojiUpdate) {
            onKojiUpdate(updatedKoji);
          }
          
          // フォームデータも更新
          const startDate = extractDateString(updatedKoji.startDate);
          const endDate = extractDateString(updatedKoji.endDate);
          
          setFormData({
            id: updatedKoji.id || '',
            companyName: updatedKoji.companyName || '',
            locationName: updatedKoji.locationName || '',
            description: updatedKoji.description || '',
            tags: Array.isArray(updatedKoji.tags) ? updatedKoji.tags.join(', ') : (updatedKoji.tags || ''),
            startDate: startDate,
            endDate: endDate
          });
        } else {
          // 文字列配列が返された場合（後方互換性）
          // 工事データを再取得
          if (currentKoji.name) {
            const updatedKojiResponse = await getKojiesByPath({
              query: {
                path: currentKoji.name
              }
            });
            
            if (updatedKojiResponse.data && onKojiUpdate) {
              // ローカルステートも更新
              setCurrentKoji(updatedKojiResponse.data);
              // 工事データを更新
              onKojiUpdate(updatedKojiResponse.data);
              // フォームデータも更新
              const startDate = extractDateString(updatedKojiResponse.data.startDate);
              const endDate = extractDateString(updatedKojiResponse.data.endDate);
              
              setFormData({
                id: updatedKojiResponse.data.id || '',
                companyName: updatedKojiResponse.data.companyName || '',
                locationName: updatedKojiResponse.data.locationName || '',
                description: updatedKojiResponse.data.description || '',
                tags: Array.isArray(updatedKojiResponse.data.tags) ? updatedKojiResponse.data.tags.join(', ') : (updatedKojiResponse.data.tags || ''),
                startDate: startDate,
                endDate: endDate
              });
            }
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
                  name="companyName"
                  value={formData.companyName}
                  onChange={handleInputChange}
                  onKeyDown={(e) => handleFilenameKeyDown(e, 'companyName')}
                  onBlur={handleFilenameBlur('companyName')}
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
                  name="locationName"
                  value={formData.locationName}
                  onChange={handleInputChange}
                  onKeyDown={(e) => handleFilenameKeyDown(e, 'locationName')}
                  onBlur={handleFilenameBlur('locationName')}
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
                  value={formData.startDate || ''}
                  onChange={(dateString) => handleDaisyDateChange(dateString, 'startDate')}
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
                  value={formData.endDate || ''}
                  onChange={(dateString) => handleDaisyDateChange(dateString, 'endDate')}
                  placeholder="終了日を選択"
                  disabled={isLoading}
                  minDate={formData.startDate}
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
          {currentKoji?.managed_files && currentKoji.managed_files.length > 0 && (
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
                {currentKoji.managed_files.map((file: ModelsManagedFile, index: number) => {
                  const needsRename = file.current && file.recommended && file.current !== file.recommended;
                  return (
                    <Paper
                      key={index}
                      elevation={0}
                      sx={{
                        mb: index < currentKoji.managed_files!.length - 1 ? 2 : 0,
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
          {currentKoji && (!currentKoji.managed_files || currentKoji.managed_files.length === 0) && (
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

export default KojiDetailModal;