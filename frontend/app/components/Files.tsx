import React, { useState, useEffect, useCallback } from "react";
import { SimpleTreeView, TreeItem } from "@mui/x-tree-view";
import type { TreeItemProps } from "@mui/x-tree-view/TreeItem";
import {
  Box,
  Typography,
  IconButton,
  Toolbar,
  Paper,
  Chip,
  Breadcrumbs,
  Link,
  CircularProgress,
  Alert,
} from "@mui/material";
import {
  InsertDriveFile,
  Refresh,
  Home,
  ExpandMore,
  ChevronRight,
} from "@mui/icons-material";
import { getBusinessFiles, getBusinessBasePath } from "../api/sdk.gen";
import { FileDetailModal } from "./FileDetailModal";
import { useFileInfo } from "../contexts/FileInfoContext";

// ユーティリティ関数群（コンポーネント外で定義）

// ノードアイコンを取得する関数
const getNodeIcon = (node: TreeNode) => {
  if (node.isDirectory) {
    return null; // ディレクトリは展開マークのみ表示
  }

  // 特定のファイル名のアイコンコードを返す
  if (node.name === ".detail.yaml") {
    return <span style={{ fontSize: "16px" }}>⚙️</span>;
  }

  // ファイル拡張子によってアイコンコードを返す
  const ext = node.name?.split(".").pop()?.toLowerCase();
  switch (ext) {
    case "pdf":
      return <span style={{ fontSize: "16px" }}>📄</span>;
    case "xlsx":
    case "xls":
    case "xlsm":
      return <span style={{ fontSize: "16px" }}>📊</span>;
    case "jpg":
    case "jpeg":
    case "png":
    case "gif":
      return <span style={{ fontSize: "16px" }}>🖼️</span>;
    case "mp4":
    case "avi":
    case "mov":
      return <span style={{ fontSize: "16px" }}>🎬</span>;
    case "mp3":
    case "wav":
      return <span style={{ fontSize: "16px" }}>🎵</span>;
    default:
      return <InsertDriveFile color="action" />;
  }
};

// ファイルサイズをフォーマットする関数
const formatFileSize = (bytes?: number): string => {
  if (!bytes || bytes === 0) return "";
  const units = ["B", "KB", "MB", "GB"];
  let size = bytes;
  let unitIndex = 0;

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }

  return ` (${size.toFixed(1)} ${units[unitIndex]})`;
};

