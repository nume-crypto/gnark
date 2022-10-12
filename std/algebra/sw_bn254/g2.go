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

package sw_bn254

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/fields_bn254"
)

// G2Jac point in Jacobian coords
type G2Jac struct {
	X, Y, Z fields_bn254.E2
}

// G2Affine point in affine coords
type G2Affine struct {
	X, Y fields_bn254.E2
}

// Neg outputs -p
func (p *G2Affine) Neg(api frontend.API, p1 G2Affine) *G2Affine {
	p.Y.Neg(api, p1.Y)
	p.X = p1.X
	return p
}
