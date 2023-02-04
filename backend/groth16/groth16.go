// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package groth16 implements Groth16 Zero Knowledge Proof system  (aka zkSNARK).
//
// # See also
//
// https://eprint.iacr.org/2016/260.pdf
package groth16

import (
	"io"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/nume-crypto/gnark/backend"
	"github.com/nume-crypto/gnark/backend/witness"
	"github.com/nume-crypto/gnark/frontend"
	backend_bls12377 "github.com/nume-crypto/gnark/internal/backend/bls12-377/cs"
	backend_bls12381 "github.com/nume-crypto/gnark/internal/backend/bls12-381/cs"
	backend_bls24315 "github.com/nume-crypto/gnark/internal/backend/bls24-315/cs"
	backend_bn254 "github.com/nume-crypto/gnark/internal/backend/bn254/cs"
	backend_bw6633 "github.com/nume-crypto/gnark/internal/backend/bw6-633/cs"
	backend_bw6761 "github.com/nume-crypto/gnark/internal/backend/bw6-761/cs"

	witness_bls12377 "github.com/nume-crypto/gnark/internal/backend/bls12-377/witness"
	witness_bls12381 "github.com/nume-crypto/gnark/internal/backend/bls12-381/witness"
	witness_bls24315 "github.com/nume-crypto/gnark/internal/backend/bls24-315/witness"
	witness_bn254 "github.com/nume-crypto/gnark/internal/backend/bn254/witness"
	witness_bw6633 "github.com/nume-crypto/gnark/internal/backend/bw6-633/witness"
	witness_bw6761 "github.com/nume-crypto/gnark/internal/backend/bw6-761/witness"

	gnarkio "github.com/nume-crypto/gnark/io"

	groth16_bls12377 "github.com/nume-crypto/gnark/internal/backend/bls12-377/groth16"
	groth16_bls12381 "github.com/nume-crypto/gnark/internal/backend/bls12-381/groth16"
	groth16_bls24315 "github.com/nume-crypto/gnark/internal/backend/bls24-315/groth16"
	groth16_bn254 "github.com/nume-crypto/gnark/internal/backend/bn254/groth16"
	groth16_bw6633 "github.com/nume-crypto/gnark/internal/backend/bw6-633/groth16"
	groth16_bw6761 "github.com/nume-crypto/gnark/internal/backend/bw6-761/groth16"
)

type groth16Object interface {
	gnarkio.WriterRawTo
	io.WriterTo
	io.ReaderFrom
	CurveID() ecc.ID
}

// Proof represents a Groth16 proof generated by groth16.Prove
//
// it's underlying implementation is curve specific (see gnark/internal/backend)
type Proof interface {
	groth16Object
}

// ProvingKey represents a Groth16 ProvingKey
//
// it's underlying implementation is strongly typed with the curve (see gnark/internal/backend)
type ProvingKey interface {
	groth16Object
	gnarkio.UnsafeReaderFrom

	// NbG1 returns the number of G1 elements in the ProvingKey
	NbG1() int

	// NbG2 returns the number of G2 elements in the ProvingKey
	NbG2() int

	IsDifferent(interface{}) bool
}

// VerifyingKey represents a Groth16 VerifyingKey
//
// it's underlying implementation is strongly typed with the curve (see gnark/internal/backend)
//
// ExportSolidity is implemented for BN254 and will return an error with other curves
type VerifyingKey interface {
	groth16Object
	gnarkio.UnsafeReaderFrom

	// NbPublicWitness returns number of elements expected in the public witness
	NbPublicWitness() int

	// NbG1 returns the number of G1 elements in the VerifyingKey
	NbG1() int

	// NbG2 returns the number of G2 elements in the VerifyingKey
	NbG2() int

	// ExportSolidity writes a solidity Verifier contract from the VerifyingKey
	// this will return an error if not supported on the CurveID()
	ExportSolidity(w io.Writer) error

	IsDifferent(interface{}) bool
}

