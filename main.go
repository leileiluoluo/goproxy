package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/olzhy/goproxy/pkg/proxy"
)

var port = flag.String("serverPort", ":8080", "server port")

func main() {
	// server
	srv := http.Server{
		Addr:    *port,
		Handler: proxy.Proxy(),
	}

	// gracefully shutdown
	closed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		select {
		case e := <-sig:
			log.Printf("interrupt, err: %s", e)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); nil != err {
				log.Printf("server shutdown error, err: %s", err)
			}

			closed <- struct{}{}
		}
	}()

	if err := srv.ListenAndServe(); nil != err {
		if http.ErrServerClosed == err {
			log.Println("server shutdown gracefully")
			return
		}
		log.Printf("serve error, err: %s", err)
	}

	<-closed
}
