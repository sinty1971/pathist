'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import '../styles/navigation.css';

export function Navigation() {
  const pathname = usePathname();

  return (
    <nav className="navigation">
      <div className="nav-container">
        <div className="nav-links">
          <Link
            href="/"
            className={pathname === '/' ? 'nav-link active' : 'nav-link'}
          >
            ホーム
          </Link>
          <Link
            href="/files"
            className={pathname === '/files' ? 'nav-link active' : 'nav-link'}
          >
            ファイル一覧
          </Link>
          <Link
            href="/kojies"
            className={pathname === '/kojies' ? 'nav-link active' : 'nav-link'}
          >
            工事一覧
          </Link>
          <Link
            href="/companies"
            className={pathname === '/companies' ? 'nav-link active' : 'nav-link'}
          >
            会社一覧
          </Link>
        </div>
        <div className="nav-logo">
          <h1>Penguin フォルダー管理</h1>
        </div>
      </div>
    </nav>
  );
}