// Verify runs the groth16.Verify algorithm on provided proof with given witness
func Verify(proof Proof, vk VerifyingKey, publicWitness *witness.Witness) error {

	switch _proof := proof.(type) {
	case *groth16_bls12377.Proof:
		w, ok := publicWitness.Vector.(*witness_bls12377.Witness)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return groth16_bls12377.Verify(_proof, vk.(*groth16_bls12377.VerifyingKey), *w)
	case *groth16_bls12381.Proof:
		w, ok := publicWitness.Vector.(*witness_bls12381.Witness)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return groth16_bls12381.Verify(_proof, vk.(*groth16_bls12381.VerifyingKey), *w)
	case *groth16_bn254.Proof:
		w, ok := publicWitness.Vector.(*witness_bn254.Witness)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return groth16_bn254.Verify(_proof, vk.(*groth16_bn254.VerifyingKey), *w)
	case *groth16_bw6761.Proof:
		w, ok := publicWitness.Vector.(*witness_bw6761.Witness)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return groth16_bw6761.Verify(_proof, vk.(*groth16_bw6761.VerifyingKey), *w)
	case *groth16_bls24315.Proof:
		w, ok := publicWitness.Vector.(*witness_bls24315.Witness)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return groth16_bls24315.Verify(_proof, vk.(*groth16_bls24315.VerifyingKey), *w)
	case *groth16_bw6633.Proof:
		w, ok := publicWitness.Vector.(*witness_bw6633.Witness)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return groth16_bw6633.Verify(_proof, vk.(*groth16_bw6633.VerifyingKey), *w)
	default:
		panic("unrecognized R1CS curve type")
	}
}

// Prove runs the groth16.Prove algorithm.
//
// if the force flag is set:
//
//		will executes all the prover computations, even if the witness is invalid
//	 will produce an invalid proof
//		internally, the solution vector to the R1CS will be filled with random values which may impact benchmarking
func Prove(r1cs frontend.CompiledConstraintSystem, pk ProvingKey, fullWitness *witness.Witness, opts ...backend.ProverOption) (Proof, error) {

	// apply options
	opt, err := backend.NewProverConfig(opts...)
	if err != nil {
		return nil, err
	}

	switch _r1cs := r1cs.(type) {
	case *backend_bls12377.R1CS:
		w, ok := fullWitness.Vector.(*witness_bls12377.Witness)
		if !ok {
			return nil, witness.ErrInvalidWitness
		}
		return groth16_bls12377.Prove(_r1cs, pk.(*groth16_bls12377.ProvingKey), *w, opt)
	case *backend_bls12381.R1CS:
		w, ok := fullWitness.Vector.(*witness_bls12381.Witness)
		if !ok {
			return nil, witness.ErrInvalidWitness
		}
		return groth16_bls12381.Prove(_r1cs, pk.(*groth16_bls12381.ProvingKey), *w, opt)
	case *backend_bn254.R1CS:
		w, ok := fullWitness.Vector.(*witness_bn254.Witness)
		if !ok {
			return nil, witness.ErrInvalidWitness
		}
		return groth16_bn254.Prove(_r1cs, pk.(*groth16_bn254.ProvingKey), *w, opt)
	case *backend_bw6761.R1CS:
		w, ok := fullWitness.Vector.(*witness_bw6761.Witness)
		if !ok {
			return nil, witness.ErrInvalidWitness
		}
		return groth16_bw6761.Prove(_r1cs, pk.(*groth16_bw6761.ProvingKey), *w, opt)
	case *backend_bls24315.R1CS:
		w, ok := fullWitness.Vector.(*witness_bls24315.Witness)
		if !ok {
			return nil, witness.ErrInvalidWitness
		}
		return groth16_bls24315.Prove(_r1cs, pk.(*groth16_bls24315.ProvingKey), *w, opt)
	case *backend_bw6633.R1CS:
		w, ok := fullWitness.Vector.(*witness_bw6633.Witness)
		if !ok {
			return nil, witness.ErrInvalidWitness
		}
		return groth16_bw6633.Prove(_r1cs, pk.(*groth16_bw6633.ProvingKey), *w, opt)
	default:
		panic("unrecognized R1CS curve type")
	}
}

