package solvers

type SolverResponse int

const (
	SOLVER_OK            SolverResponse = iota // == 0
	SOLVER_AWB_INCORRECT SolverResponse = iota // == 1
	SOLVER_BAD_REQUEST   SolverResponse = iota // == 2
	SOLVER_CACHED		 SolverResponse = iota // == 3
)
