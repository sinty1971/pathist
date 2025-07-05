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
import type { ModelsCompany } from '../api/types.gen';
import "../styles/business-detail-modal.css";

interface CompanyDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  company: ModelsCompany | null;
  onUpdate?: (company: ModelsCompany) => Promise<ModelsCompany>;
  onCompanyUpdate?: (company: ModelsCompany) => void;
}

// フォーム用のデータ型
type CompanyFormData = {
  id?: string;
  name?: string;
  short_name?: string;
  business_type?: string;
  phone?: string;
  email?: string;
  website?: string;
  address?: string;
  description?: string;
  tags?: string;
};

const CompanyDetailModal: React.FC<CompanyDetailModalProps> = ({ 
  isOpen, 
  onClose, 
  company, 
  onUpdate, 
  onCompanyUpdate 
}) => {
  const [formData, setFormData] = useState<CompanyFormData>({
    id: '',
    name: '',
    short_name: '',
    business_type: '',
    phone: '',
    email: '',
    website: '',
    address: '',
    description: '',
    tags: ''
  });
  const [currentCompany, setCurrentCompany] = useState<ModelsCompany | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 会社データが変更されたときにフォームデータを更新
  useEffect(() => {
    if (company) {
      setCurrentCompany(company);
      setFormData({
        id: company.id || '',
        name: company.name || '',
        short_name: company.short_name || '',
        business_type: company.business_type || '',
        phone: company.phone || '',
        email: company.email || '',
        website: company.website || '',
        address: company.address || '',
        description: company.description || '',
        tags: Array.isArray(company.tags) ? company.tags.join(', ') : (company.tags || '')
      });
      setError(null);
    }
  }, [company]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // 更新処理（将来の実装用）
  const handleUpdate = async () => {
    if (!company || !onUpdate) {
      console.log('更新機能は未実装です');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      // 更新処理（未実装）
      const updatedCompany: ModelsCompany = {
        ...company,
        name: formData.name,
        short_name: formData.short_name,
        business_type: formData.business_type,
        phone: formData.phone,
        email: formData.email,
        website: formData.website,
        address: formData.address,
        description: formData.description,
        tags: formData.tags ? formData.tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0) : []
      };

      const savedCompany = await onUpdate(updatedCompany);
      
      if (onCompanyUpdate) {
        onCompanyUpdate(savedCompany);
      }
      
      onClose();
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
        <Typography variant="h6" component="h2">
          会社詳細情報
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
                  name="name"
                  value={formData.name}
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
                  name="short_name"
                  value={formData.short_name}
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
                <TextField
                  fullWidth
                  size="small"
                  name="business_type"
                  value={formData.business_type}
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
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                disabled={isLoading}
                placeholder="会社の詳細説明を入力してください"
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
                disabled={isLoading}
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

          {/* 更新機能未実装の注意書き */}
          <Alert severity="info" sx={{ mb: 3 }}>
            <Typography variant="body2">
              会社情報の編集機能は現在開発中です。詳細情報の表示のみ可能です。
            </Typography>
          </Alert>

          {/* 将来の更新ボタン */}
          {onUpdate && (
            <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 2 }}>
              <Button
                variant="outlined"
                onClick={onClose}
                disabled={isLoading}
              >
                閉じる
              </Button>
              <Button
                variant="contained"
                onClick={handleUpdate}
                disabled={isLoading || true} // 現在は無効化
                startIcon={isLoading ? <CircularProgress size={16} /> : null}
              >
                {isLoading ? '更新中...' : '更新（未実装）'}
              </Button>
            </Box>
          )}
        </Box>
      </DialogContent>
    </Dialog>
  );
};

export default CompanyDetailModal;