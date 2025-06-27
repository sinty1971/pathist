import { Link, useLocation } from "react-router";

export function Navigation() {
  const location = useLocation();
  
  // 現在のページタイトルを取得
  const getPageTitle = () => {
    switch (location.pathname) {
      case '/':
        return 'ファイル一覧';
      case '/projects':
        return '工程表';
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
            フォルダー一覧
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