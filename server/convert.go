package server

import (
	"errors"

	gamev1alpha1 "github.com/unmango/game/gen/dev/unmango/game/v1alpha1"

	"github.com/unmango/game/curve"
	"github.com/unmango/game/num"
)

var errNoFamily = errors.New("curve has no family set")

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
		return nil, errNoFamily
	}
}
