package solvers

const TIME_BETWEEN_REQUEST_MIN = 20

type ISolver interface {
	updateStatuses() SolverResponse

	GetStatuses() ([]string, SolverResponse)
	GetLastStatus() ([]string, SolverResponse)
	GetAwb() string
}
