package server

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	gamev1alpha1 "github.com/unmango/game/gen/dev/unmango/game/v1alpha1"
	"github.com/unmango/game/gen/dev/unmango/game/v1alpha1/gamev1alpha1connect"

	"github.com/unmango/game/identity"
)

// Identity serves IdentityService.
type Identity struct {
	identity.Identity
}

var _ gamev1alpha1connect.IdentityServiceHandler = Identity{}

// GetIdentity implements IdentityService.
func (s Identity) GetIdentity(context.Context, *connect.Request[gamev1alpha1.GetIdentityRequest]) (*connect.Response[gamev1alpha1.GetIdentityResponse], error) {
	return connect.NewResponse(&gamev1alpha1.GetIdentityResponse{
		PublicKey: s.PublicKey(),
		CreatedAt: timestamppb.New(s.CreatedAt),
	}), nil
}
