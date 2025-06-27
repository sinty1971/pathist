import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router';
import { getFileFileinfos } from '../api/sdk.gen';
import { timestampToString } from '../utils/timestamp';
import { FileInfoModal } from './FileInfoModal';

export const FileInfoGrid: React.FC = () => {
  const navigate = useNavigate();
  const [fileInfos, setFileEntries] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPath, setCurrentPath] = useState('');
  const [pathInput, setPathInput] = useState('');
  const [selectedFileInfo, setSelectedFileInfo] = useState<any | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // ファイルシステムルートからの相対パスに変換
  const convertToRelativePath = (frontendPath: string): string => {
    // 空文字列または~/penguinの場合は空文字列を返す（ルートディレクトリ）
    if (!frontendPath || frontendPath === '~/penguin' || frontendPath === '/home/shin/penguin') {
      return '';
    }
    // ~/penguin/ で始まる場合はそれ以降を取り出す
    if (frontendPath.startsWith('~/penguin/')) {
      return frontendPath.substring('~/penguin/'.length);
    }
    // /home/shin/penguin/ で始まる場合はそれ以降を取り出す
    if (frontendPath.startsWith('/home/shin/penguin/')) {
      return frontendPath.substring('/home/shin/penguin/'.length);
    }
    // その他の場合はそのまま返す（相対パスとみなす）
    return frontendPath;
  };

  // 相対パスからフロントエンド表示用パスに変換
  const convertToDisplayPath = (relativePath: string): string => {
    if (!relativePath) {
      return '~/penguin';
    }
    return `~/penguin/${relativePath}`;
  };


  const loadFileEntries = useCallback(async (path?: string) => {
    const frontendPath = path || '';
    const relativePath = convertToRelativePath(frontendPath);
    

    setLoading(true);
    setError(null);
    
    try {
      console.log('Loading file entries for frontend path:', frontendPath);
      console.log('Converted to relative path:', relativePath);
      
      // APIクライアントを使用
      console.log('Calling API with query:', relativePath ? { path: relativePath } : {});
      const response = await getFileFileinfos({
        query: relativePath ? { path: relativePath } : {}
      });
      
      console.log('API response:', response);
      
      if (response.data) {
        // APIは直接配列を返す（実際のAPIでは日付は文字列として返される）
        const data = response.data as any[];
        console.log('Received data:', data);
        setFileEntries(Array.isArray(data) ? data : []);
        setCurrentPath(frontendPath);
      } else if (response.error) {
        console.error('API returned error:', response.error);
        throw new Error('APIエラー: ' + JSON.stringify(response.error));
      }
    } catch (err) {
      console.error('Error loading file entries:', err);
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  }, [navigate]);

  useEffect(() => {
    loadFileEntries();
  }, [loadFileEntries]);

  const handleFileInfoClick = (fileInfo: any) => {
    if (fileInfo.is_directory) {
      // ディレクトリの場合は移動
      // ファイルエントリのパスは絶対パスなので、フロントエンド表示用に変換
      const displayPath = convertToDisplayPath(convertToRelativePath(fileInfo.path || ''));
      
      
      setPathInput(displayPath);
      loadFileEntries(displayPath);
    } else {
      // ファイルの場合はモーダル表示
      setSelectedFileInfo(fileInfo);
      setIsModalOpen(true);
    }
  };

  const handlePathSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    // ~/penguinより親に行かないようにバリデーション
    const minPath = '~/penguin';
    if (pathInput.startsWith(minPath) || pathInput === minPath || pathInput === '') {
      loadFileEntries(pathInput);
    } else {
      // バリデーションエラーの場合、最小パスに設定
      setPathInput(minPath);
      loadFileEntries(minPath);
    }
  };

  const handleGoBack = () => {
    // 親ディレクトリのパスを取得
    const pathParts = currentPath.split('/');
    if (pathParts.length > 1) {
      const parentPath = pathParts.slice(0, -1).join('/');
      const newPath = parentPath || '~/penguin';
      
      // ~/penguinより親に行かないようにバリデーション
      const minPath = '~/penguin';
      if (newPath.startsWith(minPath) || newPath === minPath) {
        setPathInput(newPath);
        loadFileEntries(newPath);
      }
    }
  };



  const getFileInfoIcon = (fileInfo: any) => {
    if (fileInfo.is_directory) {
      return '📁';
    }
    const ext = fileInfo.name?.split('.').pop()?.toLowerCase();
    switch (ext) {
      case 'pdf': return '📄';
      case 'jpg':
      case 'jpeg':
      case 'png':
      case 'gif': return '🖼️';
      case 'mp4':
      case 'avi':
      case 'mov': return '🎬';
      case 'mp3':
      case 'wav': return '🎵';
      default: return '📄';
    }
  };

  // デバッグ用の表示
  console.log('FileInfoGrid render:', { 
    loading, 
    error, 
    fileInfosCount: fileInfos.length,
    pathInput,
    currentPath 
  });

  return (
    <div className="folder-container">
      <div className="header">
        <form onSubmit={handlePathSubmit} className="path-form">
          <button type="button" onClick={handleGoBack} className="back-button">
            <span className="back-arrow">⮜</span>
          </button>
          <input
            type="text"
            value={pathInput}
            onChange={(e) => setPathInput(e.target.value)}
            placeholder="フォルダーパスを入力"
            className="path-input"
          />
          <button type="submit" className="load-button">読み込み</button>
        </form>
      </div>

      <div className="folder-info">
        <span className="folder-count">{fileInfos.length} 項目</span>
        <span className="current-path">{currentPath || '~/penguin'}</span>
      </div>

      {loading && <div className="loading">読み込み中...</div>}
      {error && <div className="error">{error}</div>}

      <div className="folder-list">
        {fileInfos.map((fileInfo, index) => {
          return (
            <div
              key={index}
              className="folder-item"
              onClick={() => handleFileInfoClick(fileInfo)}
            >
              <div className="folder-icon">
                {getFileInfoIcon(fileInfo)}
              </div>
              <div className="folder-info">
                <div className="folder-name">
                  {fileInfo.name}
                </div>
                <div className="folder-meta">
                  <span>{fileInfo.is_directory ? 'フォルダー' : 'ファイル'}</span>
                  <span className="folder-date">
                    {' · 更新: '}
                    {timestampToString(fileInfo.modified_time) ? new Date(timestampToString(fileInfo.modified_time)!).toLocaleDateString('ja-JP', {
                      year: 'numeric',
                      month: '2-digit',
                      day: '2-digit',
                      hour: '2-digit',
                      minute: '2-digit'
                    }) : '-'}
                  </span>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <FileInfoModal
        fileInfo={selectedFileInfo}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      />
    </div>
  );
};