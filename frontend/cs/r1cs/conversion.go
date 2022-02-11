/*
Copyright © 2020 ConsenSys

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package r1cs

import (
	"log"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/hint"
	"github.com/consensys/gnark/debug"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/schema"
	bls12377r1cs "github.com/consensys/gnark/internal/backend/bls12-377/cs"
	bls12381r1cs "github.com/consensys/gnark/internal/backend/bls12-381/cs"
	bls24315r1cs "github.com/consensys/gnark/internal/backend/bls24-315/cs"
	bn254r1cs "github.com/consensys/gnark/internal/backend/bn254/cs"
	bw6633r1cs "github.com/consensys/gnark/internal/backend/bw6-633/cs"
	bw6761r1cs "github.com/consensys/gnark/internal/backend/bw6-761/cs"
	"github.com/consensys/gnark/internal/backend/compiled"
	"github.com/consensys/gnark/internal/dag"
)

// Compile constructs a rank-1 constraint sytem
func (cs *r1CS) Compile() (frontend.CompiledConstraintSystem, error) {

	// wires = public wires  | secret wires | internal wires

	// setting up the result
	res := compiled.R1CS{
		CS:          cs.CS,
		Constraints: cs.Constraints,
	}
	res.NbPublicVariables = len(cs.Public)
	res.NbSecretVariables = len(cs.Secret)

	// for Logs, DebugInfo and hints the only thing that will change
	// is that ID of the wires will be offseted to take into account the final wire vector ordering
	// that is: public wires  | secret wires | internal wires

	// offset variable ID depeneding on visibility
	shiftVID := func(oldID int, visibility schema.Visibility) int {
		switch visibility {
		case schema.Internal:
			return oldID + res.NbPublicVariables + res.NbSecretVariables
		case schema.Public:
			return oldID
		case schema.Secret:
			return oldID + res.NbPublicVariables
		}
		return oldID
	}

	// we just need to offset our ids, such that wires = [ public wires  | secret wires | internal wires ]
	offsetIDs := func(l compiled.LinearExpression) {
		for j := 0; j < len(l); j++ {
			_, vID, visibility := l[j].Unpack()
			l[j].SetWireID(shiftVID(vID, visibility))
		}
	}

	for i := 0; i < len(res.Constraints); i++ {
		offsetIDs(res.Constraints[i].L.LinExp)
		offsetIDs(res.Constraints[i].R.LinExp)
		offsetIDs(res.Constraints[i].O.LinExp)
	}

	// we need to offset the ids in the hints
	shiftedMap := make(map[int]*compiled.Hint)

	// we need to offset the ids in the hints
HINTLOOP:
	for _, hint := range cs.MHints {
		ws := make([]int, len(hint.Wires))
		// we set for all outputs in shiftedMap. If one shifted output
		// is in shiftedMap, then all are
		for i, vID := range hint.Wires {
			ws[i] = shiftVID(vID, schema.Internal)
			if _, ok := shiftedMap[ws[i]]; i == 0 && ok {
				continue HINTLOOP
			}
		}
		inputs := make([]interface{}, len(hint.Inputs))
		copy(inputs, hint.Inputs)
		for j := 0; j < len(inputs); j++ {
			switch t := inputs[j].(type) {
			case compiled.Variable:
				tmp := make(compiled.LinearExpression, len(t.LinExp))
				copy(tmp, t.LinExp)
				offsetIDs(tmp)
				inputs[j] = tmp
			case compiled.LinearExpression:
				tmp := make(compiled.LinearExpression, len(t))
				copy(tmp, t)
				offsetIDs(tmp)
				inputs[j] = tmp
			default:
				inputs[j] = t
			}
		}
		ch := &compiled.Hint{ID: hint.ID, Inputs: inputs, Wires: ws}
		for _, vID := range ws {
			shiftedMap[vID] = ch
		}
	}
	res.MHints = shiftedMap

	// we need to offset the ids in Logs & DebugInfo
	for i := 0; i < len(cs.Logs); i++ {

		for j := 0; j < len(res.Logs[i].ToResolve); j++ {
			_, vID, visibility := res.Logs[i].ToResolve[j].Unpack()
			res.Logs[i].ToResolve[j].SetWireID(shiftVID(vID, visibility))
		}
	}
	for i := 0; i < len(cs.DebugInfo); i++ {
		for j := 0; j < len(res.DebugInfo[i].ToResolve); j++ {
			_, vID, visibility := res.DebugInfo[i].ToResolve[j].Unpack()
			res.DebugInfo[i].ToResolve[j].SetWireID(shiftVID(vID, visibility))
		}
	}

	// build levels
	start := time.Now()
	res.Levels = buildLevels(res)
	log.Println("buildLevels took", time.Since(start).Milliseconds(), "ms")

	switch cs.CurveID {
	case ecc.BLS12_377:
		return bls12377r1cs.NewR1CS(res, cs.Coeffs), nil
	case ecc.BLS12_381:
		return bls12381r1cs.NewR1CS(res, cs.Coeffs), nil
	case ecc.BN254:
		return bn254r1cs.NewR1CS(res, cs.Coeffs), nil
	case ecc.BW6_761:
		return bw6761r1cs.NewR1CS(res, cs.Coeffs), nil
	case ecc.BW6_633:
		return bw6633r1cs.NewR1CS(res, cs.Coeffs), nil
	case ecc.BLS24_315:
		return bls24315r1cs.NewR1CS(res, cs.Coeffs), nil
	default:
		panic("not implemtented")
	}
}

func buildLevels(ccs compiled.R1CS) []dag.Level {
	// Build DAG + levels
	graph := dag.New(len(ccs.Constraints))

	nbInputs := ccs.NbPublicVariables + ccs.NbSecretVariables

	mWires := make(map[int]int, ccs.NbInternalVariables) // at which node we resolved which wire.
	dependencies := make([]int, 0, len(ccs.Constraints)) // collect constraint dependencies

	for cID, c := range ccs.Constraints {
		// clear dependencies
		dependencies = dependencies[:0]

		// nodeID == cID here
		nodeID := graph.AddNode(dag.Node{CID: cID})
		if debug.Debug {
			// TODO @gbotrel too hacky.
			if nodeID != cID {
				panic("nodeID doesn't match constraint ID")
			}
		}
		weight := 0
		hasCustomHint := false

		processLE := func(l compiled.LinearExpression) {
			for _, t := range l {
				wID := t.WireID()
				weight++ // we need to account for linear expression eval, that's the min cost
				if wID < nbInputs {
					// it's a input, we ignore it
					continue
				}

				// if we know a which constraint solves this wire, then it's a dependency
				n, ok := mWires[wID]
				if ok {
					if n != nodeID { // can happen with hints...
						dependencies = append(dependencies, n)
					}
					continue
				}

				// check if it's a hint and mark all the output wires
				if h, ok := ccs.MHints[wID]; ok {
					weight += (len(h.Inputs) + len(h.Wires)) * 10
					if !(h.ID == hint.IsZero.UUID() || h.ID == hint.IthBit.UUID()) {
						// it's an unknown hint, we max the weight
						hasCustomHint = true
					}
					for _, hwid := range h.Wires {
						mWires[hwid] = nodeID
					}
					continue
				}

				// mark this wire solved by current node
				weight += 20 // may have a division
				mWires[wID] = nodeID
			}
		}

		processLE(c.L.LinExp)
		processLE(c.R.LinExp)
		processLE(c.O.LinExp)

		// note: if len(dependencies) == 0 --> it's an entry node
		if len(dependencies) != 0 {
			graph.AddEdges(nodeID, dependencies)
		} else {
			// it's an entry node
			graph.MarkEntryNode(nodeID)
		}
		graph.Nodes[nodeID].Weight = weight
		graph.Nodes[nodeID].HasCustomHint = hasCustomHint

	}

	return graph.Levels()
}

func (cs *r1CS) SetSchema(s *schema.Schema) {
	cs.Schema = s
}
