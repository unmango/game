package server_test

import (
	"context"
	"math"
	"net/http/httptest"
	"time"

	"connectrpc.com/connect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	gamev1alpha1 "github.com/unmango/game/gen/dev/unmango/game/v1alpha1"
	"github.com/unmango/game/gen/dev/unmango/game/v1alpha1/gamev1alpha1connect"

	"github.com/unmango/game/convert"
	"github.com/unmango/game/identity"
	"github.com/unmango/game/seed"
	"github.com/unmango/game/server"
)

var _ = Describe("Handler", func() {
	var (
		id         identity.Identity
		ts         *httptest.Server
		calc       gamev1alpha1connect.CalculatorServiceClient
		ident      gamev1alpha1connect.IdentityServiceClient
		ctx        = context.Background()
		expCurve   = &gamev1alpha1.Curve{Family: &gamev1alpha1.Curve_Exponential{Exponential: &gamev1alpha1.Exponential{Base: &gamev1alpha1.Number{Mantissa: 1, Exponent: 1}, Growth: 1.15}}}
		emptyCurve = &gamev1alpha1.Curve{}
	)

	BeforeEach(func() {
		var err error
		id, err = identity.Generate(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
		Expect(err).NotTo(HaveOccurred())
		ts = httptest.NewServer(server.Handler(id))
		DeferCleanup(ts.Close)
		calc = gamev1alpha1connect.NewCalculatorServiceClient(ts.Client(), ts.URL)
		ident = gamev1alpha1connect.NewIdentityServiceClient(ts.Client(), ts.URL)
	})

	It("derives a sub-seed from the identity's root", func() {
		res, err := calc.Derive(ctx, connect.NewRequest(&gamev1alpha1.DeriveRequest{Path: "character/strength"}))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Msg.GetSeed()).To(Equal(seed.Derive(id.Seed, "character/strength")))
	})

	It("evaluates a curve", func() {
		res, err := calc.Evaluate(ctx, connect.NewRequest(&gamev1alpha1.EvaluateRequest{Curve: expCurve, N: 10}))
		Expect(err).NotTo(HaveOccurred())
		got := convert.NumberFromProto(res.Msg.GetValue()).Float64()
		Expect(got).To(BeNumerically("~", 10*math.Pow(1.15, 10), 1e-9))
	})

	It("sums a curve", func() {
		res, err := calc.Cumulative(ctx, connect.NewRequest(&gamev1alpha1.CumulativeRequest{Curve: expCurve, From: 0, To: 3}))
		Expect(err).NotTo(HaveOccurred())
		Expect(convert.NumberFromProto(res.Msg.GetValue()).Float64()).To(BeNumerically("~", 34.725, 1e-9))
	})

	It("inverts a curve", func() {
		res, err := calc.Invert(ctx, connect.NewRequest(&gamev1alpha1.InvertRequest{Curve: expCurve, From: 0, Budget: &gamev1alpha1.Number{Mantissa: 4, Exponent: 1}}))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Msg.GetCount()).To(Equal(int64(3)))
	})

	It("advances a rate over a duration", func() {
		rate := &gamev1alpha1.Curve{Family: &gamev1alpha1.Curve_Linear{Linear: &gamev1alpha1.Linear{Intercept: &gamev1alpha1.Number{Mantissa: 2}}}}
		res, err := calc.Advance(ctx, connect.NewRequest(&gamev1alpha1.AdvanceRequest{Rate: rate, Level: 0, Elapsed: durationpb.New(90 * time.Second)}))
		Expect(err).NotTo(HaveOccurred())
		Expect(convert.NumberFromProto(res.Msg.GetValue()).Float64()).To(BeNumerically("~", 180, 1e-9))
	})

	It("rejects a curve with no family", func() {
		_, err := calc.Evaluate(ctx, connect.NewRequest(&gamev1alpha1.EvaluateRequest{Curve: emptyCurve}))
		Expect(connect.CodeOf(err)).To(Equal(connect.CodeInvalidArgument))
	})

	It("reports the identity", func() {
		res, err := ident.GetIdentity(ctx, connect.NewRequest(&gamev1alpha1.GetIdentityRequest{}))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.Msg.GetPublicKey()).To(Equal([]byte(id.PublicKey())))
		Expect(res.Msg.GetCreatedAt().AsTime()).To(Equal(id.CreatedAt))
	})
})
