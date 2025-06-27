import React from 'react';
import { timestampToString } from '../utils/timestamp';

interface FileInfoModalProps {
  fileInfo: any | null;
  isOpen: boolean;
  onClose: () => void;
}

export const FileInfoModal: React.FC<FileInfoModalProps> = ({ fileInfo, isOpen, onClose }) => {
  if (!isOpen || !fileInfo) return null;

  const formatSize = (bytes?: number): string => {
    if (!bytes) return '0 B';
    const units = ['B', 'KB', 'MB', 'GB'];
    let size = bytes;
    let unitIndex = 0;
    
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }
    
    return `${size.toFixed(1)} ${units[unitIndex]}`;
  };

  const formatDate = (dateString?: string): string => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString('ja-JP');
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>{fileInfo.name}</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>
        
        <div className="modal-body">
          <div className="info-row">
            <span className="label">種類:</span>
            <span className="value">{fileInfo.is_directory ? 'フォルダー' : 'ファイル'}</span>
          </div>
          
          
          <div className="info-row">
            <span className="label">サイズ:</span>
            <span className="value">{formatSize(fileInfo.size)}</span>
          </div>
          
          <div className="info-row">
            <span className="label">パス:</span>
            <span className="value">{fileInfo.path}</span>
          </div>
          
          
          <div className="info-row">
            <span className="label">更新日時:</span>
            <span className="value">{formatDate(timestampToString(fileInfo.modified_time))}</span>
          </div>
        </div>
      </div>
    </div>
  );
};