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
	results = append(results, "Could not found and AWB and link it with a courier company")

	return results, SOLVER_AWB_INCORRECT
}

func (awbsolver *unknownCourierSolver) GetLastStatus() ([]string, SolverResponse) {
	results := []string{}
	results = append(results, "Could not found and AWB and link it with a courier company")

	return results, SOLVER_AWB_INCORRECT
}


