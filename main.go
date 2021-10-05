package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/subosito/gotenv"

	log "github.com/sirupsen/logrus"
)

func main() {
	Run()
}

func Router() http.Handler {
	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/follower/{username}").Handler(Follower())
	r.Methods(http.MethodGet).Path("/{userid}/detail").Handler(Detail())

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	return loggedRouter
}

func Follower() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)

		var follow int
		switch username["username"] {
		case "SammyShark":
			follow = 987
		case "JesseOctopus":
			follow = 432
		case "DrewSquid":
			follow = 321
		case "JamieMantisShrimp":
			follow = 654
		default:
			follow = 0

		}

		var data = map[string]int{
			"Followers": follow,
		}

		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(&data)

	}
}

func Detail() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)

		var follow int
		var user string
		switch username["userid"] {
		case "sammy":
			follow = 987
			user = "SammyShark"
		case "jesse":
			follow = 432
			user = "JesseOctopus"
		case "drew":
			follow = 321
			user = "DrewSquid"
		case "jamie":
			follow = 654
			user = "JamieMantisShrimp"
		default:
			follow = 0
			user = ""

		}

		var data map[string]interface{} = map[string]interface{}{
			"Username": user,
			"Follower": follow,
		}

		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(&data)

	}
}

func Run() {
	defaultEnv := ".env"

	if err := gotenv.Load(defaultEnv); err != nil {
		log.Warning("failed load .env")
	}

	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: Router(),
	}

	log.Println(fmt.Sprintf("starting application { %v } on port :%v", "Home Work", port))

	go listenAndServe(srv)
	waitForShutdown(srv)
}

func listenAndServe(apiServer *http.Server) {
	err := apiServer.ListenAndServe()

	if err != nil {
		log.WithField("error", err.Error()).Fatal("unable to serve")
	}
}

func waitForShutdown(apiServer *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)

	<-sig

	log.Warn("shutting down")

	if err := apiServer.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}

	log.Warn("shutdown complete")
}
