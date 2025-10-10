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
  CircularProgress,
  Chip
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import BusinessIcon from '@mui/icons-material/Business';
import PhoneIcon from '@mui/icons-material/Phone';
import EmailIcon from '@mui/icons-material/Email';
import LanguageIcon from '@mui/icons-material/Language';
import TagIcon from '@mui/icons-material/Tag';
import EditIcon from '@mui/icons-material/Edit';
import SaveIcon from '@mui/icons-material/Save';
import CancelIcon from '@mui/icons-material/Cancel';
import type { Company, CompanyCategoryInfo } from '@/api/types.gen';
import { putCompany } from '@/api/sdk.gen';
import { getCategoryName, loadCategories } from '@/utils/company';
import "../styles/business-detail-modal.css";

interface CompanyDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  company: Company | null;
  onCompanyUpdate?: (company: Company) => void;
}

// tags変換ユーティリティ関数
const tagsToString = (tags?: Array<string> | null): string => {
  if (!tags) {
    return '';
  }
  return Array.isArray(tags) ? tags.join(', ') : '';
};

const stringToTags = (str: string): string[] => {
  if (!str.trim()) {
    return [];
  }
  return str.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0);
};

// デフォルト会社データ
const defaultCompanyData: Company = {
  id: '',
  legalName: '',
  shortName: '',
  category: '',
  phone: '',
  email: '',
  website: '',
  address: '',
  postalCode: '',
  tags: []
};

