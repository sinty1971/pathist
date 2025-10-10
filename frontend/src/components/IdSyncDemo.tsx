'use client';

/**
 * ID同期システムのデモンストレーション用コンポーネント
 * 
 * このコンポーネントは以下のデモを提供します：
 * 1. 工事ID生成とバリデーション
 * 2. フルパス → Len7 ID 変換
 * 3. 自動同期機能
 * 4. 半自動同期機能
 * 5. 一括変換機能
 */

import React, { useState, useCallback } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  TextField,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  Alert,
  CircularProgress,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Grid,
  Card,
  CardContent,
  CardActions,
  Divider
} from '@mui/material';
import {
  ExpandMore,
  PlayArrow,
  Sync,
  AutoFixHigh,
  List,
  Transform,
  CheckCircle,
  Error
} from '@mui/icons-material';
import { useIdSync, IdComponents } from '@/utils/idSync';
import { useAutoIdSync, useBulkIdSync, usePathIdSync } from '@/hooks/useAutoIdSync';

interface DemoResult {
  id: string;
  isValid: boolean;
  error?: string;
  timestamp: Date;
  method: string;
}

interface PathConversion {
  fullPath: string;
  len7Id: string;
  reduction: number;
}

export const IdSyncDemo: React.FC = () => {
  // 工事ID生成デモ用の状態
  const [kojiComponents, setKojiComponents] = useState<IdComponents>({
    startDate: new Date(2025, 5, 18), // 2025-06-18
    companyName: '豊田築炉',
    locationName: '名和工場'
  });

  // パス変換デモ用の状態
  const [samplePaths] = useState([
    "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場",
    "豊田築炉/2-工事/2025-0615 豊田築炉 東海工場",
    "豊田築炉/2-工事/2025-0620 豊田築炉 刈谷工場",
    "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場/工事.xlsx",
    "豊田築炉/2-工事/2025-0618 豊田築炉 名和工場/図面.pdf"
  ]);

  // 結果の状態
  const [demoResults, setDemoResults] = useState<DemoResult[]>([]);
  const [pathConversions, setPathConversions] = useState<PathConversion[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // Hooks
  const { generateKojiId, generatePathId, validateId, bulkConvertAndValidate } = useIdSync();
  
  const autoSyncResult = useAutoIdSync(
    'DEMO_ID',
    kojiComponents,
    {
      requireConfirmation: false,
      onSuccess: (newId) => {
        addResult('自動同期', newId, true, '自動同期が完了しました');
      },
      onError: (error) => {
        addResult('自動同期', '', false, error);
      }
    }
  );

  const { convertToLen7, isConverting } = usePathIdSync(samplePaths, {
    autoConvert: true,
    onMappingReady: (mapping) => {
      const conversions = Array.from(mapping.entries()).map(([fullPath, len7Id]) => ({
        fullPath,
        len7Id,
        reduction: Math.round((1 - len7Id.length / fullPath.length) * 100)
      }));
      setPathConversions(conversions);
    }
  });

  // 結果追加ヘルパー
  const addResult = useCallback((method: string, id: string, isValid: boolean, error?: string) => {
    const result: DemoResult = {
      id,
      isValid,
      error,
      timestamp: new Date(),
      method
    };
    setDemoResults(prev => [result, ...prev].slice(0, 10)); // 最新10件のみ保持
  }, []);

  // 工事ID生成デモ
  const runKojiIdDemo = useCallback(async () => {
    setIsLoading(true);
    try {
      const id = generateKojiId(kojiComponents);
      addResult('工事ID生成', id, true);

      // バリデーションも実行
      const validation = await validateId(id, kojiComponents);
      if (validation.isValid) {
        addResult('ID検証', id, true, '検証成功');
      } else {
        addResult('ID検証', id, false, validation.error || '検証失敗');
      }
    } catch (error) {
      addResult('工事ID生成', '', false, error instanceof Error ? error.message : '不明なエラー');
    } finally {
      setIsLoading(false);
    }
  }, [kojiComponents, generateKojiId, validateId, addResult]);

  // パスID変換デモ
  const runPathIdDemo = useCallback(() => {
    setIsLoading(true);
    try {
      const testPath = samplePaths[0];
      const len7Id = generatePathId(testPath);
      addResult('パスID変換', len7Id, true, `${testPath} → ${len7Id}`);
    } catch (error) {
      addResult('パスID変換', '', false, error instanceof Error ? error.message : '不明なエラー');
    } finally {
      setIsLoading(false);
    }
  }, [samplePaths, generatePathId, addResult]);

  // 一括変換デモ
  const runBulkConversionDemo = useCallback(async () => {
    setIsLoading(true);
    try {
      const testItems = [
        { id: 'TEST1', components: kojiComponents },
        { id: 'TEST2', components: { ...kojiComponents, locationName: '東海工場' } },
        { id: 'TEST3', components: { ...kojiComponents, locationName: '刈谷工場' } }
      ];

      const results = await bulkConvertAndValidate(testItems);
      const successCount = results.filter(r => r.isValid).length;
      const needsUpdateCount = results.filter(r => r.needsUpdate).length;

      addResult(
        '一括変換', 
        `${successCount}/${results.length}件成功`, 
        true,
        `${needsUpdateCount}件が更新を必要としています`
      );
    } catch (error) {
      addResult('一括変換', '', false, error instanceof Error ? error.message : '不明なエラー');
    } finally {
      setIsLoading(false);
    }
  }, [kojiComponents, bulkConvertAndValidate, addResult]);

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        🔄 ID同期システム デモ
      </Typography>
      
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
        バックエンドとフロントエンドの自動ID同期機能をテストできます
      </Typography>

      {/* 設定セクション */}
      <Paper sx={{ p: 2, mb: 2 }}>
        <Typography variant="h6" gutterBottom>
          ⚙️ 工事データ設定
        </Typography>
        
        <Grid container spacing={2}>
          <Grid item xs={12} md={4}>
            <TextField
              label="開始日"
              type="date"
              value={kojiComponents.startDate.toISOString().slice(0, 10)}
              onChange={(e) => setKojiComponents(prev => ({
                ...prev,
                startDate: new Date(e.target.value)
              }))}
              fullWidth
              InputLabelProps={{ shrink: true }}
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <TextField
              label="会社名"
              value={kojiComponents.companyName}
              onChange={(e) => setKojiComponents(prev => ({
                ...prev,
                companyName: e.target.value
              }))}
              fullWidth
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <TextField
              label="場所名"
              value={kojiComponents.locationName}
              onChange={(e) => setKojiComponents(prev => ({
                ...prev,
                locationName: e.target.value
              }))}
              fullWidth
            />
          </Grid>
        </Grid>
      </Paper>

      {/* デモ実行ボタン */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="primary">
                <PlayArrow /> 工事ID生成
              </Typography>
              <Typography variant="body2" color="text.secondary">
                設定データから工事IDを生成し、検証します
              </Typography>
            </CardContent>
            <CardActions>
              <Button 
                onClick={runKojiIdDemo} 
                disabled={isLoading}
                startIcon={<PlayArrow />}
                variant="contained"
                fullWidth
              >
                実行
              </Button>
            </CardActions>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="secondary">
                <Transform /> パスID変換
              </Typography>
              <Typography variant="body2" color="text.secondary">
                フルパスからLen7 IDに変換します
              </Typography>
            </CardContent>
            <CardActions>
              <Button 
                onClick={runPathIdDemo} 
                disabled={isLoading}
                startIcon={<Transform />}
                variant="contained"
                color="secondary"
                fullWidth
              >
                実行
              </Button>
            </CardActions>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="success.main">
                <Sync /> 自動同期
              </Typography>
              <Typography variant="body2" color="text.secondary">
                ID不整合を自動検出・修正します
              </Typography>
            </CardContent>
            <CardActions>
              <Button 
                onClick={autoSyncResult.sync} 
                disabled={autoSyncResult.isSyncing}
                startIcon={autoSyncResult.isSyncing ? <CircularProgress size={16} /> : <Sync />}
                variant="contained"
                color="success"
                fullWidth
              >
                {autoSyncResult.isSyncing ? '同期中...' : '実行'}
              </Button>
            </CardActions>
          </Card>
        </Grid>

        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Typography variant="h6" color="warning.main">
                <List /> 一括変換
              </Typography>
              <Typography variant="body2" color="text.secondary">
                複数データを一括で変換・検証します
              </Typography>
            </CardContent>
            <CardActions>
              <Button 
                onClick={runBulkConversionDemo} 
                disabled={isLoading}
                startIcon={<List />}
                variant="contained"
                color="warning"
                fullWidth
              >
                実行
              </Button>
            </CardActions>
          </Card>
        </Grid>
      </Grid>

      {/* 現在の同期状態 */}
      <Paper sx={{ p: 2, mb: 2 }}>
        <Typography variant="h6" gutterBottom>
          📊 現在の同期状態
        </Typography>
        
        <Grid container spacing={2}>
          <Grid item xs={12} md={6}>
            <Alert severity="info">
              <Typography variant="body2">
                <strong>現在のID:</strong> {autoSyncResult.currentId}
              </Typography>
              <Typography variant="body2">
                <strong>同期状態:</strong> {autoSyncResult.isSyncing ? '同期中' : '待機中'}
              </Typography>
              {autoSyncResult.lastSyncTime && (
                <Typography variant="body2">
                  <strong>最終同期:</strong> {autoSyncResult.lastSyncTime.toLocaleString()}
                </Typography>
              )}
            </Alert>
          </Grid>
          
          <Grid item xs={12} md={6}>
            {autoSyncResult.syncError ? (
              <Alert severity="error">
                <Typography variant="body2">
                  <strong>同期エラー:</strong> {autoSyncResult.syncError}
                </Typography>
              </Alert>
            ) : (
              <Alert severity="success">
                <Typography variant="body2">
                  同期システムは正常に動作しています
                </Typography>
              </Alert>
            )}
          </Grid>
        </Grid>
      </Paper>

      {/* 結果表示 */}
      <Accordion defaultExpanded>
        <AccordionSummary expandIcon={<ExpandMore />}>
          <Typography variant="h6">
            📋 実行結果 ({demoResults.length}件)
          </Typography>
        </AccordionSummary>
        <AccordionDetails>
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>時刻</TableCell>
                  <TableCell>メソッド</TableCell>
                  <TableCell>ID</TableCell>
                  <TableCell>状態</TableCell>
                  <TableCell>詳細</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {demoResults.map((result, index) => (
                  <TableRow key={index}>
                    <TableCell>
                      {result.timestamp.toLocaleTimeString()}
                    </TableCell>
                    <TableCell>
                      <Chip label={result.method} size="small" />
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                        {result.id || 'N/A'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      {result.isValid ? (
                        <CheckCircle color="success" />
                      ) : (
                        <Error color="error" />
                      )}
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" color="text.secondary">
                        {result.error || '成功'}
                      </Typography>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
          
          {demoResults.length === 0 && (
            <Alert severity="info">
              デモを実行すると結果がここに表示されます
            </Alert>
          )}
        </AccordionDetails>
      </Accordion>

      {/* パス変換テーブル */}
      <Accordion>
        <AccordionSummary expandIcon={<ExpandMore />}>
          <Typography variant="h6">
            🔄 パス変換テーブル
          </Typography>
          {isConverting && <CircularProgress size={20} sx={{ ml: 1 }} />}
        </AccordionSummary>
        <AccordionDetails>
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>フルパス</TableCell>
                  <TableCell>Len7 ID</TableCell>
                  <TableCell>短縮率</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {pathConversions.map((conversion, index) => (
                  <TableRow key={index}>
                    <TableCell>
                      <Typography variant="body2" sx={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>
                        {conversion.fullPath}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>
                        {conversion.len7Id}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip 
                        label={`${conversion.reduction}%`} 
                        color={conversion.reduction > 80 ? 'success' : 'warning'}
                        size="small"
                      />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>

          {pathConversions.length === 0 && (
            <Alert severity="info">
              パス変換処理中...
            </Alert>
          )}
        </AccordionDetails>
      </Accordion>
    </Box>
  );
};