import { Link, useLocation } from "react-router";
import { useFileInfo } from "../contexts/FileInfoContext";

interface NavigationProps {
  projectCount?: number;
}

export function Navigation({ projectCount }: NavigationProps = {}) {
  const location = useLocation();
  const { fileCount, currentPath } = useFileInfo();
  
  // 現在のページタイトルを取得
  const getPageTitle = () => {
    switch (location.pathname) {
      case '/':
        return 'ホーム';
      case '/files':
        return `ファイル一覧${fileCount > 0 ? ` (${fileCount}項目)` : ''}${currentPath ? ` - ${currentPath}` : ''}`;
      case '/projects':
        return `工程表${projectCount !== undefined ? ` (${projectCount}件)` : ''}`;
      case '/projects/gantt':
        return 'ガントチャート';
      default:
        return 'ファイル管理';
    }
  };
  
  return (
    <nav className="navigation">
      <div className="nav-container">
        <div className="nav-logo">
          <h1>Penguin フォルダー管理</h1>
          <div className="page-title">{getPageTitle()}</div>
        </div>
        <div className="nav-links">
          <Link 
            to="/" 
            className={location.pathname === '/' ? 'nav-link active' : 'nav-link'}
          >
            ホーム
          </Link>
          <Link 
            to="/files" 
            className={location.pathname === '/files' ? 'nav-link active' : 'nav-link'}
          >
            ファイル一覧
          </Link>
          <Link 
            to="/projects" 
            className={location.pathname === '/projects' ? 'nav-link active' : 'nav-link'}
          >
            工程表
          </Link>
        </div>
      </div>
    </nav>
  );
}