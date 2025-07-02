import React, { useState, useEffect, useCallback } from 'react';
import { 
  SimpleTreeView,
  TreeItem
} from '@mui/x-tree-view';
import type { TreeItemProps } from '@mui/x-tree-view/TreeItem';
import { 
  Box, 
  Typography, 
  TextField, 
  IconButton, 
  Toolbar,
  Paper,
  Chip,
  Breadcrumbs,
  Link,
  CircularProgress,
  Alert
} from '@mui/material';
import {
  Folder,
  FolderOpen,
  InsertDriveFile,
  Refresh,
  Home,
  Search
} from '@mui/icons-material';
import { getFileFileinfos } from '../api/sdk.gen';
import { FileDetailModal } from './FileDetailModal';
import { useFileInfo } from '../contexts/FileInfoContext';

// TreeNode型の定義
interface TreeNode {
  id: string;
  name: string;
  path: string;
  isDirectory: boolean;
  size?: number;
  modifiedTime?: any;
  children?: TreeNode[];
  isLoaded?: boolean;
  isLoading?: boolean;
}

// カスタムTreeItemのProps
interface CustomTreeItemProps extends Omit<TreeItemProps, 'itemId'> {
  itemId: string;
  node: TreeNode;
  onNodeClick: (node: TreeNode) => void;
  onNodeExpand: (nodeId: string, node: TreeNode) => void;
}

// カスタムTreeItemコンポーネント
const CustomTreeItem: React.FC<CustomTreeItemProps> = React.memo(({ 
  itemId, 
  node, 
  onNodeClick,
  onNodeExpand,
  ...props 
}) => {
  const getNodeIcon = (node: TreeNode, expanded = false) => {
    if (node.isDirectory) {
      return expanded ? <FolderOpen color="primary" /> : <Folder color="primary" />;
    }
    
    // 特定のファイル名をチェック
    if (node.name === '.detail.yaml') {
      return <span style={{ fontSize: '16px' }}>⚙️</span>;
    }
    
    const ext = node.name?.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf': return <span style={{ fontSize: '16px' }}>📄</span>;
      case 'xlsx':
      case 'xls':
      case 'xlsm': return <span style={{ fontSize: '16px' }}>📊</span>;
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif': return <span style={{ fontSize: '16px' }}>🖼️</span>;
      case 'mp4':
      case 'avi':
      case 'mov': return <span style={{ fontSize: '16px' }}>🎬</span>;
      case 'mp3':
      case 'wav': return <span style={{ fontSize: '16px' }}>🎵</span>;
      default: return <InsertDriveFile color="action" />;
    }
  };

  const formatSize = (bytes?: number): string => {
    if (!bytes || bytes === 0) return '';
    const units = ['B', 'KB', 'MB', 'GB'];
    let size = bytes;
    let unitIndex = 0;
    
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }
    
    return ` (${size.toFixed(1)} ${units[unitIndex]})`;
  };

  const handleClick = React.useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onNodeClick(node);
  }, [node, onNodeClick]);

  return (
    <TreeItem
      itemId={itemId}
      onClick={handleClick}
      label={
        <Box sx={{ display: 'flex', alignItems: 'center', py: 0.5 }}>
          {getNodeIcon(node)}
          <Typography variant="body2" sx={{ flexGrow: 1, mr: 1, ml: 1 }}>
            {node.name}
            {!node.isDirectory && formatSize(node.size)}
          </Typography>
          {node.isLoading && <CircularProgress size={16} />}
          {node.name === '.detail.yaml' && (
            <Chip 
              label="詳細" 
              size="small" 
              color="info" 
              variant="outlined"
              sx={{ ml: 1, fontSize: '10px', height: '20px' }}
            />
          )}
        </Box>
      }
      {...props}
    >
      {node.children?.map((child) => (
        <CustomTreeItem
          key={child.id}
          itemId={child.id}
          node={child}
          onNodeClick={onNodeClick}
          onNodeExpand={onNodeExpand}
        />
      ))}
    </TreeItem>
  );
});

