package solvers


type ISolver interface {
	updateStatuses() SolverResponse

	GetStatuses() ([]string, SolverResponse)
	GetLastStatus() ([]string, SolverResponse)
}
