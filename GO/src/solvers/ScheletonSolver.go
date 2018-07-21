package solvers

type ScheletonSolver struct {
	awb string
}

func Scheleton(awb string) ISolver {
	sds := ScheletonSolver{awb}
	return &sds
}

func (solver *ScheletonSolver) updateStatuses() SolverResponse {
	return SOLVER_NOT_IMPLEMENTED
}

func (solver *ScheletonSolver) GetStatuses() ([]string, SolverResponse) {
	return nil, SOLVER_NOT_IMPLEMENTED
}

func (solver *ScheletonSolver) GetLastStatus() ([]string, SolverResponse) {
	return nil, SOLVER_NOT_IMPLEMENTED
}

func (solver *ScheletonSolver) GetAwb() string {
	return ""
}
