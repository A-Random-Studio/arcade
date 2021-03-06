package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/bseto/arcade/backend/game/gamefactory"
	"github.com/bseto/arcade/backend/hub/hubmanager"
	"github.com/bseto/arcade/backend/log"
	"github.com/bseto/arcade/backend/websocket"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var port *int = flag.Int("port", 8081, "defines the port to listen and serve on")

func main() {
	flag.Parse()
	initializeLogging()
	initializeRoutes()
}

func initializeRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))

	hubManager := hubmanager.GetHubManager()
	hubManager.SetupRoutes(r)
	gameFactory := gamefactory.GetGameFactory()

	r.PathPrefix("/ws/{hubID}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hubInstance, err := hubManager.GetHub(req, gameFactory)
		if err != nil {
			log.Errorf("unable to ")
		}

		wsClient := websocket.GetClient(hubInstance)
		err = wsClient.Upgrade(w, req)
		if err != nil {
			log.Errorf("unable to upgrade websocket: %v", err)
			return
		}
		wsClient.RegisterCloseListener(hubManager)
	})

	address := fmt.Sprintf(":%v", *port)
	log.Infof("starting server on: %v", address)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"}) // we need to add our domain name here one day
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	err := http.ListenAndServe(address, handlers.CORS(originsOk, headersOk, methodsOk)(r))
	log.Fatalf("unable to listen and serve: %v", err)
}

func initializeLogging() {
	config := log.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      log.Info,
		ConsoleJSONFormat: false,
		EnableFile:        true,
		FileLevel:         log.Debug,
		FileJSONFormat:    true,
		FileLocation:      "log.log",
	}
	err := log.NewLogger(config, log.InstanceLogrusLogger)
	if err != nil {
		log.Fatalf("Could not instantiate log: %v", err.Error())
	}
	log.Infof("logging initialized")
}
