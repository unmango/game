// Package server implements the ConnectRPC services over the pure packages.
package server

import (
	"context"

	"connectrpc.com/connect"

	gamev1alpha1 "github.com/unmango/game/gen/dev/unmango/game/v1alpha1"
	"github.com/unmango/game/gen/dev/unmango/game/v1alpha1/gamev1alpha1connect"

	"github.com/unmango/game/curve"
	"github.com/unmango/game/seed"
)

// Calculator serves CalculatorService for one root seed.
type Calculator struct {
	Seed int64
}

var _ gamev1alpha1connect.CalculatorServiceHandler = Calculator{}

// Derive implements CalculatorService.
func (c Calculator) Derive(_ context.Context, req *connect.Request[gamev1alpha1.DeriveRequest]) (*connect.Response[gamev1alpha1.DeriveResponse], error) {
	return connect.NewResponse(&gamev1alpha1.DeriveResponse{
		Seed: seed.Derive(c.Seed, req.Msg.GetPath()),
	}), nil
}

// Evaluate implements CalculatorService.
func (c Calculator) Evaluate(_ context.Context, req *connect.Request[gamev1alpha1.EvaluateRequest]) (*connect.Response[gamev1alpha1.EvaluateResponse], error) {
	cv, err := CurveFromProto(req.Msg.GetCurve())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&gamev1alpha1.EvaluateResponse{
		Value: NumberToProto(cv.Evaluate(req.Msg.GetN())),
	}), nil
}

// Cumulative implements CalculatorService.
func (c Calculator) Cumulative(_ context.Context, req *connect.Request[gamev1alpha1.CumulativeRequest]) (*connect.Response[gamev1alpha1.CumulativeResponse], error) {
	cv, err := CurveFromProto(req.Msg.GetCurve())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&gamev1alpha1.CumulativeResponse{
		Value: NumberToProto(curve.Cumulative(cv, req.Msg.GetFrom(), req.Msg.GetTo())),
	}), nil
}

// Invert implements CalculatorService.
func (c Calculator) Invert(_ context.Context, req *connect.Request[gamev1alpha1.InvertRequest]) (*connect.Response[gamev1alpha1.InvertResponse], error) {
	cv, err := CurveFromProto(req.Msg.GetCurve())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&gamev1alpha1.InvertResponse{
		Count: curve.Invert(cv, req.Msg.GetFrom(), NumberFromProto(req.Msg.GetBudget())),
	}), nil
}

// Advance implements CalculatorService.
func (c Calculator) Advance(_ context.Context, req *connect.Request[gamev1alpha1.AdvanceRequest]) (*connect.Response[gamev1alpha1.AdvanceResponse], error) {
	cv, err := CurveFromProto(req.Msg.GetRate())
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&gamev1alpha1.AdvanceResponse{
		Value: NumberToProto(curve.Advance(cv, req.Msg.GetLevel(), req.Msg.GetElapsed().AsDuration())),
	}), nil
}
