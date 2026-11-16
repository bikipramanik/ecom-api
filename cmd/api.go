package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	// A good base middleware stack

	/*	Suppose a user reports an API error.

		You can search your server logs using that request ID and trace the request.

		It is also useful when monitoring APIs and debugging production issues.

		Note: RequestID does not implement rate limiting itself. It helps identify individual requests
	*/
	r.Use(middleware.RequestID)              //important for rate imiting
	r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
	r.Use(middleware.Logger)                 //print logs

	/*Without Recoverer, a panic can crash the HTTP server's serving goroutine.

	With Recoverer, Chi recovers from the panic, logs it, and returns an HTTP 500 response.

	It helps prevent one bad request from crashing your server.*/
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})
	// http.ListenAndServe(":3333", r)
	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s",app.config.addr)
	return srv.ListenAndServe()
}

type config struct {
	addr string
	db   dbconfig
}

type dbconfig struct {
	dsn string
}
