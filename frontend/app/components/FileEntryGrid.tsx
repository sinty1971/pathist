import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router';
import type { FileEntry } from '../services/api';
import { folderService } from '../services/api';
import { FileEntryModal } from './FileEntryModal';

export const FileEntryGrid: React.FC = () => {
  const navigate = useNavigate();
  const [fileEntries, setFileEntries] = useState<FileEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPath, setCurrentPath] = useState('');
  const [pathInput, setPathInput] = useState('');
  const [selectedFileEntry, setSelectedFileEntry] = useState<FileEntry | null>(null);
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

  // 工事プロジェクトディレクトリかどうかをチェック
  const isKoujiProjectPath = (path: string) => {
    const normalizedPath = path.replace(/\\/g, '/');
    return normalizedPath.includes('/豊田築炉/2-工事') || 
           normalizedPath.endsWith('/2-工事') ||
           normalizedPath.includes('2-工事');
  };

  const loadFileEntries = async (path?: string) => {
    const frontendPath = path || '';
    const relativePath = convertToRelativePath(frontendPath);
    
    // 工事プロジェクトディレクトリの場合は工事プロジェクトページにリダイレクト
    if (isKoujiProjectPath(frontendPath)) {
      navigate('/kouji');
      return;
    }

    setLoading(true);
    setError(null);
    
    try {
      console.log('Loading file entries for frontend path:', frontendPath);
      console.log('Converted to relative path:', relativePath);
      
      // 直接fetchでテスト
      const directResponse = await fetch(`http://localhost:8080/api/file-entries${relativePath ? `?path=${encodeURIComponent(relativePath)}` : ''}`);
      console.log('Direct fetch response status:', directResponse.status);
      console.log('Direct fetch response ok:', directResponse.ok);
      
      if (!directResponse.ok) {
        const errorText = await directResponse.text();
        console.error('Direct fetch error:', errorText);
        throw new Error(`Direct fetch failed: ${directResponse.status} ${errorText}`);
      }
      
      const directData = await directResponse.json();
      console.log('Direct fetch data:', directData);
      
      setFileEntries(directData.file_entries || []);
      setCurrentPath(frontendPath);
    } catch (err) {
      console.error('Error loading file entries:', err);
      setError(err instanceof Error ? err.message : 'エラーが発生しました');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadFileEntries();
  }, []);

  const handleFileEntryClick = (fileEntry: FileEntry) => {
    if (fileEntry.is_directory) {
      // ディレクトリの場合は移動
      // ファイルエントリのパスは絶対パスなので、フロントエンド表示用に変換
      const displayPath = convertToDisplayPath(convertToRelativePath(fileEntry.path));
      
      // 工事プロジェクトディレクトリの場合は工事プロジェクトページにリダイレクト
      if (isKoujiProjectPath(displayPath)) {
        navigate('/kouji');
        return;
      }
      
      setPathInput(displayPath);
      loadFileEntries(displayPath);
    } else {
      // ファイルの場合はモーダル表示
      setSelectedFileEntry(fileEntry);
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


  // 特別なフォルダーかどうかをチェック
  const isSpecialFileEntry = (fileEntry: FileEntry) => {
    if (!fileEntry.is_directory) return false;
    return isKoujiProjectPath(fileEntry.path) || fileEntry.name === '2-工事';
  };

  const getFileEntryIcon = (fileEntry: FileEntry) => {
    if (fileEntry.is_directory) {
      // 工事プロジェクトフォルダーの場合は特別なアイコン
      if (isSpecialFileEntry(fileEntry)) {
        return '🏗️';
      }
      return '📁';
    }
    const ext = fileEntry.name.split('.').pop()?.toLowerCase();
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
  console.log('FileEntryGrid render:', { 
    loading, 
    error, 
    fileEntriesCount: fileEntries.length,
    pathInput,
    currentPath 
  });

  return (
    <div className="folder-container">
      <div className="header">
        <h1>フォルダー管理システム</h1>
        
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
        <span className="folder-count">{fileEntries.length} 項目</span>
        <span className="current-path">{currentPath || '~/penguin'}</span>
      </div>

      {loading && <div className="loading">読み込み中...</div>}
      {error && <div className="error">{error}</div>}

      <div className="folder-list">
        {fileEntries.map((fileEntry, index) => {
          const isSpecial = isSpecialFileEntry(fileEntry);
          return (
            <div
              key={index}
              className={`folder-item ${isSpecial ? 'folder-item--special' : ''}`}
              onClick={() => handleFileEntryClick(fileEntry)}
            >
              <div className={`folder-icon ${isSpecial ? 'folder-icon--special' : ''}`}>
                {getFileEntryIcon(fileEntry)}
              </div>
              <div className="folder-info">
                <div className={`folder-name ${isSpecial ? 'folder-name--special' : ''}`}>
                  {fileEntry.name}
                  {isSpecial && <span className="special-badge">工事一覧</span>}
                </div>
                <div className="folder-meta">
                  <span>{fileEntry.is_directory ? 'フォルダー' : 'ファイル'}</span>
                  <span className="folder-date">
                    {' · 更新: '}
                    {new Date(fileEntry.modified_time).toLocaleDateString('ja-JP', {
                      year: 'numeric',
                      month: '2-digit',
                      day: '2-digit',
                      hour: '2-digit',
                      minute: '2-digit'
                    })}
                  </span>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <FileEntryModal
        fileEntry={selectedFileEntry}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      />
    </div>
  );
};