export const Files: React.FC = () => {
  const [treeData, setTreeData] = useState<TreeNode[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPath, setCurrentPath] = useState('~/penguin');
  const [pathInput, setPathInput] = useState('~/penguin');
  const [selectedNode, setSelectedNode] = useState<TreeNode | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [expanded, setExpanded] = useState<string[]>([]);
  const { setFileCount, setCurrentPath: setContextPath } = useFileInfo();

  // パス変換関数
  const convertToRelativePath = (frontendPath: string): string => {
    if (!frontendPath || frontendPath === '~/penguin' || frontendPath === '/home/shin/penguin') {
      return '';
    }
    if (frontendPath.startsWith('~/penguin/')) {
      return frontendPath.substring('~/penguin/'.length);
    }
    if (frontendPath.startsWith('/home/shin/penguin/')) {
      return frontendPath.substring('/home/shin/penguin/'.length);
    }
    return frontendPath;
  };


  // ファイルデータをTreeNode形式に変換
  const convertToTreeNode = (fileInfo: any, basePath: string): TreeNode => {
    return {
      id: `${basePath}/${fileInfo.name}`,
      name: fileInfo.name,
      path: fileInfo.path || `${basePath}/${fileInfo.name}`,
      isDirectory: fileInfo.is_directory,
      size: fileInfo.size,
      modifiedTime: fileInfo.modified_time,
      children: fileInfo.is_directory ? [] : undefined,
      isLoaded: !fileInfo.is_directory,
      isLoading: false
    };
  };

  // ファイル一覧を読み込み
  const loadFiles = useCallback(async (path?: string, isRefresh = false) => {
    const frontendPath = path || '~/penguin';
    const relativePath = convertToRelativePath(frontendPath);
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await getFileFileinfos({
        query: relativePath ? { path: relativePath } : {}
      });
      
      if (response.data) {
        const data = response.data as any[];
        const nodes = data.map(fileInfo => convertToTreeNode(fileInfo, frontendPath));
        
        if (path === '~/penguin' || !path || isRefresh) {
          // ルートレベルの読み込みまたはリフレッシュの場合のみツリーデータを置き換え
          setTreeData(nodes);
          // コンテキストを更新（ルートレベルまたはリフレッシュの場合のみ）
          setFileCount(data.length);
          setContextPath(frontendPath);
          setCurrentPath(frontendPath);
        }
        
        return nodes;
      } else if (response.error) {
        throw new Error('APIエラー: ' + JSON.stringify(response.error));
      }
    } catch (err) {
      console.error('Error loading files:', err);
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
      return [];
    } finally {
      setLoading(false);
    }
  }, [setFileCount, setContextPath]);

  // 初期読み込み
  useEffect(() => {
    loadFiles();
  }, [loadFiles]);

  // ノード展開時の処理
  const handleNodeExpand = useCallback(async (nodeId: string, node: TreeNode) => {
    if (!node.isDirectory || node.isLoaded || node.isLoading) return;

    // ノードをローディング状態に設定
    setTreeData(prevData => updateNodeInTree(prevData, nodeId, { ...node, isLoading: true }));

    try {
      const children = await loadFiles(node.path, false);
      
      // 子ノードを設定
      setTreeData(prevData => updateNodeInTree(prevData, nodeId, {
        ...node,
        children,
        isLoaded: true,
        isLoading: false
      }));
    } catch (err) {
      // エラー時はローディング状態を解除
      setTreeData(prevData => updateNodeInTree(prevData, nodeId, { ...node, isLoading: false }));
    }
  }, [loadFiles]);

  // ツリー内のノードを更新するヘルパー関数
  const updateNodeInTree = (nodes: TreeNode[], targetId: string, updatedNode: TreeNode): TreeNode[] => {
    return nodes.map(node => {
      if (node.id === targetId) {
        return updatedNode;
      }
      if (node.children) {
        return {
          ...node,
          children: updateNodeInTree(node.children, targetId, updatedNode)
        };
      }
      return node;
    });
  };

  // ツリー内のノードを検索するヘルパー関数
  const findNodeById = (nodes: TreeNode[], targetId: string): TreeNode | null => {
    for (const node of nodes) {
      if (node.id === targetId) {
        return node;
      }
      if (node.children) {
        const found = findNodeById(node.children, targetId);
        if (found) return found;
      }
    }
    return null;
  };

  // ノードクリック処理
  const handleNodeClick = useCallback((node: TreeNode) => {
    if (!node.isDirectory) {
      // ファイルの場合はモーダル表示
      setSelectedNode(node);
      setIsModalOpen(true);
    } else {
      // ディレクトリの場合は展開/縮小
      setExpanded(prev => {
        const isExpanded = prev.includes(node.id);
        
        if (isExpanded) {
          // 縮小
          return prev.filter(id => id !== node.id);
        } else {
          // 展開
          const newExpanded = [...prev, node.id];
          
          // まだロードされていない場合は子要素を読み込み
          if (!node.isLoaded && !node.isLoading) {
            // 非同期処理は状態更新後に実行
            setTimeout(() => handleNodeExpand(node.id, node), 0);
          }
          
          return newExpanded;
        }
      });
    }
  }, [handleNodeExpand]);

  // パス移動
  const handlePathSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // パスが変わる場合は新しいツリーを読み込み
    if (pathInput !== currentPath) {
      setTreeData([]);
      setExpanded([]);
      loadFiles(pathInput, true);
    }
  };

  // リフレッシュ
  const handleRefresh = () => {
    setTreeData([]);
    setExpanded([]);
    loadFiles(currentPath, true);
  };

  // ホームに戻る
  const handleGoHome = () => {
    setPathInput('~/penguin');
    setTreeData([]);
    setExpanded([]);
    loadFiles('~/penguin', true);
  };

  // パスのブレッドクラム
  const getBreadcrumbs = () => {
    const parts = currentPath.replace('~/penguin', '').split('/').filter(Boolean);
    const breadcrumbs = [
      { label: 'penguin', path: '~/penguin' }
    ];
    
    let accumulatedPath = '~/penguin';
    parts.forEach(part => {
      accumulatedPath += `/${part}`;
      breadcrumbs.push({ label: part, path: accumulatedPath });
    });
    
    return breadcrumbs;
  };

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      {/* ツールバー */}
      <Paper elevation={1} sx={{ mb: 1 }}>
        <Toolbar variant="dense">
          <IconButton onClick={handleGoHome} size="small" title="ホームに戻る">
            <Home />
          </IconButton>
          <IconButton onClick={handleRefresh} size="small" title="リフレッシュ">
            <Refresh />
          </IconButton>
          
          <Box component="form" onSubmit={handlePathSubmit} sx={{ flexGrow: 1, mx: 2 }}>
            <TextField
              size="small"
              fullWidth
              value={pathInput}
              onChange={(e) => setPathInput(e.target.value)}
              placeholder="パスを入力"
              slotProps={{
                input: {
                  startAdornment: <Search sx={{ mr: 1, color: 'action.active' }} />
                }
              }}
            />
          </Box>
        </Toolbar>
        
        {/* ブレッドクラム */}
        <Box sx={{ px: 2, pb: 1 }}>
          <Breadcrumbs>
            {getBreadcrumbs().map((crumb, index) => (
              <Link
                key={index}
                component="button"
                variant="body2"
                color={index === getBreadcrumbs().length - 1 ? 'text.primary' : 'inherit'}
                onClick={() => {
                  setPathInput(crumb.path);
                  if (crumb.path !== currentPath) {
                    setTreeData([]);
                    setExpanded([]);
                    loadFiles(crumb.path, true);
                  }
                }}
                underline="hover"
              >
                {crumb.label}
              </Link>
            ))}
          </Breadcrumbs>
        </Box>
      </Paper>

      {/* エラー表示 */}
      {error && (
        <Alert severity="error" sx={{ mb: 1 }}>
          {error}
        </Alert>
      )}

      {/* ツリービュー */}
      <Paper sx={{ flex: 1, overflow: 'auto', p: 1 }}>
        {loading && treeData.length === 0 ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : (
          <SimpleTreeView
            expandedItems={expanded}
            onExpandedItemsChange={() => {
              // MUIのTreeViewによる自動展開/縮小を無効化
              // 手動でクリック処理を制御するため
            }}
            sx={{
              flexGrow: 1,
              overflowY: 'auto',
              '& .MuiTreeItem-content': {
                '&:hover': {
                  backgroundColor: 'action.hover',
                },
              },
            }}
          >
            {treeData.map((node) => (
              <CustomTreeItem
                key={node.id}
                itemId={node.id}
                node={node}
                onNodeClick={handleNodeClick}
                onNodeExpand={handleNodeExpand}
              />
            ))}
          </SimpleTreeView>
        )}
        
        {treeData.length === 0 && !loading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <Typography color="text.secondary">
              フォルダーが空です
            </Typography>
          </Box>
        )}
      </Paper>

      {/* ファイル詳細モーダル */}
      <FileDetailModal
        fileInfo={selectedNode}
        isOpen={isModalOpen}
        onClose={() => {
          setIsModalOpen(false);
          setSelectedNode(null);
        }}
      />
    </Box>
  );
};