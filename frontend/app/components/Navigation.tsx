import { Link, useLocation } from "react-router";

export function Navigation() {
  const location = useLocation();
  
  return (
    <nav className="navigation">
      <div className="nav-container">
        <div className="nav-logo">
          <h1>Penguin フォルダー管理</h1>
        </div>
        <div className="nav-links">
          <Link 
            to="/" 
            className={location.pathname === '/' ? 'nav-link active' : 'nav-link'}
          >
            フォルダー一覧
          </Link>
          <Link 
            to="/kouji" 
            className={location.pathname === '/kouji' ? 'nav-link active' : 'nav-link'}
          >
            工事一覧
          </Link>
        </div>
      </div>
    </nav>
  );
}