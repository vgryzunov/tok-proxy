package cmd

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	originUrl := viper.GetString(OriginUrlFlag)
	log.Printf("Origin URL: %s", originUrl)

	origin, parseErr := url.Parse(originUrl)
	if parseErr != nil {
		log.Fatalln(parseErr)
	}
	log.Printf("Origin URL: %s", origin.String())

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			log.Println("**** Director ****")

			headers := req.Header
			xAuthToken := headers.Get("X-Auth-Token")
			log.Printf("X-Auth-Token: %s", xAuthToken)

			xAuthUserId := headers.Get("X-Auth-Userid")
			log.Printf("X-Auth-Usedid: %s", xAuthUserId)

			xAuthEmail := headers.Get("X-Auth-Email")
			log.Printf("X-Auth-Usedid: %s", xAuthEmail)

			req.SetBasicAuth(xAuthEmail, xAuthToken)

			for name := range headers {
				if strings.HasPrefix(name, "X-Auth-") {
					headers.Del(name)
				}
			}

			headers.Add("X-Forwarded-Host", req.Host)
			headers.Add("X-Origin-Host", origin.Host)
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
			log.Printf("**** ModifyResponse ****")

			headers := r.Header
			for name, val := range headers {
				log.Printf("%s: %s", name, val)
			}
			return nil
		},
		ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
			log.Printf("error was: %+v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte(err.Error()))
		},
	}

	router := chi.NewRouter()
	// Base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(tokenAuthOnly)

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	r := chi.NewRouter()
	r.Route("/foo", func(r chi.Router) {

	})

	log.Printf("Reverse proxy is listening to the port %s", port)
	httpErr := http.ListenAndServe(":"+port, router)
	log.Fatal(httpErr)
	return
}

func tokenAuthOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("*** tokenAuthOnly handler ****")
		headers := r.Header
		xAuthToken := headers.Get("X-Auth-Token")
		if xAuthToken == "" {
			log.Printf("Not Authenticated. Missing authentication X-Auth-Token")
			http.Error(w, http.StatusText(403), 403)
			return
		}
		next.ServeHTTP(w, r)
	})
}