const normalizeCompany = (data: Company): Company => ({
  ...data,
  id: data.id || '',
  legalName: data.legalName || '',
  shortName: data.shortName || '',
  category:
    data.category !== undefined &&
    data.category !== null &&
    String(data.category).trim() !== ''
      ? String(data.category)
      : '',
  phone: data.phone || '',
  email: data.email || '',
  website: data.website || '',
  address: data.address || '',
  postalCode: data.postalCode || '',
  tags: Array.isArray(data.tags) ? data.tags : data.tags ? [data.tags] : [],
});
// 会社詳細情報モーダル
const CompanyDetailModal: React.FC<CompanyDetailModalProps> = ({ 
  isOpen, 
  onClose, 
  company, 
  onCompanyUpdate 
}) => {
  const [formData, setFormData] = useState<Company>(defaultCompanyData);
  const [currentCompany, setCurrentCompany] = useState<Company | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);
  const [categories, setCategories] = useState<CompanyCategoryInfo[]>([]);
  const [categoriesLoading, setCategoriesLoading] = useState(false);

  // カテゴリーデータを初期化
  useEffect(() => {
    const fetchCategories = async () => {
      setCategoriesLoading(true);
      try {
        const allCategories = await loadCategories();
        setCategories(allCategories);
      } catch (error) {
        console.error('Failed to load categories:', error);
      } finally {
        setCategoriesLoading(false);
      }
    };

    if (isOpen) {
      fetchCategories();
    }
  }, [isOpen]);

  // 会社データが変更されたときにフォームデータを更新
  useEffect(() => {
    if (company) {
      const normalized = normalizeCompany(company);
      setCurrentCompany(normalized);
      setFormData(normalized);
      setError(null);
      setIsEditing(false);
      setHasChanges(false);
    }
  }, [company]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData((prev: Company) => ({
      ...prev,
      [name]: value
    }));
    setHasChanges(true);
  };

  // 編集モード切り替え
  const handleEditToggle = () => {
    if (isEditing && hasChanges) {
      // 編集中で変更がある場合は確認
      if (confirm('変更を破棄しますか？')) {
        // 元のデータに戻す
        if (company) {
          setFormData(normalizeCompany(company));
        }
        setIsEditing(false);
        setHasChanges(false);
      }
    } else {
      setIsEditing(!isEditing);
      setHasChanges(false);
    }
  };

  // 更新処理
  const handleUpdate = async () => {
    if (!company) {
      setError('会社データがありません');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      // 更新用の会社データを作成
      const baseCompany = normalizeCompany(company);
      const updatedCompany: Company = {
        ...baseCompany,
        ...formData,
        category:
          formData.category !== undefined &&
          formData.category !== null
            ? String(formData.category)
            : '',
        tags: formData.tags || []
      };


      // API呼び出し
      const response = await putCompany({
        body: updatedCompany
      });

      if (response.data) {
        // ファイル名に影響するフィールドが変更されたかチェック
        const normalizedResponse = normalizeCompany(response.data);
        const isIdentityChanged = 
          baseCompany.shortName !== normalizedResponse.shortName || 
          baseCompany.category !== normalizedResponse.category;
        
        // 成功時の処理
        setCurrentCompany(normalizedResponse);
        setFormData(normalizedResponse);
        setIsEditing(false);
        setHasChanges(false);
        
        if (onCompanyUpdate) {
          onCompanyUpdate(normalizedResponse);
        }
        
        // 成功メッセージ（一時的に表示）
        setError(null);
        
        // ファイル名が変更された場合はモーダルを閉じる
        if (isIdentityChanged) {
          onClose();
        }
      } else {
        throw new Error('更新レスポンスが無効です');
      }
    } catch (err) {
      console.error('Error updating company:', err);
      setError(err instanceof Error ? err.message : '会社情報の更新に失敗しました');
    } finally {
      setIsLoading(false);
    }
  };

  // タグの配列を表示用に変換
  const displayTags = currentCompany?.tags && Array.isArray(currentCompany.tags)
    ? currentCompany.tags
    : [];

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
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Typography variant="h6" component="h2">
            会社詳細情報
          </Typography>
          {isEditing && (
            <Chip 
              label="編集中" 
              color="warning" 
              size="small"
              icon={<EditIcon />}
            />
          )}
        </Box>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            variant={isEditing ? "outlined" : "contained"}
            size="small"
            onClick={handleEditToggle}
            startIcon={isEditing ? <CancelIcon /> : <EditIcon />}
            color={isEditing ? "secondary" : "primary"}
          >
            {isEditing ? 'キャンセル' : '編集'}
          </Button>
          <IconButton
            onClick={onClose}
            size="small"
            sx={{
              color: 'grey.500'
            }}
          >
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent dividers>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Box sx={{ mt: 1 }}>
          {/* ID・会社名 */}
          <Box sx={{ display: 'flex', gap: 2, mb: 3, flexDirection: { xs: 'column', sm: 'row' } }}>
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
                  ID
                </Typography>
                <Typography
                  variant="body2"
                  sx={{
                    fontFamily: 'monospace',
                    padding: '8px 12px',
                    backgroundColor: 'white',
                    border: '1px solid',
                    borderColor: 'grey.300',
                    borderRadius: 1,
                    color: 'grey.700'
                  }}
                >
                  {formData.id || 'ID未設定'}
                </Typography>
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
                    mb: 1,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1
                  }}
                >
                  <BusinessIcon fontSize="small" />
                  会社名
                </Typography>
                <TextField
                  fullWidth
                  size="small"
                  name="legalName"
                  value={formData.legalName}
                  onChange={handleInputChange}
                  disabled={isLoading || !isEditing}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      backgroundColor: isEditing ? 'white' : 'grey.50'
                    }
                  }}
                />
              </Paper>
            </Box>
          </Box>

          {/* 業種・略称 */}
          <Box sx={{ display: 'flex', gap: 2, mb: 3, flexDirection: { xs: 'column', sm: 'row' } }}>
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
                  業種
                </Typography>
                {!isEditing ? (
                  <Typography variant="body1" sx={{ mt: 1 }}>
                    {getCategoryName(formData.category)}
                  </Typography>
                ) : (
                  <TextField
                    fullWidth
                    size="small"
                    name="category"
                    value={formData.category || ''}
                    onChange={(e) => {
                      const value = e.target.value;
                      setFormData((prev: Company) => ({
                        ...prev,
                        category: value
                      }));
                      setHasChanges(true);
                    }}
                    disabled={isLoading || categoriesLoading}
                    select
                    slotProps={{
                      select: {
                        native: true,
                      }
                    }}
                    sx={{
                      '& .MuiOutlinedInput-root': {
                        backgroundColor: 'white'
                      }
                    }}
                  >
                    <option value="">選択してください</option>
                    {categoriesLoading ? (
                      <option value="" disabled>カテゴリーを読み込み中...</option>
                    ) : (
                      categories.map((cat) => (
                        <option key={cat.code} value={cat.code}>
                          {cat.label}
                        </option>
                      ))
                    )}
                  </TextField>
                )}
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
                  略称
                </Typography>
                <TextField
                  fullWidth
                  size="small"
                  name="shortName"
                  value={formData.shortName}
                  onChange={handleInputChange}
                  disabled={isLoading || !isEditing}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      backgroundColor: isEditing ? 'white' : 'grey.50'
                    }
                  }}
                />
              </Paper>
            </Box>
          </Box>

          {/* 連絡先情報 */}
          <Box sx={{ display: 'flex', gap: 2, mb: 3, flexDirection: { xs: 'column', sm: 'row' } }}>
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
                    mb: 1,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1
                  }}
                >
                  <PhoneIcon fontSize="small" />
                  電話番号
                </Typography>
                <TextField
                  fullWidth
                  size="small"
                  name="phone"
                  value={formData.phone}
                  onChange={handleInputChange}
                  disabled={isLoading || !isEditing}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      backgroundColor: isEditing ? 'white' : 'grey.50'
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
                    mb: 1,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1
                  }}
                >
                  <EmailIcon fontSize="small" />
                  メールアドレス
                </Typography>
                <TextField
                  fullWidth
                  size="small"
                  name="email"
                  type="email"
                  value={formData.email}
                  onChange={handleInputChange}
                  disabled={isLoading || !isEditing}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      backgroundColor: isEditing ? 'white' : 'grey.50'
                    }
                  }}
                />
              </Paper>
            </Box>
          </Box>

          {/* ウェブサイト・住所 */}
          <Box sx={{ display: 'flex', gap: 2, mb: 3, flexDirection: { xs: 'column', sm: 'row' } }}>
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
                    mb: 1,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1
                  }}
                >
                  <LanguageIcon fontSize="small" />
                  ウェブサイト
                </Typography>
                <TextField
                  fullWidth
                  size="small"
                  name="website"
                  type="url"
                  value={formData.website}
                  onChange={handleInputChange}
                  disabled={isLoading || !isEditing}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      backgroundColor: isEditing ? 'white' : 'grey.50'
                    }
                  }}
                />
              </Paper>
            </Box>
          </Box>

          {/* 住所 */}
          <Box sx={{ mb: 3 }}>
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
                住所
              </Typography>
              <TextField
                fullWidth
                size="small"
                name="address"
                value={formData.address}
                onChange={handleInputChange}
                disabled={isLoading}
                sx={{
                  '& .MuiOutlinedInput-root': {
                    backgroundColor: 'white'
                  }
                }}
              />
            </Paper>
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
                name="tags"
                value={tagsToString(formData.tags)}
                onChange={(e) => {
                  const tagsArray = stringToTags(e.target.value);
                  setFormData((prev: Company) => ({ ...prev, tags: tagsArray }));
                  setHasChanges(true);
                }}
                disabled={isLoading || !isEditing}
                placeholder="会社の詳細説明を入力してください"
                sx={{
                  '& .MuiOutlinedInput-root': {
                    backgroundColor: isEditing ? 'white' : 'grey.50'
                  }
                }}
              />
            </Box>
            <Box>
              <Typography
                variant="subtitle2"
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  mb: 1,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 1
                }}
              >
                <TagIcon fontSize="small" />
                タグ
              </Typography>
              <TextField
                fullWidth
                name="tags"
                value={tagsToString(formData.tags)}
                onChange={(e) => {
                  const tagsArray = stringToTags(e.target.value);
                  setFormData((prev: Company) => ({ ...prev, tags: tagsArray }));
                  setHasChanges(true);
                }}
                placeholder="カンマ区切りで入力"
                disabled={isLoading || !isEditing}
                sx={{
                  '& .MuiOutlinedInput-root': {
                    backgroundColor: isEditing ? 'white' : 'grey.50'
                  }
                }}
              />
              
              {/* 現在のタグ表示 */}
              {displayTags.length > 0 && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                    現在のタグ:
                  </Typography>
                  <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                    {displayTags.map((tag, index) => (
                      <Chip 
                        key={index} 
                        label={tag} 
                        size="small" 
                        color="primary" 
                        variant="outlined"
                      />
                    ))}
                  </Box>
                </Box>
              )}
            </Box>
          </Box>

          {/* 編集時の保存ボタン */}
          {isEditing && (
            <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 2, mt: 3 }}>
              <Button
                variant="outlined"
                onClick={handleEditToggle}
                disabled={isLoading}
                startIcon={<CancelIcon />}
              >
                キャンセル
              </Button>
              <Button
                variant="contained"
                onClick={handleUpdate}
                disabled={isLoading || !hasChanges}
                startIcon={isLoading ? <CircularProgress size={16} /> : <SaveIcon />}
                color="primary"
              >
                {isLoading ? '保存中...' : '保存'}
              </Button>
            </Box>
          )}

          {/* 更新成功メッセージ */}
          {!isEditing && hasChanges === false && currentCompany && (
            <Alert severity="success" sx={{ mt: 2 }}>
              会社情報が正常に更新されました。
            </Alert>
          )}
        </Box>
      </DialogContent>
    </Dialog>
  );
};

export default CompanyDetailModal;
