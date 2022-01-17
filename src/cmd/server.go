package main

import (
	"bytes"
	"crypto/tls"
	"dongpham/config"
	//"dongpham/crons"
	"dongpham/repository"
	"dongpham/rest"
	"dongpham/utils"
	"dongpham/version"
	b64 "encoding/base64"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
)

const BYTE_SPLIT string = "@@"

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%s: %s", err, debug.Stack())
				log.Printf("Recovered in service %+v", err)
				utils.ResponseError(utils.ERROR_UNKNOWN_ERROR, w)
			}
		}()

		startTime := time.Now().UnixNano()
		h.ServeHTTP(w, r)
		if config.Verbose {
			code := r.Header.Get("CF-Ipcountry")
			if len(code) > 0 {
				log.Printf("%s %s \t%dms \t%sb \t: %s\t--> %s", code, utils.GetRemoteIp(r), (time.Now().UnixNano()-startTime)/1000000, w.Header().Get("Expected-Size"), r.Method, "https://"+r.Host+r.URL.Path+"?"+r.URL.RawQuery)
			} else {
				log.Printf("%s \t%dms \t%sb \t: %s\t--> %s", utils.GetRemoteIp(r), (time.Now().UnixNano()-startTime)/1000000, w.Header().Get("Expected-Size"), r.Method, "https://"+r.Host+r.URL.Path+"?"+r.URL.RawQuery)
			}
		}
	})
}

func initHTTPServer(httpPort int) {

	router := mux.NewRouter().StrictSlash(true)
	router = rest.RegisterRoutes(router)
	router.HandleFunc("/ping", ping)

	srv := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		Addr:         ":" + strconv.Itoa(httpPort),
		Handler:      middleware(router),
	}

	hostname, _ := os.Hostname()
	log.Printf("%s : (%s) Starting HTTP server at %d. isMaster=%v, verbose=%v", hostname, version.Version, httpPort, config.IsMaster, config.Verbose)
	log.Printf("%s : (%s) Starting HTTP server at %d. isMaster=%v, verbose=%v",
		hostname, version.Version, httpPort, config.IsMaster, config.Verbose)
	log.Fatal(srv.ListenAndServe())
}

func initHttpsServer(httpsPort int) {
	if len(config.SSLKeyBase64) > 0 {
		//certPEMBlock, err := os.ReadFile("../cert/server.crt")
		//if err != nil {
		//	return
		//}
		//keyPEMBlock, err := os.ReadFile("../cert/server.private.key")
		//if err != nil {
		//	return
		//}

		sEnc, err := b64.StdEncoding.DecodeString(config.SSLKeyBase64)
		if err != nil {
			return
		}
		certs := bytes.Split(sEnc, []byte(BYTE_SPLIT))

		if len(certs) < 2 {
			return
		}
		cert, err := tls.X509KeyPair(certs[0], certs[1])
		if err != nil {
			log.Fatalf("Error creating x509 keypair from client cert file  and client key file %v", err)
			return
		}
		router := mux.NewRouter().StrictSlash(true)
		router = rest.RegisterRoutes(router)
		router.HandleFunc("/ping", ping)
		//router.Use("/", func(w http.ResponseWriter, req *http.Request) {
		//	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		//	w.Write([]byte("This is an example server.\n"))
		//})

		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			Certificates: []tls.Certificate{cert},
		}
		srv := &http.Server{
			Addr:         ":" + strconv.Itoa(httpsPort),
			Handler:      middleware(router),
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}

		ln, err := net.Listen("tcp", srv.Addr)
		if err != nil {
			return
		}

		defer ln.Close()
		hostname, _ := os.Hostname()
		tlsListener := tls.NewListener(ln, cfg)
		log.Printf("%s : (%s) Starting HTTPs server at %d. isMaster=%v, verbose=%v", hostname, version.Version, httpsPort, config.IsMaster, config.Verbose)
		//log.Printf("%s : (%s) Starting HTTP server at %d. isMaster=%v, verbose=%v")
		log.Fatal(srv.Serve(tlsListener))
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		repository.CloseConn()
		os.Exit(0)
	}()
}

// ping use to test if the server is still alive
func ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("PONG"))
	if err != nil {
		log.Printf("Ping %v", err)
	}

}

func main() {
	httpPort := flag.Int("port", config.HTTPPort, "which http port that server will be listening")
	config.Init()
	// init randome seed
	rand.Seed(time.Now().UTC().UnixNano())

	requireLoop := false
	if *httpPort > 0 {
		requireLoop = true
		go initHTTPServer(*httpPort)
		go initHttpsServer(config.HTTPSPort)
	}

	if requireLoop {
		SetupCloseHandler()
		select {}
	}

}
