import { Link, useLocation } from "react-router";
import "../styles/navigation.css";

export function Navigation() {
  const location = useLocation();
  
  
  return (
    <nav className="navigation">
      <div className="nav-container">
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
            工事一覧
          </Link>
        </div>
        <div className="nav-logo">
          <h1>Penguin フォルダー管理</h1>
        </div>
      </div>
    </nav>
  );
}