package dag

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDAGReduction(t *testing.T) {
	assert := require.New(t)

	// we start with
	// ┌────A
	// │    │
	// │    ▼
	// │    B
	// │    │
	// │    ▼
	// └───►C
	const (
		A Node = iota
		B
		C
		nbNodes
	)
	dag := New(int(nbNodes))
	// virtually adds A and B
	a := dag.AddNode(A)
	b := dag.AddNode(B)

	dag.AddEdges(b, []int{a})

	// virtuall adds C
	c := dag.AddNode(C)
	dag.AddEdges(c, []int{a, b})

	// we should get
	// 		A
	// 		│
	// 		▼
	// 		B
	// 		│
	// 		▼
	// 		C
	assert.Equal(0, len(dag.parents[a]))
	assert.Equal(1, len(dag.parents[b]))
	assert.Equal(1, len(dag.parents[c]))

	assert.Equal(a, dag.parents[b][0])
	assert.Equal(b, dag.parents[c][0])

	assert.Equal(1, len(dag.children[a]))
	assert.Equal(1, len(dag.children[b]))
	assert.Equal(0, len(dag.children[c]))

	assert.Equal(b, dag.children[a][0])
	assert.Equal(c, dag.children[b][0])

}

func TestDAGReductionFork(t *testing.T) {
	assert := require.New(t)

	// we start with this
	// ┌─────D◄───┐
	// │     ▲    │
	// │     │    │
	// │ A   B    C
	// │ │   │    │
	// │ │   ▼    │
	// │ └──►E ◄──┘
	// │     ▲
	// └─────┘
	const (
		A Node = iota
		B
		C
		D
		E
		nbNodes
	)

	dag := New(int(nbNodes))
	// virtually adds A,B,C,D
	a := dag.AddNode(A)
	b := dag.AddNode(B)
	c := dag.AddNode(C)
	d := dag.AddNode(D)

	dag.AddEdges(d, []int{b, c})

	// virtuall adds E
	e := dag.AddNode(E)
	dag.AddEdges(e, []int{a, b, c, d})

	// we should get
	// A     B     C
	// │     │     │
	// │     ▼     │
	// │     D ◄───┘
	// │     │
	// │     ▼
	// └────►E
	assert.Equal(0, len(dag.parents[a]))
	assert.Equal(0, len(dag.parents[b]))
	assert.Equal(0, len(dag.parents[c]))
	assert.Equal(2, len(dag.parents[d]))
	assert.Equal(2, len(dag.parents[e]))

	assert.Equal(c, dag.parents[d][0])
	assert.Equal(b, dag.parents[d][1])

	assert.Equal(d, dag.parents[e][0])
	assert.Equal(a, dag.parents[e][1])

	assert.Equal(1, len(dag.children[a]))
	assert.Equal(1, len(dag.children[b]))
	assert.Equal(1, len(dag.children[c]))
	assert.Equal(1, len(dag.children[d]))
	assert.Equal(0, len(dag.children[e]))

	assert.Equal(e, dag.children[a][0])
	assert.Equal(d, dag.children[b][0])
	assert.Equal(d, dag.children[c][0])
	assert.Equal(e, dag.children[d][0])

	// Check that levels are coherent
	levels := dag.Levels(nil)

	// we should have 3 levels:
	// [A,B,C]
	// [D]
	// [E]
	assert.Equal(3, len(levels))
	assert.Equal(3, len(levels[0]))
	assert.Equal(1, len(levels[1]))
	assert.Equal(1, len(levels[2]))

	// level 0
	assert.Equal(A, levels[0][0])
	assert.Equal(B, levels[0][1])
	assert.Equal(C, levels[0][2])

	// level 1
	assert.Equal(D, levels[1][0])

	// level 2
	assert.Equal(E, levels[2][0])
}

func BenchmarkDAGReduction(b *testing.B) {
	rand.Seed(42)
	const nbNodes = 100000
	parents := make([]int, 0, nbNodes)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		dag := New(nbNodes)
		for j := 0; j < nbNodes/1000; j++ {
			dag.AddNode(Node(j)) // initial nodes
		}
		b.StartTimer()
		for j := nbNodes / 1000; j < nbNodes; j++ {
			parents = parents[:0]
			for k := 0; k < 10; k++ {
				parents = append(parents, rand.Intn(j-1))
			}
			dag.AddNode(Node(j))
			dag.AddEdges(j, parents)
		}
		_ = dag.Levels(nil)
	}
}
