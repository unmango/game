package server

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/unmango/game/gen/dev/unmango/game/v1alpha1/gamev1alpha1connect"

	"github.com/unmango/game/identity"
)

// Handler returns an HTTP handler serving every service for id.
// It accepts HTTP/1.1, HTTP/2 over cleartext, and therefore gRPC.
func Handler(id identity.Identity) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(gamev1alpha1connect.NewCalculatorServiceHandler(Calculator{Seed: id.Seed}))
	mux.Handle(gamev1alpha1connect.NewIdentityServiceHandler(Identity{id}))
	return h2c.NewHandler(mux, &http2.Server{})
}
