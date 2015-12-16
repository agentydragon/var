package bitvecsat

import (
	"github.com/MichalPokorny/var/sat"
)

type PlusConstrain struct {
	AIndex int
	BIndex int
	SumIndex int

	carryIndex int
	// A *Vector
	// B *Vector
	// Sum *Vector

	// carry *Vector
}

func AddBitSumClause(a int, b int, cIn int, sum int) []sat.Clause {
	return []sat.Clause{
		sat.NewClause(true, a, true, b, true, cIn, false, sum),
		sat.NewClause(true, a, true, b, false, cIn, true, sum),
		sat.NewClause(true, a, false, b, true, cIn, true, sum),
		sat.NewClause(true, a, false, b, false, cIn, false, sum),
		sat.NewClause(false, a, true, b, true, cIn, true, sum),
		sat.NewClause(false, a, true, b, false, cIn, false, sum),
		sat.NewClause(false, a, false, b, true, cIn, false, sum),
		sat.NewClause(false, a, false, b, false, cIn, true, sum),
	}
}

func AddBitCarryClause(a int, b int, cIn int, cOut int) []sat.Clause {
	return []sat.Clause{
		sat.NewClause(true, a, true, b, true, cIn, false, cOut),
		sat.NewClause(true, a, true, b, false, cIn, false, cOut),
		sat.NewClause(true, a, false, b, true, cIn, false, cOut),
		sat.NewClause(true, a, false, b, false, cIn, true, cOut),
		sat.NewClause(false, a, true, b, true, cIn, false, cOut),
		sat.NewClause(false, a, true, b, false, cIn, true, cOut),
		sat.NewClause(false, a, false, b, true, cIn, true, cOut),
		sat.NewClause(false, a, false, b, false, cIn, true, cOut),
	}
}

func (constrain PlusConstrain) Materialize(problem *Problem) []sat.Clause {
	// returns sat.formula
	// all vectors have proper width (8)

	a := problem.Vectors[constrain.AIndex]
	b := problem.Vectors[constrain.BIndex]
	carry := problem.Vectors[constrain.carryIndex]
	sum := problem.Vectors[constrain.SumIndex]

	width := a.Width()

	if (width != b.Width()) || (width != sum.Width()) || (width != carry.Width()) {
		panic("unequal bit widths")
	}

	clauses := make([]sat.Clause, 0) // TODO: exact width
	clauses = append(clauses, sat.NewClause(false, carry.SatVarIndices[0]))
	for i := 0; i < width; i++ {
		aIn := a.SatVarIndices[i]
		bIn := b.SatVarIndices[i]
		carryIn := carry.SatVarIndices[i]
		sumOut := sum.SatVarIndices[i]
		clauses = append(clauses, AddBitSumClause(aIn, bIn, carryIn, sumOut)...)

		if i + 1 < width {
			carryOut := carry.SatVarIndices[i + 1]
			clauses = append(clauses, AddBitCarryClause(aIn, bIn, carryIn, carryOut)...)
		}
	}
	return clauses
}

func (constrain PlusConstrain) AddToProblem(problem *Problem) {
	constrain.carryIndex = problem.AddNewVector()
	problem.AddNewConstrain(constrain)
}