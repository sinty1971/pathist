package main

import (
	"backend-grpc/gen/grpc/v1/grpcv1connect"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"backend-grpc/internal/services"

	"connectrpc.com/grpcreflect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// コマンドライン引数
var (
	httpAddr  = flag.String("http-addr", ":9090", "HTTP (h2c) の待受アドレス")
	httpsAddr = flag.String("https-addr", ":9443", "HTTPS の待受アドレス")
	enableTLS = flag.Bool("enable-tls", false, "true の場合は HTTPS も起動します")
	certPath  = flag.String("cert", "cert.pem", "TLS 証明書のパス")
	keyPath   = flag.String("key", "key.pem", "TLS 秘密鍵のパス")
)

func main() {
	// コマンドライン引数の解析
	flag.Parse()

	// シグナルの受信を待機するコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	services, err := services.NewServices()
	if err != nil {
		log.Fatalf("サービスの初期化に失敗しました: %v", err)
	}
	defer services.Cleanup()

	if services.KojiService == nil {
		log.Fatal("KojiService が初期化されていません")
	}
	if services.CompanyService == nil {
		log.Fatal("CompanyService が初期化されていません")
	}
	if services.FileService == nil {
		log.Fatal("FileService が初期化されていません")
	}

	mux := http.NewServeMux()

	filePath, fileConnectHandler := grpcv1connect.NewFileServiceHandler(services.FileService)
	mux.Handle(filePath, fileConnectHandler)

	kojiPath, kojiConnectHandler := grpcv1connect.NewKojiServiceHandler(services.KojiService)
	mux.Handle(kojiPath, kojiConnectHandler)

	companyPath, companyConnectHandler := grpcv1connect.NewCompanyServiceHandler(services.CompanyService)
	mux.Handle(companyPath, companyConnectHandler)

	reflector := grpcreflect.NewStaticReflector(
		grpcv1connect.FileServiceName,
		grpcv1connect.CompanyServiceName,
		grpcv1connect.KojiServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	httpServer := &http.Server{
		Addr:    *httpAddr,
		Handler: h2c.NewHandler(cors(mux), &http2.Server{}),
	}

	go func() {
		log.Printf("HTTP gRPC サーバーを %s で起動します", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP サーバーでエラーが発生しました: %v", err)
		}
	}()

	var httpsServer *http.Server
	if *enableTLS {
		if err := ensureCertificate(*certPath, *keyPath); err != nil {
			log.Fatalf("TLS 証明書の準備に失敗しました: %v", err)
		}

		httpsServer = &http.Server{
			Addr:      *httpsAddr,
			Handler:   cors(mux),
			TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		}

		go func() {
			log.Printf("HTTPS gRPC サーバーを %s で起動します", httpsServer.Addr)
			if err := httpsServer.ListenAndServeTLS(*certPath, *keyPath); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS サーバーでエラーが発生しました: %v", err)
			}
		}()
	}

	<-ctx.Done()
	log.Printf("停止シグナルを受信しました。サーバーをシャットダウンします。")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP サーバーの停止に失敗しました: %v", err)
	}

	if httpsServer != nil {
		if err := httpsServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTPS サーバーの停止に失敗しました: %v", err)
		}
	}

	log.Printf("gRPC サーバーを停止しました。")
}

func ensureCertificate(certFile, keyFile string) error {
	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(certFile), 0o755); err != nil {
		return err
	}

	log.Printf("自己署名証明書を %s に生成します", certFile)

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName:   "Penguin Backend gRPC",
			Organization: []string{"Penguin Backend"},
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	if err := writePEM(certFile, "CERTIFICATE", certDER); err != nil {
		return err
	}

	keyDER := x509.MarshalPKCS1PrivateKey(priv)
	if err := writePEM(keyFile, "RSA PRIVATE KEY", keyDER); err != nil {
		return err
	}

	return nil
}

func writePEM(path, typ string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, &pem.Block{Type: typ, Bytes: data})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
