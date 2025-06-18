import React from 'react';
import type { FileEntry } from '../services/api';

interface FileEntryModalProps {
  fileEntry: FileEntry | null;
  isOpen: boolean;
  onClose: () => void;
}

export const FileEntryModal: React.FC<FileEntryModalProps> = ({ fileEntry, isOpen, onClose }) => {
  if (!isOpen || !fileEntry) return null;

  const formatSize = (bytes: number): string => {
    const units = ['B', 'KB', 'MB', 'GB'];
    let size = bytes;
    let unitIndex = 0;
    
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }
    
    return `${size.toFixed(1)} ${units[unitIndex]}`;
  };

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleString('ja-JP');
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>{fileEntry.name}</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>
        
        <div className="modal-body">
          <div className="info-row">
            <span className="label">種類:</span>
            <span className="value">{fileEntry.is_directory ? 'フォルダー' : 'ファイル'}</span>
          </div>
          
          <div className="info-row">
            <span className="label">ID:</span>
            <span className="value">{fileEntry.id}</span>
          </div>
          
          <div className="info-row">
            <span className="label">サイズ:</span>
            <span className="value">{formatSize(fileEntry.size)}</span>
          </div>
          
          <div className="info-row">
            <span className="label">パス:</span>
            <span className="value">{fileEntry.path}</span>
          </div>
          
          
          <div className="info-row">
            <span className="label">更新日時:</span>
            <span className="value">{formatDate(fileEntry.modified_time)}</span>
          </div>
        </div>
      </div>
    </div>
  );
};