// Setup runs groth16.Setup with provided R1CS and outputs a key pair associated with the circuit.
//
// Note that careful consideration must be given to this step in production environment.
// groth16.Setup uses some randomness to precompute the Proving and Verifying keys. If the process
// or machine leaks this randomness, an attacker could break the ZKP protocol.
//
// Two main solutions to this deployment issues are: running the Setup through a MPC (multi party computation)
// or using a ZKP backend like PLONK where the per-circuit Setup is deterministic.
func Setup(r1cs frontend.CompiledConstraintSystem) (ProvingKey, VerifyingKey, error) {

	switch _r1cs := r1cs.(type) {
	case *backend_bls12377.R1CS:
		var pk groth16_bls12377.ProvingKey
		var vk groth16_bls12377.VerifyingKey
		if err := groth16_bls12377.Setup(_r1cs, &pk, &vk); err != nil {
			return nil, nil, err
		}
		return &pk, &vk, nil
	case *backend_bls12381.R1CS:
		var pk groth16_bls12381.ProvingKey
		var vk groth16_bls12381.VerifyingKey
		if err := groth16_bls12381.Setup(_r1cs, &pk, &vk); err != nil {
			return nil, nil, err
		}
		return &pk, &vk, nil
	case *backend_bn254.R1CS:
		var pk groth16_bn254.ProvingKey
		var vk groth16_bn254.VerifyingKey
		if err := groth16_bn254.Setup(_r1cs, &pk, &vk); err != nil {
			return nil, nil, err
		}
		return &pk, &vk, nil
	case *backend_bw6761.R1CS:
		var pk groth16_bw6761.ProvingKey
		var vk groth16_bw6761.VerifyingKey
		if err := groth16_bw6761.Setup(_r1cs, &pk, &vk); err != nil {
			return nil, nil, err
		}
		return &pk, &vk, nil
	case *backend_bls24315.R1CS:
		var pk groth16_bls24315.ProvingKey
		var vk groth16_bls24315.VerifyingKey
		if err := groth16_bls24315.Setup(_r1cs, &pk, &vk); err != nil {
			return nil, nil, err
		}
		return &pk, &vk, nil
	case *backend_bw6633.R1CS:
		var pk groth16_bw6633.ProvingKey
		var vk groth16_bw6633.VerifyingKey
		if err := groth16_bw6633.Setup(_r1cs, &pk, &vk); err != nil {
			return nil, nil, err
		}
		return &pk, &vk, nil
	default:
		panic("unrecognized R1CS curve type")
	}
}

// DummySetup create a random ProvingKey with provided R1CS
// it doesn't return a VerifyingKey and is use for benchmarking or test purposes only.
func DummySetup(r1cs frontend.CompiledConstraintSystem) (ProvingKey, error) {
	switch _r1cs := r1cs.(type) {
	case *backend_bls12377.R1CS:
		var pk groth16_bls12377.ProvingKey
		if err := groth16_bls12377.DummySetup(_r1cs, &pk); err != nil {
			return nil, err
		}
		return &pk, nil
	case *backend_bls12381.R1CS:
		var pk groth16_bls12381.ProvingKey
		if err := groth16_bls12381.DummySetup(_r1cs, &pk); err != nil {
			return nil, err
		}
		return &pk, nil
	case *backend_bn254.R1CS:
		var pk groth16_bn254.ProvingKey
		if err := groth16_bn254.DummySetup(_r1cs, &pk); err != nil {
			return nil, err
		}
		return &pk, nil
	case *backend_bw6761.R1CS:
		var pk groth16_bw6761.ProvingKey
		if err := groth16_bw6761.DummySetup(_r1cs, &pk); err != nil {
			return nil, err
		}
		return &pk, nil
	case *backend_bls24315.R1CS:
		var pk groth16_bls24315.ProvingKey
		if err := groth16_bls24315.DummySetup(_r1cs, &pk); err != nil {
			return nil, err
		}
		return &pk, nil
	case *backend_bw6633.R1CS:
		var pk groth16_bw6633.ProvingKey
		if err := groth16_bw6633.DummySetup(_r1cs, &pk); err != nil {
			return nil, err
		}
		return &pk, nil
	default:
		panic("unrecognized R1CS curve type")
	}
}

