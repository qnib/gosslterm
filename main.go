package main

import (
	"log"
	"net/http"
	"os"
	"crypto/tls"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/wunderlist/moxy"
	"github.com/codegangsta/cli"
	"github.com/zpatrick/go-config"

)

func AddSecurityHeaders(request *http.Request, response *http.Response) {
	//response.Header.Del("X-Powered-By")
	//response.Header.Set("X-Super-Secure", "Yes!!")
}


func main() {
	app := cli.NewApp()
	app.Name = "Reverse Proxy to terminate SSL for arbitrary http service"
	app.Usage = "gosslterm [options]"
	app.Version = "0.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "cert,c",
			Value: "cert.pem",
			Usage: "Path to SSL certificate.",
			EnvVar: "GOSSLTERM_CERT",
		},
		cli.StringFlag{
			Name:  "key,k",
			Value: "key.pem",
			Usage: "Path to SSL key.",
			EnvVar: "GOSSLTERM_KEY",
		},
		cli.StringFlag{
			Name:  "backend-addr",
			Value: "127.0.0.1:80",
			Usage: "Address to proxy to",
			EnvVar: "GOSSLTERM_BACKEND_ADDR",
		},
		cli.StringFlag{
			Name:  "frontend-addr",
			Value: ":8080",
			Usage: "Address to service proxy at",
			EnvVar: "GOSSLTERM_FRONTEND_ADDR",
		},
		cli.StringFlag{
			Name:   "log-level",
			Value:  "warn",
			Usage:  "Log level (warn: silent)",
			EnvVar: "LOG_LEVEL",
		},
	}
	app.Action = Run
	app.Run(os.Args)
}

func getTLSConfig(cPath, kPath string) tls.Config {
	log.Printf("Load cert '%s' and key '%s'", cPath, kPath)
	cert, err := tls.LoadX509KeyPair(cPath, kPath)
	if err != nil {
		log.Fatalf("error in tls.LoadX509KeyPair: %s", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	return config
}

func Run(ctx *cli.Context) {
	cfg := config.NewConfig([]config.Provider{config.NewCLI(ctx, true)})
	cPath, _ := cfg.String("cert")
	kPath, _ := cfg.String("key")
	tlsConfig := getTLSConfig(cPath, kPath)
	bAddr, _ := cfg.String("backend-addr")
	hosts := []string{bAddr}
	filters := []moxy.FilterFunc{}
	proxy := moxy.NewReverseProxy(hosts, filters)

	router := mux.NewRouter()
	router.HandleFunc("/", proxy.ServeHTTP)
	router.HandleFunc("/pi", proxy.ServeHTTP)
	router.HandleFunc("/pi/{{num}}", proxy.ServeHTTP)

	app := negroni.New(negroni.NewLogger())
	app.UseHandler(router)

	fAddr, _ := cfg.String("frontend-addr")
	log.Printf("Create http.Server on '%s'", fAddr)
	srv := &http.Server{
		Addr:      fAddr,
		Handler:   app,
		TLSConfig: &tlsConfig,
	}
	log.Fatal(srv.ListenAndServeTLS(cPath, kPath))
}
