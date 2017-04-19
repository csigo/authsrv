package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"gopkg.in/session.v1"

	"github.com/csigo/authsrv"
)

var (
	globalSessions *session.Manager
	users          map[string]*authsrv.User
	clients        map[string]*authsrv.Client
	clientStore    *store.ClientStore
	manager        *manage.Manager
)

var (
	hostAddr   = flag.String("host.addr", ":8898", "Specify server's host:port")
	cbAddr     = flag.String("callback.addr", "localhost:8899", "Specify callback host:port")
	clientFile = flag.String("client.file", "client.json", "predefined client")
)

func init() {
	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid","gclifetime":3600}`)
	go globalSessions.GC()
}

func main() {
	flag.Parse()
	manager = manage.NewDefaultManager()
	// token store
	tstore, err := store.NewMemoryTokenStore()
	csiToken := models.NewToken()
	manager.MustTokenStorage(tstore, err)

	clientStore = store.NewClientStore()
	if _, err := os.Stat(*clientFile); err == nil {
		log.Println("has predefined client file")
		file, err := ioutil.ReadFile(*clientFile)
		if err != nil {
			fmt.Printf("File error: %v\n", err)
			os.Exit(1)
		}
		json.Unmarshal(file, &clients)
	} else {
		clients = map[string]*authsrv.Client{}
	}
	for _, u := range clients {
		var domain string
		if u.AuthURI != "" {
			url, err := url.Parse(u.AuthURI)
			if err != nil {
				log.Fatal(err)
			}
			domain = url.Hostname()
		} else {
			domain = fmt.Sprintf("http://%s", *hostAddr)
		}

		clientStore.Set(u.ClientID, &models.Client{
			ID:     u.ClientID,
			Secret: u.ClientSecret,
			Domain: domain,
		})
		// the following only for CSI internal test
		if csiToken.ClientID == "" {
			csiToken.ClientID = u.ClientID
			csiToken.UserID = "csiuser"
			csiToken.RedirectURI = u.RedirectURI
			csiToken.Scope = "read"
			csiToken.Access = "87654321"
			csiToken.AccessCreateAt = time.Now()
			csiToken.AccessExpiresIn = 1 * time.Hour
			tstore.Create(csiToken)
		}
	}
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		ti, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		user := &authsrv.User{
			Id:    ti.GetUserID(),
			Email: "csiuser@csigo.com",
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(*user)
	})

	http.HandleFunc("/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	_, port, err := net.SplitHostPort(*hostAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server is running at %s port.\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
