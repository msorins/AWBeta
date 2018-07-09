package solvers

import (
	"wit"
)

type unknownCourierSolver struct {
	awb                string
	url                string
	Statuses           []AWbFanCourierCheckpoint
	LastSolverResponse SolverResponse
}


func UnknownFanCourierSolverBuilder(awb string, entities map[string][]wit.WitEntity) ISolver{
	awbSolver := unknownCourierSolver{}
	awbSolver.awb = awb
	return &awbSolver
}

func (solver *unknownCourierSolver) updateStatuses() SolverResponse  {
	solver.LastSolverResponse = SOLVER_AWB_INCORRECT
	return SOLVER_AWB_INCORRECT
}

func (awbsolver *unknownCourierSolver) GetStatuses() ([]string, SolverResponse) {
	results := []string{}
	results = append(results, "Unknown: Could not found and AWB in your message, please specify one")

	return results, SOLVER_AWB_INCORRECT
}

func (awbsolver *unknownCourierSolver) GetLastStatus() ([]string, SolverResponse) {
	results := []string{}
	results = append(results, "Unknown: Could not found and AWB in your message, please specify one")

	return results, SOLVER_AWB_INCORRECT
}

func (awbsolver *unknownCourierSolver) GetLastSolverResponse() SolverResponse {
	return SOLVER_AWB_INCORRECT
}

func (awbsolver *unknownCourierSolver) GetAwb() string {
	return awbsolver.awb
}
