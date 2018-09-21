package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"testRouterAPI/go_recipe/insert"
	"testRouterAPI/go_recipe/logging"
	"testRouterAPI/go_recipe/router"

	"github.com/sirupsen/logrus"
)

var (
	httpPort        = flag.String("http_port", "80", "the port for http traffic") // TODO: add TLS.
	shutdownTimeout = flag.Duration("shutdown_timeout", time.Second*8, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
)

func main() {
	flag.Parse()
	logging.Init()

	r := router.NewRouter()
	
	adddatas.View()
	qq := router.GetRecipe()
	fmt.Printf("%s", qq)
	// Log not-founds as well.
	r.NotFoundHandler = logging.Middleware(http.NotFoundHandler())

	srv := &http.Server{
		Addr:         "0.0.0.0:" + *httpPort,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Println("serving at port", *httpPort)
		// TODO: we should use TLS instead.
		if err := srv.ListenAndServe(); err != nil {
			logrus.Errorln(err)
		}
	}()

	// Wait for SIGINT (Ctrl+C) or SIGTERM.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // Block until we receive our signal.
	log.Println("shutting down")

	// Create the shutdown deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}
