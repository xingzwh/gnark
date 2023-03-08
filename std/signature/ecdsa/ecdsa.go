/*
Package ecdsa implements ECDSA signature verification over any elliptic curve.

The package depends on the [weierstrass] package for elliptic curve group
operations using non-native arithmetic. Thus we can verify ECDSA signatures over
any curve. The cost for a single secp256k1 signature verification is
approximately 4M constraints in R1CS and 10M constraints in PLONKish.

See [ECDSA] for the signature verification algorithm.

[ECDSA]:
https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm
*/
package ecdsa

import (
	"github.com/xingzwh/gnark/frontend"
	"github.com/xingzwh/gnark/std/algebra/weierstrass"
	"github.com/xingzwh/gnark/std/math/emulated"
)

// Signature represents the signature for some message.
type Signature[Scalar emulated.FieldParams] struct {
	R, S emulated.Element[Scalar]
}

// PublicKey represents the public key to verify the signature for.
type PublicKey[Base, Scalar emulated.FieldParams] weierstrass.AffinePoint[Base]

// Verify asserts that the signature sig verifies for the message msg and public
// key pk. The curve parameters params define the elliptic curve.
//
// We assume that the message msg is already hashed to the scalar field.
func (pk PublicKey[T, S]) Verify(api frontend.API, params weierstrass.CurveParams, msg *emulated.Element[S], sig *Signature[S]) {
	cr, err := weierstrass.New[T, S](api, params)
	if err != nil {
		// TODO: softer handling.
		panic(err)
	}
	scalarApi, err := emulated.NewField[S](api)
	if err != nil {
		panic(err)
	}
	baseApi, err := emulated.NewField[T](api)
	if err != nil {
		panic(err)
	}
	pkpt := weierstrass.AffinePoint[T](pk)
	sInv := scalarApi.Inverse(&sig.S)
	msInv := scalarApi.MulMod(msg, sInv)
	rsInv := scalarApi.MulMod(&sig.R, sInv)

	qa := cr.ScalarMul(cr.Generator(), msInv)
	qb := cr.ScalarMul(&pkpt, rsInv)
	q := cr.Add(qa, qb)
	qx := baseApi.Reduce(&q.X)
	qxBits := baseApi.ToBits(qx)
	rbits := scalarApi.ToBits(&sig.R)
	if len(rbits) != len(qxBits) {
		panic("non-equal lengths")
	}
	for i := range rbits {
		api.AssertIsEqual(rbits[i], qxBits[i])
	}
}