// ツリー内のノードを更新するヘルパー関数
const updateNodeInTree = (
  nodes: TreeNode[],
  targetId: string,
  updatedNode: TreeNode
): TreeNode[] => {
  return nodes.map((node) => {
    if (node.id === targetId) {
      return updatedNode;
    }
    if (node.children) {
      return {
        ...node,
        children: updateNodeInTree(node.children, targetId, updatedNode),
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

// TreeNode型の定義
interface TreeNode {
  id: string;
  name: string;
  path: string;
  relativePath?: string; // 相対パス部分を保存
  isDirectory: boolean;
  size?: number;
  modifiedTime?: any;
  children?: TreeNode[];
  isLoaded?: boolean;
  isLoading?: boolean;
}

// カスタムTreeItemのProps
interface CustomTreeItemProps extends Omit<TreeItemProps, "itemId"> {
  itemId: string;
  node: TreeNode;
  onNodeClick: (node: TreeNode) => void;
  onNodeExpand: (nodeId: string, node: TreeNode) => void;
  isExpanded: boolean;
  expanded: string[];
}

// カスタムTreeItemコンポーネント
const CustomTreeItem: React.FC<CustomTreeItemProps> = React.memo(
  ({
    itemId,
    node,
    onNodeClick,
    onNodeExpand,
    isExpanded,
    expanded,
    ...props
  }) => {
    // ノードクリック処理
    const handleClick = React.useCallback(
      (e: React.MouseEvent) => {
        e.stopPropagation();
        onNodeClick(node);
      },
      [node, onNodeClick]
    );

    return (
      <TreeItem
        itemId={itemId}
        onClick={handleClick}
        sx={{
          // ローディング中のノードの視覚的フィードバックのみ
          ...(node.isLoading && {
            "& .MuiTreeItem-content": {
              opacity: 0.7,
            },
          }),
        }}
        label={
          <Box sx={{ display: "flex", alignItems: "center", py: 0.5, pr: 2 }}>
            {node.isDirectory ? (
              isExpanded ? (
                <ExpandMore sx={{ mr: 0.5 }} />
              ) : (
                <ChevronRight sx={{ mr: 0.5 }} />
              )
            ) : (
              getNodeIcon(node)
            )}
            <Typography variant="body2" sx={{ flexGrow: 1, mr: 1, ml: 1 }}>
              {node.name}
              {!node.isDirectory && formatFileSize(node.size)}
            </Typography>
            {node.name === ".detail.yaml" && (
              <Chip
                label="詳細"
                size="small"
                color="info"
                variant="outlined"
                sx={{ ml: 1, fontSize: "10px", height: "20px" }}
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
            isExpanded={expanded.includes(child.id)}
            expanded={expanded}
          />
        ))}
      </TreeItem>
    );
  }
);

export const Files: React.FC = () => {
  const [treeData, setTreeData] = useState<TreeNode[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPath, setCurrentPath] = useState("");
  const [basePath, setBasePath] = useState<string>("");
  const [basePathError, setBasePathError] = useState(false);
  const [selectedNode, setSelectedNode] = useState<TreeNode | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [expanded, setExpanded] = useState<string[]>([]);
  const { setFileCount, setCurrentPath: setContextPath } = useFileInfo();

  // ベースパスを取得
  const loadBasePath = useCallback(async () => {
    try {
      const response = await getBusinessBasePath();
      if (response.data && response.data.businessBasePath) {
        setBasePath(response.data.businessBasePath);
        setBasePathError(false);
        return response.data.businessBasePath;
      }
      throw new Error("ベースパスが取得できませんでした");
    } catch (err) {
      setBasePathError(true);
      setError("基準となるパスを取得できないため、一覧を表示できません");
      return null;
    }
  }, []);


  // ファイルデータをTreeNode形式に変換
  const convertToTreeNode = (fileInfo: any): TreeNode => {
    const node = {
      id: fileInfo.path, // バックエンドからのpathをIDとして使用
      name: fileInfo.name,
      path: fileInfo.path, // バックエンドからのpathをそのまま使用
      relativePath: fileInfo.path.replace(basePath + '/', '') || '', // 相対パス部分を保存
      isDirectory: fileInfo.is_directory,
      size: fileInfo.size,
      modifiedTime: fileInfo.modified_time,
      children: fileInfo.is_directory ? [] : undefined,
      isLoaded: !fileInfo.is_directory,
      isLoading: false,
    };

    return node;
  };

  // ファイル一覧を読み込み
  const loadFiles = useCallback(
    async (relativePath?: string, isRefresh = false) => {
      const requestPath = relativePath || "";

      setLoading(true);
      setError(null);

      try {
        const response = await getBusinessFiles({
          query: requestPath ? { path: requestPath } : {},
        });

        if (response.data) {
          const data = response.data as any[];

          const nodes = data.map((fileInfo) => {
            return convertToTreeNode(fileInfo);
          });

          if (!relativePath || isRefresh) {
            // ルートレベルの読み込みまたはリフレッシュの場合のみツリーデータを置き換え
            setTreeData(nodes);
            // コンテキストを更新
            setFileCount(data.length);
            setContextPath(basePath);
            setCurrentPath(basePath);
          }

          return nodes;
        } else if (response.error) {
          throw new Error("APIエラー: " + JSON.stringify(response.error));
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : "エラーが発生しました");
        return [];
      } finally {
        setLoading(false);
      }
    },
    [basePath, setFileCount, setContextPath]
  );

  // 初期読み込み
  useEffect(() => {
    const initialize = async () => {
      setLoading(true);
      const basePathResult = await loadBasePath();
      if (basePathResult) {
        await loadFiles();
      }
      setLoading(false);
    };
    initialize();
  }, [loadBasePath, loadFiles]);

  // ノード展開時の処理
  const handleNodeExpand = useCallback(
    async (nodeId: string, node: TreeNode) => {
      if (!node.isDirectory || node.isLoaded || node.isLoading) return;

      // ノードをローディング状態に設定
      setTreeData((prevData) =>
        updateNodeInTree(prevData, nodeId, { ...node, isLoading: true })
      );

      try {
        // nodeの相対パス部分を直接使用
        const targetRelativePath = node.relativePath || '';

        const children = await loadFiles(targetRelativePath, false);

        // 子ノードを設定
        setTreeData((prevData) =>
          updateNodeInTree(prevData, nodeId, {
            ...node,
            children,
            isLoaded: true,
            isLoading: false,
          })
        );
      } catch (err) {
        console.error("handleNodeExpand error:", err);
        // エラー時はローディング状態を解除
        setTreeData((prevData) =>
          updateNodeInTree(prevData, nodeId, { ...node, isLoading: false })
        );
      }
    },
    [loadFiles]
  );

  // ノードクリック処理
  const handleNodeClick = useCallback(
    (node: TreeNode) => {
      if (!node.isDirectory) {
        // ファイルの場合はモーダル表示
        setSelectedNode(node);
        setIsModalOpen(true);
      } else {
        // ディレクトリの場合は展開/縮小
        setExpanded((prev) => {
          const isExpanded = prev.includes(node.id);

          if (isExpanded) {
            // 縮小
            return prev.filter((id) => id !== node.id);
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
    },
    [handleNodeExpand]
  );

  // リフレッシュ
  const handleRefresh = () => {
    setTreeData([]);
    setExpanded([]);
    loadFiles(currentPath, true);
  };

  // ホームに戻る
  const handleGoHome = () => {
    setTreeData([]);
    setExpanded([]);
    loadFiles("", true);
  };

  // パスのブレッドクラム
  const getBreadcrumbs = () => {
    const currentBasePath = currentPath || basePath;
    const parts = currentBasePath
      .replace(basePath, "")
      .split("/")
      .filter(Boolean);
    const breadcrumbs = [
      { label: basePath.split("/").pop() || "ホーム", path: "" },
    ];

    let accumulatedPath = "";
    parts.forEach((part) => {
      if (accumulatedPath === "") {
        accumulatedPath = `${basePath}/${part}`;
      } else {
        accumulatedPath += `/${part}`;
      }
      breadcrumbs.push({ label: part, path: accumulatedPath });
    });

    return breadcrumbs;
  };

  return (
    <Box sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
      {/* ツールバー */}
      <Paper elevation={1} sx={{ mb: 1 }}>
        <Toolbar variant="dense">
          <IconButton onClick={handleGoHome} size="small" title="ホームに戻る">
            <Home />
          </IconButton>
          <IconButton onClick={handleRefresh} size="small" title="リフレッシュ">
            <Refresh />
          </IconButton>
        </Toolbar>

        {/* ブレッドクラム */}
        <Box sx={{ px: 2, pb: 1 }}>
          <Breadcrumbs>
            {getBreadcrumbs().map((crumb, index) => (
              <Link
                key={index}
                component="button"
                variant="body2"
                color={
                  index === getBreadcrumbs().length - 1
                    ? "text.primary"
                    : "inherit"
                }
                onClick={() => {
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
      <Paper
        sx={{
          flex: 1,
          overflow: "auto",
          p: 1,
          position: "relative",
          minHeight: 0,
        }}
      >
        {loading && treeData.length === 0 ? (
          <Box
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
              p: 3,
            }}
          >
            <Box sx={{ position: "relative", display: "inline-flex" }}>
              <CircularProgress
                size={60}
                thickness={1}
                sx={{
                  color: "primary.main",
                  animationDuration: "2s",
                }}
              />
              <CircularProgress
                variant="determinate"
                size={60}
                thickness={2}
                value={25}
                sx={{
                  color: "grey.300",
                  position: "absolute",
                  top: 0,
                  left: 0,
                  zIndex: 0,
                }}
              />
            </Box>
            <Typography color="text.secondary" sx={{ mt: 2 }}>
              データ取得中...
            </Typography>
            <Typography
              variant="caption"
              color="text.secondary"
              sx={{ mt: 0.5 }}
            >
              大量のファイルがある場合は時間がかかることがあります
            </Typography>
          </Box>
        ) : basePathError ? (
          <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
            <Typography color="error">
              基準となるパスを取得できないため、一覧を表示できません
            </Typography>
          </Box>
        ) : (
          <>
            <SimpleTreeView
              expandedItems={expanded}
              onExpandedItemsChange={() => {
                // MUIのTreeViewによる自動展開/縮小を無効化
                // 手動でクリック処理を制御するため
              }}
              sx={{
                flexGrow: 1,
                overflowY: "auto",
                "& .MuiTreeItem-content": {
                  "&:hover": {
                    backgroundColor: "action.hover",
                  },
                },
                "& .MuiTreeItem-iconContainer": {
                  display: "none", // MUIのデフォルト展開アイコンを非表示（カスタムアイコンを使用）
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
                  isExpanded={expanded.includes(node.id)}
                  expanded={expanded}
                />
              ))}
            </SimpleTreeView>

            {/* ノード展開時のローディングオーバーレイ */}
            {treeData.some(
              (node) =>
                node.isLoading ||
                (node.children &&
                  node.children.some((child) => child.isLoading))
            ) && (
              <Box
                sx={{
                  position: "absolute",
                  top: "50%",
                  left: "50%",
                  transform: "translate(-50%, -50%)",
                  backgroundColor: "rgba(255, 255, 255, 0.95)",
                  borderRadius: 2,
                  p: 3,
                  boxShadow: 3,
                  zIndex: 1000,
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                  minWidth: 200,
                }}
              >
                <CircularProgress
                  size={48}
                  thickness={3}
                  sx={{
                    color: "primary.main",
                  }}
                />
                <Typography
                  color="text.primary"
                  sx={{ mt: 2, fontWeight: 500 }}
                >
                  フォルダーを読み込んでいます...
                </Typography>
                <Typography
                  variant="caption"
                  color="text.secondary"
                  sx={{ mt: 0.5 }}
                >
                  大量のファイルがある場合は時間がかかります
                </Typography>
              </Box>
            )}
          </>
        )}

        {treeData.length === 0 && !loading && !basePathError && (
          <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
            <Typography color="text.secondary">フォルダーが空です</Typography>
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