// NewProvingKey instantiates a curve-typed ProvingKey and returns an interface object
// This function exists for serialization purposes
func NewProvingKey(curveID ecc.ID) ProvingKey {
	var pk ProvingKey
	switch curveID {
	case ecc.BN254:
		pk = &groth16_bn254.ProvingKey{}
	case ecc.BLS12_377:
		pk = &groth16_bls12377.ProvingKey{}
	case ecc.BLS12_381:
		pk = &groth16_bls12381.ProvingKey{}
	case ecc.BW6_761:
		pk = &groth16_bw6761.ProvingKey{}
	case ecc.BLS24_315:
		pk = &groth16_bls24315.ProvingKey{}
	case ecc.BW6_633:
		pk = &groth16_bw6633.ProvingKey{}
	default:
		panic("not implemented")
	}
	return pk
}

// NewVerifyingKey instantiates a curve-typed VerifyingKey and returns an interface
// This function exists for serialization purposes
func NewVerifyingKey(curveID ecc.ID) VerifyingKey {
	var vk VerifyingKey
	switch curveID {
	case ecc.BN254:
		vk = &groth16_bn254.VerifyingKey{}
	case ecc.BLS12_377:
		vk = &groth16_bls12377.VerifyingKey{}
	case ecc.BLS12_381:
		vk = &groth16_bls12381.VerifyingKey{}
	case ecc.BW6_761:
		vk = &groth16_bw6761.VerifyingKey{}
	case ecc.BLS24_315:
		vk = &groth16_bls24315.VerifyingKey{}
	case ecc.BW6_633:
		vk = &groth16_bw6633.VerifyingKey{}
	default:
		panic("not implemented")
	}

	return vk
}

// NewProof instantiates a curve-typed Proof and returns an interface
// This function exists for serialization purposes
func NewProof(curveID ecc.ID) Proof {
	var proof Proof
	switch curveID {
	case ecc.BN254:
		proof = &groth16_bn254.Proof{}
	case ecc.BLS12_377:
		proof = &groth16_bls12377.Proof{}
	case ecc.BLS12_381:
		proof = &groth16_bls12381.Proof{}
	case ecc.BW6_761:
		proof = &groth16_bw6761.Proof{}
	case ecc.BLS24_315:
		proof = &groth16_bls24315.Proof{}
	case ecc.BW6_633:
		proof = &groth16_bw6633.Proof{}
	default:
		panic("not implemented")
	}

	return proof
}

// NewCS instantiate a concrete curved-typed R1CS and return a R1CS interface
// This method exists for (de)serialization purposes
func NewCS(curveID ecc.ID) frontend.CompiledConstraintSystem {
	var r1cs frontend.CompiledConstraintSystem
	switch curveID {
	case ecc.BN254:
		r1cs = &backend_bn254.R1CS{}
	case ecc.BLS12_377:
		r1cs = &backend_bls12377.R1CS{}
	case ecc.BLS12_381:
		r1cs = &backend_bls12381.R1CS{}
	case ecc.BW6_761:
		r1cs = &backend_bw6761.R1CS{}
	case ecc.BLS24_315:
		r1cs = &backend_bls24315.R1CS{}
	case ecc.BW6_633:
		r1cs = &backend_bw6633.R1CS{}
	default:
		panic("not implemented")
	}
	return r1cs
}
