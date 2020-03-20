package cmd

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func singleJoiningSlash(a, b string) string {
	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")
	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		return a + "/" + b
	}
	return a + b
}

func Server(cmd *cobra.Command, agrs []string) {

	log.Printf(cmd.Short)

	log.Print("Config file used: ", viper.ConfigFileUsed())
	port := viper.GetString(HttpPortFlag)
	log.Printf("Using HTTP Port: %s", port)

	path := "/*catchall"
	origin, parseErr := url.Parse("http://localhost:9091/")
	if parseErr != nil {
		log.Fatalln(parseErr)
		return
	}
	log.Printf("Origin URL: %s", origin.String())

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			log.Println("*** Entering Director...")
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = origin.Scheme
			req.URL.Host = origin.Host

			wildcardIndex := strings.IndexAny(path, "*")
			proxyPath := singleJoiningSlash(origin.Path, req.URL.Path[wildcardIndex:])
			if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
				proxyPath = proxyPath[:len(proxyPath)-1]
			}
			req.URL.Path = proxyPath
		},
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).Dial,
		},
		ModifyResponse: func(r *http.Response) error {
			return nil
		},
		ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
			fmt.Printf("error was: %+v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		},
	}

	router := httprouter.New()
	router.Handle("GET", path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Reverse proxy is listening to the port %s", port)
	httpErr := http.ListenAndServe(":"+port, router)
	log.Fatal(httpErr)
	return
}
