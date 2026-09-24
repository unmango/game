// Package convert translates between the wire types in gen and the Go types
// in num and curve. Servers and clients both use it.
package convert

import (
	"errors"

	gamev1alpha1 "github.com/unmango/game/gen/dev/unmango/game/v1alpha1"

	"github.com/unmango/game/curve"
	"github.com/unmango/game/num"
)

// ErrNoFamily is returned for a Curve with no family set.
var ErrNoFamily = errors.New("curve has no family set")

// ErrUnknownCurve is returned for a Go curve type with no wire form.
var ErrUnknownCurve = errors.New("curve type has no wire form")

// NumberFromProto converts a wire Number, treating nil as zero.
func NumberFromProto(n *gamev1alpha1.Number) num.Number {
	if n == nil {
		return num.Zero()
	}
	return num.New(n.GetMantissa(), n.GetExponent())
}

// NumberToProto converts a Number for the wire.
func NumberToProto(n num.Number) *gamev1alpha1.Number {
	return &gamev1alpha1.Number{Mantissa: n.Mantissa, Exponent: n.Exponent}
}

// CurveFromProto converts a wire Curve to its Go form.
func CurveFromProto(c *gamev1alpha1.Curve) (curve.Curve, error) {
	switch f := c.GetFamily().(type) {
	case *gamev1alpha1.Curve_Linear:
		return curve.Linear{
			Intercept: NumberFromProto(f.Linear.GetIntercept()),
			Slope:     NumberFromProto(f.Linear.GetSlope()),
		}, nil
	case *gamev1alpha1.Curve_Exponential:
		return curve.Exponential{
			Base:   NumberFromProto(f.Exponential.GetBase()),
			Growth: f.Exponential.GetGrowth(),
		}, nil
	case *gamev1alpha1.Curve_Polynomial:
		return curve.Polynomial{
			Scale:  NumberFromProto(f.Polynomial.GetScale()),
			Degree: f.Polynomial.GetDegree(),
			Offset: NumberFromProto(f.Polynomial.GetOffset()),
		}, nil
	case *gamev1alpha1.Curve_Logistic:
		return curve.Logistic{
			Max:       NumberFromProto(f.Logistic.GetMax()),
			Steepness: f.Logistic.GetSteepness(),
			Midpoint:  f.Logistic.GetMidpoint(),
		}, nil
	default:
		return nil, ErrNoFamily
	}
}

// CurveToProto converts a Go curve to its wire form.
func CurveToProto(c curve.Curve) (*gamev1alpha1.Curve, error) {
	switch v := c.(type) {
	case curve.Linear:
		return &gamev1alpha1.Curve{Family: &gamev1alpha1.Curve_Linear{Linear: &gamev1alpha1.Linear{
			Intercept: NumberToProto(v.Intercept),
			Slope:     NumberToProto(v.Slope),
		}}}, nil
	case curve.Exponential:
		return &gamev1alpha1.Curve{Family: &gamev1alpha1.Curve_Exponential{Exponential: &gamev1alpha1.Exponential{
			Base:   NumberToProto(v.Base),
			Growth: v.Growth,
		}}}, nil
	case curve.Polynomial:
		return &gamev1alpha1.Curve{Family: &gamev1alpha1.Curve_Polynomial{Polynomial: &gamev1alpha1.Polynomial{
			Scale:  NumberToProto(v.Scale),
			Degree: v.Degree,
			Offset: NumberToProto(v.Offset),
		}}}, nil
	case curve.Logistic:
		return &gamev1alpha1.Curve{Family: &gamev1alpha1.Curve_Logistic{Logistic: &gamev1alpha1.Logistic{
			Max:       NumberToProto(v.Max),
			Steepness: v.Steepness,
			Midpoint:  v.Midpoint,
		}}}, nil
	default:
		return nil, ErrUnknownCurve
	}
}
