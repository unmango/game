// Command game serves one player's root: identity and calculator.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/unmango/game/identity"
	"github.com/unmango/game/server"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dataDir := flag.String("data-dir", defaultDataDir(), "directory holding the identity file")
	flag.Parse()

	id, err := identity.LoadOrGenerate(*dataDir, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("identity %x created %s", id.PublicKey(), id.CreatedAt.Format(time.RFC3339))

	srv := &http.Server{
		Addr:              *addr,
		Handler:           server.Handler(id),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()

	log.Printf("listening on %s", *addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func defaultDataDir() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return fmt.Sprintf("%s/unmango-game", dir)
	}
	return ".game"
}
