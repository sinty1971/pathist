import { Link } from "react-router";

export default function HomePage() {
  return (
    <div style={{
      height: "100vh",
      overflowY: "auto"
    }}>
      <div style={{
        padding: "3rem",
        maxWidth: "1200px",
        margin: "0 auto"
      }}>
      <div style={{
        textAlign: "center",
        marginBottom: "3rem"
      }}>
        <h1 style={{
          fontSize: "3rem",
          color: "#333",
          marginBottom: "1rem",
          fontWeight: "600"
        }}>
          Penguin フォルダー管理システム
        </h1>
        <p style={{
          fontSize: "1.2rem",
          color: "#666",
          fontWeight: "400"
        }}>
          ファイル管理と工事プロジェクト管理を効率的に
        </p>
      </div>

      <div style={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))",
        gap: "1rem",
        marginTop: "2rem"
      }}>
        {/* ファイル一覧カード */}
        <Link to="/files" style={{ textDecoration: "none" }}>
          <div className="feature-card" style={{
            background: "rgba(255, 255, 255, 0.95)",
            borderRadius: "12px",
            padding: "1.5rem",
            boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
            transition: "all 0.2s ease",
            cursor: "pointer",
            height: "180px",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center"
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.transform = "translateY(-3px)";
            e.currentTarget.style.boxShadow = "0 8px 20px rgba(0, 0, 0, 0.15)";
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.transform = "translateY(0)";
            e.currentTarget.style.boxShadow = "0 4px 12px rgba(0, 0, 0, 0.1)";
          }}>
            <div style={{
              fontSize: "2.5rem",
              marginBottom: "0.5rem"
            }}>
              📁
            </div>
            <h3 style={{
              fontSize: "1.1rem",
              color: "#333",
              marginBottom: "0.5rem",
              textAlign: "center"
            }}>
              ファイル一覧
            </h3>
            <p style={{
              color: "#666",
              fontSize: "0.85rem",
              textAlign: "center",
              lineHeight: "1.4"
            }}>
              ファイルとフォルダーの閲覧・管理
            </p>
          </div>
        </Link>

        {/* 工程表カード */}
        <div className="feature-card" style={{
          background: "rgba(255, 255, 255, 0.95)",
          borderRadius: "12px",
          padding: "1.5rem",
          boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
          transition: "all 0.2s ease",
          height: "220px",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          position: "relative"
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.transform = "translateY(-3px)";
          e.currentTarget.style.boxShadow = "0 8px 20px rgba(0, 0, 0, 0.15)";
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.transform = "translateY(0)";
          e.currentTarget.style.boxShadow = "0 4px 12px rgba(0, 0, 0, 0.1)";
        }}>
          <Link to="/projects" style={{ 
            textDecoration: "none", 
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            flex: 1,
            width: "100%"
          }}>
            <div style={{
              fontSize: "2.5rem",
              marginBottom: "0.5rem"
            }}>
              📊
            </div>
            <h3 style={{
              fontSize: "1.1rem",
              color: "#333",
              marginBottom: "0.5rem",
              textAlign: "center"
            }}>
              工程表
            </h3>
            <p style={{
              color: "#666",
              fontSize: "0.85rem",
              textAlign: "center",
              lineHeight: "1.4",
              marginBottom: "0.5rem"
            }}>
              工事プロジェクトの一覧管理
            </p>
          </Link>
          
          {/* ガントチャートリンク */}
          <Link 
            to="/projects/gantt" 
            style={{
              color: "#667eea",
              fontSize: "0.8rem",
              textDecoration: "none",
              cursor: "pointer",
              transition: "all 0.2s",
              position: "absolute",
              bottom: "12px",
              left: "50%",
              transform: "translateX(-50%)",
              display: "flex",
              alignItems: "center",
              gap: "4px",
              borderBottom: "1px solid transparent"
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.color = "#4f46e5";
              e.currentTarget.style.borderBottomColor = "#4f46e5";
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.color = "#667eea";
              e.currentTarget.style.borderBottomColor = "transparent";
            }}
            onClick={(e) => {
              e.stopPropagation();
            }}
          >
            📈 ガントチャートを見る
          </Link>
        </div>

      </div>

      <div style={{
        marginTop: "4rem",
        padding: "2rem",
        background: "#f8f9fa",
        borderRadius: "12px",
        textAlign: "center",
        border: "1px solid #e9ecef"
      }}>
        <h3 style={{
          color: "#333",
          marginBottom: "1rem",
          fontWeight: "600"
        }}>
          システム情報
        </h3>
        <p style={{
          color: "#666",
          lineHeight: "1.8"
        }}>
          データ保存場所: ~/penguin<br />
          工事プロジェクト: ~/penguin/豊田築炉/2-工事
        </p>
      </div>
      </div>
    </div>
  );
}