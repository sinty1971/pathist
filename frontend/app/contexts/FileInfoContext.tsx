import React, { createContext, useContext, useState } from 'react';

interface FileInfoContextType {
  fileCount: number;
  currentPath: string;
  setFileCount: (count: number) => void;
  setCurrentPath: (path: string) => void;
}

const FileInfoContext = createContext<FileInfoContextType | undefined>(undefined);

export const FileInfoProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [fileCount, setFileCount] = useState(0);
  const [currentPath, setCurrentPath] = useState('');

  return (
    <FileInfoContext.Provider value={{ fileCount, currentPath, setFileCount, setCurrentPath }}>
      {children}
    </FileInfoContext.Provider>
  );
};

export const useFileInfo = () => {
  const context = useContext(FileInfoContext);
  if (!context) {
    throw new Error('useFileInfo must be used within a FileInfoProvider');
  }
  return context;
};