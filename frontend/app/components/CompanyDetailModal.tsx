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
import type { ModelsCompany } from '../api/types.gen';
import { putBusinessCompanies } from '../api/sdk.gen';
import { getCategoryName } from '../utils/company';
import "../styles/business-detail-modal.css";

interface CompanyDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  company: ModelsCompany | null;
  onCompanyUpdate?: (company: ModelsCompany) => void;
}

// フォーム用のデータ型
type CompanyFormData = {
  id?: string;
  legalName?: string;
  shortName?: string;
  category?: number;
  phone?: string;
  email?: string;
  website?: string;
  address?: string;
  postalCode?: string;
  tags?: string;
};

const CompanyDetailModal: React.FC<CompanyDetailModalProps> = ({ 
  isOpen, 
  onClose, 
  company, 
  onCompanyUpdate 
}) => {
  const [formData, setFormData] = useState<CompanyFormData>({
    id: '',
    legalName: '',
    shortName: '',
    category: undefined,
    phone: '',
    email: '',
    website: '',
    address: '',
    postalCode: '',
    tags: ''
  });
  const [currentCompany, setCurrentCompany] = useState<ModelsCompany | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);

  // 会社データが変更されたときにフォームデータを更新
  useEffect(() => {
    if (company) {
      setCurrentCompany(company);
      const newFormData: CompanyFormData = {
        id: company.id || '',
        legalName: company.legalName || '',
        shortName: company.shortName || '',
        category: company.category,
        phone: company.phone || '',
        email: company.email || '',
        website: company.website || '',
        address: company.address || '',
        postalCode: company.postalCode || '',
        tags: Array.isArray(company.tags) ? company.tags.join(', ') : (company.tags || '')
      };
      setFormData(newFormData);
      setError(null);
      setIsEditing(false);
      setHasChanges(false);
    }
  }, [company]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
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
          setFormData({
            id: company.id || '',
            legalName: company.legalName || '',
            shortName: company.shortName || '',
            category: company.category,
            phone: company.phone || '',
            email: company.email || '',
            website: company.website || '',
            address: company.address || '',
            postalCode: company.postalCode || '',
            tags: Array.isArray(company.tags) ? company.tags.join(', ') : (company.tags || '')
          });
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
      // タグをカンマ区切りから配列に変換
      const tagsArray = formData.tags 
        ? formData.tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0) 
        : [];

      // 更新用の会社データを作成
      const updatedCompany: ModelsCompany = {
        ...company,
        legalName: formData.legalName,
        shortName: formData.shortName,
        category: formData.category,
        phone: formData.phone,
        email: formData.email,
        website: formData.website,
        address: formData.address,
        tags: tagsArray
      };

      // API呼び出し
      const response = await putBusinessCompanies({
        body: updatedCompany
      });

      if (response.data) {
        // 成功時の処理
        setCurrentCompany(response.data);
        setIsEditing(false);
        setHasChanges(false);
        
        if (onCompanyUpdate) {
          onCompanyUpdate(response.data);
        }
        
        // 成功メッセージ（一時的に表示）
        setError(null);
        
        // モーダルを閉じる（オプション）
        // onClose();
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
  const displayTags = currentCompany?.tags ? 
    (Array.isArray(currentCompany.tags) ? currentCompany.tags : [currentCompany.tags]) : [];

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
          {/* 基本情報セクション */}
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

          {/* 業種・ID */}
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
                      const value = e.target.value === '' ? undefined : parseInt(e.target.value, 10);
                      setFormData(prev => ({ ...prev, category: value }));
                      setHasChanges(true);
                    }}
                    disabled={isLoading}
                    select
                    SelectProps={{
                      native: true,
                    }}
                    sx={{
                      '& .MuiOutlinedInput-root': {
                        backgroundColor: 'white'
                      }
                    }}
                  >
                    <option value="">選択してください</option>
                    <option value="0">特別</option>
                    <option value="1">下請会社</option>
                    <option value="2">築炉会社</option>
                    <option value="3">一人親方</option>
                    <option value="4">元請け</option>
                    <option value="5">リース会社</option>
                    <option value="6">販売会社</option>
                    <option value="7">販売会社</option>
                    <option value="8">求人会社</option>
                    <option value="9">その他</option>
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
                value={formData.tags}
                onChange={handleInputChange}
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
                value={formData.tags}
                onChange={handleInputChange}
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